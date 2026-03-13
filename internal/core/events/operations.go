package events

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

// Operation is the operation that an event is related to, with a unique ID, scope and target
type Operation struct {
	ID       string
	Scope    Scope
	Intent   string
	ParentID string
	Depth    int
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
	return Operation{
		ID:     newOpID(),
		Scope:  e.Scope,
		Intent: intent,
		Depth:  0,
	}
}

// newChildOp is a helper that creates a new chil operation with its parentID and a unique ID
func (e *EmitterBase) newChildOp(intent string, parentOp Operation) Operation {
	return Operation{
		ID:       newOpID(),
		Scope:    e.Scope,
		Intent:   intent,
		ParentID: parentOp.ID,
		Depth:    parentOp.Depth + 1,
	}
}

func newOpID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

type contextKey string

const opKey contextKey = "op"

// WithOp puts Operation in ctx
func WithOp(ctx context.Context, op Operation) context.Context {
	return context.WithValue(ctx, opKey, op)
}

// OpFromCtx gets Operation from ctx
func OpFromCtx(ctx context.Context) (Operation, bool) {
	op, ok := ctx.Value(opKey).(Operation)
	return op, ok
}
