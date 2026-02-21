package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

func rootDir() string {
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
