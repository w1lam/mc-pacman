// Package events handles app activity for ux
package events

import (
	"time"
)

// EventType is the type of an event
type EventType string

const (
	EventStart    EventType = "START"
	EventComplete EventType = "COMPLETE"
	EventError    EventType = "ERROR"
	EventInfo     EventType = "INFO"
	EventPayload  EventType = "PAYLOAD"
)

// Scope is the scope of an operation
type Scope string

const (
	ScopeDownloader    Scope = "DOWNLOADER"
	ScopeGetter        Scope = "GETTER"
	ScopeInstaller     Scope = "INSTALELR"
	ScopeUninstaller   Scope = "UNINSTALLER"
	ScopeUpdater       Scope = "UPDATER"
	ScopeResolver      Scope = "RESOLVER"
	ScopeVerifier      Scope = "VERIFIER"
	ScopeActivator     Scope = "ACTIVATOR"
	ScopeBackup        Scope = "BACKUP"
	ScopeList          Scope = "LIST"
	ScopeRemoteRepo    Scope = "REMOTEREPO"
	ScopeInstalledRepo Scope = "INSTALLEDREPO"
)

// ErrorLvl is the level of an error
type ErrorLvl string

const (
	// WARN a warning error, function can continue
	WARN ErrorLvl = "WARN"
	// ERROR function error, fcuntion ends
	ERROR ErrorLvl = "ERROR"
	// FATAL fatal error, program ends
	FATAL ErrorLvl = "FATAL"
)

// Event is an event payload emitted by services for ux to consume
type Event struct {
	Type    EventType
	Op      Operation
	Message string

	Payload *Payload

	Error    error
	ErrorLvl ErrorLvl

	Timestamp time.Time
}
