package app

import (
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/events"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/installer"
	"github.com/w1lam/mc-pacman/internal/services/lister"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
	"github.com/w1lam/mc-pacman/internal/services/updater"
	"github.com/w1lam/mc-pacman/internal/services/verifier"
)

// Services holds all services strored in app
type Services struct {
	Installer *installer.Installer
	Updater   *updater.Updater
	Verifier  *verifier.Verifier
	Lister    *lister.Lister
}

func initServices(
	p *paths.Paths,
	repo manifest.Repo,
	view events.Emitter,
) *Services {
	l := lister.New(p, repo)
	l.SetEmitter(view)

	d := downloader.New()
	d.SetEmitter(view)

	r := resolver.New()
	r.SetEmitter(view)

	i := installer.New(p, repo, d, r)
	i.SetEmitter(view)

	v := verifier.New(p, repo)
	v.SetEmitter(view)

	return &Services{
		Installer: i,
		Updater:   nil,
		Verifier:  v,
		Lister:    l,
	}
}
