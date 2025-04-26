package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Constants and global maps are now expected to be in globals.go

// printBanner displays the tool's banner
func printBanner() {
	fmt.Println(ColorBanner)
	fmt.Println(`
██╗  ██╗██╗  ██╗      ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
██║  ██║╚██╗██╔╝      ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗ ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝ ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗      ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝      ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
	`)
	fmt.Println(ColorAccent + "         HyperScanner v1.4+CORS (IP/Domain/URL Scanner)" + ColorReset)
	fmt.Println()
}

// colorPrint displays a single primary scan result (status or error).
func colorPrint(target string, code int, desc string, err error, quiet bool, isRescan bool) {
	if quiet && err == nil {
		return
	}

	rescanPrefix := ""
	if isRescan {
		rescanPrefix = "[RESCAN] "
	}

	if err != nil {
		fmt.Printf("%s%s[X]%s %s -> %sERROR: %v%s\n", rescanPrefix, ColorError, ColorReset, target, ColorError, err, ColorReset)
		return
	}

	category := code / 100
	emoji := statusEmojis[category]
	if emoji == "" {
		emoji = "❓"
	}

	color, ok := statusColors[code]
	if !ok {
		switch category {
		case 1:
			color = ColorInfo
		case 2:
			color = ColorSuccess
		case 3:
			color = ColorWarning
		case 4:
			color = ColorError
		case 5:
			color = ColorError
		default:
			color = ColorReset
		}
	}

	successMark := "[✓]"
	if isRescan {
		successMark = "[✓✓]"
	}

	fmt.Printf("%s%s%s%s %s -> %s%d%s %s %s%s\n",
		rescanPrefix,
		ColorSuccess, successMark, ColorReset,
		target,
		color, code, ColorReset,
		emoji,
		desc,
		ColorReset)
}

// createProgressBar initializes a new progress bar using settings from globals.go
func createProgressBar(total int, description string) *progressbar.ProgressBar {
	if total <= 0 {
		total = 1
	}
	return progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		// progressbar.OptionShowElapsedTime(true), // <-- Removed/Commented out this line
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        fmt.Sprintf("%s█%s", ColorSuccess, ColorReset),
			SaucerHead:    fmt.Sprintf("%s>%s", ColorAccent, ColorReset),
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

// printSummary displays the scan summary table using final counts
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
	finalSuccess := atomic.LoadInt64(successfulScans)
	finalFailed := atomic.LoadInt64(failedScans)

	fmt.Printf("\n--- %s Summary (%s) ---\n", title, time.Since(startTime).Round(time.Millisecond))
	fmt.Printf("Total Targets: %d\n", totalTargets)
	fmt.Printf("%sSuccessful Scans: %d%s\n", ColorSuccess, finalSuccess, ColorReset)
	fmt.Printf("%sFailed Scans:     %d%s\n", ColorError, finalFailed, ColorReset)

	if finalSuccess > 0 {
		// Use the helper function from main.go
		printStatusBreakdown(statusCounts, statusCountsMutex)
	}

	fmt.Printf("\n%s[*]%s Output saved to: %s%s%s\n", ColorInfo, ColorReset, ColorAccent, outputDir, ColorReset)
}

// Note: printStatusBreakdown function is now in main.go
