// Package lister holds list service
package lister

import (
	"context"
	"fmt"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/usecases"
)

type Lister struct {
	usecases.Base

	installedRepo packages.InstalledRepo
	remoteRepo    packages.RemoteRepo
}

func New(base usecases.Base, iRepo packages.InstalledRepo, rRepo packages.RemoteRepo) *Lister {
	l := Lister{
		Base:          base,
		installedRepo: iRepo,
		remoteRepo:    rRepo,
	}

	return &l
}

// SearchAll searches all packages from github repo
func (l *Lister) SearchAll(ctx context.Context) error {
	op := l.StartOp(events.Operation{}, "search_remote_all")
	l.EmitStart(op, "searching for packages")
	defer l.EmitEnd(op)

	pkgs, err := l.remoteRepo.GetAll(events.WithOp(ctx, op))
	if err != nil {
		return err
	}

	ps := make([]packages.Package, 0, len(pkgs))
	for _, p := range pkgs {
		ps = append(ps, p)
	}
	l.EmitPackages(op, ps)

	return nil
}

// SearchPkg searches for a package from github repo
func (l *Lister) SearchPkg(ctx context.Context, pkgID packages.PkgID) error {
	op := l.StartOp(events.Operation{}, fmt.Sprintf("search_remote_%s", pkgID))
	l.EmitStart(op, fmt.Sprintf("searching for package: %s", pkgID))
	defer l.EmitEnd(op)

	pkg, err := l.remoteRepo.GetByID(events.WithOp(ctx, op), pkgID)
	if err != nil {
		return err
	}

	l.EmitPackage(op, pkg)

	return nil
}

// ListAll lists all local/installed packages
func (l *Lister) ListAll(ctx context.Context) error {
	op := l.StartOp(events.Operation{}, "list_installed_all")
	l.EmitStart(op, "listing pacakages")
	defer l.EmitEnd(op)

	pkgs, err := l.installedRepo.GetAll(events.WithOp(ctx, op))
	if err != nil {
		return err
	}

	ps := make([]packages.Package, 0, len(pkgs))
	for _, p := range pkgs {
		ps = append(ps, p)
	}
	l.EmitPackages(op, ps)

	return nil
}

// NOTE: maybe make ListPkg actually list entries, or atleast add the feature?

// ListPkg lits a local/installed package
func (l *Lister) ListPkg(ctx context.Context, pkgID packages.PkgID) error {
	op := l.StartOp(events.Operation{}, fmt.Sprintf("list_installed_%s", pkgID))
	l.EmitStart(op, fmt.Sprintf("listing pacakage %s", pkgID))
	defer l.EmitEnd(op)

	pkg, err := l.installedRepo.GetByID(events.WithOp(ctx, op), pkgID)
	if err != nil {
		return err
	}

	l.EmitPackage(op, pkg)

	return nil
}
