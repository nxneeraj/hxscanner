package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	// Progressbar is used via ui.go
)

// runScanPhase executes either the initial scan or the re-scan
// Added corsCheck flag
func runScanPhase(
	targets []string,
	description string,
	isRescan bool,
	client *http.Client,
	workersCount int,
	outputDir string,
	quiet bool,
	corsCheck bool, // <-- Added parameter
	successfulScans *int64,
	failedScans *int64,
	statusCounts map[int]int64,
	statusCountsMutex *sync.Mutex,
) { // [source: 29]
	totalScanTargets := len(targets)
	if totalScanTargets == 0 {
		if isRescan {
			fmt.Printf("%s[*] No targets needed re-scanning.%s\n", ColorInfo, ColorReset)
		}
		return // Nothing to do in this phase
	}

	fmt.Printf("\n%s[*] Starting %s for %d targets...%s\n", ColorInfo, description, totalScanTargets, ColorReset)

	// Setup channels and waitgroup for this scan phase
	jobs := make(chan string, workersCount)
	results := make(chan scanResult, workersCount*2) // Increase buffer slightly for results+cors
	var wg sync.WaitGroup                            // For workers in this phase
	var resultWg sync.WaitGroup                      // For result processor in this phase

	// Start workers for this phase (worker function is in scanner.go)
	// Pass corsCheck flag to worker
	for w := 1; w <= workersCount; w++ {
		wg.Add(1)
		// Pass corsCheck flag here
		go worker(w, &wg, client, jobs, results, isRescan, corsCheck) // <-- Pass corsCheck
	} // [source: 30]

	// Progress bar for this phase (createProgressBar is in ui.go)
	bar := createProgressBar(totalScanTargets, fmt.Sprintf("%s[*] %s%s", ColorInfo, description, ColorReset))

	// Start results processor for this phase (processResults is in results.go)
	resultWg.Add(1)
	// Pass corsCheck flag down to results processor to potentially influence output
	go processResults(results, &resultWg, bar, quiet, corsCheck, outputDir,
		successfulScans, failedScans,
		statusCounts, statusCountsMutex,
		!isRescan, // trackFailures is true only for initial scan (!isRescan)
	)

	// Send jobs for this phase
	for _, target := range targets {
		jobs <- target
	}
	close(jobs) // Done sending jobs for this phase

	// Wait for completion of this phase
	wg.Wait()       // Wait for workers
	close(results)  // Done sending results
	resultWg.Wait() // Wait for result processor [source: 31]

	// Explicitly add a small delay to ensure progress bar finishes drawing
	// especially if result processing finishes very quickly.
	time.Sleep(100 * time.Millisecond)
	bar.Finish() // Cleanly finish progress bar
	fmt.Printf("%s[*] %s phase complete.%s\n", ColorInfo, description, ColorReset)
}

func main() {
	startTime := time.Now()
	printBanner() // From ui.go

	// --- Command Line Flags ---
	targetInput := flag.String("i", "", "Input file containing IPs, Domains, or URLs (one per line)")
	fileInput := flag.String("f", "", "Alias for -i (optional)")
	helpFlag := flag.Bool("h", false, "Show help")
	defaultWorkers := runtime.NumCPU()
	if defaultWorkers < 4 {
		defaultWorkers = 4
	}
	workers := flag.Int("w", defaultWorkers, fmt.Sprintf("Number of concurrent workers (default: %d)", defaultWorkers))
	timeout := flag.Duration("t", 5*time.Second, "HTTP request timeout (e.g., 3s, 10s)")
	quiet := flag.Bool("q", false, "Quiet mode: suppress individual results, show only progress and summary")
	corsCheck := flag.Bool("cors", false, "Perform a basic CORS vulnerability check on successful targets") // <-- Added CORS flag

	flag.Parse() // [source: 32]

	// --- Help Flag Handling ---
	if *helpFlag {
		fmt.Println("Usage: hxscanner [options]")
		fmt.Println("\nScans IPs, Domains, or full URLs from an input file via HTTP/S.")
		fmt.Println("If no scheme (http:// or https://) is provided for a domain/IP, http:// is assumed.")
		fmt.Println("Offers an option to re-scan failed targets and check for CORS misconfigurations.")
		fmt.Println("\nOptions:")
		fmt.Println("  -i <file>     Input file with targets (IPs/Domains/URLs), one per line (required if -f not used)")
		fmt.Println("  -f <file>     Alias for -i")
		fmt.Println("  -w <number>   Number of concurrent scanning workers (default: number of CPU cores)") // [source: 33]
		fmt.Println("  -t <duration> HTTP request timeout (default: 5s)")                                   // [source: 33]
		fmt.Println("  -q            Quiet mode: suppress individual results (except errors/warnings)")     // [source: 33]
		fmt.Println("  --cors        Perform basic CORS vulnerability check on successful targets")         // <-- Added help text
		fmt.Println("  -h            Show this help message")                                               // [source: 33]
		os.Exit(0)
	}

	// --- Input File Validation ---
	if *targetInput == "" && *fileInput == "" {
		fmt.Printf("%sError: No input file provided. Use -i <file> or -f <file>%s\n", ColorError, ColorReset) // [source: 34]
		flag.Usage()
		os.Exit(1)
	}
	targetListPath := *targetInput
	if *fileInput != "" {
		if targetListPath != "" && targetListPath != *fileInput {
			fmt.Printf("%sWarning: Both -i and -f provided, using -f value: %s%s\n", ColorWarning, *fileInput, ColorReset)
		}
		targetListPath = *fileInput
	}

	// --- Read all targets into memory first ---
	fmt.Printf("%s[*] Reading targets from %s...%s\n", ColorInfo, targetListPath, ColorReset)
	initialTargets, err := readTargetsFromFile(targetListPath) // From utils.go
	if err != nil {
		fmt.Printf("%sError reading input file %s: %v%s\n", ColorError, targetListPath, err, ColorReset)
		os.Exit(1)
	}
	totalTargets := len(initialTargets)
	if totalTargets == 0 {
		fmt.Printf("%sWarning: Input file %s appears to be empty or contains no valid targets.%s\n", ColorWarning, targetListPath, ColorReset)
		os.Exit(0)
	}
	fmt.Printf("%s[*] Found %d targets to scan.%s\n", ColorInfo, totalTargets, ColorReset)

	// --- Output Directory Setup ---
	outputDir := strings.TrimSuffix(filepath.Base(targetListPath), filepath.Ext(targetListPath)) + "_output"
	err = createOutputStructure(outputDir) // From output.go
	if err != nil {
		fmt.Printf("%sError creating output structure in %s: %v%s\n", ColorError, outputDir, err, ColorReset)
		os.Exit(1)
	}
	fmt.Printf("%s[*] Output will be saved to: %s%s%s\n", ColorInfo, ColorAccent, outputDir, ColorReset) // [source: 35]

	// --- Shared HTTP Client Setup ---
	// Note: CORS check uses its own client settings within checkCORS for specific needs
	sharedClient := setupHTTPClient(*timeout, *workers) // From scanner.go

	// --- Overall Statistics Setup ---
	var successfulScans int64
	var failedScans int64
	statusCounts := make(map[int]int64)
	var statusCountsMutex sync.Mutex // [source: 35]

	// --- Run Initial Scan ---
	// Pass corsCheck flag value here
	runScanPhase(initialTargets, "Initial Scan", false, /* isRescan = false */
		sharedClient, *workers, outputDir, *quiet, *corsCheck, // <-- Pass corsCheck value
		&successfulScans, &failedScans,
		statusCounts, &statusCountsMutex,
	)

	// --- Initial Summary ---
	printSummary("Initial Scan", startTime, totalTargets, &successfulScans, &failedScans, statusCounts, &statusCountsMutex, outputDir) // [source: 37]

	// --- Prompt and Run Re-scan ---
	initialFailCount := atomic.LoadInt64(&failedScans)
	if initialFailCount > 0 {
		fmt.Printf("\n%s[*] %d targets failed initially. Do you want to re-scan them? (y/N): %s", ColorWarning, initialFailCount, ColorReset) // [source: 37]
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" { // [source: 38]
			failedTargetsMutex.Lock()
			targetsToRescan := make([]string, len(failedTargets))
			copy(targetsToRescan, failedTargets)
			// Optional: Clear global list if memory is a concern
			// failedTargets = nil
			failedTargetsMutex.Unlock() // [source: 38]

			// Run the rescan phase
			// Pass corsCheck flag value here too
			runScanPhase(targetsToRescan, "Re-scan", true, /* isRescan = true */
				sharedClient, *workers, outputDir, *quiet, *corsCheck, // <-- Pass corsCheck value
				&successfulScans, &failedScans,
				statusCounts, &statusCountsMutex,
			)
			// --- Final Summary (after re-scan) ---
			printSummary("Final", startTime, totalTargets, &successfulScans, &failedScans, statusCounts, &statusCountsMutex, outputDir)

		} else {
			fmt.Printf("%s[*] Skipping re-scan.%s\n", ColorInfo, ColorReset)
			// Print the initial summary again as the final summary if no re-scan
			fmt.Println("\n--- Final Summary (No Re-scan) ---")
			fmt.Printf("Total Targets: %d\n", totalTargets)
			fmt.Printf("%sSuccessful: %d%s\n", ColorSuccess, atomic.LoadInt64(&successfulScans), ColorReset)
			fmt.Printf("%sFailed: %d%s\n", ColorError, atomic.LoadInt64(&failedScans), ColorReset)
			// Re-print breakdown if needed, using same variables
			printStatusBreakdown(statusCounts, &statusCountsMutex)                                                       // Extracted breakdown logic
			fmt.Printf("\n%s[*]%s Output saved to: %s%s%s\n", ColorInfo, ColorReset, ColorAccent, outputDir, ColorReset) // [source: 39]
			fmt.Printf("%s[*]%s Scan complete.%s\n", ColorInfo, ColorReset, ColorReset)
		}
	} else {
		// If no initial failures, the initial summary is the final one.
		fmt.Printf("\n%s[*] No targets failed initial scan. Scan complete.%s\n", ColorInfo, ColorReset) // [source: 40]
	}
}

// Helper function extracted from printSummary to avoid duplication
func printStatusBreakdown(statusCounts map[int]int64, statusCountsMutex *sync.Mutex) {
	statusCountsMutex.Lock()
	mapLen := len(statusCounts)
	if mapLen > 0 {
		fmt.Println("\nStatus Code Breakdown:")
		codes := make([]int, 0, mapLen)
		for code := range statusCounts {
			codes = append(codes, code)
		}
		statusCountsMutex.Unlock() // Unlock early after getting keys
		sort.Ints(codes)

		for _, code := range codes {
			category := code / 100
			if category < 1 || category > 5 { // [source: 57]
				continue
			}

			// Use statusColors defined in globals.go
			color, ok := statusColors[code]
			if !ok {
				color = ColorReset // Default color if specific one not found
			}

			// Use statusEmojis defined in globals.go
			emoji := statusEmojis[category]
			// Use statusCodes defined in globals.go
			desc, descOk := statusCodes[code]
			if !descOk {
				desc = "(Unknown Status)" // Handle unknown codes gracefully
			}

			statusCountsMutex.Lock() // Lock again to safely read count
			count := statusCounts[code]
			statusCountsMutex.Unlock()

			// Adjusted formatting for potentially longer descriptions
			fmt.Printf("  %s%d%s %s %-25s : %d\n", color, code, ColorReset, emoji, desc, count)
		}
	} else {
		statusCountsMutex.Unlock() // Ensure unlock if mapLen was 0
	}
}
