package core

import (
	"fmt"

	packages "github.com/w1lam/mc-pacman/internal/core/packages"
)

// InitManifest initializes the manifest loading it if it exists and creating blank if not
func InitManifest(path string) (*Manifest, error) {
	m, err := Load(path)
	if err != nil {
		m, err = blankManifest(path)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// BlankManifest builds a blank manifest and saves it to path
func blankManifest(path string) (*Manifest, error) {
	fmt.Printf(" * Building Initial Manifest...\n")

	m := Manifest{
		SchemaVersion: 1,

		EnabledPackages: EnabledPackages{
			Modpack:        "",
			ResourceBundle: "",
			ShaderBundle:   "",
			DatapackBundle: "",
		},

		InstalledPackages: InstalledPackages{
			Modpacks:        make(map[string]packages.InstalledPackage),
			ResourceBundles: make(map[string]packages.InstalledPackage),
			ShaderBundles:   make(map[string]packages.InstalledPackage),
			DatapackBundles: make(map[string]packages.InstalledPackage),
		},

		InstalledLoaders: []LoaderInfo{},

		Backups:     []BackupEntry{},
		Initialized: false,
		Path:        path,
	}

	return &m, m.Save()
}
