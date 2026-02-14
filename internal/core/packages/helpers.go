package packages

// GetHashes gets hashes of entries in an installed package
func (p *InstalledPackage) GetHashes() []string {
	var hashes []string
	for _, mod := range p.Entries {
		hashes = append(hashes, mod.Hash)
	}
	return hashes
}
