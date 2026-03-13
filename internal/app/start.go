package app

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/installed"
	"github.com/w1lam/mc-pacman/internal/infra/logger"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
)

// Start starts the app
func (a *App) Start() error {
	p := paths.New(paths.RootDir(), "")

	// EnsureDirectories
	if err := p.Ensure(); err != nil {
		return err
	}

	// Create logger
	logger := logger.New(p.LogFile())

	// repos
	sRepo := state.NewStateRepo(p.StateFile())
	iRepo := installed.New(a.view, p.PackagesDir())
	rRepo := remote.New(a.view)

	// state init
	if err := sRepo.Ensure(); err != nil {
		return err
	}

	// load state and resolve mcDir
	st, err := sRepo.Load()
	if err != nil {
		return err
	}
	mcDir := resolveMincraftDir(st, "")

	// update state
	if err := sRepo.Update(func(s *state.State) error {
		s.McDir = mcDir
		return nil
	}); err != nil {
		return err
	}

	// update path with correct mcDir
	p = paths.New(paths.RootDir(), mcDir)

	// initialize useCases
	uc := initUseCases(p, sRepo, iRepo, rRepo, a.view, logger)

	a.logger = logger
	a.useCases = uc
	a.Paths = p
	a.StateRepo = sRepo
	a.RemoteRepo = rRepo
	a.InstalledRepo = iRepo

	return nil
}
