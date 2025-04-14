package main

import (
	"fmt"
	"os"
	"time"
	"flag"
	"net/http"
	"io/ioutil"
	"strings"
)

// Global variables
var (
	urlFile     string
	concurrency int
	timeout     int
	showAll     bool
	delay       int
	userAgent   string
)

// init function to parse flags and setup
func init() {
	// Command-line flags
	flag.StringVar(&urlFile, "f", "", "📁 Path to file containing target URLs")
	flag.IntVar(&concurrency, "c", 50, "🚀 Number of concurrent requests")
	flag.IntVar(&timeout, "t", 10, "⏱️ Request timeout in seconds")
	flag.BoolVar(&showAll, "a", false, "🧾 Show all status codes (including non-2xx/3xx)")
	flag.IntVar(&delay, "d", 0, "🕒 Delay between requests in milliseconds")
	flag.StringVar(&userAgent, "ua", "HyperScanner/1.1", "🕵️ Custom User-Agent header")

	// Custom Usage banner
	flag.Usage = func() {
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

USAGE:
    hxscanner -f urls.txt [options]

OPTIONS:
`
		fmt.Fprintln(os.Stderr, banner)
		flag.PrintDefaults()
	}

	// Small delay for aesthetic effect before showing usage
	time.Sleep(200 * time.Millisecond)
}

// Function to start scanning URLs
func scanURLs() {
	// Read file containing URLs
	fileContent, err := ioutil.ReadFile(urlFile)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", urlFile, err)
		return
	}

	// Split file content into URLs (assuming each URL is on a new line)
	urls := strings.Split(string(fileContent), "\n")

	// Loop through URLs and scan them
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url != "" {
			// Start HTTP request
			statusCode := scanURL(url)
			if showAll || (statusCode >= 200 && statusCode < 400) {
				fmt.Printf("URL: %s - Status: %d\n", url, statusCode)
			}
		}
	}
}

// Function to send an HTTP request and get the status code
func scanURL(url string) int {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for URL %s: %v\n", url, err)
		return -1
	}

	// Set custom User-Agent if provided
	req.Header.Set("User-Agent", userAgent)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to URL %s: %v\n", url, err)
		return -1
	}
	defer resp.Body.Close()

	// Return the status code
	return resp.StatusCode
}

// main function to handle program execution
func main() {
	// Parse the flags
	flag.Parse()

	// Validate flags
	if urlFile == "" {
		fmt.Println("Error: Please provide a valid URL file using the -f flag.")
		return
	}

	// Show UI before scanning
	showUI()

	// Start scanning URLs
	scanURLs()
}

// Function to show a cool loading UI
func showUI() {
	clearScreen()
	fmt.Println("Starting the HyperScanner... Please wait...")
	time.Sleep(2 * time.Second) // simulate loading

	// Displaying a loading bar (optional)
	fmt.Print("Loading: [")
	for i := 0; i < 10; i++ {
		fmt.Print("=")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("] Done!")
}

// Function to clear the screen (cross-platform)
func clearScreen() {
	if os.Getenv("OS") == "Windows_NT" {
		// Windows
		fmt.Print("\x0c")
	} else {
		// Unix-like systems
		fmt.Print("\033[H\033[2J")
	}
}
