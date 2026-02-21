// Package tui holds tui features
package tui

import (
	ansi "github.com/w1lam/Packages/tui"
	"github.com/w1lam/mc-pacman/internal/events"
)

type View struct {
	ansi bool
}

func New() *View {
	return &View{
		ansi: ansi.SupportsANSI(),
	}
}

func (v *View) Emit(event events.Event) {}
