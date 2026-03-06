// Package getter epic funny hahah
package getter

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/infra/downloader"
	"github.com/w1lam/mc-pacman/internal/infra/filesystem"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
	"github.com/w1lam/mc-pacman/internal/infra/resolver"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// Getter gets a package and installs it to the correct location
type Getter struct {
	events.EmitterBase

	paths *paths.Paths
	iRepo packages.InstalledRepo
	rRepo packages.RemoteRepo

	downloader *downloader.Downloader
	resolver   *resolver.Resolver
}

// New creates a new getter useCase
func New(view ux.View, p *paths.Paths, iRepo packages.InstalledRepo, rRepo packages.RemoteRepo, d *downloader.Downloader, r *resolver.Resolver) *Getter {
	g := Getter{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeGetter,
		},
		paths: p,

		iRepo: iRepo,
		rRepo: rRepo,

		downloader: d,
		resolver:   r,
	}
	g.SetEmitter(view)
	return &g
}

// Get downloads a package and stores it in PkgID/ dir
func (g *Getter) Get(ctx context.Context, pkgID string) error {
	pOp, _ := events.OpFromCtx(ctx)
	op := g.StartOp(pOp, fmt.Sprintf("get %s", pkgID))
	g.EmitStart(op, fmt.Sprintf("starting installation of: %s", pkgID))
	defer g.EmitEnd(op)

	exists, err := g.iRepo.Exists(packages.PkgID(pkgID))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("package already installed")
	}

	// get remote package
	pkg, err := g.rRepo.GetByID(events.WithOp(ctx, op), packages.PkgID(pkgID))
	if err != nil {
		return err
	}
	if pkg.ID == "" || len(pkg.Entries) == 0 {
		return errors.New("invalid package")
	}

	// resolve remote package
	resolved, err := g.resolver.Resolve(events.WithOp(ctx, op), pkg)
	if err != nil {
		return err
	}

	// temp dir
	tmp, err := os.MkdirTemp(g.paths.PackagesDir(), "download.tmp")
	if err != nil {
		return err
	}

	// download entries
	resultFiles, err := g.downloader.Download(events.WithOp(ctx, op), tmp, buildFileRequests(resolved))
	if err != nil {
		return err
	}

	downloadedPackage := buildDownloadedPackage(resolved, resultFiles)

	// compute full hash
	hash, err := filesystem.ComputeDirHash(tmp)
	if err != nil {
		return err
	}

	// build installed package
	installed, err := buildInstalledPackage(downloadedPackage, hash)
	if err != nil {
		return err
	}

	if err := g.iRepo.Add(installed, tmp); err != nil {
		return err
	}

	g.EmitComplete(op, fmt.Sprintf("sucessfully installed %s", pkg.ID))

	return nil
}
