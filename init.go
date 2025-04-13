package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"path/filepath"
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get executable path:", err)
		return
	}

	var destPath string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		destPath = filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "hxscanner.exe")
	case "darwin": // macOS
		destPath = "/usr/local/bin/hxscanner"
	case "linux":
		if isTermux() {
			home := os.Getenv("HOME")
			destPath = filepath.Join(home, ".termux", "bin", "hxscanner")
		} else if isArch() {
			destPath = "/usr/bin/hxscanner"
		} else {
			destPath = "/usr/local/bin/hxscanner"
		}
	default:
		fmt.Println("Unsupported OS")
		return
	}

	err = copyFile(execPath, destPath)
	if err != nil {
		fmt.Println("Failed to copy binary:", err)
		return
	}

	err = os.Chmod(destPath, 0755)
	if err != nil {
		fmt.Println("Failed to make binary executable:", err)
	}

	fmt.Println("✅ HyperScanner installed globally at:", destPath)
}

func isTermux() bool {
	return strings.Contains(os.Getenv("PREFIX"), "com.termux")
}

func isArch() bool {
	cmd := exec.Command("uname", "-a")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "arch")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0755)
}
