package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/schollz/progressbar/v3"
)

// processResults handles incoming results from workers for both initial and re-scan phases.
// Added corsCheck bool parameter to know if CORS details should be expected/printed.
func processResults(
	results <-chan scanResult,
	resultWg *sync.WaitGroup,
	bar *progressbar.ProgressBar,
	quiet bool,
	corsCheck bool, // <-- Added parameter
	outputDir string,
	successfulScans *int64,
	failedScans *int64,
	statusCounts map[int]int64,
	statusCountsMutex *sync.Mutex,
	trackFailures bool, // Flag to control adding to failedTargets list [source: 43]
) {
	defer resultWg.Done()
	// Get file paths using constants from output.go
	logPath := filepath.Join(outputDir, logFileName)
	existPath := filepath.Join(outputDir, existFileName)
	invalidPath := filepath.Join(outputDir, invalidFileName)
	corsVulnPath := filepath.Join(outputDir, corsVulnerableFileName) // <-- Added CORS vuln file path

	for res := range results {
		// Safely increment progress bar for each processed result
		bar.Add(1) // [source: 43]

		// Get status code description, handle unknown codes
		desc, descOk := statusCodes[res.statusCode]
		if res.err == nil && !descOk {
			desc = "(Unknown Status Code)"
		} else if res.err != nil {
			desc = "" // No description needed if there's an error
		}

		// Display primary result (status code or error) unless in quiet mode
		colorPrint(res.target, res.statusCode, desc, res.err, quiet, res.isRescan) // Call ui function [source: 43]

		// --- Handle CORS Result Processing (if check was enabled) ---
		if corsCheck && res.cors != nil { // Check if CORS check was performed (res.cors is not nil)
			if res.cors.err != nil {
				// Log CORS check error separately
				logMsg := fmt.Sprintf("[!] CORS Check Error for %s: %v", res.cors.target, res.cors.err)
				appendToFile(logPath, logMsg)
				// Print CORS error distinctly, even in quiet mode as it's an operational error
				// Use specific colors defined in globals.go
				fmt.Printf("%s[CORS ERR]%s %s -> %s%v%s\n", ColorCorsErr, ColorReset, res.cors.target, ColorError, res.cors.err, ColorReset)

			} else if res.cors.vulnerable {
				// Log CORS vulnerability to dedicated file and main log
				vulnMsg := fmt.Sprintf("%s (%s)", res.cors.target, res.cors.details)
				appendToFile(corsVulnPath, vulnMsg)                                    // Save to cors_vulnerable.txt
				appendToFile(logPath, fmt.Sprintf("[!] CORS VULNERABLE: %s", vulnMsg)) // Log clearly

				// Print vulnerability indication, even in quiet mode as it's a finding
				// Use specific colors defined in globals.go
				fmt.Printf("%s[CORS VULN]%s %s -> %s%s%s\n", ColorCorsVuln, ColorReset, res.cors.target, ColorWarning, res.cors.details, ColorReset)

			} else {
				// Optionally log non-vulnerable CORS checks (can be verbose)
				// appendToFile(logPath, fmt.Sprintf("[✓] CORS OK: %s (%s)", res.cors.target, res.cors.details))

				// Print non-vulnerable status only if NOT in quiet mode
				if !quiet {
					// Indicate CORS check was OK inline or below main result line
					// Using a less prominent color like Info or just default reset
					// fmt.Printf("%s      ↳ CORS OK: %s%s\n", ColorInfo, res.cors.details, ColorReset)
				}
			}
		}
		// --- End CORS Result Handling ---

		// Process primary scan result logic: update counters, manage failures, write files
		if res.err != nil { // Handle Primary Scan Failure [source: 44]
			logMsg := ""
			if !res.isRescan {
				// --- Initial Scan Failure ---
				atomic.AddInt64(failedScans, 1) // Increment overall fail count
				if trackFailures {              // Only track initial failures for potential rescan
					failedTargetsMutex.Lock()
					failedTargets = append(failedTargets, res.target) // Add to global list
					failedTargetsMutex.Unlock()
					// Write to the invalid list only on the first failure
					appendToFile(invalidPath, res.target) // [source: 44]
				}
				logMsg = fmt.Sprintf("[!] FAIL %s -> ERROR: %v", res.target, res.err)
			} else {
				// --- Re-scan Failure ---
				// Failure persists after rescan
				logMsg = fmt.Sprintf("[!!] RESCAN FAIL %s -> ERROR: %v", res.target, res.err) // [source: 45]
				// Do not increment failedScans again, it was already counted during initial fail
			}
			appendToFile(logPath, logMsg) // Log the failure

		} else { // Handle Primary Scan Success
			logMsg := ""
			if res.isRescan {
				// --- Success during Re-scan ---
				// Target failed initially but succeeded on rescan. Adjust overall counts.
				atomic.AddInt64(failedScans, -1)                                                          // Decrease overall fail count
				logMsg = fmt.Sprintf("[✓✓] RESCAN SUCCESS %s -> %d %s", res.target, res.statusCode, desc) // [source: 45]
			} else {
				// --- Success during Initial Scan ---
				logMsg = fmt.Sprintf("[✓] SUCCESS %s -> %d %s", res.target, res.statusCode, desc) // [source: 45]
			}
			// Increment overall success count regardless of initial/rescan success
			atomic.AddInt64(successfulScans, 1) // [source: 45]

			// Update status code counts safely
			statusCountsMutex.Lock()
			statusCounts[res.statusCode]++
			statusCountsMutex.Unlock() // [source: 45]

			// Write to common success files (exist list and log)
			// Only write to ip_exist.txt if it succeeded at least once (initial or rescan)
			appendToFile(existPath, res.target) // [source: 46]
			appendToFile(logPath, logMsg)

			// Write to specific status code file based on category
			catDigit := res.statusCode / 100
			catName, catOk := statusCategories[catDigit] // Assumes statusCategories is populated globally [source: 46]
			if !catOk {
				catName = "unknown_category"                                // Handle unexpected category
				os.MkdirAll(filepath.Join(outputDir, catName), os.ModePerm) // Create dir if needed
			}

			// Check if the status code itself is known (from globals map)
			_, codeKnown := statusCodes[res.statusCode] // [source: 46]

			if codeKnown && catOk { // Write to status file if code and category are known
				targetFile := filepath.Join(outputDir, catName, fmt.Sprintf("%d.txt", res.statusCode))
				appendToFile(targetFile, res.target) // [source: 46]
			} else { // Log if code was unknown or category was unknown even on success
				unknownLogMsg := fmt.Sprintf("[?] UNKNOWN STATUS %s -> %d (Desc Known: %t, Cat Known: %t)", res.target, res.statusCode, codeKnown, catOk)
				appendToFile(logPath, unknownLogMsg)
				// Optionally write to a dedicated "unknown_status.txt" file
				unknownFile := filepath.Join(outputDir, "unknown_status.txt")
				appendToFile(unknownFile, fmt.Sprintf("%s -> %d", res.target, res.statusCode))
			}
		}
	}
}
