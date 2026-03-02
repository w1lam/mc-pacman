package packages

import (
	"context"
)

// RemotePackage is a resolved package that gets passed to the installer
type RemotePackage struct {
	Name string `json:"name"`
	ID   PkgID  `json:"id"`

	ListVersion string `json:"listVersion"`
	McVersion   string `json:"mcVersion"`
	Loader      string `json:"loader"`
	Env         string `json:"env"`
	Description string `json:"description"`

	Type       PkgTypeID `json:"pkgType"`
	ListSource string    `json:"-"`

	Entries []RemoteEntry `json:"entries"`
}

// RemoteEntry is an entry in a remote package
type RemoteEntry struct {
	ID        EntryID `json:"id"`
	PinnedVer string  `json:"pinnedVer"`
}

// RemoteRepo is the interface for fetching packages from remote sources (e.g. github)
type RemoteRepo interface {
	GetAll(ctx context.Context) (map[PkgID]RemotePackage, error)
	GetByID(ctx context.Context, id PkgID) (RemotePackage, error)
}
