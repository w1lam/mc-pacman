// Package app holds the App struct
package app

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// App is the main struct of application holding all core services and state
type App struct {
	view   ux.View
	logger events.Logger

	*useCases

	Paths         *paths.Paths
	StateRepo     state.Repo
	RemoteRepo    packages.RemoteRepo
	InstalledRepo packages.InstalledRepo
}

// New creates a new App
func New(view ux.View) *App {
	return &App{
		view: view,
	}
}

const (
	name    = "mc-pacman"
	version = "0.1"
)

// UserAgent returns the user agent string for http requests
func UserAgent() string {
	return name + "-v" + version
}
