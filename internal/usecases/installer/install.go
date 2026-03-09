package installer

import (
	"context"
	"fmt"

	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/usecases"
	"github.com/w1lam/mc-pacman/internal/usecases/getter"
)

// Installer installer useCase
type Installer struct {
	usecases.Base

	paths  *paths.Paths
	sRepo  state.Repo
	getter *getter.Getter
}

func New(base usecases.Base, p *paths.Paths, repo state.Repo, g *getter.Getter) *Installer {
	i := Installer{
		Base:  base,
		paths: p,

		sRepo: repo,

		getter: g,
	}

	return &i
}

func (i *Installer) Install(ctx context.Context, pkgID string) error {
	if err := i.paths.Validate(); err != nil {
		return err
	}

	op := i.StartOp(events.Operation{}, fmt.Sprintf("install_package_%s", pkgID))
	i.EmitStart(op, fmt.Sprintf("starting installation of %s", pkgID))
	defer i.EmitEnd(op)

	err := i.getter.Get(events.WithOp(ctx, op), pkgID)
	if err != nil {
		return err
	}

	return nil
}
