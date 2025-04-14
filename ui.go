package ui

import (
	"fmt"
)

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Cyan    = "\033[36m"
	White   = "\033[97m"
	Magenta = "\033[35m"
	Reset   = "\033[0m"
	Bold    = "\033[1m"
)

func ShowBanner() {
	banner := `
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

HyperScanner v1.1 🔥  -  Ultra Fast HTTP Status Scanner
`
	fmt.Println(string(Cyan) + Bold + banner + Reset)
}

func ShowHelp() {
	fmt.Println(string(Yellow) + Bold + "\nUsage:")
	fmt.Println(string(White) + "  hyperscanner -i input.txt -json results.json -csv results.csv -c 100 -r 2")
	fmt.Println(string(Yellow) + Bold + "\nOptions:")
	fmt.Println(string(Green) + "  -i, -f" + string(White) + "      Input file containing list of targets (IP or domain)")
	fmt.Println(string(Green) + "  -json" + string(White) + "     Output file for JSON results (default: output.json)")
	fmt.Println(string(Green) + "  -csv" + string(White) + "      Output file for CSV results (default: output.csv)")
	fmt.Println(string(Green) + "  -c" + string(White) + "         Concurrency level (default: 100)")
	fmt.Println(string(Green) + "  -r" + string(White) + "         Max retries per request (default: 1)")
	fmt.Println(string(Green) + "  -h" + string(White) + "         Show help info")
	fmt.Println(string(Reset))
}

func PrintStatus(url string, status int) {
	switch {
	case status >= 200 && status < 300:
		fmt.Printf("%s[+]%s %s => %s%d OK%s\n", Green, Reset, url, Green, status, Reset)
	case status >= 300 && status < 400:
		fmt.Printf("%s[~]%s %s => %s%d Redirect%s\n", Cyan, Reset, url, Cyan, status, Reset)
	case status >= 400 && status < 500:
		fmt.Printf("%s[!]%s %s => %s%d Client Error%s\n", Yellow, Reset, url, Yellow, status, Reset)
	case status >= 500:
		fmt.Printf("%s[✖]%s %s => %s%d Server Error%s\n", Red, Reset, url, Red, status, Reset)
	default:
		fmt.Printf("%s[?]%s %s => %s%d Unknown%s\n", Magenta, Reset, url, Magenta, status, Reset)
	}
}
