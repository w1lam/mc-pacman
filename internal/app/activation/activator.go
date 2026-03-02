// Package activation handles package enabling/disabling
package activation

import (
	"github.com/w1lam/mc-pacman/internal/app/base"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
)

type Activator struct {
	base.UseCase
	paths *paths.Paths
	repo  state.Repo
}

func New(repo state.Repo, paths *paths.Paths) *Activator {
	return &Activator{
		paths: paths,
		repo:  repo,
	}
}

// Enable enables the specified package
func (a *Activator) Enable(pkgID string) error {
	m, err := a.repo.Load()
	if err != nil {
		return err
	}

	_ = m
	return nil
}

// Disable disables a package of specified package type. packages types "modpack", "resourcebundle", "shaderbundle"
func (a *Activator) Disable(pkgType string) error {
	m, err := a.repo.Load()
	if err != nil {
		return err
	}

	_ = m
	return nil
}
