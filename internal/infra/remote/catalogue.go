package remote

import "github.com/w1lam/mc-pacman/internal/core/packages"

type Catalogue struct {
	SchemaVersion int                `json:"schemaVersion"`
	Packages      []CataloguePackage `json:"packages"`
}

type CataloguePackage struct {
	ID   packages.PkgID   `json:"id"`
	Type packages.PkgType `json:"type"`
}
