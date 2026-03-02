package state

// // Installed queries the state for installed packageIDs
// func (m *State) Installed() []packages.PkgID {
// 	return m.InstalledPackageIDs
// }
//
// // Installed queries the manifest for installed packages
// func (m *Manifest) IsInstalled(id packages.PkgID) bool {
// 	ok := false
// 	for _, pkgs := range m.InstalledPackages {
// 		_, ok = pkgs[id]
// 		if ok {
// 			break
// 		}
// 	}
// 	return ok
// }
//
// // EnabledOfType returns enabled package of specified type, and false if it does not exist
// func (m *Manifest) EnabledOfType(t packages.PkgType) (packages.InstalledPackage, bool) {
// 	id, ok := m.EnabledPackages[t]
// 	if !ok {
// 		return packages.InstalledPackage{}, false
// 	}
//
// 	pkgs, ok := m.InstalledPackages[t]
// 	if !ok {
// 		return packages.InstalledPackage{}, false
// 	}
//
// 	pkg, ok := pkgs[id]
// 	return pkg, ok
// }
//
// // SetEnabled sets the specified package as enabled
// func (m *Manifest) SetEnabled(id packages.PkgID) error {
// 	for t, pkgs := range m.InstalledPackages {
// 		_, ok := pkgs[id]
// 		if ok {
// 			m.EnabledPackages[t] = id
// 			return nil
// 		}
// 	}
//
// 	return fmt.Errorf("specified package not installed: %s", id)
// }
