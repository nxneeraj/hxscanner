package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var statusCategories = map[int]string{
	1: "1xx", 2: "2xx", 3: "3xx", 4: "4xx", 5: "5xx",
}

var statusCodes = map[int]string{
	100: "Continue", 101: "Switching Protocols", 102: "Processing", 103: "Early Hints",
	200: "OK", 201: "Created", 202: "Accepted", 203: "Non-Authoritative Info", 204: "No Content",
	205: "Reset Content", 206: "Partial Content", 207: "Multi-Status", 208: "Already Reported", 226: "IM Used",
	300: "Multiple Choices", 301: "Moved Permanently", 302: "Found", 303: "See Other", 304: "Not Modified",
	305: "Use Proxy", 306: "Unused", 307: "Temporary Redirect", 308: "Permanent Redirect",
	400: "Bad Request", 401: "Unauthorized", 402: "Payment Required", 403: "Forbidden", 404: "Not Found",
	405: "Method Not Allowed", 406: "Not Acceptable", 407: "Proxy Authentication Required", 408: "Request Timeout",
	409: "Conflict", 410: "Gone", 411: "Length Required", 412: "Precondition Failed", 413: "Payload Too Large",
	414: "URI Too Long", 415: "Unsupported Media Type", 416: "Range Not Satisfiable", 417: "Expectation Failed",
	418: "I'm a Teapot ☕", 421: "Misdirected Request", 422: "Unprocessable Entity", 423: "Locked",
	424: "Failed Dependency", 425: "Too Early", 426: "Upgrade Required", 428: "Precondition Required",
	429: "Too Many Requests", 431: "Request Header Fields Too Large", 451: "Unavailable for Legal Reasons",
	500: "Internal Server Error", 501: "Not Implemented", 502: "Bad Gateway", 503: "Service Unavailable",
	504: "Gateway Timeout", 505: "HTTP Version Not Supported", 506: "Variant Also Negotiates",
	507: "Insufficient Storage", 508: "Loop Detected", 510: "Not Extended", 511: "Network Authentication Required",
}

func scanIP(ip string) (int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://" + ip)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func appendToFile(path, line string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.ModePerm, 0644)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	defer f.Close()
	f.WriteString(line + "\n")
}

func processIPs(ipFilePath string, outputDir string) {
	ipFile, err := os.Open(ipFilePath)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	defer ipFile.Close()

	scanner := bufio.NewScanner(ipFile)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip == "" {
			continue
		}

		status, err := scanIP(ip)
		logPath := filepath.Join(outputDir, "log.txt")
		if err != nil {
			appendToFile(filepath.Join(outputDir, "ip_invalid.txt"), ip)
			appendToFile(logPath, fmt.Sprintf("[!] %s -> ERROR: %v", ip, err))
			colorPrint(ip, 0, "", err)
		} else {
			appendToFile(filepath.Join(outputDir, "ip_exist.txt"), ip)
			cat := statusCategories[status/100]
			target := filepath.Join(outputDir, cat, fmt.Sprintf("%d.txt", status))
			appendToFile(target, ip)
			appendToFile(logPath, fmt.Sprintf("[✓] %s -> %d %s", ip, status, statusCodes[status]))
			colorPrint(ip, status, statusCodes[status], nil)
		}
	}
}

func createOutputStructure(base string) error {
	for _, cat := range statusCategories {
		for code := range statusCodes {
			if strings.HasPrefix(fmt.Sprint(code), cat[:1]) {
				dir := filepath.Join(base, cat)
				os.MkdirAll(dir, os.ModePerm)
				f, _ := os.Create(filepath.Join(dir, fmt.Sprintf("%d.txt", code)))
				f.Close()
			}
		}
	}
	extras := []string{"ip_exist.txt", "ip_invalid.txt", "log.txt"}
	for _, name := range extras {
		f, _ := os.Create(filepath.Join(base, name))
		f.Close()
	}
	return nil
}

func parseFlags() {
	ipInput := flag.String("i", "", "Scan IP file only")
	fileInput := flag.String("f", "", "Scan file with any link, domain, or URL")
	helpFlag := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *helpFlag || (*ipInput == "" && *fileInput == "") {
		fmt.Println("Usage:")
		fmt.Println("  -i <ip-file>       Scan IP addresses")
		fmt.Println("  -f <target-file>   Scan domains, URLs or anything")
		fmt.Println("  -h                 Show help")
		os.Exit(0)
	}

	inputPath := *ipInput
	if *fileInput != "" {
		inputPath = *fileInput
	}

	outputDir := strings.TrimSuffix(inputPath, ".txt") + " output"
	os.MkdirAll(outputDir, os.ModePerm)

	err := createOutputStructure(outputDir)
	if err != nil {
		fmt.Println("Error creating folders:", err)
		os.Exit(1)
	}

	fastProcess(inputPath, outputDir)
}
