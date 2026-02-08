package state

import (
	"fmt"
	"log"

	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
	packages "github.com/w1lam/mc-pacman/internal/core/packages"
	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

// InitState initializes the program state
func InitState(paths *paths.Paths) error {
	if paths == nil {
		return fmt.Errorf("InitState: paths cannot be nil")
	}

	m, err := manifest.InitManifest(paths.ManifestPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init manifest: %w", err))
	}

	initState := blankState()

	initState.paths = paths
	initState.manifest = m

	globalState.mu.Lock()
	defer globalState.mu.Unlock()
	SetState(initState)

	return nil
}

func blankState() *State {
	return &State{
		ProgramVersion: ProgramVersion,
		paths:          nil,

		manifest: nil,

		availablePackages: &packages.AvailablePackages{
			ModPacks:            make(map[string]packages.ResolvedPackage),
			ResourcePackBundles: make(map[string]packages.ResolvedPackage),
			ShaderPackBundles:   make(map[string]packages.ResolvedPackage),
			DatapackBundles:     make(map[string]packages.ResolvedPackage),
		},
	}
}
