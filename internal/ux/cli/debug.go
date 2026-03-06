package cli

import (
	"fmt"
	"strings"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

type DebugView struct{}

func NewDebugView() *DebugView {
	return &DebugView{}
}

func (v *DebugView) Emit(e events.Event) {
	sep := ""
	if e.Op.ParentID != "" {
		sep = " > "
	} else {
		fmt.Println()
	}
	fmt.Printf("%s %s[%s:%s] (%s) -> %s %s : %s %s",
		e.Timestamp.Format("15:04:05"),
		sep,
		strings.ToUpper(string(e.Op.Scope)),
		e.Op.ID,
		strings.ToUpper(string(e.Type)),
		e.Op.Intent,
		e.Error,
		e.Message,
		e.FileName,
	)
	if e.Payload != nil {
		fmt.Printf("[PAYLOAD]: %+v", e.Payload)
	}

	fmt.Println()
}
