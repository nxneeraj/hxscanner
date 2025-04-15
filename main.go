package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nxneeraj/hxscanner/init"
	"github.com/nxneeraj/hxscanner/ui"
	"github.com/nxneeraj/hxscanner/installer"
	"github.com/fatih/color"
)

// ScanHTTPStatus performs the scanning for HTTP status codes and categorizes them.
func ScanHTTPStatus(inputFile string, outputDir string) error {
	// Open the input file containing IPs or URLs
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("could not open input file: %v", err)
	}
	defer file.Close()

	// Prepare to write the output files
	var wg sync.WaitGroup
	ips := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Assuming input file is a list of IPs or URLs
		// For simplicity, we categorize the status codes manually in this example
		statusCode := getStatusCode(line) // This should be a function that gets the status code from the URL/IP
		ips[statusCode] = append(ips[statusCode], line)
		wg.Add(1)
	}

	// Wait for all statuses to be processed
	wg.Wait()

	// Handle output
	return writeOutput(outputDir, ips)
}

// getStatusCode is a stub function that simulates getting a status code from an IP or URL.
// In a real-world application, this would involve making HTTP requests to each URL/IP.
func getStatusCode(ip string) string {
	// Dummy logic for status code assignment
	return "200" // Here, we simply assume all IPs have a 200 status for demonstration purposes
}

// writeOutput writes categorized output based on HTTP status codes.
func writeOutput(outputDir string, ips map[string][]string) error {
	// Create the output directory
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create output directory: %v", err)
	}

	// Writing the categorized status codes into separate files
	for statusCode, ipList := range ips {
		fileName := fmt.Sprintf("%s/%s.txt", outputDir, statusCode)
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("could not create file %s: %v", fileName, err)
		}
		defer file.Close()

		for _, ip := range ipList {
			_, err := file.WriteString(ip + "\n")
			if err != nil {
				return fmt.Errorf("could not write to file %s: %v", fileName, err)
			}
		}
	}

	return nil
}

// ShowBanner displays the initial ASCII banner for the tool.
func ShowBanner() {
	ui.DisplayBanner()
}

// HandleArguments processes command-line arguments for file input and output paths.
func HandleArguments() (string, string) {
	if len(os.Args) < 3 {
		color.Red("Usage: hxscanner <input_file> <output_dir>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputDir := os.Args[2]

	return inputFile, outputDir
}

// Main function to tie everything together
func main() {
	// Display tool's banner
	ShowBanner()

	// Handle command-line arguments
	inputFile, outputDir := HandleArguments()

	// Run the HTTP status scanning process
	err := ScanHTTPStatus(inputFile, outputDir)
	if err != nil {
		color.Red("Error scanning HTTP status codes: %v", err)
		return
	}

	color.Green("Scanning complete. Output written to: %s", outputDir)
}
