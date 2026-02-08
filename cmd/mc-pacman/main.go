package main

import (
	"log"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/app/cli"
	"github.com/w1lam/mc-pacman/internal/app/init"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates

func main() {
	if len(os.Args) > 1 {
		err := init.InitCore()
		if err != nil {
			log.Fatal(err)
		}

		cli.Execute()
		return
	}

	if err := init.InitCore(); err != nil {
		log.Fatal(err)
	}

	app.InitTUI()

	app.Run()
}
