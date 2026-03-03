package packages

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

type Package interface {
	GetBase() PackageBase
	IsInstalled() bool
}
