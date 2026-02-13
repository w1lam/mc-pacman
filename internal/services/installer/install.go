package installer

import (
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
)

type Installer struct {
	paths      *paths.Paths
	manifest   *manifest.Repository
	downloader *downloader.Service
}

func (i *Installer) Install(pkg packages.RemotePackage) error {
	if err := validate(pkg); err != nil {
		return err
	}

	m.InstalledPackages[pkg.Type][pkg.ID] = installed
	return i.repo.Save(m)
}
