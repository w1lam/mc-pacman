package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type Repo interface {
	Load() (*State, error)
	Save(*State) error
	Update(func(*State) error) error
}

type stateRepo struct {
	path string
	mu   sync.Mutex
}

func NewStateRepo(path string) *stateRepo {
	return &stateRepo{
		path: path,
	}
}

func (r *stateRepo) Update(fn func(*State) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, err := r.loadFromDisk()
	if err != nil {
		return err
	}

	if err := fn(m); err != nil {
		return err
	}

	return r.saveToDisk(m)
}

// Load loads the manifest
func (r *stateRepo) Load() (*State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.loadFromDisk()
}

// loadFromDisk loads from disk safe
func (r *stateRepo) loadFromDisk() (*State, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
	}

	var m State
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// Save saves the given state
func (r *stateRepo) Save(s *State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tmp := r.path + ".tmp"

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall manifest: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest temp file: %s", err)
	}

	return os.Rename(tmp, r.path)
}

// saveToDisk safe save to disk no mutex locking
func (r *stateRepo) saveToDisk(m *State) error {
	tmp := r.path + ".tmp"

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall manifest: %w", err)
	}

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write manifest temp file: %w", err)
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("failed to write manifest temp file: %w", err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("failed to write manifest temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, r.path)
}

// Ensureensures the state is initialized
func (r *stateRepo) Ensure() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := os.Stat(r.path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	s := &State{
		SchemaVersion:     1,
		EnabledPackageIDs: make(map[packages.PkgTypeID]packages.PkgID),
		InstalledLoaders:  []LoaderInfo{},
		Backups:           []BackupEntry{},
	}

	return r.saveToDisk(s)
}
