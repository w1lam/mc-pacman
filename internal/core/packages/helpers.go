package packages

// GetHashes gets hashes of entries in an installed package
func (p *InstalledPackage) GetHashes() []string {
	var hashes []string
	for _, mod := range p.Entries {
		hashes = append(hashes, mod.Hash)
	}
	return hashes
}

func NewInstalledPackage(
	name string,
	id PkgID,
	instVer,
	mcVer,
	loader string,
	pkgType PkgTypeID,
	listSource,
	hash string,
	entries []InstalledPackageEntry,
) *InstalledPackage {
	e := make(map[EntryID]InstalledPackageEntry)
	return &InstalledPackage{
		Name:             name,
		ID:               id,
		InstalledVersion: instVer,
		McVersion:        mcVer,
		Loader:           loader,
		Type:             pkgType,
		ListSource:       listSource,
		Hash:             hash,
		Entries:          e,
	}
}
