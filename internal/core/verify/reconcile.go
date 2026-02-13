package verify

import (
	"fmt"

	errors "github.com/w1lam/mc-pacman/internal/core/errors"
	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
	packages "github.com/w1lam/mc-pacman/internal/core/packages"
)

func ReconcileActiveWithManifest(m *manifest.Manifest, found []packages.InstalledPackage) {
	for _, pkg := range found {
		_, ok := m.InstalledPackages[pkg.Type.PackageType][pkg.ID]

		if !ok {
			errors.ReportCtx("startup.reconile",
				fmt.Errorf("active package not in manifest"),
				map[string]string{
					"name": pkg.Name,
					"type": string(pkg.Type.TypeName),
				},
			)
			m.InstalledPackages[pkg.Type.PackageType][pkg.ID] = pkg
		}

		// LAST SCANNED WINS MIGHT NEED CHANGING IN THE FUTURE
		m.EnabledPackages[pkg.Type.PackageType] = pkg.ID
	}
}
