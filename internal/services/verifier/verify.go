// Package verifier hold verifier service
package verifier

import (
	"fmt"
	"path/filepath"

	"github.com/w1lam/mc-pacman/internal/core/errors"
	"github.com/w1lam/mc-pacman/internal/core/filesystem"
	"github.com/w1lam/mc-pacman/internal/core/manifest"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/paths"
	"github.com/w1lam/mc-pacman/internal/events"
	"github.com/w1lam/mc-pacman/internal/packageid"
	"github.com/w1lam/mc-pacman/internal/services"
)

type Verifier struct {
	services.Base
	paths *paths.Paths
	repo  manifest.Repo
}

func New(p *paths.Paths, r manifest.Repo) *Verifier {
	return &Verifier{
		Base: services.Base{
			Scope:   events.ScopeVerifier,
			Emitter: nil,
		},
		paths: p,
		repo:  r,
	}
}

func (v *Verifier) Verify() error {
	op := events.NewOperation(v.Scope, "verify")
	v.Emit(events.Event{
		Type: events.EventStart,
		Op:   op,
	})

	m, err := v.repo.Load()
	if err != nil {
		return err
	}

	v.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: "scanning packages",
	})
	found, err := v.scanActitvePackages()
	if err != nil {
		errors.Report("startup.scan", err)
	}

	v.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: "reconciling with manifest",
	})
	v.reconcile(m, found)

	v.Emit(events.Event{
		Type:    events.EventUpdate,
		Op:      op,
		Message: "verifying manifest packages",
	})
	v.verifyManifestPackages(m)

	v.Emit(events.Event{
		Type: events.EventComplete,
		Op:   op,
	})
	return v.repo.Save(m)
}

func (v *Verifier) verifyManifestPackages(m *manifest.Manifest) {
	for pkgType, pkgs := range m.InstalledPackages {
		for id, ip := range pkgs {
			ok, err := v.verifyPackageIntegrity(m, packages.Pkg{
				Title: ip.Name,
				ID:    id,
				Type:  pkgType,
			})
			if err != nil {
				errors.ReportCtx(
					"startup.verify",
					err,
					map[string]string{
						"name": ip.Name,
						"type": ip.Type.TypeName,
					},
				)
				continue
			}

			if !ok {
				errors.ReportCtx(
					"startup.verify",
					fmt.Errorf("package integrity mismath"),
					map[string]string{
						"name": ip.Name,
						"type": ip.Type.TypeName,
					},
				)
			}
		}
	}
}

// VerifyPackageIntegrity verifies packages integrity by computing hash and comparing it to hash in *PkgID*.id.json
func (v *Verifier) verifyPackageIntegrity(m *manifest.Manifest, pkg packages.Pkg) (bool, error) {
	iPkg, found := m.InstalledPackages[pkg.Type][pkg.ID]
	if !found {
		return false, fmt.Errorf("package not installed: %s", pkg.Title)
	}

	enabled := m.EnabledPackages[pkg.Type] == pkg.ID

	var path string
	if enabled {
		path = iPkg.FullActivePath
	} else {
		path = iPkg.FullStoragePath
	}

	hash, err := filesystem.ComputeDirHash(path)
	if err != nil {
		return false, err
	}

	pkgIDPath := filepath.Join(path, string(pkg.ID)+".id.json")
	pkgID, err := packageid.ReadPackageIDFile(pkgIDPath)
	if err != nil {
		return false, err
	}

	return hash == pkgID.Hash, nil
}

func (v *Verifier) reconcile(m *manifest.Manifest, found []packages.InstalledPackage) {
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
