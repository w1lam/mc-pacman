package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/w1lam/Packages/utils"
)

func RootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", ".mcpacman")
	case "linux":
		return filepath.Join(home, ".mcpacman")
	case "darwin":
		return filepath.Join("Library", "Application Support", "mcpacman")
	default:
		return filepath.Join(home, ".mcpacman")
	}
}

func DefaultMinecraftDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		dir := filepath.Join(home, "AppData", "Roaming", ".minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	case "linux":
		dir := filepath.Join(home, ".minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	case "darwin":
		dir := filepath.Join("Library", "Application Support", "minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	}

	return ""
}
