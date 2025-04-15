package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Installer handles moving the binary to a global path.
func Installer(binaryName string) error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary path: %v", err)
	}

	var destDir string

	switch runtime.GOOS {
	case "linux", "darwin":
		if os.Geteuid() == 0 {
			destDir = "/usr/local/bin"
		} else {
			home := os.Getenv("HOME")
			destDir = filepath.Join(home, ".local", "bin")
		}
	case "windows":
		// On Windows, move to GOPATH\bin or add to PATH manually
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
		}
		destDir = filepath.Join(gopath, "bin")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	dest := filepath.Join(destDir, binaryName)

	// Ensure destination directory exists
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if mkerr := os.MkdirAll(destDir, 0755); mkerr != nil {
			return fmt.Errorf("failed to create dir: %v", mkerr)
		}
	}

	// Copy the binary
	if err := copyBinary(binaryPath, dest); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}

	// Add to PATH permanently (Termux / Windows)
	addToPathIfNeeded(destDir)

	fmt.Printf("✅ Installed globally to: %s\n", dest)
	return nil
}

// copyBinary moves the built binary to the global destination.
func copyBinary(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0755)
}

// addToPathIfNeeded attempts to add the dir to PATH for Termux / Windows users.
func addToPathIfNeeded(dir string) {
	if runtime.GOOS == "android" {
		// For Termux
		bashrc := filepath.Join(os.Getenv("HOME"), ".bashrc")
		appendPathExport(bashrc, dir)
	} else if runtime.GOOS == "windows" {
		// Windows requires setting PATH via PowerShell
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`[Environment]::SetEnvironmentVariable("Path", "$env:Path;%s", "User")`, dir))
		_ = cmd.Run()
	}
}

func appendPathExport(filePath, dir string) {
	content := fmt.Sprintf("\nexport PATH=\"$PATH:%s\"\n", dir)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(content)
	}
}
