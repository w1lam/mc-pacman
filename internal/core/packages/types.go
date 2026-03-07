// Package packages
package packages

// PkgID is the ID of a package
type PkgID string

// PkgTypeID is the type of a package
type PkgTypeID string

const (
	// modpack type constant
	PackageTypeModPack PkgTypeID = "modpack"
	// resourcebundle type constant
	PackageTypeResourceBundle PkgTypeID = "resourcebundle"
	// shaderbundle type constant
	PackageTypeShaderBundle PkgTypeID = "shaderbundle"
	// datapackbundle type constant
	PackageTypeDataPack PkgTypeID = "datapackbundle"
)

// EntryID is the id of an entry
type EntryID string

// EntryTypeID is the id of an entry
type EntryTypeID string

const (
	// mod type constant
	EntryTypeMod EntryTypeID = "mod"
	// resourcepack type constant
	EntryTypeResourcepack EntryTypeID = "resourcepack"
	// shaderpack type constant
	EntryTypeShaderpack EntryTypeID = "shaderpack"
	// datapack type constant
	EntryTypeDatapack EntryTypeID = "datapack"
)
