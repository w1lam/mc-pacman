package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "PANIC: %v\n", r)
			debug.PrintStack()
			os.Exit(2)
		}
	}()

	a := app.New(cli.NewDebugView())

	if err := a.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL STARTUP ERROR: %v\n", err)
		os.Exit(1)
	}
	defer a.Exit()

	cli.Run(a)
}
