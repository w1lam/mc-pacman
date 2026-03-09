// Package updater holds updater service
package updater

import (
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/usecases"
)

// TODO: IMPLEMENT UPDATER

type Updater struct {
	usecases.Base

	iRepo packages.InstalledRepo
	rRepo packages.RemoteRepo
}

func New(base usecases.Base, iRepo packages.InstalledRepo, rRepo packages.RemoteRepo) *Updater {
	u := Updater{
		Base:  base,
		iRepo: iRepo,
		rRepo: rRepo,
	}

	return &u
}
