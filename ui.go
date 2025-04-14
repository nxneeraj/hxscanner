package ui

import (
	"fmt"
	"time"
	"os"
)

// showWelcomeBanner displays the banner with cool ASCII art and program name
func showWelcomeBanner() {
	clearScreen()

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
	fmt.Println(banner)
	fmt.Println("Initializing HyperScanner... Please Wait...")
	time.Sleep(2 * time.Second) // simulate the loading time
}

// showLoadingAnimation displays a loading animation during the scanning process
func showLoadingAnimation() {
	fmt.Print("Loading: [")
	for i := 0; i < 10; i++ {
		fmt.Print("=")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("] Done!")
}

// showScanStart displays a professional-looking start message before the scanning begins
func showScanStart() {
	clearScreen()
	fmt.Println("\n\n-------------------------------------------------")
	fmt.Println("HyperScanner v1.1")
	fmt.Println("-------------------------------------------------")
	fmt.Println("Scanning started... 🚀")
	fmt.Println("-------------------------------------------------")
}

// showScanProgress shows progress with a real-time status update on the URL being scanned
func showScanProgress(url string, statusCode int) {
	// Customize status output to include URL and Status Code
	fmt.Printf("Scanning URL: %-50s Status: %-4d\n", url, statusCode)
}

// showScanComplete shows a cool summary at the end of the scan
func showScanComplete(totalURLs, successfulScans, failedScans int) {
	fmt.Println("\n-------------------------------------------------")
	fmt.Println("Scan Complete! 🎉")
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Total URLs scanned: %d\n", totalURLs)
	fmt.Printf("Successful scans: %d\n", successfulScans)
	fmt.Printf("Failed scans: %d\n", failedScans)
	fmt.Println("-------------------------------------------------")
	fmt.Println("Thanks for using HyperScanner! 🔥")
	fmt.Println("-------------------------------------------------\n")
}

// clearScreen clears the console screen for a neat UI experience
func clearScreen() {
	if os.Getenv("OS") == "Windows_NT" {
		// Windows
		fmt.Print("\x0c")
	} else {
		// Unix-based systems
		fmt.Print("\033[H\033[2J")
	}
}

