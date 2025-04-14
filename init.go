package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	outputFolder     string
	categorizedPaths map[string]string
	allCodes         = []int{
		100, 101, 102, 103,
		200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
		300, 301, 302, 303, 304, 305, 306, 307, 308,
		400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410,
		411, 412, 413, 414, 415, 416, 417, 418, 421, 422,
		423, 424, 425, 426, 428, 429, 431, 451,
		500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511,
	}
)

// SetupEnv creates output folder structure based on input file
func SetupEnv() {
	ts := time.Now().Format("2006-01-02_15-04-05")
	outputFolder = "hx_output_" + ts
	os.MkdirAll(outputFolder, 0755)

	categorizedPaths = map[string]string{
		"1xx": filepath.Join(outputFolder, "1xx"),
		"2xx": filepath.Join(outputFolder, "2xx"),
		"3xx": filepath.Join(outputFolder, "3xx"),
		"4xx": filepath.Join(outputFolder, "4xx"),
		"5xx": filepath.Join(outputFolder, "5xx"),
	}

	// Create main folders
	for _, folder := range categorizedPaths {
		os.MkdirAll(folder, 0755)
	}

	// Create individual status code files
	for _, code := range allCodes {
		codeStr := strconv.Itoa(code)
		codeDir := filepath.Join(outputFolder, string(codeStr[0])+"xx")
		codeFile := filepath.Join(codeDir, codeStr+".txt")
		os.WriteFile(codeFile, []byte{}, 0644)
	}

	// Create special logs
	os.WriteFile(filepath.Join(outputFolder, "ip_exist.txt"), []byte{}, 0644)
	os.WriteFile(filepath.Join(outputFolder, "ip_invalid.txt"), []byte{}, 0644)
	os.WriteFile(filepath.Join(outputFolder, "log.txt"), []byte{}, 0644)
}

// StoreResult categorizes and writes the URL based on status code
func StoreResult(url string, status int) {
	statusStr := strconv.Itoa(status)
	mainCat := string(statusStr[0]) + "xx"
	categoryFolder := categorizedPaths[mainCat]

	codeFile := filepath.Join(categoryFolder, statusStr+".txt")
	appendLine(codeFile, url)

	// Additional centralized logs
	appendLine(filepath.Join(outputFolder, "log.txt"), fmt.Sprintf("%s -> %d", url, status))

	if status == 0 || status > 599 {
		appendLine(filepath.Join(outputFolder, "ip_invalid.txt"), url)
	} else {
		appendLine(filepath.Join(outputFolder, "ip_exist.txt"), url)
	}
}

func appendLine(path string, line string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Failed to write to", path, ":", err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(strings.TrimSpace(line) + "\n")
}
