package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// corsCheckResult holds the outcome of a CORS check
type corsCheckResult struct {
	target     string
	vulnerable bool
	details    string
	err        error
}

// checkCORS performs the actual CORS vulnerability check
// It sends an OPTIONS request with a potentially malicious Origin.
func checkCORS(target string, client *http.Client) corsCheckResult {
	result := corsCheckResult{target: target}
	// Use a distinct origin for testing that's unlikely to be whitelisted by chance
	checkOrigin := "https://evil-cors-test.com"

	// Ensure the target has a scheme (http/https)
	urlToScan := target
	if !strings.Contains(target, "://") {
		urlToScan = "http://" + target // Default to http if no scheme
	}

	// Validate the URL structure again before the OPTIONS request
	parsedURL, err := url.ParseRequestURI(urlToScan)
	if err != nil {
		result.err = fmt.Errorf("invalid target format for CORS check '%s': %w", target, err)
		return result
	}
	urlToScan = parsedURL.String()

	// --- Create OPTIONS request ---
	req, err := http.NewRequest("OPTIONS", urlToScan, nil)
	if err != nil {
		result.err = fmt.Errorf("failed to create OPTIONS request for %s: %w", urlToScan, err)
		return result
	}

	// --- Set Headers for CORS Check ---
	// Use a specific user agent for CORS checks
	req.Header.Set("User-Agent", "HyperScanner/1.4+CORSCheck")
	req.Header.Set("Origin", checkOrigin)
	// Common methods often allowed via CORS
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "X-Requested-With") // Common header

	// Use a dedicated client for CORS checks to avoid interference and manage settings like redirects
	corsClient := &http.Client{
		Timeout: client.Timeout, // Use the same timeout or a specific one for CORS
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          10, // Less need for connection pooling for OPTIONS
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,  // OPTIONS is usually a one-off check per target
			ForceAttemptHTTP2:     false, // HTTP/1.1 is sufficient for OPTIONS
			// Add TLSClientConfig if needed (e.g., for InsecureSkipVerify based on a main flag)
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Example if needed
		},
		// Explicitly prevent following redirects for OPTIONS requests
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// --- Perform the OPTIONS request ---
	resp, err := corsClient.Do(req)
	if err != nil {
		// Network errors (timeouts, connection refused, DNS issues) are not CORS vulns per se
		result.err = fmt.Errorf("OPTIONS request network error for %s: %w", urlToScan, err)
		return result
	}
	defer resp.Body.Close()

	// --- Analyze Response Headers ---
	// Standard CORS headers
	acaoHeader := resp.Header.Get("Access-Control-Allow-Origin")
	acacHeader := resp.Header.Get("Access-Control-Allow-Credentials")

	// Check 1: Reflects arbitrary Origin?
	if acaoHeader == checkOrigin {
		result.vulnerable = true
		result.details = fmt.Sprintf("Reflects Origin: ACAO='%s'", acaoHeader)
		// Check 1a: If reflects origin AND allows credentials, it's highly sensitive
		if acacHeader == "true" {
			result.details += ", ACAC='true' (CRITICAL)"
		}
		return result // Found vulnerability
	}

	// Check 2: Allows wildcard Origin?
	if acaoHeader == "*" {
		result.vulnerable = true
		result.details = "Wildcard Origin: ACAO='*'"
		// Check 2a: If wildcard AND allows credentials (SPEC VIOLATION but servers might do it)
		// Browsers should block this, but indicates server misconfiguration.
		if acacHeader == "true" {
			result.details += ", ACAC='true' (Misconfiguration/Severe)"
		}
		return result // Found vulnerability
	}

	// Check 3: Allows "null" Origin (less common but possible vector)
	if acaoHeader == "null" {
		// Vulnerability depends on context, but worth noting
		result.vulnerable = true // Treat 'null' as potentially problematic
		result.details = "Null Origin: ACAO='null'"
		if acacHeader == "true" {
			result.details += ", ACAC='true' (Potentially Problematic)"
		}
		return result
	}

	// If none of the above conditions were met
	result.vulnerable = false
	if acaoHeader == "" {
		result.details = "ACAO header missing or empty"
	} else {
		result.details = fmt.Sprintf("ACAO header: '%s' (Not matching test origin, wildcard, or null)", acaoHeader)
	}

	return result
}
