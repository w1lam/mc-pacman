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

	if !utils.CheckFileExists(path.RootDir) {
		err := os.MkdirAll(path.RootDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.DataDir) {
		err := os.MkdirAll(path.DataDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.PackagesDir) {
		err := os.MkdirAll(path.PackagesDir, 0o755)
		if err != nil {
			return err
		}
	}

	if !utils.CheckFileExists(path.BackupsDir) {
		err := os.MkdirAll(path.BackupsDir, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}
