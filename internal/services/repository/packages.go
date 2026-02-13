package repository

import (
	"fmt"
	"sync"

	"github.com/w1lam/mc-pacman/internal/core/netcfg"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type RemotePackageRepository struct {
	mu     sync.RWMutex
	cache  packages.RemotePackageIndex
	loaded bool
}

// GetSpecificPackage searches for the specified package in githubrepo
func GetSpecificPackage(pkgID packages.PkgID) (packages.RemotePackage, error) {
	folders, err := ScanPackagesFolder()
	if err != nil {
		return packages.RemotePackage{}, err
	}

	for _, folder := range folders {
		url := fmt.Sprintf("%scontents/packages/%s/%s.json", netcfg.GithubRepo, folder, pkgID)

		var resp GithubContentResponse
		if err := GithubGetJSON(url, &resp); err != nil {
			continue
		}

		if resp.Type != "file" {
			continue
		}

		return ResolvePackageJSON(resp.RawURL)
	}

	return packages.RemotePackage{}, fmt.Errorf("package not found: %s", pkgID)
}

func GetAllAvailablePackages() (packages.RemotePackageIndex, error) {
	folders, err := ScanPackagesFolder()
	if err != nil {
		return nil, err
	}

	out := packages.BlankRemotePackageIndex()

	for _, folder := range folders {
		pkgs, err := GetPackagesFromFolder(folder)
		if err != nil {
			return nil, err
		}

		for _, pkg := range pkgs {
			out[pkg.Type.PackageType][pkg.ID] = pkg
		}
	}

	return out, nil
}
