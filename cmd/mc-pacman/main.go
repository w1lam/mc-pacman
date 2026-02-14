package main

import (
	"log"
	"os"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
)

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
