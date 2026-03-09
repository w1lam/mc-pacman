package main

import (
	"fmt"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
)

func main() {
	view := cli.NewDebugView()
	a, err := app.New(view)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
		os.Exit(1)
	}
	cli.Run(a)
}
