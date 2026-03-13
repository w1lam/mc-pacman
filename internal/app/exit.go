package app

import (
	"fmt"
	"os"
)

func (a *App) Exit() {
	if a.logger != nil {
		if err := a.logger.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		}
	}
}
