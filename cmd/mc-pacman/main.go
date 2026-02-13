package main

import (
	"log"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/app/cli"
	"github.com/w1lam/mc-pacman/internal/app/tui"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		cli.Execute(app)
		return
	}

	tui.InitTUI()
}
