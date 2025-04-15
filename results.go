package main

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/schollz/progressbar/v3"
)

// processResults handles incoming results from workers for both initial and re-scan phases.
func processResults(
	results <-chan scanResult,
	resultWg *sync.WaitGroup,
	bar *progressbar.ProgressBar, // Pass the appropriate bar
	quiet bool,
	outputDir string,
	successfulScans *int64,
	failedScans *int64, // Pointer needed to adjust based on rescan success
	statusCounts map[int]int64,
	statusCountsMutex *sync.Mutex,
	trackFailures bool, // Flag to control adding to failedTargets list
) {
	defer resultWg.Done()
	// Get file paths using constants from output.go
	logPath := filepath.Join(outputDir, logFileName)
	existPath := filepath.Join(outputDir, existFileName)
	invalidPath := filepath.Join(outputDir, invalidFileName)

	for res := range results {
		// Safely increment progress bar for each processed result
		bar.Add(1)

		// Display result unless in quiet mode
		desc, _ := statusCodes[res.statusCode] // Assumes statusCodes is populated globally
		colorPrint(res.target, res.statusCode, desc, res.err, quiet, res.isRescan) // Call ui function

		// Process result logic: update counters, manage failures, write files
		if res.err != nil { // Handle Failure
			logMsg := ""
			if !res.isRescan {
				// --- Initial Scan Failure ---
				atomic.AddInt64(failedScans, 1) // Increment overall fail count
				if trackFailures {              // Only track initial failures for potential rescan
					failedTargetsMutex.Lock()
					failedTargets = append(failedTargets, res.target) // Add to global list
					failedTargetsMutex.Unlock()
					// Write to the invalid list only on the first failure
					appendToFile(invalidPath, res.target)
				}
				logMsg = fmt.Sprintf("[!] %s -> ERROR: %v", res.target, res.err)
			} else {
				// --- Re-scan Failure ---
				logMsg = fmt.Sprintf("[!! RESCAN FAIL] %s -> ERROR: %v", res.target, res.err)
				// Do not increment failedScans again, it was already counted
			}
			// Log the failure (initial or persistent)
			appendToFile(logPath, logMsg)

		} else { // Handle Success
			logMsg := ""
			if res.isRescan {
				// --- Success during Re-scan ---
				// Adjust overall counters: decrease failures, increase successes
				atomic.AddInt64(failedScans, -1)
				logMsg = fmt.Sprintf("[✓✓ RESCAN SUCCESS] %s -> %d %s", res.target, res.statusCode, desc)
			} else {
				// --- Success during Initial Scan ---
				logMsg = fmt.Sprintf("[✓] %s -> %d %s", res.target, res.statusCode, desc)
			}
			// Increment overall success count regardless of initial/rescan
			atomic.AddInt64(successfulScans, 1)

			// Update status code counts safely
			statusCountsMutex.Lock()
			statusCounts[res.statusCode]++
			statusCountsMutex.Unlock()

			// Write to common success files (exist list and log)
			appendToFile(existPath, res.target) // Target succeeded at least once
			appendToFile(logPath, logMsg)

			// Write to specific status code file
			catDigit := res.statusCode / 100
			catName, ok := statusCategories[catDigit] // Assumes statusCategories is populated globally
			if !ok {
				catName = "unknown" // Should ideally not happen if maps are complete
			}
			_, codeKnown := statusCodes[res.statusCode] // Assumes statusCodes is populated globally

			if codeKnown && catName != "unknown" {
				targetFile := filepath.Join(outputDir, catName, fmt.Sprintf("%d.txt", res.statusCode))
				appendToFile(targetFile, res.target)
			} else if !codeKnown { // Log if code was unknown even on success
				// This case might indicate a server returning a non-standard code
				appendToFile(logPath, fmt.Sprintf("[?] %s -> %d (Unknown Desc/Cat on Success)", res.target, res.statusCode))
			}
		}
	}
}
