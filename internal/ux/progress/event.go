// Package progress handles progress emits
package progress

import "time"

type ProgressEventType string

const (
	ProgressStart    ProgressEventType = "start"
	ProgressUpdate   ProgressEventType = "update"
	ProgressSuccess  ProgressEventType = "success"
	ProgressFailure  ProgressEventType = "failure"
	ProgressComplete ProgressEventType = "complete"
)

type ProgressEvent struct {
	Type       ProgressEventType
	Context    string
	PackageID  string
	FileName   string
	Percentage float64
	Bytes      int
	Message    string
	Error      error
	ExtraData  any
	Timestamp  time.Time
}

type ProgressEmitter interface {
	Emit(event ProgressEvent)
}
