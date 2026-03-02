package app

import (
	"github.com/w1lam/mc-pacman/internal/app/getter"
	"github.com/w1lam/mc-pacman/internal/app/installer"
	"github.com/w1lam/mc-pacman/internal/app/lister"
	"github.com/w1lam/mc-pacman/internal/app/updater"
	"github.com/w1lam/mc-pacman/internal/app/verifier"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/infra/downloader"
	"github.com/w1lam/mc-pacman/internal/infra/errors"
	"github.com/w1lam/mc-pacman/internal/infra/installed"
	"github.com/w1lam/mc-pacman/internal/infra/paths"
	"github.com/w1lam/mc-pacman/internal/infra/remote"
	"github.com/w1lam/mc-pacman/internal/infra/resolver"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// TODO: CLEAN UPP APP NEW, MAKE SURE EVERYTHING FLOWS DOWNWARDS AND NO DUPLICATES

// App is the main struct of application holding all core services and state
type App struct {
	view ux.View

	UseCases *useCases

	paths *paths.Paths

	stateRepo     state.Repo
	remoteRepo    packages.RemoteRepo
	installedRepo packages.InstalledRepo
}

// New creates a new App initializing core of app
func New(view ux.View, cfg Config) (*App, error) {
	p := paths.New(paths.RootDir(), "")

	// EnsureDirectories
	if err := p.Ensure(); err != nil {
		return nil, err
	}

	// Start error handler
	if err := errors.Start(p.LogFile()); err != nil {
		return nil, err
	}

	// repos
	sRepo := state.NewStateRepo(p.StateFile())
	iRepo := installed.New(p.PackagesDir())
	rRepo := remote.New()

	// state init
	if err := sRepo.Ensure(); err != nil {
		return nil, err
	}

	// load state and resolve mcDir
	st, err := sRepo.Load()
	if err != nil {
		return nil, err
	}
	mcDir := resolveMincraftDir(st, cfg.McDir)

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
	uc := initUseCases(p, sRepo, iRepo, rRepo, view)

	return &App{
		view: view,

		UseCases: uc,

		paths:     p,
		stateRepo: sRepo,

		remoteRepo:    rRepo,
		installedRepo: iRepo,
	}, nil
}

// useCases holds all UseCases strored in app
type useCases struct {
	Installer *installer.Installer
	Getter    *getter.Getter
	Updater   *updater.Updater
	Verifier  *verifier.Verifier
	Lister    *lister.Lister
}

// initUseCases initializes all use cases with their dependencies and returns a useCases struct
func initUseCases(
	p *paths.Paths,
	sRepo state.Repo,
	iRepo packages.InstalledRepo,
	rRepo packages.RemoteRepo,
	view events.Emitter,
) *useCases {
	l := lister.New(p, sRepo)
	l.SetEmitter(view)

	d := downloader.New()
	d.SetEmitter(view)

	r := resolver.New(UserAgent())
	r.SetEmitter(view)

	g := getter.New(p, iRepo, d, r)
	g.SetEmitter(view)

	v := verifier.New(p, sRepo)
	v.SetEmitter(view)

	i := installer.New(p, sRepo, g)
	i.SetEmitter(view)

	return &useCases{
		Installer: i,
		Getter:    g,
		Updater:   nil,
		Verifier:  v,
		Lister:    l,
	}
}
