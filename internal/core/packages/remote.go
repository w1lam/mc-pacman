package packages

import (
	"context"

	"github.com/w1lam/Packages/modrinth"
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

	Type       PackageType `json:"pkgType"`
	ListSource string      `json:"listSource"`
	Hash       string      `json:"hash"` // sha512

	Entries []modrinth.Entry `json:"entries"`
}

// BlankRemotePackageIndex creates an empty RemotePackagesIndex and safely initiates all maps
func BlankRemotePackageIndex() RemotePackageIndex {
	index := make(RemotePackageIndex)

	for pkgType := range PackageTypeIndex {
		index[pkgType] = make(map[PkgID]RemotePackage)
	}

	return index
}

// RemotePackageIndex are all available packages found in github repo
type RemotePackageIndex map[PkgType]map[PkgID]RemotePackage

// RemoteRepo is the interface for fetching packages from remote sources (e.g. github)
type RemoteRepo interface {
	GetAll(ctx context.Context) (RemotePackageIndex, error)
	GetByID(ctx context.Context, id PkgID) (RemotePackage, error)
}
