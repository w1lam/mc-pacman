package core

import (
	"time"

	packages "github.com/w1lam/mc-pacman/internal/core/packages"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion int `json:"schemaVersion"`

	EnabledPackages EnabledPackages `json:"enabledPackages"`

	InstalledPackages InstalledPackages `json:"installedPackages"`

	InstalledLoaders []LoaderInfo `json:"installedLoader"`

	Backups     []BackupEntry `json:"backups"`
	Initialized bool          `json:"initialized"`
	Path        string        `json:"path"`
}

// EnabledPackages type
type EnabledPackages struct {
	Modpack        string `json:"modpack"`
	ResourceBundle string `json:"resourcebundle"`
	ShaderBundle   string `json:"shaderbundle"`
	DatapackBundle string `json:"datapackbundle"` // not really used
}

// InstalledPackages type
type InstalledPackages struct {
	Modpacks        map[string]packages.InstalledPackage `json:"modpacks"`
	ResourceBundles map[string]packages.InstalledPackage `json:"resourceBundles"`
	ShaderBundles   map[string]packages.InstalledPackage `json:"shaderBundles"`
	DatapackBundles map[string]packages.InstalledPackage `json:"datapackBundles"` // not really used
}

// BackupEntry
type BackupEntry struct {
	PkgID       packages.PkgID       `json:"pkgID"`
	Time        time.Time            `json:"time"`
	Type        packages.PackageType `json:"pkgType"`
	Path        string               `json:"path"`
	GeneratedID string               `json:"generatedID"`
}

// LoaderInfo is the information about a mod loader
type LoaderInfo struct {
	Loader        string `json:"loader"`
	McVersion     string `json:"mcVersion"`
	LoaderVersion string `json:"version"`
}
