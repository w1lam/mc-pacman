package filesystem

import (
	"os"

	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

// EnsureDirectories ensures all program directories exists
func EnsureDirectories(path *paths.Paths) error {
	dirs := []string{
		path.RootDir,
		path.DataDir,
		path.PackagesDir,
		path.BackupsDir,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
	}
	return nil
}
