package ui

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

// ShowBanner prints ASCII art banner
func ShowBanner() {
	banner := `
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

HyperScanner v1.1 🔥  -  Ultra Fast HTTP Status Scanner

`
	fmt.Println(color.HiMagentaString(banner))
}

// ShowHelp prints help
func ShowHelp() {
	fmt.Println(`Usage:
  -i <file>      Input file with IPs/URLs
  -f <file>      Alias for -i
  -c <int>       Concurrency (default: 100)
  -r <int>       Retries (default: 1)
  -json <file>   Save output to JSON file (default: output.json)
  -csv <file>    Save output to CSV file (default: output.csv)
  -h             Show this help

Example:
  hxscanner -i list.txt -c 200 -r 2
`)
}

// PrintStatus shows colored status result
func PrintStatus(url string, status int) {
	var statusColor *color.Color

	switch {
	case status >= 200 && status < 300:
		statusColor = color.New(color.FgGreen)
	case status >= 300 && status < 400:
		statusColor = color.New(color.FgCyan)
	case status >= 400 && status < 500:
		statusColor = color.New(color.FgYellow)
	case status >= 500:
		statusColor = color.New(color.FgRed)
	default:
		statusColor = color.New(color.FgWhite)
	}

	statusStr := strconv.Itoa(status)
	statusColor.Printf("[ %s ] ", statusStr)
	fmt.Println(url)
}
