// Package services holds all services
package services

import (
	"time"

	"github.com/w1lam/mc-pacman/internal/events"
)

type Base struct {
	Scope   events.Scope
	Emitter events.Emitter
}

func (b *Base) SetEmitter(e events.Emitter) {
	b.Emitter = e
}

func (b *Base) Emit(e events.Event) {
	if b.Emitter == nil {
		return
	}

	e.Timestamp = time.Now()
	b.Emitter.Emit(e)
}
