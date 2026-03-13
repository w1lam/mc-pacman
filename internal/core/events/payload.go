package events

import "github.com/w1lam/mc-pacman/internal/core/packages"

// Payload is a payload of data
type Payload struct {
	Package  *PackageItem
	Packages []PackageItem

	Progress *Progress
}

// PackageItem is a package item payload
type PackageItem struct {
	packages.PackageBase
	Installed bool
}

// Progress is a progress payload
type Progress struct {
	Label      string
	Current    int64
	Total      int64
	Percentage float64
}
