// Package app is da app bruh
package app

import (
	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/paths"
)

type App struct {
	Paths *paths.Paths
	Repo  manifest.Repository

	Services *Services
}

// New creates a new App initializing core of app
func New() (*App, error) {
	// Resolve paths
	p, err := paths.Resolve()
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
	repo := manifest.NewFileRepository(p.ManifestPath)

	// manifest init
	if err := repo.EnsureInitialized(); err != nil {
		return nil, err
	}

	services := InitServices(p, repo)

	return &App{
		Repo:  repo,
		Paths: p,

		Services: services,
	}, nil
}
