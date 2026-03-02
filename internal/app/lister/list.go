// Package lister holds list service
package lister

import (
	"context"
	"fmt"

	"github.com/w1lam/mc-pacman/internal/app/base"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
)

type Lister struct {
	base.UseCase
	path       *paths.Paths
	repo       state.Repo
	remoteRepo packages.RemoteRepo
}

func New(path *paths.Paths, repo state.Repo) *Lister {
	return &Lister{
		path:       path,
		repo:       repo,
		remoteRepo: remote.New(),
	}
}

func (l *Lister) ListAllRemote(ctx context.Context) error {
	op := events.NewOperation(events.ScopeList, "remote")

	l.Emit(events.Event{
		Op:   op,
		Type: events.EventStart,
	})

	index, err := l.remoteRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	l.Emit(events.Event{
		Op:        op,
		Type:      events.EventComplete,
		ExtraData: index,
	})

	return nil
}

func (l *Lister) SearchPkg(ctx context.Context, id packages.PkgID) error {
	op := events.NewOperation(events.ScopeList, "remote")

	l.Emit(events.Event{
		Op:      op,
		Type:    events.EventStart,
		Message: "searching",
	})

	pkg, err := l.remoteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	l.Emit(events.Event{
		Op:      op,
		Type:    events.EventComplete,
		Message: fmt.Sprintf("found %s", pkg.ID),
	})

	return nil
}
