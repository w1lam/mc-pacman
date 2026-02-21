// Package events handles app activity for ux
package events

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type EventType string

const (
	EventStart    EventType = "start"
	EventUpdate   EventType = "update"
	EventSuccess  EventType = "success"
	EventFailure  EventType = "failure"
	EventComplete EventType = "complete"
)

type Scope string

const (
	ScopeDownloader        Scope = "downloader"
	ScopeDownloaderPerFile Scope = "perFile"
	ScopeInstaller         Scope = "installer"
	ScopeResolver          Scope = "resolver"
	ScopeVerifier          Scope = "verifier"
	ScopeBackup            Scope = "backup"
	ScopeList              Scope = "list"
)

type Event struct {
	Type     EventType
	Op       Operation
	SubScope string

	Message  string
	FileName string

	Percentage float64
	Bytes      int

	ExtraData any

	Error     error
	Timestamp time.Time
}

type Operation struct {
	ID     string
	Scope  Scope
	Target string
}

func NewOperation(scope Scope, target string) Operation {
	return Operation{
		ID:     newOpID(),
		Scope:  scope,
		Target: target,
	}
}

func newOpID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

type Emitter interface {
	Emit(Event)
}
