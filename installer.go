package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	binaryName := "hxscanner"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fmt.Println("[+] Building the hxscanner tool...")
	build := exec.Command("go", "build", "-o", binaryName, "main.go", "init.go", "ui.go", "hxs.go")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	err := build.Run()
	if err != nil {
		fmt.Println("[-] Build failed:", err)
		os.Exit(1)
	}

	destDir := "/usr/local/bin"
	if runtime.GOOS == "windows" {
		destDir = filepath.Join(os.Getenv("SystemDrive"), "\\Windows\\System32")
	}
	destPath := filepath.Join(destDir, binaryName)

	srcFile, err := os.Open(binaryName)
	if err != nil {
		fmt.Println("[-] Error opening built binary:", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("[-] Cannot write to %s. Try running with sudo.\n", destPath)
		os.Exit(1)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("[-] Installation failed:", err)
		os.Exit(1)
	}

	fmt.Println("[✔] Installed successfully! Now you can run \"hxscanner\" from anywhere.")
}
