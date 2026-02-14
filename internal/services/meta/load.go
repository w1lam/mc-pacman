package meta

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadMetaData loads metadata
func Load(path string) *MetaData {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var m MetaData
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}

	return &m
}

// Save metadata
func (md *MetaData) Save(path string) error {
	tmp := path + ".tmp"

	data, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall metadata: %s", err)
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write metadata temp file: %s", err)
	}

	return os.Rename(tmp, path)
}
