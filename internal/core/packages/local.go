package packages

import "github.com/w1lam/Packages/modrinth"

// InstalledPackage is an installed package which holds all information about the package
type InstalledPackage struct {
	Name             string `json:"name"`
	ID               PkgID  `json:"id"`
	InstalledVersion string `json:"version"`
	McVersion        string `json:"mcVersion"`
	Loader           string `json:"loader"`

	Type PackageType `json:"type"`

	ListSource string                                `json:"listSource"`
	Hash       string                                `json:"hash"`
	Entries    map[modrinth.ID]InstalledPackageEntry `json:"installedEntries"`

	FullActivePath  string `json:"activePath"`
	FullStoragePath string `json:"storagePath"`
}

// InstalledPackageEntry is a mod entry in the manifest that holds all information about an entry
type InstalledPackageEntry struct {
	ID               modrinth.ID `json:"id"` // id or slug
	InstalledVersion string      `json:"InstalledVersion"`

	FileName string `json:"fileName"`
	Hash     string `json:"hash"`
	Algo     string `json:"hashAlgo"`
}

// InstalledPackageIndex is an index of all packages
type InstalledPackageIndex map[PkgType]map[PkgID]InstalledPackage

// BlankInstalledPackageIndex creates an empty InstalledPackagesIndex and safely initiates all maps
func BlankInstalledPackageIndex() InstalledPackageIndex {
	index := make(InstalledPackageIndex)

	for pkgType := range PackageTypeIndex {
		index[pkgType] = make(map[PkgID]InstalledPackage)
	}

	return index
}
