package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// setupHTTPClient creates the shared HTTP client for initial GET requests
func setupHTTPClient(timeout time.Duration, workers int) *http.Client {
	// Transport settings optimized for potentially many connections
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          workers * 2,      // Allow more idle connections globally
		MaxIdleConnsPerHost:   10,               // Allow more idle connections per host
		IdleConnTimeout:       90 * time.Second, // Keep idle connections longer
		TLSHandshakeTimeout:   timeout,          // Use main timeout for handshake
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false, // Keep connections alive for efficiency
		ForceAttemptHTTP2:     true,  // Try HTTP/2 [source: 47]
		// Add TLSClientConfig if needed globally, e.g., for InsecureSkipVerify
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
		// Important: Prevent following redirects automatically to capture 3xx codes correctly
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Return the 3xx response itself
		}, // [source: 47]
	}
}

// scanTarget performs the primary HTTP GET request for a target
func scanTarget(target string, client *http.Client) (int, error) {
	urlToScan := target
	// Prepend http:// if no scheme is present (handles IPs and domains)
	if !strings.Contains(target, "://") {
		// Check if it looks like an IPv6 address that needs brackets
		if strings.Contains(target, ":") && !strings.HasPrefix(target, "[") && !strings.HasSuffix(target, "]") {
			// Basic check, might need refinement for edge cases
			isIPv6 := false
			parts := strings.Split(target, ":")
			if len(parts) > 2 { // Simple heuristic for IPv6
				isIPv6 = true
			}
			// Check if it's likely host:port
			if len(parts) == 2 {
				_, portErr := url.Parse("http://dummy:" + parts[1]) // Check if part after : is a valid port
				if portErr != nil {                                 // If not a valid port, assume IPv6
					isIPv6 = true
				}
			}

			if isIPv6 {
				urlToScan = "http://[" + target + "]"
			} else {
				urlToScan = "http://" + target
			}
		} else {
			urlToScan = "http://" + target
		}
	} // [source: 48]

	// Validate the final URL structure before making the request
	parsedURL, err := url.ParseRequestURI(urlToScan)
	if err != nil {
		// Return an error if the URL format is fundamentally invalid after scheme prepending
		return 0, fmt.Errorf("invalid target format '%s' -> '%s': %w", target, urlToScan, err) // [source: 48]
	}
	urlToScan = parsedURL.String() // Use the validated URL string

	// Create request (defaulting to GET)
	req, err := http.NewRequest("GET", urlToScan, nil)
	if err != nil {
		// This error is less likely if url.Parse succeeded, but check anyway
		return 0, fmt.Errorf("failed to create GET request for %s: %w", urlToScan, err) // [source: 49]
	}
	// Set a distinct user agent for the main scanner
	req.Header.Set("User-Agent", "HyperScanner/1.4") // [source: 49]

	// Perform the request using the provided shared client
	resp, err := client.Do(req)
	if err != nil {
		// Handle network errors, timeouts, DNS errors, connection refused etc.
		// Try to provide more context if possible (e.g., timeout vs connection refused)
		// urlErr, ok := err.(*url.Error)
		// if ok && urlErr.Timeout() {
		//  return 0, fmt.Errorf("timeout reaching %s: %w", urlToScan, err)
		// }
		return 0, fmt.Errorf("request failed for %s: %w", urlToScan, err) // Return wrapped error
	}
	// Ensure the response body is always closed to free up resources
	defer resp.Body.Close() // [source: 49]

	// Return the status code from the response
	return resp.StatusCode, nil // [source: 49]
}

// worker executes scan jobs received from the jobs channel
// Added corsCheck flag parameter
func worker(id int, wg *sync.WaitGroup, client *http.Client, jobs <-chan string, results chan<- scanResult, isRescan bool, corsCheck bool) {
	defer wg.Done() // Signal completion when channel is closed and loop finishes
	for target := range jobs {
		if target == "" { // Skip empty lines [source: 50]
			continue
		}

		status, err := scanTarget(target, client) // Perform the initial GET scan

		// Prepare the basic result struct
		result := scanResult{
			target:     target,
			statusCode: status,
			err:        err,
			isRescan:   isRescan,
			cors:       nil, // Initialize CORS result pointer to nil
		}

		// --- Perform CORS Check if enabled AND initial scan was successful ---
		// Also ensure status code is not 0 (which indicates an error in scanTarget itself)
		if corsCheck && err == nil && status != 0 {
			// Perform the CORS check (function defined in cors.go)
			// Pass the same target and the main client (checkCORS uses its own internal client settings)
			corsResult := checkCORS(target, client)
			result.cors = &corsResult // Store the pointer to the CORS result in the main scanResult
		}
		// --- End CORS Check ---

		// Send the combined result (including potential CORS info) back to the results processor
		results <- result // [source: 50]
	}
}
