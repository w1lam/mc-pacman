// Package installer has the installer
package installer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
	"github.com/w1lam/mc-pacman/internal/ux/progress"
)

// Service Installer service
type Service struct {
	paths      *paths.Paths
	repo       manifest.Repository
	downloader *downloader.Service
	emitter    progress.ProgressEmitter
	// modrinth   *modrinth.Client
}

func New(p *paths.Paths, repo manifest.Repository, d *downloader.Service) *Service {
	return &Service{
		paths:      p,
		repo:       repo,
		downloader: d,
	}
}

// Install installs a package
func (s *Service) Install(ctx context.Context, pkg packages.RemotePackage) error {
	if pkg.ID == "" || len(pkg.Entries) == 0 {
		return errors.New("invalid package")
	}

	s.emitter.Emit(progress.ProgressEvent{
		Type:      progress.ProgressStart,
		Context:   "installer",
		PackageID: string(pkg.ID),
		Message:   fmt.Sprintf("installing %s", pkg.Name),
		Timestamp: time.Now(),
	})

	// grab manifest
	m, err := s.repo.Load()
	if err != nil {
		return err
	}

	// check if pkg already installed
	if _, exists := m.InstalledPackages[pkg.Type.PackageType][pkg.ID]; exists {
		return fmt.Errorf("package already installed")
	}

	// resolve remote package ie fetch correct versions
	resolved, err := resolver.New().Resolve(ctx, pkg)
	if err != nil {
		return err
	}

	// temp dir
	tmp, err := os.MkdirTemp(s.paths.PackagesDir, "download.tmp")
	if err != nil {
		return err
	}

	// download entries
	results, err := s.downloader.DownloadBatch(ctx, tmp, buildFileRequests(resolved))
	if err != nil {
		return err
	}

	meta := packages.PackageTypeIndex[pkg.Type.PackageType]

	fullActivePath := filepath.Join(s.paths.MinecraftDir, meta.ActivePath)
	fullStoragePath := filepath.Join(s.paths.PackagesDir, meta.StorageDir, string(pkg.ID))

	// install to final dir
	if err := s.installToStorage(tmp, fullStoragePath); err != nil {
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

	s.emitter.Emit(progress.ProgressEvent{
		Type:      progress.ProgressSuccess,
		Context:   "installer",
		PackageID: string(pkg.ID),
		Message:   fmt.Sprintf("%s installed!", pkg.Name),
		Timestamp: time.Now(),
	})

	return s.repo.Save(m)
}

func (s *Service) SetEmitter(e progress.ProgressEmitter) {
	s.emitter = e
}
