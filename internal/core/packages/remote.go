package packages

import (
	"context"
)

// RemotePackage is a resolved package that gets passed to the installer
type RemotePackage struct {
	PackageBase

	Entries []RemoteEntry `json:"entries"`
}

// RemoteEntry is an entry in a remote package
type RemoteEntry struct {
	ID        EntryID     `json:"id"`
	PinnedVer string      `json:"pinnedVer"`
	Type      EntryTypeID `json:"type"`
}

func (p RemotePackage) GetBase() PackageBase {
	return p.PackageBase
}

func (p RemotePackage) IsInstalled() bool {
	return false
}

// RemoteRepo is the interface for fetching packages from remote sources (e.g. github)
type RemoteRepo interface {
	GetAll(ctx context.Context) ([]RemotePackage, error)
	GetByID(ctx context.Context, id PkgID) (RemotePackage, error)
}
