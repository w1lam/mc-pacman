package verifier

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/packageid"
)

// ScanActitvePackages scans all active packages directories looking for a *PkgID*.id.json
func (v *Verifier) scanActitvePackages() ([]packages.InstalledPackage, error) {
	var found []packages.InstalledPackage

	for _, dir := range []string{
		v.paths.ModsDir,
		v.paths.ResourcePacksDir,
		v.paths.ShaderPacksDir,
	} {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			pkg, err := ScanDirForPackageID(filepath.Join(dir, e.Name()))
			if err != nil {
				errors.ReportCtx(
					"startup.scan.active",
					err,
					map[string]string{"dir": dir},
				)
				continue
			}

			found = append(found, pkg)
		}
	}

	return found, nil
}

// ScanDirForPackageID scans a directory looking for a *PkgID*.id.json file
func ScanDirForPackageID(dir string) (packages.InstalledPackage, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return packages.InstalledPackage{}, nil
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if strings.HasSuffix(e.Name(), ".id.json") {
			idPath := filepath.Join(dir, e.Name())
			return packageid.ReadPackageIDFile(idPath)
		}
	}

	return packages.InstalledPackage{}, fmt.Errorf("no id file found in: %s", dir)
}
