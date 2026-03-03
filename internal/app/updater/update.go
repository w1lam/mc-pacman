// Package updater holds updater service
package updater

import (
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// TODO: IMPLEMENT UPDATER

type Updater struct {
	events.EmitterBase

	iRepo packages.InstalledRepo
	rRepo packages.RemoteRepo
}

func New(view ux.View, iRepo packages.InstalledRepo, rRepo packages.RemoteRepo) *Updater {
	u := Updater{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeUpdater,
		},
		iRepo: iRepo,
		rRepo: rRepo,
	}

	u.SetEmitter(view)
	return &u
}
