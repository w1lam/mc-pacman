package logger

import (
	"fmt"
	"io"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

func (l *logger) formatEvent(w io.Writer, e events.Event) {
	sep := ""
	if e.Op.ParentID != "" {
		sep = " > "
	} else {
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "%s %s[%s:%s] (%s) -> %s",
		e.Timestamp.Format("15:04:05"),
		sep,
		e.Op.Scope,
		e.Op.ID,
		e.Type,
		e.Op.Intent,
	)

	if e.Error != nil {
		fmt.Fprintf(w, "\n error: %v", e.Error)
	}
	if e.Message != "" {
		fmt.Fprintf(w, "\n message: %s", e.Message)
	}
	if e.ErrorLvl != "" {
		fmt.Fprintf(w, "\n level: %s", e.ErrorLvl)
	}

	fmt.Fprintln(w)
}
