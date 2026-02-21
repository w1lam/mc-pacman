package installer

import (
	"os"
	"path/filepath"
)

func (s *Installer) installToStorage(
	tempDir string,
	finalDir string,
) error {
	if err := os.MkdirAll(filepath.Dir(finalDir), 0755); err != nil {
		return err
	}

	if err := os.Rename(tempDir, finalDir); err != nil {
		return err
	}

	return nil
}
