package ui

import (
	"fmt"
	"github.com/fatih/color"
)

// DisplayBanner displays a simple ASCII banner for hxscanner.
func DisplayBanner() {
	banner := `
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗ 
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

        HyperScanner v1.1 🔥 - Ultra Fast HTTP Status Scanner by Neeraj Sah
        GitHub: https://github.com/nxneeraj/hxscanner
--------------------------------------------------------------------------------
`
	color.Cyan(banner)
	fmt.Println("Welcome to hxscanner! The HTTP Status Code Scanner")
	fmt.Println("---------------------------------------------------")
}

// DisplayScanningProgress shows the current status of the scan.
func DisplayScanningProgress(ip string, statusCode string) {
	color.Yellow("Scanning IP/URL: %s - Status Code: %s", ip, statusCode)
}

// DisplayScanCompletion shows a message when the scan is completed.
func DisplayScanCompletion(outputDir string) {
	color.Green("Scan completed successfully. Results are saved in: %s", outputDir)
}

// DisplayError shows an error message in red.
func DisplayError(errMsg string) {
	color.Red("Error: %s", errMsg)
}

// DisplayHelp displays basic help instructions for the tool.
func DisplayHelp() {
	helpText := `
Usage:
  hxscanner <input_file> <output_directory>

Options:
  <input_file>      Path to the file containing a list of IPs or URLs to scan.
  <output_directory> Directory where the output (categorized status codes) will be saved.

Examples:
  hxscanner ips.txt output_folder
`
	color.Cyan(helpText)
}
