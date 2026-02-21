package main

import (
	"log"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/ux"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
	"github.com/w1lam/mc-pacman/internal/ux/tui"
)

// TODO: PRIO 1: finsih remaining services. enable/disable, uninstall, updater, meta, loaders
//       PRIO 2: add cli commands and flags to all services
//       PRIO 3: add basic commands and flags like help/version etc
//       PRIO 4: reintroduce TUI frontend

func main() {
	useTUI := len(os.Args) == 0
	var view ux.View

	if useTUI {
		view = tui.New()
	} else {
		view = cli.New()
	}

	app, err := app.New(view)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		cli.Run(app)
		return
	}
}
