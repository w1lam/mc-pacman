package state

import (
	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
	packages "github.com/w1lam/mc-pacman/internal/core/packages"
	paths "github.com/w1lam/mc-pacman/internal/core/paths"
)

// Get gets the state only read or edit inside Read or Write funcs
func Get() *State {
	if globalState == nil {
		panic("globalState not initialized")
	}
	return globalState
}

// AvailablePackages safe packages accessor
func (s *State) AvailablePackages() *packages.AvailablePackages {
	return s.availablePackages
}

// Manifest safe manifest accessor
func (s *State) Manifest() *manifest.Manifest {
	return s.manifest
}

// Paths small paths reader DO NOT USE INSIDE READ OF WRITE
func Paths() *paths.Paths {
	globalState.mu.RLock()
	defer globalState.mu.RUnlock()
	return globalState.paths
}
