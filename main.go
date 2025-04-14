package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	. "hxscanner/ui"
	. "hxscanner/init"
	. "hxscanner/installer"
)

// Result holds individual scan results
type Result struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

var (
	inputFile     string
	outputJSON    string
	outputCSV     string
	concurrency   int
	maxRetries    int
	showHelp      bool
	runSetup      bool

	results     []Result
	resultMutex sync.Mutex
	wg          sync.WaitGroup
	sem         chan struct{}
)

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file with IPs or URLs")
	flag.StringVar(&inputFile, "f", "", "Alias for -i (input file)")
	flag.StringVar(&outputJSON, "json", "output.json", "Save results to JSON file")
	flag.StringVar(&outputCSV, "csv", "output.csv", "Save results to CSV file")
	flag.IntVar(&concurrency, "c", 100, "Number of concurrent scans")
	flag.IntVar(&maxRetries, "r", 1, "Number of retries for failed URLs")
	flag.BoolVar(&showHelp, "h", false, "Show help and usage")
	flag.BoolVar(&runSetup, "setup", false, "Run installer to make hxscanner global")
}

func main() {
	flag.Parse()

	if runSetup {
		RunInstaller() // Run global installation (installer.go)
		return
	}

	if showHelp || inputFile == "" {
		ShowBanner()
		ShowHelp()
		return
	}

	ShowBanner()
	SetupEnv() // Create folders, output structure

	urls, err := readInput(inputFile)
	if err != nil {
		fmt.Println("❌ Error reading input file:", err)
		return
	}

	sem = make(chan struct{}, concurrency)

	start := time.Now()
	for _, url := range urls {
		wg.Add(1)
		go scanURL(url)
	}
	wg.Wait()
	elapsed := time.Since(start)

	saveJSON(outputJSON)
	saveCSV(outputCSV)

	fmt.Printf("\n✅ Scan completed in %s\n", elapsed)
}

func readInput(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if !strings.HasPrefix(line, "http") {
				line = "http://" + line
			}
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

func scanURL(url string) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	var resp *http.Response
	var err error

	for i := 0; i <= maxRetries; i++ {
		resp, err = http.Get(url)
		if err == nil && resp != nil {
			defer resp.Body.Close()
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	status := 0
	if resp != nil {
		status = resp.StatusCode
	}

	PrintStatus(url, status)
	StoreResult(url, status) // Categorize + log file (init.go)

	resultMutex.Lock()
	results = append(results, Result{URL: url, Status: status})
	resultMutex.Unlock()
}

func saveJSON(filename string) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("❌ Failed to save JSON:", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("❌ Error writing JSON file:", err)
	}
}

func saveCSV(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("❌ Failed to save CSV:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"URL", "Status"})
	for _, r := range results {
		writer.Write([]string{r.URL, fmt.Sprintf("%d", r.Status)})
	}
}
