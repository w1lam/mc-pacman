// Package app is the high level core of the program
package app

import (
	"fmt"
	"log"

	"github.com/w1lam/Raw-Mod-Installer/internal/filesystem"
	"github.com/w1lam/Raw-Mod-Installer/internal/manifest"
	"github.com/w1lam/Raw-Mod-Installer/internal/meta"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
	"github.com/w1lam/Raw-Mod-Installer/internal/state"
	"github.com/w1lam/Raw-Mod-Installer/internal/verify"
	errors "github.com/w1lam/mc-pacman/internal/core/errors"
)

// InitCore core functionality initiation
func InitCore() error {
	path, err := paths.Resolve()
	if err != nil {
		return err
	}

	if err := filesystem.EnsureDirectories(path); err != nil {
		return err
	}

	// Start error handler
	if err := errors.Start(path.LogPath); err != nil {
		log.Fatal(fmt.Errorf("failed to start error handler: %w", err))
	}

	m, err := manifest.Load(path)
	if err != nil {
		m, err = manifest.BuildInitialManifest(state.ProgramVersion, path)
		if err != nil {
			return err
		}
	}

	metaD := meta.LoadMetaData(path)
	if metaD == nil {
		emptyMd := &meta.MetaData{
			SchemaVersion: 1,
			Mods:          make(map[string]meta.ModMetaData),
		}
		metaD = emptyMd
	}

	// Sets global state
	state.SetState(state.NewState(m, metaD))

	// Verify packages
	verify.VerifyAndReconcile(path)

	return nil
}
