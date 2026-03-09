package usecases

import (
	"github.com/w1lam/mc-pacman/internal/core/events"
)

type Base struct {
	events.EmitterBase
	events.Logger
}
