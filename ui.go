package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	gray    = color.New(color.FgHiBlack).SprintFunc()
	red     = color.New(color.FgHiRed).SprintFunc()
	green   = color.New(color.FgHiGreen).SprintFunc()
	yellow  = color.New(color.FgHiYellow).SprintFunc()
	cyan    = color.New(color.FgHiCyan).SprintFunc()
	blue    = color.New(color.FgHiBlue).SprintFunc()
	magenta = color.New(color.FgHiMagenta).SprintFunc()
	white   = color.New(color.FgHiWhite).SprintFunc()
)

func PrintBanner() {
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
	fmt.Println(cyan(banner))
}

func LogResult(url string, code int) {
	var status string

	switch {
	case code >= 100 && code < 200:
		status = blue(fmt.Sprintf("[%d]", code))
	case code >= 200 && code < 300:
		status = green(fmt.Sprintf("[%d]", code))
	case code >= 300 && code < 400:
		status = yellow(fmt.Sprintf("[%d]", code))
	case code >= 400 && code < 500:
		status = red(fmt.Sprintf("[%d]", code))
	case code >= 500:
		status = magenta(fmt.Sprintf("[%d]", code))
	default:
		status = gray("[???]")
	}

	fmt.Printf("➡️  %s -> %s\n", url, status)
}

func LogInfo(msg string) {
	fmt.Println(white("[*] "), msg)
}

func LogError(msg string) {
	fmt.Println(red("[!] "), msg)
}

func LogSuccess(msg string) {
	fmt.Println(green("[✓] "), msg)
}
