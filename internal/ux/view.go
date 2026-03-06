package ux

import "github.com/w1lam/mc-pacman/internal/core/events"

type View interface {
	Emit(events.Event)
}
