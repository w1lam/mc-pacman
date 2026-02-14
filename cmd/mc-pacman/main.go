package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/services/remote"
	"github.com/w1lam/mc-pacman/internal/ux/cli"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	p, err := remote.GetAllAvailablePackages()
	if err != nil {
		panic(err)
	}

	fmt.Println(p)

	time.Sleep(time.Hour * 1)

	if len(os.Args) > 1 {
		cli.Execute(app)
		return
	}
}
