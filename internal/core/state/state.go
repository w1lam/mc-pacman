// Package state holds the state of app
package state

//
// import (
// 	"fmt"
// 	"sync"
//
// 	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
// )
//
// // State is the global state struct
// type State struct {
// 	mu sync.RWMutex
//
// 	manifest *manifest.Manifest
// 	repo     manifest.Repository
// }
//
// // New initializes app state
// func New(m *manifest.Manifest, repo manifest.Repository) (*State, error) {
// 	if m == nil || repo == nil {
// 		return nil, fmt.Errorf("input data is nil")
// 	}
//
// 	return &State{
// 		manifest: m,
// 		repo:     repo,
// 	}, nil
// }
//
// // Manifest accessor
// func (s *State) Manifest() *manifest.Manifest {
// 	return s.manifest
// }
//
// // Repo manifest repo accessor
// func (s *State) Repo() manifest.Repository {
// 	return s.repo
// }
