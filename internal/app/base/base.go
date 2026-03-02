// Package base provides base struct for usecases directly accessed by user
package base

import (
	"time"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

type UseCase struct {
	Scope   events.Scope
	Emitter events.Emitter
}

func (b *UseCase) SetEmitter(e events.Emitter) {
	b.Emitter = e
}

func (b *UseCase) Emit(e events.Event) {
	if b.Emitter == nil {
		return
	}

	e.Timestamp = time.Now()
	b.Emitter.Emit(e)
}
