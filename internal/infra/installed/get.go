package installed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

// Exists returns true if package exists
func (r *repo) Exists(pkgID packages.PkgID) (bool, error) {
	pkgDir := filepath.Join(r.path, string(pkgID))
	if _, err := os.Stat(pkgDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Add creates package dir, moves entries to entries folder and writes pkg.json
func (r *repo) Add(p packages.InstalledPackage, entriesSrcDir string) error {
	pkgDir := filepath.Join(r.path, string(p.ID))
	entriesDir := filepath.Join(pkgDir, "entries")

	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		return fmt.Errorf("failed to create entries dir: %w", err)
	}

	if err := os.Rename(entriesSrcDir, entriesDir); err != nil {
		return fmt.Errorf("failed to move entries dir: %w", err)
	}

	data, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal pkg.json: %w", err)
	}

	if err := os.WriteFile(filepath.Join(pkgDir, "pkg.json"), data, 0o644); err != nil {
		return fmt.Errorf("failed to write pkg.json: %w", err)
	}

	return nil
}

// Remove removes the a package completly from packages/
func (r *repo) Remove(pkgID packages.PkgID) error {
	pkgDir := filepath.Join(r.path, string(pkgID))

	if _, err := os.Stat(pkgDir); err != nil {
		return fmt.Errorf("package %s does not exist", pkgID)
	}

	return os.RemoveAll(pkgDir)
}

// Update updates a pkg.json file overwriting it
func (r *repo) Update(p packages.InstalledPackage) error {
	pkgDir := filepath.Join(r.path, string(p.ID))
	tmpFile := filepath.Join(pkgDir, "pkg.json.tmp")
	finalFile := filepath.Join(pkgDir, "pkg.json")

	data, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tmpFile, finalFile)
}

// GetAll gets all installed packages
func (r *repo) GetAll(ctx context.Context) ([]packages.InstalledPackage, error) {
	parentOp, _ := events.OpFromCtx(ctx)
	op := r.StartOp(parentOp, "get_installed_packages")
	r.EmitStart(op, "")
	defer r.EmitEnd(op)

	paths, err := r.scanDir()
	if err != nil {
		return nil, err
	}

	pkgs := make([]packages.InstalledPackage, 0, len(paths))

	for _, path := range paths {
		p := filepath.Join(path, "pkg.json")

		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}

		var installedPkg packages.InstalledPackage
		if err := json.Unmarshal(data, &installedPkg); err != nil {
			continue
		}

		pkgs = append(pkgs, installedPkg)
	}

	return pkgs, nil
}

// GetByID gets an installed package with given PkgID
func (r *repo) GetByID(ctx context.Context, pkgID packages.PkgID) (packages.InstalledPackage, error) {
	parentOp, _ := events.OpFromCtx(ctx)
	op := r.StartOp(parentOp, fmt.Sprintf("get_installed_%s", pkgID))
	r.EmitStart(op, "")
	defer r.EmitEnd(op)

	pkgPath := filepath.Join(r.path, string(pkgID), "pkg.json")

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return packages.InstalledPackage{}, err
	}

	var installedPkg packages.InstalledPackage
	if err := json.Unmarshal(data, &installedPkg); err != nil {
		return packages.InstalledPackage{}, err
	}

	return installedPkg, nil
}

// scanDir scans repo packages dir and returns a slice of the dirs found
func (r *repo) scanDir() ([]string, error) {
	entries, err := os.ReadDir(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(r.path, entry.Name()))
		}
	}

	return dirs, nil
}
