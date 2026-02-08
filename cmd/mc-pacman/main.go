package main

import (
	"log"
	"os"
)

// NOTES:
// Add independent mod update checking and updating and only update mods that have new versions
// Add version checking for program updates

// initiation
func init() {}

func main() {
	// packages.CompareLists("C:\\Users\\William\\PROG-GO\\RAWMODINSTALLER\\client_list.txt", "C:\\Users\\William\\PROG-GO\\RAWMODINSTALLER\\server_list.txt")
	// time.Sleep(time.Hour * 1)

	if len(os.Args) > 1 {
		err := app.InitCore()
		if err != nil {
			log.Fatal(err)
		}

		cli.Execute()
		return
	}

	if err := app.InitCore(); err != nil {
		log.Fatal(err)
	}

	app.InitTUI()

	app.Run()
}
