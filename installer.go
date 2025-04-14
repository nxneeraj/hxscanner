package installer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CheckAndSetup() {
	binName := "hxscanner"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	path, err := exec.LookPath(binName)
	if err == nil {
		fmt.Println("✅ Already in PATH:", path)
		return
	}

	fmt.Println("⚙️  Installing globally...")

	binPath := "./" + binName
	globalPath := getInstallPath() + "/" + binName

	err = os.Rename(binPath, globalPath)
	if err != nil {
		fmt.Println("❌ Move failed:", err)
		return
	}

	fmt.Println("✅ Installed at", globalPath)
	fmt.Println("🔁 Please restart terminal or add to PATH manually if needed.")
}

func getInstallPath() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("USERPROFILE") + "\\AppData\\Local\\Microsoft\\WindowsApps"
	case "darwin":
		return "/usr/local/bin"
	default:
		return "/usr/local/bin"
	}
}
