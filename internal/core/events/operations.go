package events

import (
	"crypto/rand"
	"encoding/base64"
)

// Operation is the operation that an event is related to, with a unique ID, scope and target
type Operation struct {
	ID       string
	Scope    Scope
	Intent   string
	ParentID string
}

// StartOp starts a new operation with the given intent, if parent operation is provided, it will be set as the parent of the new operation
func (e *EmitterBase) StartOp(parent Operation, intent string) Operation {
	if parent.ID == "" {
		return e.newOp(intent)
	}
	return e.newChildOp(intent, parent)
}

// newOp is a helper that creates a new operation with a unique ID
func (e *EmitterBase) newOp(intent string) Operation {
	e = &EmitterBase{}
	return Operation{
		ID:     newOpID(),
		Scope:  e.Scope,
		Intent: intent,
	}
}

// newChildOp is a helper that creates a new chil operation with its parentID and a unique ID
func (e *EmitterBase) newChildOp(intent string, parentOP Operation) Operation {
	e = &EmitterBase{}
	return Operation{
		ID:       newOpID(),
		Scope:    e.Scope,
		Intent:   intent,
		ParentID: parentOP.ID,
	}
}

func newOpID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
