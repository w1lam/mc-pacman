package events

import (
	"time"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

// Emitter is the interface for emitting events, implemented by views and services
type Emitter interface {
	Emit(Event)
}

// EmitterBase provides base struct for emitters, with auto timestamp and emitter setter
type EmitterBase struct {
	Scope   Scope
	emitter Emitter
}

// SetEmitter sets emitter for EmitterBase
func (b *EmitterBase) SetEmitter(e Emitter) {
	b.emitter = e
}

// Emit base emit helper, auto set timestamp
func (b *EmitterBase) Emit(e Event) {
	if b.emitter == nil {
		return
	}

	e.Timestamp = time.Now()
	b.emitter.Emit(e)
}

// EmitStart operation start helper
func (b *EmitterBase) EmitStart(op Operation, msg string) {
	if b.emitter == nil {
		return
	}

	e := Event{
		Type:      EventStart,
		Op:        op,
		Message:   msg,
		Timestamp: time.Now(),
	}

	b.emitter.Emit(e)
}

// EmitComplete operation end helper
func (b *EmitterBase) EmitComplete(op Operation, msg string) {
	if b.emitter == nil {
		return
	}

	e := Event{
		Type:      EventEnd,
		Op:        op,
		Message:   msg,
		Timestamp: time.Now(),
	}

	b.emitter.Emit(e)
}

// EmitEnd operation end helper
func (b *EmitterBase) EmitEnd(op Operation) {
	if b.emitter == nil {
		return
	}

	e := Event{
		Type:      EventEnd,
		Op:        op,
		Timestamp: time.Now(),
	}

	b.emitter.Emit(e)
}

// EmitError error emitter helper
func (b *EmitterBase) EmitError(op Operation, err error) {
	if b.emitter == nil {
		return
	}

	e := Event{
		Type:      EventError,
		Op:        op,
		Error:     err,
		Timestamp: time.Now(),
	}

	b.emitter.Emit(e)
}

// EmitPackage emit a payload package
func (b *EmitterBase) EmitPackage(op Operation, p packages.Package) {
	if b.emitter == nil {
		return
	}

	e := Event{
		Type:           EventInfo,
		Op:             op,
		PackagePayload: NewPackagePayload(p.GetBase(), p.IsInstalled()),
		Timestamp:      time.Now(),
	}

	b.emitter.Emit(e)
}

// EmitPackages emit a batch of payload packages
func (b *EmitterBase) EmitPackages(op Operation, ps []packages.Package) {
	if b.emitter == nil {
		return
	}

	pp := make([]PackagePayload, 0, len(ps))
	for _, p := range ps {
		pp = append(pp, NewPackagePayload(p.GetBase(), p.IsInstalled()))
	}

	e := Event{
		Type:            EventInfo,
		Op:              op,
		PackagePayloads: pp,
		Timestamp:       time.Now(),
	}

	b.emitter.Emit(e)
}
