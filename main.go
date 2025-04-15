package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/nxneeraj/hxscanner/ui"
)

func main() {
	uix := ui.NewUI()
	uix.PrintBanner()

	if len(os.Args) < 2 {
		uix.ShowHelp()
		return
	}

	inputFile := ""

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-i" || arg == "-f" {
			if i+1 < len(os.Args) {
				inputFile = os.Args[i+1]
				i++
			} else {
				uix.PrintError("Missing value for " + arg)
				return
			}
		} else if arg == "-h" || arg == "--help" {
			uix.ShowHelp()
			return
		}
	}

	if inputFile == "" {
		uix.PrintError("No input file provided. Use -i or -f.")
		return
	}

	if !fileExists(inputFile) {
		uix.PrintError("File not found: " + inputFile)
		return
	}

	targets, err := readLines(inputFile)
	if err != nil {
		uix.PrintError("Failed to read input file: " + err.Error())
		return
	}

	uix.PrintInfo(fmt.Sprintf("Loaded %d URLs/IPs for scanning...", len(targets)))

	RunScanner(targets, uix)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readLines(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned, nil
}
