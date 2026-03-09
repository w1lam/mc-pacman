// Package app holds the App struct
package app

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/errors"
	"github.com/w1lam/mc-pacman/internal/infra/installed"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// App is the main struct of application holding all core services and state
type App struct {
	view ux.View

	*useCases

	Paths *paths.Paths

	StateRepo     state.Repo
	RemoteRepo    packages.RemoteRepo
	InstalledRepo packages.InstalledRepo
}

// New creates a new App initializing core of app
func New(view ux.View) (*App, error) {
	p := paths.New(paths.RootDir(), "")

	// EnsureDirectories
	if err := p.Ensure(); err != nil {
		return nil, err
	}

	// Start error handler
	if err := errors.Start(p.LogFile()); err != nil {
		return nil, err
	}
	logger := errors.New()

	// repos
	sRepo := state.NewStateRepo(p.StateFile())
	iRepo := installed.New(view, p.PackagesDir())
	rRepo := remote.New(view)

	// state init
	if err := sRepo.Ensure(); err != nil {
		return nil, err
	}

	// load state and resolve mcDir
	st, err := sRepo.Load()
	if err != nil {
		return nil, err
	}
	mcDir := resolveMincraftDir(st, "")

	// update state
	if err := sRepo.Update(func(s *state.State) error {
		s.McDir = mcDir
		return nil
	}); err != nil {
		return nil, err
	}

	// update path with correct mcDir
	p = paths.New(paths.RootDir(), mcDir)

	// initialize useCases
	uc := initUseCases(p, sRepo, iRepo, rRepo, view, logger)

	return &App{
		view: view,

		useCases: uc,

		Paths: p,

		StateRepo:     sRepo,
		RemoteRepo:    rRepo,
		InstalledRepo: iRepo,
	}, nil
}
