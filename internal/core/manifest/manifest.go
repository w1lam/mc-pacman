// Package manifest contains manifest
package manifest

import (
	"time"

	packages "github.com/w1lam/mc-pacman/internal/core/packages"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion     int                            `json:"schemaVersion"`
	EnabledPackages   EnabledPackages                `json:"enabledPackages"`
	InstalledPackages packages.InstalledPackageIndex `json:"installedPackages"`
	InstalledLoaders  []LoaderInfo                   `json:"installedLoader"`
	Backups           []BackupEntry                  `json:"backups"`
	Initialized       bool                           `json:"initialized"`
}

// EnabledPackages asda
type EnabledPackages map[packages.PkgType]packages.PkgID

// BackupEntry asda
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
