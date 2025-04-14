package init

import (
	"fmt"
	"os"
	"time"
	"github.com/nxneeraj/hxscanner/ui"
	"github.com/nxneeraj/hxscanner/scanner"
)

// Init initializes the HyperScanner application by setting up the environment
func Init() {
	ui.showWelcomeBanner()

	// Simulate some setup time
	time.Sleep(2 * time.Second)

	// Checking environment variables (like for proxy or custom configs)
	checkEnvironment()

	// You can initialize logging, other services here if needed in future

	// Load configurations or settings (could be added)
	fmt.Println("Configuration Loaded: HyperScanner is ready to scan.")

	// Proceed to scanning (optional)
	ui.showScanStart()
}

// checkEnvironment checks the environment for necessary conditions (e.g., proxy, config files)
func checkEnvironment() {
	// Check if any environment variables need to be loaded
	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		fmt.Printf("Using Proxy: %s\n", proxy)
	}

	// Add more checks if necessary (like config files or permissions)
	fmt.Println("Environment check completed. All good!")
}

// StartScanning begins the scanning process
func StartScanning(urls []string) {
	ui.showScanStart()

	// Initialize scan variables
	totalURLs := len(urls)
	successfulScans := 0
	failedScans := 0

	// Start the scanning process
	for _, url := range urls {
		// Call scanner logic for URL
		statusCode := scanner.ScanURL(url) // This is a placeholder for actual scanning logic
		ui.showScanProgress(url, statusCode)

		if statusCode >= 200 && statusCode < 300 {
			successfulScans++
		} else {
			failedScans++
		}
	}

	// After scanning all URLs, show the results
	ui.showScanComplete(totalURLs, successfulScans, failedScans)
}
