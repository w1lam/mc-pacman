package core

import (
	"github.com/w1lam/Packages/modrinth"
)

// PackageType is the type of a package modpack/resourcebundle/shaderbundle?
type PackageType struct {
	EntryType  modrinth.EntryType `json:"entryType"`
	ActivePath string             `json:"activePath"`
	StorageDir string             `json:"storageDir"` // package types storage dir ie .mc-pacman/modpacks
}

// AvailablePackages are all available packages found in github repo
type AvailablePackages struct {
	ModPacks            map[string]ResolvedPackage
	ResourcePackBundles map[string]ResolvedPackage
	ShaderPackBundles   map[string]ResolvedPackage
}

// Pkg is a small pacakge struct used for passing around packages
type Pkg struct {
	Name string
	ID   string
	Type PackageType
}

// ResolvedPackage is a resolved package that gets passed to the installer
type ResolvedPackage struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	ListVersion string `json:"listVersion"`
	McVersion   string `json:"mcVersion"`
	Loader      string `json:"loader"`
	Env         string `json:"env"`
	Description string `json:"description"`

	Type       PackageType `json:"pkgType"`
	ListSource string      `json:"listSource"`
	Hash       string      `json:"hash"` // sha512

	Entries []modrinth.Entry `json:"entries"`
}

// InstalledPackage is an installed package which holds all information about the package
type InstalledPackage struct {
	Name             string                           `json:"name"`
	ID               string                           `json:"id"`
	Type             PackageType                      `json:"type"`
	ListSource       string                           `json:"listSource"`
	InstalledVersion string                           `json:"version"`
	McVersion        string                           `json:"mcVersion"`
	Loader           string                           `json:"loader"`
	Hash             string                           `json:"hash"`
	Entries          map[string]InstalledPackageEntry `json:"installedEntries"`

	ActivePath  string `json:"activePath"`
	StoragePath string `json:"storagePath"`
}

// PackageEntry is a mod entry in the manifest that holds all information about an entry
type InstalledPackageEntry struct {
	Name             string `json:"name"`
	ID               string `json:"id"` // id or slug
	FileName         string `json:"fileName"`
	Sha512           string `json:"sha512"`
	Sha1             string `json:"sha1,omitempty"`
	InstalledVersion string `json:"InstalledVersion"`
}
