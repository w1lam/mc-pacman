package installer

import (
	"context"
	"fmt"

	"github.com/w1lam/mc-pacman/internal/app/getter"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// Installer installer useCase
type Installer struct {
	events.EmitterBase

	paths *paths.Paths

	repo state.Repo

	getter *getter.Getter
}

func New(view ux.View, p *paths.Paths, repo state.Repo, g *getter.Getter) *Installer {
	i := Installer{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeInstaller,
		},
		paths: p,

		repo: repo,

		getter: g,
	}

	i.SetEmitter(view)
	return &i
}

func (i *Installer) Install(ctx context.Context, pkgID string) error {
	op := i.StartOp(events.Operation{}, fmt.Sprintf("install %s", pkgID))
	i.EmitStart(op, fmt.Sprintf("starting installation of %s", pkgID))
	defer i.EmitEnd(op)

	err := i.getter.Get(ctx, op, pkgID)
	if err != nil {
		return err
	}

	return nil
}
