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
	logger  Logger
}

// SetEmitter sets the emitter for the EmitterBase
func (b *EmitterBase) SetEmitter(e Emitter) {
	b.emitter = e
}

// SetLogger sets the logger for the EmitterBase
func (b *EmitterBase) SetLogger(l Logger) {
	b.logger = l
}

// emit is the base emit helper that automatically adds timestamp and checks if emitter is nil
func (b *EmitterBase) emit(e Event) {
	if b.emitter == nil {
		return
	}
	e.Timestamp = time.Now()
	b.emitter.Emit(e)
}

// EmitStart operation start helper
func (b *EmitterBase) EmitStart(op Operation, msg string) {
	b.emit(Event{
		Type:    EventStart,
		Op:      op,
		Message: msg,
	})
}

// EmitComplete operation end helper
func (b *EmitterBase) EmitComplete(op Operation, msg string) {
	b.emit(Event{
		Type:    EventComplete,
		Op:      op,
		Message: msg,
	})
}

// EmitInfo emits an info event
func (b *EmitterBase) EmitInfo(op Operation, msg string) {
	b.emit(Event{
		Type:    EventInfo,
		Op:      op,
		Message: msg,
	})
}

// ERROR EMITTERs

// EmitWarn emits a warning, operation continues
func (b *EmitterBase) EmitWarn(op Operation, err error, msg string) {
	b.emitError(op, err, msg, WARN)
}

// EmitError emits an error, operation failed
func (b *EmitterBase) EmitError(op Operation, err error, msg string) {
	b.emitError(op, err, msg, WARN)
}

// EmitFatal emits a fatal error, program should exit
func (b *EmitterBase) EmitFatal(op Operation, err error, msg string) {
	b.emitError(op, err, msg, WARN)
}

// emitError emit an error
func (b *EmitterBase) emitError(op Operation, err error, msg string, errLvl ErrorLvl) {
	e := Event{
		Type:     EventError,
		Op:       op,
		Message:  msg,
		Error:    err,
		ErrorLvl: errLvl,
	}

	if b.logger != nil {
		b.logger.Log(e)
	}

	b.emit(e)
}

// PAYLOADS

// emitPayload base payload emitter emits a payload event
func (b *EmitterBase) emitPayload(op Operation, payload *Payload) {
	b.emit(Event{
		Type:    EventPayload,
		Op:      op,
		Payload: payload,
	})
}

// EmitProgress emits an progress payload event
func (b *EmitterBase) EmitProgress(op Operation, progress Progress) {
	b.emitPayload(
		op,
		&Payload{
			Progress: &progress,
		},
	)
}

// EmitPackage emits a package
func (b *EmitterBase) EmitPackage(op Operation, p packages.Package) {
	b.emitPayload(
		op,
		&Payload{
			Package: &PackageItem{
				PackageBase: p.GetBase(),
				Installed:   p.IsInstalled(),
			},
		})
}

// EmitPackages emit a batch of packages
func (b *EmitterBase) EmitPackages(op Operation, ps []packages.Package) {
	pi := make([]PackageItem, 0, len(ps))
	for _, p := range ps {
		pi = append(pi, PackageItem{
			PackageBase: p.GetBase(),
			Installed:   p.IsInstalled(),
		})
	}

	b.emitPayload(
		op,
		&Payload{
			Packages: pi,
		})
}
