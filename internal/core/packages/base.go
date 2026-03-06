package packages

// PackageBase is the base of a package
type PackageBase struct {
	Name        string    `json:"name"`
	ID          PkgID     `json:"id"`
	Type        PkgTypeID `json:"type"`
	PkgVersion  string    `json:"pkgVersion"`
	McVersion   string    `json:"mcVersion"`
	Loader      string    `json:"loader"`
	Env         string    `json:"env"`
	Description string    `json:"description"`
	ListSource  string    `json:"listSource"`
}

// Package is the base interface of remote and installed packages
type Package interface {
	GetBase() PackageBase
	IsInstalled() bool
}
