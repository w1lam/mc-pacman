// Package ux holds user facing/visual part of app
package ux

import "github.com/w1lam/mc-pacman/internal/events"

type View interface {
	Emit(events.Event)
}
