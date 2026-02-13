package app

import (
	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/core/verify"
)

type App struct {
	Paths *paths.Paths
	State *state.State
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
	m, err := repo.Init()
	if err != nil {
		return nil, err
	}

	// Init State
	st, err := state.New(m, repo)
	if err != nil {
		return nil, err
	}

	// Verify packages
	verify.VerifyAndReconcile(p)

	return &App{
		State: st,
	}, nil
}
