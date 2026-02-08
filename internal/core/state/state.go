package state

import (
	"sync"

	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
	packages "github.com/w1lam/mc-pacman/internal/core/packages"
	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

var ProgramVersion string = "0.1a"

var (
	globalState *State
	once        sync.Once
)

// State is the global state struct
type State struct {
	ProgramVersion string
	mu             sync.RWMutex

	paths    *paths.Paths
	manifest *manifest.Manifest

	availablePackages *packages.AvailablePackages
}

// SetState sets the global state
func SetState(s *State) {
	once.Do(func() {
		globalState = s
	})
}

func SetAvailablePackages(pkgs *packages.AvailablePackages) {
	globalState.mu.Lock()
	defer globalState.mu.Unlock()

	globalState.availablePackages = pkgs
}
