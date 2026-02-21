package cli

import (
	"fmt"
	"strings"

	ansi "github.com/w1lam/Packages/tui"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/events"
)

type View struct {
	ansi bool
}

func New() *View {
	return &View{
		ansi: ansi.SupportsANSI(),
	}
}

func (v *View) Emit(event events.Event) {
	scope := strings.ToUpper(string(event.Op.Scope))
	target := event.Op.Target

	switch v.ansi {

	// ANSI SUPPORT
	case true:
		switch event.Op.Scope {
		// FULL DOWNLOADER
		case events.ScopeDownloader:
			switch event.Type {
			case events.EventStart:
				fmt.Printf("\n[%s] downloading %s...                                \n", scope, target)
			case events.EventUpdate:
				fmt.Printf("\r[%s] %s %.0f%% %s                                    ", scope, target, event.Percentage, event.Message)
			case events.EventSuccess:
				fmt.Printf("\n[%s] %s finished                                     ", scope, target)
			case events.EventFailure:
				fmt.Printf("\n[%s] %s failed                                  ", scope, target)
			case events.EventComplete:
				fmt.Printf("\n\n[%s] %s Download Complete!                     \n", scope, target)
			}

		// PER FILE DOWNLOADER
		case events.ScopeDownloaderPerFile:
			switch event.Type {
			case events.EventStart:
				fmt.Printf("\n downloading %s...                                \n", target)
			case events.EventUpdate:
				fmt.Printf("\r %s %.0f%% %s                                    ", target, event.Percentage, event.Message)
			case events.EventSuccess:
				fmt.Printf("\n %s finished                                     ", target)
			case events.EventFailure:
				fmt.Printf("\n %s failed                                  ", target)
			case events.EventComplete:
				fmt.Printf("\n\n %s Download Complete!                     \n", target)
			}

		// INSTALLER
		case events.ScopeInstaller:
			switch event.Type {
			case events.EventStart:
				fmt.Printf("\n[%s] installing %s...                                \n", scope, target)
			case events.EventUpdate:
				fmt.Printf("\r[%s] %s %.0f%% %s                                    ", scope, target, event.Percentage, event.Message)
			case events.EventSuccess:
				fmt.Printf("\n[%s] %s finished                                     ", scope, target)
			case events.EventFailure:
				fmt.Printf("\n[%s] %s failed                                  ", scope, target)
			case events.EventComplete:
				fmt.Printf("\n\n[%s] %s Installation Complete!                     \n", scope, target)
			}

		// LIST
		case events.ScopeList:
			switch event.Type {
			case events.EventStart:
				fmt.Println(event.Message)
			case events.EventComplete:
				index := event.ExtraData.(packages.RemotePackageIndex)
				for pType, pkgs := range index {
					fmt.Println()
					fmt.Println(strings.ToUpper(packages.PackageTypeIndex[pType].TypeName) + "S:")
					fmt.Println(strings.Repeat("-", len(packages.PackageTypeIndex[pType].TypeName)) + "---")
					for _, pkg := range pkgs {
						fmt.Println(" " + pkg.Name)
						fmt.Println(" * ID:", pkg.ID)
						fmt.Println(" - List Version:", pkg.ListVersion)
						fmt.Println(" - Minecraft Version:", pkg.McVersion)
						if pkg.Description != "" {
							fmt.Println(" - Desctription:", pkg.Description)
						}
						fmt.Println(" - Env:", pkg.Env)
						if pkg.Loader != "" {
							fmt.Println(" - Loader:", pkg.Loader)
						}
						fmt.Println()
					}
				}
			}
		}

	// NO ANSI SUPPORT
	case false:
		switch event.Op.Scope {
		// FULL DOWNLOADER
		case events.ScopeDownloader:
			switch event.Type {
			case events.EventStart:
				fmt.Printf(" downloading %s...\n", target)
			case events.EventSuccess:
				fmt.Printf(" %s downloaded!\n", target)
			case events.EventComplete:
				fmt.Printf("[%s] %s Download Complete!\n", scope, target)
			}

		// PER FILE DOWNLOADER
		case events.ScopeDownloaderPerFile:
			switch event.Type {
			case events.EventStart:
				fmt.Printf(" downloading %s...\n", target)
			case events.EventSuccess:
				fmt.Printf(" %s finished\n", target)
			case events.EventFailure:
				fmt.Printf(" %s failed\n", target)
			case events.EventComplete:
				fmt.Printf(" %s downloaded\n", target)
			}

		// INSTALLER
		case events.ScopeInstaller:
			switch event.Type {
			case events.EventStart:
				fmt.Printf("\n[%s] installing %s...                                \n", scope, target)
			case events.EventComplete:
				fmt.Printf("\n\n[%s] %s Installation Complete!                     \n", scope, target)
			}
		}
	}
}

func CliPackageListRenderer(pkgIndex packages.RemotePackageIndex) {
	for pType, pkgs := range pkgIndex {
		fmt.Println()
		fmt.Println(strings.ToUpper(packages.PackageTypeIndex[pType].TypeName) + "S:")
		fmt.Println(strings.Repeat("-", len(packages.PackageTypeIndex[pType].TypeName)) + "---")
		for _, pkg := range pkgs {
			fmt.Println(" " + pkg.Name)
			fmt.Println(" * ID:", pkg.ID)
			fmt.Println(" - List Version:", pkg.ListVersion)
			fmt.Println(" - Minecraft Version:", pkg.McVersion)
			if pkg.Description != "" {
				fmt.Println(" - Desctription:", pkg.Description)
			}
			fmt.Println(" - Env:", pkg.Env)
			if pkg.Loader != "" {
				fmt.Println(" - Loader:", pkg.Loader)
			}
			fmt.Println()
		}
	}
}
