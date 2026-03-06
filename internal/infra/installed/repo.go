// Package installed holds installed repo implementation
package installed

import (
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/ux"
)

type repo struct {
	events.EmitterBase
	path string
}

func New(view ux.View, path string) *repo {
	r := repo{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeInstalledRepo,
		},
		path: path,
	}
	r.SetEmitter(view)
	return &r
}
