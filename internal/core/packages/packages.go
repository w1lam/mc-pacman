// Package packages
package packages

import "github.com/w1lam/Packages/modrinth"

// PkgID is the ID of a package
type PkgID string

// PkgType is the type of a package
type PkgType int

const (
	// modpack type constant
	PackageTypeModPack PkgType = iota
	// resourcebundle type constant
	PackageTypeResourceBundle
	// shaderbundle type constant
	PackageTypeShaderBundle
	// datapackbundle type constant
	PackageTypeDataPack
)

// PackageType is the type of a package modpack/resourcebundle/shaderbundle?
type PackageType struct {
	PackageType PkgType            `json:"pkgType"`
	TypeName    string             `json:"typeName"`
	EntryType   modrinth.EntryType `json:"entryType"`
	ActivePath  string             `json:"-"`
	StorageDir  string             `json:"-"` // package types storage dir ie .mc-pacman/packages/modpacks
}

// Pkg is a small pacakge struct used for passing around packages
type Pkg struct {
	Title string
	ID    PkgID
	Type  PkgType
}

// PackageTypeIndex is the index of all PackageTypes with PkgType as key
var PackageTypeIndex = map[PkgType]PackageType{
	PackageTypeModPack: {
		PackageType: PackageTypeModPack,
		TypeName:    "Modpack",
		EntryType:   modrinth.Mod,
		ActivePath:  "mods",     // relative from mc dir
		StorageDir:  "modpacks", // relative from mcDir/.mc-pacman/packages dir
	},
	PackageTypeResourceBundle: {
		PackageType: PackageTypeResourceBundle,
		TypeName:    "Resource Bundle",
		EntryType:   modrinth.Resourcepack,
		ActivePath:  "resourcepacks",   // relative from mc dir
		StorageDir:  "resourcebundles", // relative from mcDir/.mc-pacman/packages dir
	},
	PackageTypeShaderBundle: {
		PackageType: PackageTypeShaderBundle,
		TypeName:    "Shader Bundle",
		EntryType:   modrinth.Shaderpack,
		ActivePath:  "shaderpacks",   // relative from mc dir
		StorageDir:  "shaderbundles", // relative from mcDir/.mc-pacman/packages dir
	},
	// not really used not high on prio
	PackageTypeDataPack: {
		PackageType: PackageTypeDataPack,
		TypeName:    "Datapack",
		EntryType:   modrinth.Datapack,
		ActivePath:  "",                // relative from mc dir
		StorageDir:  "datapackbundles", // relative from mcDir/.mc-pacman/packages dir
	},
}
