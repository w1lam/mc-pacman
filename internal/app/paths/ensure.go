package paths

import (
	"errors"
	"fmt"
	"os"
)

// Ensure ensures all directories exist
func (p *Paths) Ensure() error {
	dirs := []string{
		p.RootDir(),
		p.BinDir(),
		p.DataDir(),
		p.PackagesDir(),
		p.BackupsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}

var ErrMcDirNotConfigured error = errors.New("mcDir not configured")

// Validate validates the minecraft directory in Paths
func (p *Paths) Validate() error {
	if p.minecraftDir == "" {
		return ErrMcDirNotConfigured
	}
	if _, err := os.Stat(p.minecraftDir); err != nil {
		return fmt.Errorf("minecraft directory does not exist: %s", p.minecraftDir)
	}
	return nil
}
