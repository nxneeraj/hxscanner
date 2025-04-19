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

// Terminal color codes
const (
	ColorReset   = "\033[0m"
	ColorAccent  = "\033[36m"
	ColorError   = "\033[31m"
	ColorSuccess = "\033[32m"
	ColorInfo    = "\033[34m"
	ColorBanner  = "\033[35m"
)

// Emoji for each HTTP status code class
var statusEmojis = map[int]string{
	1: "🟦",
	2: "🟩",
	3: "🟨",
	4: "🟥",
	5: "🟥",
}

// Color associated with each specific status code
var statusColors = map[int]string{
	200: ColorSuccess,
	201: ColorSuccess,
	301: ColorAccent,
	302: ColorAccent,
	400: ColorError,
	401: ColorError,
	403: ColorError,
	404: ColorError,
	500: ColorError,
}

// Short descriptions for common HTTP status codes
var statusCodes = map[int]string{
	200: "OK",
	201: "Created",
	301: "Moved Permanently",
	302: "Found",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	500: "Internal Server Error",
}

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
	fmt.Println(ColorAccent + "         HyperScanner v1.4 (IP/Domain/URL Scanner w/ Rescan)" + ColorReset)
	fmt.Println()
}

// colorPrint displays a single scan result, respecting quiet mode and indicating rescans
func colorPrint(target string, code int, desc string, err error, quiet bool, isRescan bool) {
	if quiet {
		return
	}

	rescanPrefix := ""
	if isRescan {
		rescanPrefix = "[RESCAN] "
	}

	if err != nil {
		fmt.Printf("%s%s[x]%s %s -> %sERROR: %v%s\n", rescanPrefix, ColorError, ColorReset, target, ColorError, err, ColorReset)
		return
	}

	category := code / 100
	if category < 1 || category > 5 {
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
		successMark = "✓✓"
	}

	fmt.Printf("%s%s%s%s %s -> %s%d%s %s%s%s\n",
		rescanPrefix, ColorSuccess, successMark, ColorReset,
		target, color, code, ColorReset,
		emoji, desc, ColorReset)
}

// createProgressBar initializes a new progress bar
func createProgressBar(total int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
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

	if atomic.LoadInt64(successfulScans) > 0 {
		statusCountsMutex.Lock()
		mapLen := len(statusCounts)
		if mapLen > 0 {
			fmt.Println("\nStatus Code Breakdown:")
			codes := make([]int, 0, mapLen)
			for code := range statusCounts {
				codes = append(codes, code)
			}
			statusCountsMutex.Unlock()
			sort.Ints(codes)

			for _, code := range codes {
				category := code / 100
				if category < 1 || category > 5 {
					continue
				}

				color, ok := statusColors[code]
				if !ok {
					color = ColorReset
				}

				emoji := statusEmojis[category]
				desc, _ := statusCodes[code]

				statusCountsMutex.Lock()
				count := statusCounts[code]
				statusCountsMutex.Unlock()

				fmt.Printf("  %s%d%s %s%-25s : %d\n", color, code, ColorReset, emoji, desc, count)
			}
		} else {
			statusCountsMutex.Unlock()
		}
	}

	// Optional completion message
	// fmt.Printf("\n%s[*]%s Output saved to: %s%s%s\n", ColorInfo, ColorReset, ColorAccent, outputDir, ColorReset)
	// fmt.Printf("%s[*]%s Scan complete.%s\n", ColorInfo, ColorReset, ColorReset)
}
