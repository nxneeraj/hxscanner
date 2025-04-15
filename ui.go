package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

func printBanner() {
	fmt.Println(ColorBanner)
	fmt.Println(`
██╗  ██╗██╗  ██╗      ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
██║  ██║╚██╗██╔╝      ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗ ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝ ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗      ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝      ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
	`) // Using backticks for multi-line raw string
	fmt.Println(ColorAccent + "         HyperScanner v1.4 (IP/Domain/URL Scanner w/ Rescan)" + ColorReset)
	fmt.Println()
}

// colorPrint displays a single scan result, respecting quiet mode and indicating rescans
func colorPrint(target string, code int, desc string, err error, quiet bool, isRescan bool) {
	if quiet { // Suppress all individual output if quiet mode is enabled
		return
	}

	rescanPrefix := ""
	if isRescan {
		rescanPrefix = "[RESCAN] "
	}

	if err != nil {
		fmt.Printf("%s%s%s[x]%s %s -> %sERROR: %v%s\n", rescanPrefix, ColorError, ColorReset, target, ColorError, err, ColorReset)
		return
	}
	category := code / 100
	// Handle cases where code might be 0 (e.g., before error checked) - though less likely now
	if category < 1 || category > 5 {
		// If code is 0, it's likely an error state already handled.
		// If code is non-zero but category is invalid, print basic info.
		if code != 0 {
		    fmt.Printf("%s[?] %s -> %d %s\n", rescanPrefix, target, code, desc)
		}
		return
	}
	emoji := statusEmojis[category]
	color, ok := statusColors[code]
	if !ok {
		color = ColorReset
	}
	successMark := "✓"
	if isRescan {
		successMark = "✓✓" // Double check indicates success on rescan
	}
	fmt.Printf("%s%s%s%s %s -> %s%d%s %s%s%s\n", rescanPrefix, ColorSuccess, successMark, ColorReset, target, color, code, ColorReset, emoji, desc, ColorReset)
}


// createProgressBar initializes a new progress bar
func createProgressBar(total int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr), // Write bar to stderr
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        fmt.Sprintf("%s=%s", ColorSuccess, ColorReset),
			SaucerHead:    fmt.Sprintf("%s>%s", ColorSuccess, ColorReset),
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

// printSummary displays the scan summary table
func printSummary(
	title string,
	startTime time.Time,
	totalTargets int,
	successfulScans *int64,
	failedScans *int64,
	statusCounts map[int]int64,
	statusCountsMutex *sync.Mutex,
	outputDir string,
) {
	fmt.Printf("\n--- %s Summary (%s) ---\n", title, time.Since(startTime).Round(time.Millisecond))
	fmt.Printf("Total Targets: %d\n", totalTargets)
	fmt.Printf("%sSuccessful: %d%s\n", ColorSuccess, atomic.LoadInt64(successfulScans), ColorReset)
	fmt.Printf("%sFailed: %d%s\n", ColorError, atomic.LoadInt64(failedScans), ColorReset)

	// Print status code breakdown only if there are successful scans
	if atomic.LoadInt64(successfulScans) > 0 {
		statusCountsMutex.Lock() // Lock map before reading
		mapLen := len(statusCounts)
		if mapLen > 0 {
			fmt.Println("\nStatus Code Breakdown:")
			codes := make([]int, 0, mapLen)
			for code := range statusCounts {
				codes = append(codes, code)
			}
			statusCountsMutex.Unlock() // Unlock map immediately after reading keys
			sort.Ints(codes)           // Sort keys

			for _, code := range codes {
				category := code / 100
                // Prevent panic if category is somehow invalid (e.g., code 0)
                if category < 1 || category > 5 { continue }

				color, ok := statusColors[code]
				if !ok { color = ColorReset }

				emoji := statusEmojis[category]
				desc, _ := statusCodes[code] // Desc might be empty if code not in map

				// Lock map again just to read the count for this code
				statusCountsMutex.Lock()
				count := statusCounts[code]
				statusCountsMutex.Unlock()
				fmt.Printf("  %s%d%s %s%-25s : %d\n", color, code, ColorReset, emoji, desc, count)
			}
		} else {
			statusCountsMutex.Unlock() // Ensure unlock even if map was empty
		}
	}

    // Only print completion message once at the very end in main.go
	// fmt.Printf("\n%s[*]%s Output saved to: %s%s%s\n", ColorInfo, ColorReset, ColorAccent, outputDir, ColorReset)
	// fmt.Printf("%s[*]%s Scan complete.%s\n", ColorInfo, ColorReset, ColorReset)
}
