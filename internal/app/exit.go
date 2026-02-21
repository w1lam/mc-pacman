package app

import (
	"fmt"
	"os"
	"time"

	ansi "github.com/w1lam/Packages/tui"
)

func Exit() {
	ansi.ClearScreenRaw()
	fmt.Printf("Exiting...")
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
}
