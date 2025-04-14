package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"hxscanner/ui"
	"hxscanner/installer"
)

var (
	outputBase = ""
	client     = &http.Client{Timeout: 5 * time.Second}
)

func main() {
	ui.PrintBanner()

	if len(os.Args) < 2 {
		ui.LogError("Usage: hxscanner <input-file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputBase = strings.TrimSuffix(inputFile, ".txt") + "_output"
	os.MkdirAll(outputBase, 0755)

	urls, err := readLines(inputFile)
	if err != nil {
		ui.LogError("Failed to read input file")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	for _, url := range urls {
		url := normalize(url)

		if url == "" {
			appendToFile("log.txt", "Empty or invalid line\n")
			continue
		}

		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			sem <- struct{}{}
			scanURL(link)
			<-sem
		}(url)
	}

	wg.Wait()
	ui.LogSuccess("Scanning complete.")
	installer.CheckAndSetup()
}

func scanURL(url string) {
	resp, err := client.Get(url)
	if err != nil {
		ui.LogError(fmt.Sprintf("Invalid: %s", url))
		appendToFile("ip_invalid.txt", url)
		return
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	ui.LogResult(url, code)
	appendToFile("ip_exist.txt", url)

	codeStr := fmt.Sprintf("%d", code)
	writeToCategoryFolder(codeStr, url)
}

func normalize(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "http") {
		return line
	}
	if strings.Contains(line, ".") {
		return "http://" + line
	}
	return ""
}

func writeToCategoryFolder(code string, url string) {
	folder := outputBase + "/" + code[:1] + "xx"
	os.MkdirAll(folder, 0755)
	appendToFile(fmt.Sprintf("%s/%s.txt", folder, code), url)
}

func appendToFile(filename, line string) {
	f, _ := os.OpenFile(outputBase+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(line + "\n")
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
