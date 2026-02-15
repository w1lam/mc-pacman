package main

import (
	"log"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
)

// TODO: SERVICES: enable/disable, uninstall, updater, meta, loaders, verifier
//       UX: progress, renderers
//       CLI: commands for enable/disable, uninstall, update, loader, meta, verify

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		cli.Execute(app)
		return
	}
}
