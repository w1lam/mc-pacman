package tui

// // InitTUI initializes the tui
// func InitTUI() {
// 	tui.EnableANSI()
// 	tui.HideCursor()
// 	tui.ClearScreenRaw()
//
// 	fmt.Println("* Starting up...")
// 	var m *manifest.Manifest
// 	var p *paths.Paths
// 	state.Get().Read(func(s *state.State) {
// 		m = s.Manifest()
// 		p = s.Paths()
// 	})
//
// 	// Setting Program Exit Function
// 	menu.SetProgramExitFunc(app.Exit)
//
// 	// Start menu workers
// 	menu.StartWorkers(4)
//
// 	// Start input checker
// 	if err := menu.StartInput(); err != nil {
// 		log.Fatal(fmt.Errorf("failed to start menu workers: %w", err))
// 	}
//
// 	// Backup if first run
// 	if !m.Initialized {
// 		go func() {
// 			res, err := services.PerformInitialBackup(p)
// 			if err != nil {
// 				errors.Report("startup.backup", err)
// 				return
// 			}
//
// 			if err := state.Get().Write(func(s *state.State) error {
// 				if err := services.ApplyInitialBackup(s.Manifest(), res, s.Manifest().Paths); err != nil {
// 					errors.Report("startup.backup.apply", err)
// 				}
// 				return nil
// 			}); err != nil {
// 				log.Fatal(fmt.Errorf("OH NOOOO ): failed to write initial backup to state: %w", err))
// 			}
// 		}()
// 	}
//
// 	// refresh meta data of installed package entries
// 	go refreshMetaData(state.Get().Manifest().Paths, m, state.Get().MetaData())
//
// 	// Initialize menus
// 	minit.InitializeMenus(state.Get().Manifest())
// }
