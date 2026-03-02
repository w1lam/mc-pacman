// Package paths builds all paths
package paths

import (
	"path/filepath"
)

// Paths is the object holding root paths.
type Paths struct {
	minecraftDir string
	rootDir      string
}

// New creates a new Paths object.
func New(root, mcDir string) *Paths {
	return &Paths{
		minecraftDir: mcDir,
		rootDir:      root,
	}
}

// --- DIRECTORIES ---

// MC

// McDir returns the oaths to minecraft directory.
func (p *Paths) McDir() string {
	return p.minecraftDir
}

// ModsDir resturns the path to mods directory.
func (p *Paths) ModsDir() string {
	return filepath.Join(p.minecraftDir, "mods")
}

// ResourcepackDir returns the path to resourcepacks directory.
func (p *Paths) ResourcepackDir() string {
	return filepath.Join(p.minecraftDir, "resourcepacks")
}

// ShaderpackDir returns the path to shaderpacks directory.
func (p *Paths) ShaderpackDir() string {
	return filepath.Join(p.minecraftDir, "shaderpacks")
}

// APP

// RootDir returns the path to root directory of app.
func (p *Paths) RootDir() string {
	return p.rootDir
}

// BinDir returns the path to bin directory of app.
func (p *Paths) BinDir() string {
	return filepath.Join(p.rootDir, "bin")
}

// DataDir returns the path to data directory of app.
func (p *Paths) DataDir() string {
	return filepath.Join(p.rootDir, "data")
}

// PackagesDir returns the path to packages directory of app.
func (p *Paths) PackagesDir() string {
	return filepath.Join(p.rootDir, "packages")
}

// BackupsDir returns the path to backups directory of app.
func (p *Paths) BackupsDir() string {
	return filepath.Join(p.rootDir, "backups")
}

// --- FILES ---

// LogFile returns the path to log file of app.
func (p *Paths) LogFile() string {
	return filepath.Join(p.rootDir, "logs.log")
}

// StateFile returns the path to state file of app.
func (p *Paths) StateFile() string {
	return filepath.Join(p.DataDir(), "state.json")
}
