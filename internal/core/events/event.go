// Package events handles app activity for ux
package events

import (
	"time"
)

// EventType is the type of an event
type EventType string

const (
	EventStart    EventType = "start"
	EventUpdate   EventType = "update"
	EventSuccess  EventType = "success"
	EventFailure  EventType = "failure"
	EventComplete EventType = "complete"
	EventEnd      EventType = "end"
	EventError    EventType = "error"
	EventInfo     EventType = "info"
)

// Scope is the scope of an operation
type Scope string

const (
	ScopeDownloader  Scope = "downloader"
	ScopeGetter      Scope = "getter"
	ScopeInstaller   Scope = "installer"
	ScopeUninstaller Scope = "uninstaller"
	ScopeUpdater     Scope = "updater"
	ScopeResolver    Scope = "resolver"
	ScopeVerifier    Scope = "verifier"
	ScopeBackup      Scope = "backup"
	ScopeList        Scope = "list"
)

// Event is an event payload emitted by services for ux to consume
type Event struct {
	Type     EventType
	Op       Operation
	SubScope string

	Message string

	FileName   string
	Percentage float64
	Bytes      int

	PackagePayload  PackagePayload
	PackagePayloads []PackagePayload

	Error error

	Timestamp time.Time
}
