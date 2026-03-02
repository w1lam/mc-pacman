// Package getter epic funny hahah
package getter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/w1lam/mc-pacman/internal/app/base"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/infra/downloader"
	"github.com/w1lam/mc-pacman/internal/infra/filesystem"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
	"github.com/w1lam/mc-pacman/internal/infra/resolver"
)

// Getter gets a package and installs it to the correct location
type Getter struct {
	base.UseCase
	paths *paths.Paths

	iRepo packages.InstalledRepo
	rRepo packages.RemoteRepo

	downloader *downloader.Downloader
	resolver   *resolver.Resolver
}

// New creates a new getter useCase
func New(p *paths.Paths, iRepo packages.InstalledRepo, d *downloader.Downloader, r *resolver.Resolver) *Getter {
	return &Getter{
		UseCase: base.UseCase{
			Scope:   events.ScopeGetter,
			Emitter: nil,
		},
		paths: p,

		iRepo: iRepo,
		rRepo: remote.New(),

		downloader: d,
		resolver:   r,
	}
}

// Get downloads a package and stores it in PkgID/ dir
func (g *Getter) Get(ctx context.Context, pkgID string) error {
	exists, err := g.iRepo.Exists(packages.PkgID(pkgID))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("package already installed")
	}

	op := events.NewOperation(g.Scope, pkgID)

	g.Emit(events.Event{
		Type:    events.EventStart,
		Op:      op,
		Message: fmt.Sprintf("starting installer for: %s", pkgID),
	})

	// get remote package
	pkg, err := g.rRepo.GetByID(ctx, packages.PkgID(pkgID))
	if err != nil {
		return err
	}
	if pkg.ID == "" || len(pkg.Entries) == 0 {
		return errors.New("invalid package")
	}

	g.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: fmt.Sprintf("[%s] starting resolver for: %s", strings.ToUpper(string(g.Scope)), pkg.Name),
	})

	// resolve remote package
	resolved, err := g.resolver.Resolve(ctx, pkg)
	if err != nil {
		return err
	}

	g.Emit(events.Event{
		Type:    events.EventSuccess,
		Op:      op,
		Message: fmt.Sprintf("[%s] resolver success! for: %s", strings.ToUpper(string(g.Scope)), pkg.Name),
	})

	// temp dir
	tmp, err := os.MkdirTemp(g.paths.PackagesDir(), "download.tmp")
	if err != nil {
		return err
	}

	g.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: fmt.Sprintf("[%s] starting downloader for: %s", strings.ToUpper(string(g.Scope)), pkg.Name),
	})

	// download entries
	resultFiles, err := g.downloader.Download(ctx, tmp, buildFileRequests(resolved))
	if err != nil {
		return err
	}

	downloadedPackage := buildDownloadedPackage(resolved, resultFiles)

	g.Emit(events.Event{
		Type:    events.EventSuccess,
		Op:      op,
		Message: fmt.Sprintf("[%s] downloader success! for: %s", strings.ToUpper(string(g.Scope)), pkg.Name),
	})

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

	g.Emit(events.Event{
		Type:    events.EventComplete,
		Op:      op,
		Message: fmt.Sprintf("%s install complete!", pkg.Name),
	})

	return nil
}
