package app

import (
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/installer"
	"github.com/w1lam/mc-pacman/internal/services/updater"
	"github.com/w1lam/mc-pacman/internal/services/verifier"
	"github.com/w1lam/mc-pacman/internal/ux/progress"
)

// Services holds all services strored in app
type Services struct {
	Installer *installer.Service
	Updater   *updater.Service
	Verifier  *verifier.Service
}

func InitServices(p *paths.Paths, repo manifest.Repository) *Services {
	d := downloader.New()
	d.SetEmitter(&progress.CLIProgress{})

	i := installer.New(p, repo, d)
	i.SetEmitter(&progress.CLIProgress{})

	return &Services{
		Installer: i,
		Updater:   nil,
		Verifier:  nil,
	}
}
