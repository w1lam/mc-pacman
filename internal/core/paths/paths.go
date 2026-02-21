// Package paths builds all paths
package paths

import (
	"path/filepath"
)

// Paths is the object holding all paths
type Paths struct {
	MinecraftDir     string `json:"mcDir"`
	ModsDir          string `json:"modsDir"`
	ResourcePacksDir string `json:"resourcePacksDir"`
	ShaderPacksDir   string `json:"shaderPacksDir"`

	RootDir   string `json:"rootDir"`
	PathsPath string `json:"pathsPath"`

	DataDir      string `json:"dataDir"`
	ManifestPath string `json:"manifestPath"`
	MetaDataPath string `json:"metaDatapath"`

	PackagesDir string `json:"packagesDir"`
	BackupsDir  string `json:"backupsDir"`
	LogPath     string `json:"logPath"`
}

// New creates a new Paths object
func New(root, mcDir string) *Paths {
	dataDir := filepath.Join(root, "data")
	return &Paths{
		MinecraftDir:     mcDir,
		ModsDir:          filepath.Join(mcDir, "mods"),
		ResourcePacksDir: filepath.Join(mcDir, "resourcepacks"),
		ShaderPacksDir:   filepath.Join(mcDir, "shaderpacks"),

		RootDir:   root,
		PathsPath: filepath.Join(root, "paths.json"),

		DataDir:      dataDir,
		ManifestPath: filepath.Join(dataDir, "manifest.json"),
		MetaDataPath: filepath.Join(dataDir, "meta.json"),

		PackagesDir: filepath.Join(root, "packages"),
		BackupsDir:  filepath.Join(root, "backups"),
		LogPath:     filepath.Join(root, "logs.log"),
	}
}

// Init initializes paths with its repo
func (r *JSONPathsRepository) Init() (*Paths, error) {
	p, err := r.Load()
	if err != nil {
		return nil, err
	}

	if p != nil {
		return p, nil
	}

	mcDir, err := detectOrPromptMinecraftDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Dir(r.file)
	p = New(root, mcDir)

	if err := r.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
