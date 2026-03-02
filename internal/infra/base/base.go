// Package base has base service for "sub services" services not directly accessed by user
package base

import (
	"time"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

type Service struct {
	Scope   events.Scope
	Emitter events.Emitter
}

func (b *Service) SetEmitter(e events.Emitter) {
	b.Emitter = e
}

func (b *Service) Emit(e events.Event) {
	if b.Emitter == nil {
		return
	}

	e.Timestamp = time.Now()
	b.Emitter.Emit(e)
}
