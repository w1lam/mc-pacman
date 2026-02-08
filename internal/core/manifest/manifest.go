package core

import (
	"time"

	"github.com/w1lam/Packages/modrinth"
	packages "github.com/w1lam/mc-pacman/internal/core/packages"
	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

// Manifest is the manifest for all global information required by the program
type Manifest struct {
	SchemaVersion  int    `json:"schemaVersion"`
	ProgramVersion string `json:"programVersion"`

	EnabledPackages map[packages.PackageType]string `json:"enabledPackages"`

	InstalledPackages map[packages.PackageType]map[string]InstalledPackage `json:"installedPackages"`

	InstalledLoaders map[string]LoaderInfo `json:"installedLoader"`
	Paths            *paths.Paths          `json:"-"`
	Backups          []BackupEntry
	Initialized      bool `json:"initialized"`
}

type BackupEntry struct {
	Time time.Time
	Type packages.PackageType
	Path string
	ID   string
}

// LoaderInfo is the information about a mod loader
type LoaderInfo struct {
	Loader        string `json:"loader"`
	McVersion     string `json:"mcVersion"`
	LoaderVersion string `json:"version"`
}

// Updates is all info on available updates
type Updates struct {
	ModListUpdate map[string]bool                   `json:"-"`
	ModUpdates    map[string][]modrinth.UpdateEntry `json:"-"`
}
