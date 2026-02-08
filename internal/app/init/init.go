// Package app is the high level core of the program
package init

import (
	"fmt"
	"log"

	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/core/verify"
)

// InitCore core functionality initiation
func InitCore() error {
	paths, err := paths.Resolve()
	if err != nil {
		return err
	}

	if err := filesystem.EnsureDirectories(paths); err != nil {
		return err
	}

	// Start error handler
	if err := errors.Start(paths.LogPath); err != nil {
		log.Fatal(fmt.Errorf("failed to start error handler: %w", err))
	}

	// Init State
	if err := state.InitState(paths); err != nil {
		log.Fatal(fmt.Errorf("failed to init state: %w", err))
	}

	// Verify packages
	verify.VerifyAndReconcile(paths)

	return nil
}
