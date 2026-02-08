package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/w1lam/Packages/utils"
)

// Load loads the manifest
func Load(path string) (*Manifest, error) {
	if !utils.CheckFileExists(path) {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	if m.Path == "" {
		m.Path = path
	}

	return &m, nil
}

// Save saves the manifest to the specified path atomically.
func (m *Manifest) Save() error {
	tmp := m.Path + ".tmp"

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall manifest: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest temp file: %s", err)
	}

	return os.Rename(tmp, m.Path)
}
