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

// type packageType struct {
// 	PackageType PkgTypeID          `json:"pkgType"`
// 	TypeName    string             `json:"typeName"`
// 	EntryType   modrinth.EntryType `json:"entryType"`
// }
//
// // PackageTypeIndex is the index of all PackageTypes with PkgType as key
// var PackageTypeIndex = map[PkgTypeID]packageType{
// 	PackageTypeModPack: {
// 		PackageType: PackageTypeModPack,
// 		TypeName:    "Modpack",
// 		EntryType:   modrinth.Mod,
// 	},
// 	PackageTypeResourceBundle: {
// 		PackageType: PackageTypeResourceBundle,
// 		TypeName:    "Resource Bundle",
// 		EntryType:   modrinth.Resourcepack,
// 	},
// 	PackageTypeShaderBundle: {
// 		PackageType: PackageTypeShaderBundle,
// 		TypeName:    "Shader Bundle",
// 		EntryType:   modrinth.Shaderpack,
// 	},
// 	// not really used not high on prio
// 	PackageTypeDataPack: {
// 		PackageType: PackageTypeDataPack,
// 		TypeName:    "Datapack",
// 		EntryType:   modrinth.Datapack,
// 	},
// }
