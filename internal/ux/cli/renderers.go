package cli

import (
	"fmt"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

func downloaderRenderer(e events.Event, ansi bool) {}
func installerRenderer(e events.Event, ansi bool)  {}
func getterRenderer(e events.Event, ansi bool)     {}
func uninstallRenderer(e events.Event, ansi bool)  {}
func updaterRenderer(e events.Event, ansi bool)    {}
func resolverRenderer(e events.Event, ansi bool)   {}
func verifierRenderer(e events.Event, ansi bool)   {}
func backupRenderer(e events.Event, ansi bool)     {}
func listRenderer(e events.Event, ansi bool) {
	switch e.Type {
	case events.EventStart:
		fmt.Println("listing packages...")
	default:
		if e.PackagePayload.ID != "" {
			for _, p := range e.PackagePayloads {
				fmt.Println(p.Name)
				fmt.Println(" " + p.ID)
				fmt.Println(" " + p.Version)
				fmt.Println(" " + p.McVersion)
				fmt.Println()
			}
		}
	}
}
