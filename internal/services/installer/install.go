// Package installer has the installer
package installer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/events"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
	"github.com/w1lam/mc-pacman/internal/services"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
)

// Service Installer service
type Installer struct {
	services.Base
	paths *paths.Paths

	repo       manifest.Repo
	RemoteRepo packages.RemoteRepo

	downloader *downloader.Downloader
	resolver   *resolver.Resolver
	// modrinth   *modrinth.Client
}

func New(p *paths.Paths, repo manifest.Repo, d *downloader.Downloader, r *resolver.Resolver) *Installer {
	return &Installer{
		Base: services.Base{
			Scope:   events.ScopeInstaller,
			Emitter: nil,
		},
		paths: p,

		repo:       repo,
		RemoteRepo: remote.New(),

		downloader: d,
		resolver:   r,
	}
}

// Install installs a package
func (i *Installer) Install(ctx context.Context, pkg packages.RemotePackage) error {
	if pkg.ID == "" || len(pkg.Entries) == 0 {
		return errors.New("invalid package")
	}

	op := events.NewOperation(i.Scope, string(pkg.ID))

	i.Emit(events.Event{
		Type:    events.EventStart,
		Op:      op,
		Message: fmt.Sprintf("starting installer for: %s", pkg.Name),
	})

	// grab manifest
	m, err := i.repo.Load()
	if err != nil {
		return err
	}

	// check if pkg already installed
	if _, exists := m.InstalledPackages[pkg.Type.PackageType][pkg.ID]; exists {
		return fmt.Errorf("package already installed")
	}

	i.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: fmt.Sprintf("[%s] starting resolver for: %s", strings.ToUpper(string(i.Scope)), pkg.Name),
	})
	// resolve remote package ie fetch correct versions
	resolved, err := i.resolver.Resolve(ctx, pkg)
	if err != nil {
		return err
	}
	i.Emit(events.Event{
		Type:    events.EventSuccess,
		Op:      op,
		Message: fmt.Sprintf("[%s] resolver success! for: %s", strings.ToUpper(string(i.Scope)), pkg.Name),
	})

	// temp dir
	tmp, err := os.MkdirTemp(i.paths.PackagesDir, "download.tmp")
	if err != nil {
		return err
	}

	i.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: fmt.Sprintf("[%s] starting downloader for: %s", strings.ToUpper(string(i.Scope)), pkg.Name),
	})
	// download entries
	results, err := i.downloader.DownloadBatch(ctx, tmp, buildFileRequests(resolved))
	if err != nil {
		return err
	}
	i.Emit(events.Event{
		Type:    events.EventSuccess,
		Op:      op,
		Message: fmt.Sprintf("[%s] downloader success! for: %s", strings.ToUpper(string(i.Scope)), pkg.Name),
	})

	meta := packages.PackageTypeIndex[pkg.Type.PackageType]

	fullActivePath := filepath.Join(i.paths.MinecraftDir, meta.ActivePath)
	fullStoragePath := filepath.Join(i.paths.PackagesDir, meta.StorageDir, string(pkg.ID))

	// install to final dir
	if err := i.installToStorage(tmp, fullStoragePath); err != nil {
		return err
	}

	// compute full hash
	fullHash, err := filesystem.ComputeDirHash(fullStoragePath)
	if err != nil {
		return err
	}

	// build installed package
	installed, err := buildInstalledPackage(pkg, resolved, results, fullStoragePath, fullActivePath, fullHash)
	if err != nil {
		return err
	}

	// set installed package
	m.InstalledPackages[pkg.Type.PackageType][pkg.ID] = installed

	i.Emit(events.Event{
		Type:    events.EventComplete,
		Op:      op,
		Message: fmt.Sprintf("%s install complete!", pkg.Name),
	})

	return i.repo.Save(m)
}
