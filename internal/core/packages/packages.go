package packages

import (
	"github.com/w1lam/Packages/modrinth"
)

// AvailablePackages are all available packages found in github repo
type AvailablePackages struct {
	ModPacks            map[string]ResolvedPackage
	ResourcePackBundles map[string]ResolvedPackage
	ShaderPackBundles   map[string]ResolvedPackage
	DatapackBundles     map[string]ResolvedPackage
}

// Pkg is a small pacakge struct used for passing around packages
type Pkg struct {
	Name string
	ID   PkgID
	Type PackageType
}

// PkgID is the ID of a package
type PkgID string

// ResolvedPackage is a resolved package that gets passed to the installer
type ResolvedPackage struct {
	Name string `json:"name"`
	ID   PkgID  `json:"id"`

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
	Name             string `json:"name"`
	ID               PkgID  `json:"id"`
	InstalledVersion string `json:"version"`
	McVersion        string `json:"mcVersion"`
	Loader           string `json:"loader"`

	Type PackageType `json:"type"`

	ListSource string                           `json:"listSource"`
	Hash       string                           `json:"hash"`
	Entries    map[string]InstalledPackageEntry `json:"installedEntries"`

	FullActivePath  string `json:"activePath"`
	FullStoragePath string `json:"storagePath"`
}

// PackageEntry is a mod entry in the manifest that holds all information about an entry
type InstalledPackageEntry struct {
	Name             string `json:"name"`
	ID               string `json:"id"` // id or slug
	InstalledVersion string `json:"InstalledVersion"`

	FileName string `json:"fileName"`
	Sha512   string `json:"sha512"`
	Sha1     string `json:"sha1,omitempty"`
}
