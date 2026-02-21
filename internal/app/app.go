// Package app is da app bruh
package app

import (
	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/ux"
)

type App struct {
	View ux.View

	Services *Services

	Paths        *paths.Paths
	PathsRepo    paths.PathsRepository
	ManifestRepo manifest.Repo
}

// New creates a new App initializing core of app
func New(view ux.View) (*App, error) {
	pRepo := paths.NewPathRepo()

	p, err := pRepo.Init()
	if err != nil {
		return nil, err
	}

	// EnsureDirectories
	if err := filesystem.EnsureDirectories(p); err != nil {
		return nil, err
	}

	// Start error handler
	if err := errors.Start(p.LogPath); err != nil {
		return nil, err
	}

	// manifest repo
	mRepo := manifest.NewFileRepository(p.ManifestPath)

	// manifest init
	if err := mRepo.EnsureInitialized(); err != nil {
		return nil, err
	}

	s := initServices(p, mRepo, view)

	return &App{
		View: view,

		Services: s,

		Paths:        p,
		PathsRepo:    pRepo,
		ManifestRepo: mRepo,
	}, nil
}
