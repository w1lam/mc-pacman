// Package activation handles package enabling/disabling
package activation

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// TODO: IMPLEMENT ACTIVATOR IE ENABLE/DISABLE

type Activator struct {
	events.EmitterBase
	paths *paths.Paths
	iRepo packages.InstalledRepo
}

func New(view ux.View, iRepo packages.InstalledRepo, paths *paths.Paths) *Activator {
	a := Activator{
		paths: paths,
		iRepo: iRepo,
	}
	a.SetEmitter(view)
	return &a
}

// Enable enables the specified package
func (a *Activator) Enable(pkgID string) error {
	return nil
}

// Disable disables a package of specified package type. packages types "modpack", "resourcebundle", "shaderbundle"
func (a *Activator) Disable(pkgType string) error {
	return nil
}
