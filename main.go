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

var statusEmojis = map[int]string{
        1: "🔵", 2: "✅", 3: "🟡", 4: "❌", 5: "💥",
}

var statusColors = map[int]string{
        100: "\033[38;5;39m", 101: "\033[38;5;33m", 102: "\033[38;5;45m", 103: "\033[38;5;27m",
        200: "\033[38;5;82m", 201: "\033[38;5;76m", 202: "\033[38;5;46m", 203: "\033[38;5;70m", 204: "\033[38;5;40m",
        205: "\033[38;5;78m", 206: "\033[38;5;34m", 207: "\033[38;5;48m", 208: "\033[38;5;35m", 226: "\033[38;5;83m",
        300: "\033[38;5;220m", 301: "\033[38;5;190m", 302: "\033[38;5;184m", 303: "\033[38;5;214m", 304: "\033[38;5;226m",
        305: "\033[38;5;228m", 306: "\033[38;5;229m", 307: "\033[38;5;230m", 308: "\033[38;5;227m",
        400: "\033[38;5;160m", 401: "\033[38;5;196m", 402: "\033[38;5;124m", 403: "\033[38;5;203m", 404: "\033[38;5;161m",
        405: "\033[38;5;197m", 406: "\033[38;5;125m", 407: "\033[38;5;162m", 408: "\033[38;5;198m",
        409: "\033[38;5;126m", 410: "\033[38;5;199m", 411: "\033[38;5;127m", 412: "\033[38;5;200m", 413: "\033[38;5;128m",
        414: "\033[38;5;201m", 415: "\033[38;5;129m", 416: "\033[38;5;202m", 417: "\033[38;5;130m",
        418: "\033[38;5;203m", 421: "\033[38;5;131m", 422: "\033[38;5;204m", 423: "\033[38;5;132m",
        424: "\033[38;5;205m", 425: "\033[38;5;133m", 426: "\033[38;5;206m", 428: "\033[38;5;134m",
        429: "\033[38;5;207m", 431: "\033[38;5;135m", 451: "\033[38;5;208m",
        500: "\033[38;5;201m", 501: "\033[38;5;200m", 502: "\033[38;5;165m", 503: "\033[38;5;199m",
        504: "\033[38;5;164m", 505: "\033[38;5;198m", 506: "\033[38;5;163m", 507: "\033[38;5;197m",
        508: "\033[38;5;162m", 510: "\033[38;5;196m", 511: "\033[38;5;161m",
}

func printBanner() {
        fmt.Println("\033[38;5;206m")
        fmt.Println(`
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
        `)
        fmt.Println("\033[38;5;46mHyperScanner v1.1\033[0m\n")
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
        f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
        if err != nil {
                fmt.Println("Write error:", err)
                return
        }
        defer f.Close()
        f.WriteString(line + "\n")
}

func colorPrint(ip string, code int, desc string, err error) {
        if err != nil {
                fmt.Printf("\033[31m[x]\033[0m %s -> \033[31mERROR: %v\033[0m\n", ip, err)
                return
        }
        emoji := statusEmojis[code/100]
        color := statusColors[code]
        fmt.Printf("[✓] %s -> %s%d\033[0m %s%s\033[0m\n", ip, color, code, emoji, desc)
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

func main() {
        printBanner()

        ipInput := flag.String("i", "", "IP list file (example: ips.txt)")
        fileInput := flag.String("f", "", "Alias for -i (optional)")
        helpFlag := flag.Bool("h", false, "Show help")
        flag.Parse()

        if *helpFlag {
                fmt.Println("Usage:")
                fmt.Println("  -i <file>     Specify IP list file")
                fmt.Println("  -f <file>     Alias for -i")
                fmt.Println("  -h            Show help")
                os.Exit(0)
        }

        if *ipInput == "" && *fileInput == "" {
                fmt.Println("Error: No input file provided. Use -i <file> or -f <file>")
                os.Exit(1)
        }

        ipListPath := *ipInput
        if *fileInput != "" {
                ipListPath = *fileInput
        }

        outputDir := strings.TrimSuffix(ipListPath, ".txt") + " output"
        os.MkdirAll(outputDir, os.ModePerm)

        err := createOutputStructure(outputDir)
        if err != nil {
                fmt.Println("Error creating folders:", err)
                os.Exit(1)
        }

        processIPs(ipListPath, outputDir)

        fmt.Println("\n✅ Scan complete. Output saved to:", outputDir)
}