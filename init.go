package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	urlFile     string
	concurrency int
	timeout     int
	showAll     bool
	delay       int
	userAgent   string
)

func init() {
	flag.StringVar(&urlFile, "f", "", "📁 Path to file containing target URLs")
	flag.IntVar(&concurrency, "c", 50, "🚀 Number of concurrent requests")
	flag.IntVar(&timeout, "t", 10, "⏱️ Request timeout in seconds")
	flag.BoolVar(&showAll, "a", false, "🧾 Show all status codes (including non-2xx/3xx)")
	flag.IntVar(&delay, "d", 0, "🕒 Delay between requests in milliseconds")
	flag.StringVar(&userAgent, "ua", "HyperScanner/1.1", "🕵️ Custom User-Agent header")

	// Automatically called before main()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
██╗  ██╗██╗  ██╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗ 
██║  ██║╚██╗██╔╝     ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
███████║ ╚███╔╝█████╗███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║ ██╔██╗╚════╝╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║██╔╝ ██╗     ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

        HyperScanner v1.1 🔥 - Ultra Fast HTTP Status Scanner by Neeraj Sah

USAGE:
    hxscanner -f urls.txt [options]

OPTIONS:
`)
		flag.PrintDefaults()
	}

	// Small delay to ensure flag.Usage shows up cleanly
	time.Sleep(200 * time.Millisecond)
}
