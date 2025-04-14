package installe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Main function to handle global installation
func main() {
	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Failed to get executable path:", err)
		return
	}

	var destPath string

	// Set the installation path depending on the OS
	switch runtime.GOOS {
	case "windows":
		destPath = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "hxscanner.exe")
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
		fmt.Println("❌ Unsupported OS")
		return
	}

	// Check for root privileges on Linux/macOS for moving binary
	if !checkSudo() && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		fmt.Println("⚠️ This operation requires sudo privileges. Please run as root.")
		return
	}

	// Copy the binary to the destination path
	err = copyFile(execPath, destPath)
	if err != nil {
		fmt.Println("❌ Failed to copy binary:", err)
		return
	}

	// Make the binary executable
	err = os.Chmod(destPath, 0755)
	if err != nil {
		fmt.Println("❌ Failed to make binary executable:", err)
		return
	}

	fmt.Println("✅ HyperScanner successfully installed globally at:", destPath)
}

// isTermux checks if the environment is Termux (Android)
func isTermux() bool {
	return strings.Contains(os.Getenv("PREFIX"), "com.termux")
}

// isArch checks if the platform is Arch Linux
func isArch() bool {
	cmd := exec.Command("uname", "-a")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "arch")
}

// checkSudo checks if the program is running with sudo privileges
func checkSudo() bool {
	cmd := exec.Command("whoami")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Error checking sudo:", err)
		return false
	}
	return strings.TrimSpace(string(out)) == "root"
}

// copyFile copies the binary from source to destination
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0755)
}
