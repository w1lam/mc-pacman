package paths

import "os"

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
