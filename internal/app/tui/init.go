package tui

import (
	"fmt"
	"log"

	"github.com/w1lam/Packages/tui"
	initial "github.com/w1lam/mc-pacman/internal/app/init"
)

// InitTUI initializes the tui
func InitTUI() {
	tui.EnableANSI()
	tui.HideCursor()
	tui.ClearScreenRaw()

	fmt.Println("* Starting up...")
	m := state.Get().Manifest()

	// Setting Program Exit Function
	menu.SetProgramExitFunc(initial.Exit)

	// Start menu workers
	menu.StartWorkers(4)

	// Start input checker
	if err := menu.StartInput(); err != nil {
		log.Fatal(fmt.Errorf("failed to start menu workers: %w", err))
	}

	// Backup if first run
	if !m.Initialized {
		go func() {
			res, err := services.PerformInitialBackup(m.Paths)
			if err != nil {
				errors.Report("startup.backup", err)
				return
			}

			if err := state.Get().Write(func(s *state.State) error {
				if err := services.ApplyInitialBackup(s.Manifest(), res, s.Manifest().Paths); err != nil {
					errors.Report("startup.backup.apply", err)
				}
				return nil
			}); err != nil {
				log.Fatal(fmt.Errorf("OH NOOOO ): failed to write initial backup to state: %w", err))
			}
		}()
	}

	// refresh meta data of installed package entries
	go refreshMetaData(state.Get().Manifest().Paths, m, state.Get().MetaData())

	// Initialize menus
	minit.InitializeMenus(state.Get().Manifest())
}
