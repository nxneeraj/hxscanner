package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	// No progressbar import needed here if createProgressBar is in ui.go
)

// runScanPhase executes either the initial scan or the re-scan
func runScanPhase(
	targets []string,
	description string,
	isRescan bool,
	client *http.Client,
	workersCount int,
	outputDir string,
	quiet bool,
	successfulScans *int64,
	failedScans *int64,
	statusCounts map[int]int64,
	statusCountsMutex *sync.Mutex,
) {
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
	results := make(chan scanResult, workersCount)
	var wg sync.WaitGroup      // For workers in this phase
	var resultWg sync.WaitGroup // For result processor in this phase

	// Start workers for this phase (worker function is in scanner.go)
	for w := 1; w <= workersCount; w++ {
		wg.Add(1)
		go worker(w, &wg, client, jobs, results, isRescan) // Pass isRescan flag
	}

	// Progress bar for this phase (createProgressBar is in ui.go)
	bar := createProgressBar(totalScanTargets, fmt.Sprintf("%s[*] %s%s", ColorInfo, description, ColorReset))

	// Start results processor for this phase (processResults is in results.go)
	resultWg.Add(1)
	go processResults(results, &resultWg, bar, quiet, outputDir,
		successfulScans, failedScans, // Pass pointers/refs to update overall stats
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
	resultWg.Wait() // Wait for result processor

	bar.Finish() // Cleanly finish progress bar
	fmt.Printf("%s[*] %s phase complete.%s\n", ColorInfo, description, ColorReset)
}

func main() {
	startTime := time.Now()
	// initMaps() // Call if using init function for maps in globals.go
	printBanner() // From ui.go

	// --- Command Line Flags ---
	targetInput := flag.String("i", "", "Input file containing IPs, Domains, or URLs (one per line)")
	fileInput := flag.String("f", "", "Alias for -i (optional)")
	helpFlag := flag.Bool("h", false, "Show help")
	// Default workers to number of CPU cores, with a minimum
	defaultWorkers := runtime.NumCPU()
	if defaultWorkers < 4 {	defaultWorkers = 4 }
	workers := flag.Int("w", defaultWorkers, fmt.Sprintf("Number of concurrent workers (default: %d)", defaultWorkers))
	timeout := flag.Duration("t", 5*time.Second, "HTTP request timeout (e.g., 3s, 10s)")
	quiet := flag.Bool("q", false, "Quiet mode: suppress individual results, show only progress and summary")
	flag.Parse()

	// --- Help Flag Handling ---
	if *helpFlag {
		fmt.Println("Usage: hxscanner [options]") // Use desired binary name
		fmt.Println("\nScans IPs, Domains, or full URLs from an input file via HTTP/S.")
		fmt.Println("If no scheme (http:// or https://) is provided for a domain/IP, http:// is assumed.")
		fmt.Println("Offers an option to re-scan failed targets after the initial scan.")
		fmt.Println("\nOptions:")
		fmt.Println("  -i <file>     Input file with targets (IPs/Domains/URLs), one per line (required if -f not used)")
		fmt.Println("  -f <file>     Alias for -i")
		fmt.Println("  -w <number>   Number of concurrent scanning workers (default: number of CPU cores)")
		fmt.Println("  -t <duration> HTTP request timeout (default: 5s)")
		fmt.Println("  -q            Quiet mode: suppress individual results (except errors/warnings)")
		fmt.Println("  -h            Show this help message")
		os.Exit(0)
	}

	// --- Input File Validation ---
	if *targetInput == "" && *fileInput == "" {
		fmt.Printf("%sError: No input file provided. Use -i <file> or -f <file>%s\n", ColorError, ColorReset)
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
	fmt.Printf("%s[*] Output will be saved to: %s%s%s\n", ColorInfo, ColorAccent, outputDir, ColorReset)

	// --- Shared HTTP Client Setup ---
	sharedClient := setupHTTPClient(*timeout, *workers) // From scanner.go

	// --- Overall Statistics Setup ---
	var successfulScans int64
	var failedScans int64 // This will reflect initial failures, adjusted by rescan successes
	statusCounts := make(map[int]int64)
	var statusCountsMutex sync.Mutex

	// --- Run Initial Scan ---
	runScanPhase(initialTargets, "Initial Scan", false, /* isRescan = false */
		sharedClient, *workers, outputDir, *quiet,
		&successfulScans, &failedScans, // Pass pointers
		statusCounts, &statusCountsMutex,
	)

	// --- Initial Summary ---
	// Summary reflects state *after* the first pass. Failed count is accurate for initial run.
	printSummary("Initial Scan", startTime, totalTargets, &successfulScans, &failedScans, statusCounts, &statusCountsMutex, outputDir)


	// --- Prompt and Run Re-scan ---
	initialFailCount := atomic.LoadInt64(&failedScans) // Get # of failures from initial run
	if initialFailCount > 0 {
		fmt.Printf("\n%s[*] %d targets failed initially. Do you want to re-scan them? (y/N): %s", ColorWarning, initialFailCount, ColorReset)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			// Make a copy of the failed targets slice before starting rescan
			failedTargetsMutex.Lock()
			targetsToRescan := make([]string, len(failedTargets))
			copy(targetsToRescan, failedTargets)
			// Clear the global list if no longer needed (optional)
			// failedTargets = nil
			failedTargetsMutex.Unlock()

			// Run the rescan phase
			runScanPhase(targetsToRescan, "Re-scan", true, /* isRescan = true */
				sharedClient, *workers, outputDir, *quiet,
				&successfulScans, &failedScans, // Pass same pointers, processResults will adjust counts
				statusCounts, &statusCountsMutex,
			)
			// --- Final Summary (after re-scan) ---
			printSummary("Final", startTime, totalTargets, &successfulScans, &failedScans, statusCounts, &statusCountsMutex, outputDir)

		} else {
			fmt.Printf("%s[*] Skipping re-scan.%s\n", ColorInfo, ColorReset)
			// Print the initial summary again as the final summary
			fmt.Println("\n--- Final Summary (No Re-scan) ---")
			fmt.Printf("Total Targets: %d\n", totalTargets)
			fmt.Printf("%sSuccessful: %d%s\n", ColorSuccess, atomic.LoadInt64(&successfulScans), ColorReset)
			fmt.Printf("%sFailed: %d%s\n", ColorError, atomic.LoadInt64(&failedScans), ColorReset)
			// Re-print breakdown if needed, using same variables
            // ... (code similar to breakdown in printSummary) ...
			fmt.Printf("\n%s[*]%s Output saved to: %s%s%s\n", ColorInfo, ColorReset, ColorAccent, outputDir, ColorReset)
			fmt.Printf("%s[*]%s Scan complete.%s\n", ColorInfo, ColorReset, ColorReset)
		}
	} else {
		// If no initial failures, the initial summary is the final one.
		fmt.Printf("\n%s[*] No targets failed initial scan. Scan complete.%s\n", ColorInfo, ColorReset)
	}
}
