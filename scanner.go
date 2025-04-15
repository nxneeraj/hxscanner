package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// setupHTTPClient creates the shared HTTP client
func setupHTTPClient(timeout time.Duration, workers int) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        workers * 2,
			MaxIdleConnsPerHost: 10,           // Default is 2, increasing might help
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,       // Keep connections alive
			ForceAttemptHTTP2:   true,        // Try HTTP/2
			// Add TLSClientConfig if needed, e.g., for InsecureSkipVerify
		},
		// Important: Prevent following redirects automatically to capture 3xx codes
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}


// scanTarget performs the HTTP request for various target types
func scanTarget(target string, client *http.Client) (int, error) {
	urlToScan := target
	// If target doesn't contain "://", assume it's an IP or domain and prepend http://
	if !strings.Contains(target, "://") {
		urlToScan = "http://" + target
	}

	// Validate the final URL structure before making the request
	parsedURL, err := url.ParseRequestURI(urlToScan)
	if err != nil {
		// Return an error if the URL format is fundamentally invalid
		return 0, fmt.Errorf("invalid target format '%s': %w", target, err)
	}
	// Use the parsed and potentially cleaned-up URL string
	urlToScan = parsedURL.String()

	// Create request (defaulting to GET)
	req, err := http.NewRequest("GET", urlToScan, nil)
	if err != nil {
		// This error is less likely if url.Parse succeeded, but check anyway
		return 0, fmt.Errorf("failed to create request for %s: %w", urlToScan, err)
	}
	// Set a user agent
	req.Header.Set("User-Agent", "HyperScanner/1.4")

	// Perform the request using the provided client
	resp, err := client.Do(req)
	if err != nil {
		// Return network errors, timeouts, DNS errors, etc.
		return 0, err
	}
	// Ensure the response body is closed even if we don't read it
	defer resp.Body.Close()

	// Return the status code
	return resp.StatusCode, nil
}

// worker executes scan jobs received from the jobs channel
func worker(id int, wg *sync.WaitGroup, client *http.Client, jobs <-chan string, results chan<- scanResult, isRescan bool) {
	defer wg.Done() // Signal completion when channel is closed and loop finishes
	for target := range jobs {
		if target == "" { // Skip empty lines, just in case
			continue
		}
		status, err := scanTarget(target, client) // Perform the scan
		// Send the result back, including the original target and rescan status
		results <- scanResult{target: target, statusCode: status, err: err, isRescan: isRescan}
	}
}
