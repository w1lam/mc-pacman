package ux

import (
	"fmt"
	"strings"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

func CliPackageListRenderer(pkgIndex packages.RemotePackageIndex) {
	for pType, pkgs := range pkgIndex {
		fmt.Println()
		fmt.Println(strings.ToUpper(packages.PackageTypeIndex()[pType].TypeName) + "S:")
		fmt.Println(strings.Repeat("-", len(packages.PackageTypeIndex()[pType].TypeName)) + "---")
		for _, pkg := range pkgs {
			fmt.Println(" " + pkg.Name)
			fmt.Println(" - List Version:", pkg.ListVersion)
			fmt.Println(" - Minecraft Version:", pkg.McVersion)
			if pkg.Description != "" {
				fmt.Println(" - Desctription:", pkg.Description)
			}
			fmt.Println(" - Env:", pkg.Env)
			if pkg.Loader != "" {
				fmt.Println(" - Loader:", pkg.Loader)
			}
			fmt.Println()
		}
	}
}
