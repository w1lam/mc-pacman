package paths

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PathsRepository interface {
	Load() (*Paths, error)
	Save(*Paths) error
}

type JSONPathsRepository struct {
	file string
}

func NewPathRepo() *JSONPathsRepository {
	return &JSONPathsRepository{
		file: filepath.Join(rootDir(), "paths.json"),
	}
}

func (r *JSONPathsRepository) Load() (*Paths, error) {
	data, err := os.ReadFile(r.file)
	if os.IsNotExist(err) {
		return nil, nil // first run
	}
	if err != nil {
		return nil, err
	}

	var p Paths
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *JSONPathsRepository) Save(p *Paths) error {
	data, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(r.file), 0755); err != nil {
		return err
	}

	return os.WriteFile(r.file, data, 0644)
}
