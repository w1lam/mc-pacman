// Package state contains global state
package state

import (
	"time"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

// State is the global state of the program
type State struct {
	SchemaVersion     int                                   `json:"schemaVersion"`
	McDir             string                                `json:"mcDir"`
	EnabledPackageIDs map[packages.PkgTypeID]packages.PkgID `json:"enabledPackages"`
	InstalledLoaders  []LoaderInfo                          `json:"installedLoader"`
	Backups           []BackupEntry                         `json:"backups"`
}

// BackupEntry asda
type BackupEntry struct {
	PkgID       packages.PkgID     `json:"pkgID"`
	Time        time.Time          `json:"time"`
	Type        packages.PkgTypeID `json:"pkgType"`
	Path        string             `json:"path"`
	GeneratedID string             `json:"generatedID"`
}

// LoaderInfo is the information about a mod loader
type LoaderInfo struct {
	Loader        string `json:"loader"`
	McVersion     string `json:"mcVersion"`
	LoaderVersion string `json:"version"`
}
