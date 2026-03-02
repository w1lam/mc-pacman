package installer

import (
	"context"

	"github.com/w1lam/mc-pacman/internal/app/base"
	"github.com/w1lam/mc-pacman/internal/app/getter"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
)

// Installer installer useCase
type Installer struct {
	base.UseCase

	paths *paths.Paths

	repo state.Repo

	getter *getter.Getter
}

func New(p *paths.Paths, repo state.Repo, g *getter.Getter) *Installer {
	return &Installer{
		UseCase: base.UseCase{
			Scope:   events.ScopeInstaller,
			Emitter: nil,
		},
		paths: p,

		repo: repo,

		getter: g,
	}
}

func (i *Installer) Install(ctx context.Context, pkgID string) error {
	err := i.getter.Get(ctx, pkgID)
	if err != nil {
		return err
	}

	return nil
}
