package filesystem

import (
	"fmt"
	"os"

	"github.com/w1lam/Packages/utils"
	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

// EnsureDirectories ensures all program directories exists
func EnsureDirectories(path *paths.Paths) error {
	if !utils.CheckFileExists(path.MinecraftDir) {
		return fmt.Errorf("minecraft directory not found")
	}

	dirs := []string{
		path.RootDir,
		path.DataDir,
		path.PackagesDir,
		path.BackupsDir,
	}

	for _, dir := range dirs {
		if !utils.CheckFileExists(dir) {
			err := os.MkdirAll(dir, 0o755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
