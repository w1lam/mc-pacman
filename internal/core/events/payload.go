package events

import "github.com/w1lam/mc-pacman/internal/core/packages"

type Payload struct {
	Package  *PackageItem
	Packages []PackageItem
}

type PackageItem struct {
	packages.PackageBase
	Installed bool
}
