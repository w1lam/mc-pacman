// Package lister holds list service
package lister

import (
	"context"
	"fmt"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// TODO: REFACTOR OR REMOVE LISTER

type Lister struct {
	events.EmitterBase
	installedRepo packages.InstalledRepo
	remoteRepo    packages.RemoteRepo
}

func New(view ux.View, iRepo packages.InstalledRepo, rRepo packages.RemoteRepo) *Lister {
	l := Lister{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeList,
		},
		installedRepo: iRepo,
		remoteRepo:    rRepo,
	}
	l.SetEmitter(view)

	return &l
}

func (l *Lister) SearchAll(ctx context.Context) error {
	op := l.StartOp(events.Operation{}, "search")
	l.EmitStart(op, "searching for packages")
	defer l.EmitEnd(op)

	pkgs, err := l.remoteRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	ps := make([]packages.Package, 0, len(pkgs))
	for i, p := range pkgs {
		pkgs[i] = p
	}
	l.EmitPackages(op, ps)

	l.Emit(events.Event{
		Op:   op,
		Type: events.EventComplete,
	})

	return nil
}

func (l *Lister) SearchPkg(ctx context.Context, id packages.PkgID) error {
	op := l.StartOp(events.Operation{}, "remote")
	l.EmitStart(op, fmt.Sprintf("searching for package: %s", id))
	defer l.EmitEnd(op)

	pkg, err := l.remoteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	l.EmitPackage(op, pkg)

	return nil
}

func (l *Lister) ListAll(ctx context.Context) error {
	op := l.StartOp(events.Operation{}, "list")
	l.EmitStart(op, "listing pacakages")
	defer l.EmitEnd(op)

	pkgs, err := l.installedRepo.GetAll()
	if err != nil {
		return err
	}

	ps := make([]packages.Package, 0, len(pkgs))
	for i, p := range pkgs {
		pkgs[i] = p
	}
	l.EmitPackages(op, ps)

	return nil
}
