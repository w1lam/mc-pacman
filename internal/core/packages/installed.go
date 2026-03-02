package packages

// InstalledPackage is an installed package which holds all information about the package
type InstalledPackage struct {
	Name             string `json:"name"`
	ID               PkgID  `json:"id"`
	InstalledVersion string `json:"version"`
	McVersion        string `json:"mcVersion"`
	Loader           string `json:"loader"`

	Type PkgTypeID `json:"pkgType"`

	ListSource string                            `json:"listSource"`
	Hash       string                            `json:"hash"`
	Entries    map[EntryID]InstalledPackageEntry `json:"installedEntries"`

	HasConfig bool `json:"hasConfig,omitempty"`
}

// InstalledPackageEntry is a mod entry in the manifest that holds all information about an entry
type InstalledPackageEntry struct {
	ID               EntryID `json:"id"` // id or slug
	InstalledVersion string  `json:"InstalledVersion"`

	FileName string `json:"fileName"`
	Hash     string `json:"hash"`
	Algo     string `json:"hashAlgo"`
}

type InstalledRepo interface {
	// Exists returns true if package exists
	Exists(PkgID) (bool, error)

	// GetAll gets all installed packages
	GetAll() ([]InstalledPackage, error)
	// GetByID gets an installed package with given PkgID
	GetByID(PkgID) (InstalledPackage, error)

	// Add creates package dir, moves entries to entries folder and writes pkg.json
	Add(p InstalledPackage, entriesSrcDir string) error
	// Update updates a pkg.json file overwriting it
	Update(InstalledPackage) error
	// Remove removes the a package completly from packages/
	Remove(PkgID) error
}
