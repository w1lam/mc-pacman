package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type Repository interface {
	Load() (*Manifest, error)
	Save(*Manifest) error
	Update(func(*Manifest) error) error
}

type FileRepository struct {
	path string
	mu   sync.Mutex
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{path: path}
}

func (r *FileRepository) Update(fn func(*Manifest) error) error {
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
func (r *FileRepository) Load() (*Manifest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// loadFromDisk loads from disk safe
func (r *FileRepository) loadFromDisk() (*Manifest, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Manifest{
				InstalledPackages: make(map[packages.PkgType]map[packages.PkgID]packages.InstalledPackage),
			}, nil
		}
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// Save saves the manifest to the specified path atomically.
func (r *FileRepository) Save(m *Manifest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tmp := r.path + ".tmp"

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall manifest: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest temp file: %s", err)
	}

	return os.Rename(tmp, r.path)
}

// saveToDisk safe save to disk no mutex locking
func (r *FileRepository) saveToDisk(m *Manifest) error {
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

// EnsureInitialized ensures the manifest is initialized
func (r *FileRepository) EnsureInitialized() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := os.Stat(r.path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	m := &Manifest{
		SchemaVersion:     1,
		EnabledPackages:   make(EnabledPackages),
		InstalledPackages: packages.BlankInstalledPackageIndex(),
		InstalledLoaders:  []LoaderInfo{},
		Backups:           []BackupEntry{},
		Initialized:       false,
	}

	return r.saveToDisk(m)
}
