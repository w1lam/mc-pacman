package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type Repository interface {
	Load() (*Manifest, error)
	Save(*Manifest) error
}

type FileRepository struct {
	path string
	mu   sync.Mutex
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{path: path}
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

func (r *FileRepository) Init() (*Manifest, error) {
	m, err := r.Load()
	if err == nil {
		return m, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	fmt.Println(" * Building Initial Manifest...")

	m = &Manifest{
		SchemaVersion:     1,
		EnabledPackages:   make(EnabledPackages),
		InstalledPackages: packages.BlankInstalledPackageIndex(),
		InstalledLoaders:  []LoaderInfo{},
		Backups:           []BackupEntry{},
		Initialized:       false,
	}

	if err := r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}
