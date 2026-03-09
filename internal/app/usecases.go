package app

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
	"github.com/w1lam/mc-pacman/internal/usecases"
	"github.com/w1lam/mc-pacman/internal/usecases/activation"
	"github.com/w1lam/mc-pacman/internal/usecases/getter"
	"github.com/w1lam/mc-pacman/internal/usecases/installer"
	"github.com/w1lam/mc-pacman/internal/usecases/lister"
	"github.com/w1lam/mc-pacman/internal/usecases/updater"
	"github.com/w1lam/mc-pacman/internal/usecases/verifier"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// useCases holds all UseCases strored in app
type useCases struct {
	Installer *installer.Installer
	Getter    *getter.Getter
	Updater   *updater.Updater
	Verifier  *verifier.Verifier
	Lister    *lister.Lister
	Activator *activation.Activator
}

type factory struct {
	view   ux.View
	logger events.Logger
}

func (f *factory) Base(scope events.Scope) usecases.Base {
	b := usecases.Base{
		EmitterBase: events.EmitterBase{
			Scope: scope,
		},
		Logger: f.logger,
	}
	b.SetEmitter(f.view)

	return b
}

// initUseCases initializes all use cases with their dependencies and returns a useCases struct
func initUseCases(
	p *paths.Paths,
	sRepo state.Repo,
	iRepo packages.InstalledRepo,
	rRepo packages.RemoteRepo,
	view ux.View,
	logger events.Logger,
) *useCases {
	fac := factory{
		view:   view,
		logger: logger,
	}

	// SERVICES
	r := resolver.New(fac.Base(events.ScopeResolver), UserAgent())
	d := downloader.New(fac.Base(events.ScopeDownloader))

	// USECASES
	a := activation.New(fac.Base(events.ScopeActivator), iRepo, sRepo, p)
	l := lister.New(fac.Base(events.ScopeList), iRepo, rRepo)
	g := getter.New(fac.Base(events.ScopeGetter), p, iRepo, rRepo, d, r)
	v := verifier.New(fac.Base(events.ScopeVerifier), p, iRepo)
	i := installer.New(fac.Base(events.ScopeInstaller), p, sRepo, g)
	u := updater.New(fac.Base(events.ScopeUpdater), iRepo, rRepo)

	return &useCases{
		Installer: i,
		Getter:    g,
		Updater:   u,
		Verifier:  v,
		Lister:    l,
		Activator: a,
	}
}
