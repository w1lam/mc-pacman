package activation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

func (a *Activator) copyEntries(op events.Operation, pkg packages.InstalledPackage, entriesDir string) error {
	for entryID, entry := range pkg.Entries {
		destDir, err := a.entryDestDir(packages.EntryTypeID(pkg.Entries[entryID].Type))
		if err != nil {
			return err
		}

		src := filepath.Join(entriesDir, entry.FileName)
		dst := filepath.Join(destDir, entry.FileName)

		a.EmitInfo(op, fmt.Sprintf("copying %s -> %s", entry.FileName, destDir))

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy %s: %w", entry.FileName, err)
		}
	}
	return nil
}

func (a *Activator) removeEntries(op events.Operation, pkg packages.InstalledPackage) error {
	for _, entry := range pkg.Entries {
		destDir, err := a.entryDestDir(entry.Type)
		if err != nil {
			return err
		}

		dst := filepath.Join(destDir, entry.FileName)
		a.EmitInfo(op, fmt.Sprintf("removing %s", entry.FileName))

		if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", entry.FileName, err)
		}
	}
	return nil
}

func (a *Activator) entryDestDir(entryType packages.EntryTypeID) (string, error) {
	switch entryType {
	case packages.EntryTypeMod:
		return a.paths.ModsDir(), nil
	case packages.EntryTypeResourcepack:
		return a.paths.ResourcepackDir(), nil
	case packages.EntryTypeShaderpack:
		return a.paths.ShaderpackDir(), nil
	default:
		return "", fmt.Errorf("unknown entry type: %s", entryType)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
