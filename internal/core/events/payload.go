package events

import (
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type PackagePayload struct {
	ID          string
	Name        string
	Type        string
	Version     string
	McVersion   string
	Loader      string
	Description string
	ListSource  string
	Installed   bool
}

func NewPackagePayload(p packages.PackageBase, installed bool) PackagePayload {
	return PackagePayload{
		ID:          string(p.ID),
		Name:        p.Name,
		Type:        string(p.Type),
		Version:     p.PkgVersion,
		McVersion:   p.McVersion,
		Loader:      p.Loader,
		Description: p.Description,
		ListSource:  p.ListSource,
		Installed:   installed,
	}
}
