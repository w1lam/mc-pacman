package verify

import (
	"fmt"

	errors "github.com/w1lam/mc-pacman/internal/core/errors"
	manifest "github.com/w1lam/mc-pacman/internal/core/manifest"
)

func ReconcileActiveWithManifest(m *manifest.Manifest, found []manifest.InstalledPackage) {
	for _, pkg := range found {
		_, ok := m.InstalledPackages[pkg.Type][pkg.Name]

		if !ok {
			errors.ReportCtx("startup.reconile",
				fmt.Errorf("active package not in manifest"),
				map[string]string{
					"name": pkg.Name,
					"type": string(pkg.Type),
				},
			)
			m.InstalledPackages[pkg.Type][pkg.Name] = pkg
		}

		// LAST SCANNED WINS MIGHT NEED CHANGING IN THE FUTURE
		m.EnabledPackages[pkg.Type] = pkg.Name
	}
}
