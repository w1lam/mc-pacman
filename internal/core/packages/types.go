package packages

import "github.com/w1lam/Packages/modrinth"

// PackageType is the type of a package modpack/resourcebundle/shaderbundle?
type PackageType struct {
	EntryType  modrinth.EntryType `json:"entryType"`
	ActivePath string             `json:"activePath"`
	StorageDir string             `json:"storageDir"` // package types storage dir ie .mc-pacman/packages/modpacks
}

var (
	PackageModPack = PackageType{
		EntryType:  modrinth.Mod,
		ActivePath: "mods",     // relative from mc dir
		StorageDir: "modpacks", // relative from mcDir/.mc-pacman/packages dir
	}
	PackageResourceBundle = PackageType{
		EntryType:  modrinth.Resourcepack,
		ActivePath: "resourcepacks",   // relative from mc dir
		StorageDir: "resourcebundles", // relative from mcDir/.mc-pacman/packages dir
	}
	PackageShaderBundles = PackageType{
		EntryType:  modrinth.Shaderpack,
		ActivePath: "shaderpacks",   // relative from mc dir
		StorageDir: "shaderbundles", // relative from mcDir/.mc-pacman/packages dir
	}
	// not really used not high on prio
	PackageDatapackBundles = PackageType{
		EntryType:  modrinth.Datapack,
		ActivePath: "",                // relative from mc dir
		StorageDir: "datapackbundles", // relative from mcDir/.mc-pacman/packages dir
	}
)
