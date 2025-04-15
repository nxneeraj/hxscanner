package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Install function to set up the binary globally
func install() {
	// Get the Go binary directory for the current OS
	goBinDir := filepath.Join(os.Getenv("GOPATH"), "bin")
	if goBinDir == "" {
		// If GOPATH is not set, default to $HOME/go/bin
		goBinDir = filepath.Join(os.Getenv("HOME"), "go", "bin")
	}

	// Check if the binary directory exists
	if _, err := os.Stat(goBinDir); os.IsNotExist(err) {
		fmt.Println("Go binary directory not found. Creating it...")
		err := os.MkdirAll(goBinDir, 0755)
		if err != nil {
			fmt.Println("Error creating Go binary directory:", err)
			return
		}
	}

	// Ensure the binary path is in the system PATH
	if !strings.Contains(os.Getenv("PATH"), goBinDir) {
		fmt.Println("Adding Go binary path to system PATH...")

		// Add Go binary path to system's PATH
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo 'export PATH=$PATH:%s' >> ~/.bashrc", goBinDir))
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error adding Go bin path to PATH:", err)
			return
		}
		fmt.Println("Go binary path added successfully to ~/.bashrc. Please restart your terminal.")
	}

	// Check if hxscanner is already installed (old version)
	_, err := exec.LookPath("hxscanner")
	if err == nil {
		fmt.Println("Old version of hxscanner found. Replacing with the latest version...")
	} else {
		fmt.Println("No previous version of hxscanner found. Installing...")
	}

	// Perform installation: build and copy binary to global path
	cmd := exec.Command("go", "install", "github.com/nxneeraj/hxscanner@latest")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error installing hxscanner:", err)
		return
	}

	fmt.Println("Installation successful! You can now use the 'hxscanner' command globally.")
}

// Main function to call the install function
func main() {
	install()
}
