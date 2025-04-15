package ui

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// PrintBanner displays the ASCII logo and version.
func PrintBanner() {
	banner := `
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗ 
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

	`
	color.Cyan(banner)
	color.Green("HyperScanner v1.1")
	color.Yellow("Made by @nxneeraj ⚡")
	color.Blue("Start Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// PrintHelp displays help info.
func PrintHelp() {
	fmt.Println("Usage:")
	color.Green("  -i <file>\tSpecify IP/URL list file")
	color.Green("  -f <file>\tAlias for -i")
	color.Green("  -o <format>\tOutput format: json / csv / txt")
	color.Green("  -h\t\tShow help")
	color.Green("  --version\tShow version")
}

// PrintStatus prints colored output for status codes.
func PrintStatus(code int, url string) {
	switch {
	case code >= 200 && code < 300:
		color.Green("[ %d ] %s", code, url)
	case code >= 300 && code < 400:
		color.Yellow("[ %d ] %s", code, url)
	case code >= 400 && code < 500:
		color.Red("[ %d ] %s", code, url)
	case code >= 500:
		color.HiRed("[ %d ] %s", code, url)
	default:
		color.White("[ %d ] %s", code, url)
	}
}

// PrintRetryInfo shows retry attempt info.
func PrintRetryInfo(url string, attempt int) {
	color.Yellow("Retrying %s (Attempt %d)", url, attempt)
}

// PrintError displays an error message.
func PrintError(msg string) {
	color.Red("Error: %s", msg)
}
