// Package paths builds all paths
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	MinecraftDir     string
	ModsDir          string
	ResourcePacksDir string
	ShaderPacksDir   string

	RootDir string

	DataDir      string
	ManifestPath string
	MetaDataPath string

	PackagesDir string
	BackupsDir  string
	LogPath     string
}

func DefaultMinecraftDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", ".minecraft"), nil
	case "linux":
		return filepath.Join(home, ".minecraft"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "minecraft"), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Resolve() (*Paths, error) {
	mcDir, err := DefaultMinecraftDir()
	if err != nil {
		return nil, err
	}

	rootDir := filepath.Join(mcDir, ".mc-pacman")
	dataDir := filepath.Join(rootDir, "data")

	return &Paths{
		MinecraftDir:     mcDir,
		ModsDir:          filepath.Join(mcDir, "mods"),
		ResourcePacksDir: filepath.Join(mcDir, "resourcepacks"),
		ShaderPacksDir:   filepath.Join(mcDir, "shaderpacks"),

		RootDir: rootDir,

		DataDir:      dataDir,
		ManifestPath: filepath.Join(dataDir, "manifest.json"),
		MetaDataPath: filepath.Join(dataDir, "meta.json"),

		PackagesDir: filepath.Join(rootDir, "packages"),
		BackupsDir:  filepath.Join(rootDir, "backups"),
		LogPath:     filepath.Join(rootDir, "logs.log"),
	}, nil
}
