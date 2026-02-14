package verify

// import (
// 	"fmt"
// 	"path/filepath"
//
// 	"github.com/w1lam/mc-pacman/internal/core/errors"
// 	"github.com/w1lam/mc-pacman/internal/core/filesystem"
// 	"github.com/w1lam/mc-pacman/internal/core/manifest"
// 	"github.com/w1lam/mc-pacman/internal/core/packages"
// 	"github.com/w1lam/mc-pacman/internal/core/paths"
// 	"github.com/w1lam/mc-pacman/internal/core/state"
// )
//
// func VerifyAndReconcile(p *paths.Paths, s *state.State) error {
// 	var m *manifest.Manifest
//
// 	s.Read(func(s *state.State) {
// 		m = s.Manifest()
// 	})
//
// 	found, err := ScanActitvePackages(p)
// 	if err != nil {
// 		errors.Report("startup.scan", err)
// 	}
//
// 	ReconcileActiveWithManifest(m, found)
//
// 	verifyManifestPackages(m)
//
// 	return s.Write(func(s *state.State) error {
// 		return s.Repo().Save(m)
// 	})
// }
//
// func verifyManifestPackages(m *manifest.Manifest) {
// 	for pkgType, pkgs := range m.InstalledPackages {
// 		for id, ip := range pkgs {
// 			ok, err := VerifyPackageIntegrity(m, packages.Pkg{
// 				Title: ip.Name,
// 				ID:    id,
// 				Type:  pkgType,
// 			})
// 			if err != nil {
// 				errors.ReportCtx(
// 					"startup.verify",
// 					err,
// 					map[string]string{
// 						"name": ip.Name,
// 						"type": ip.Type.TypeName,
// 					},
// 				)
// 				continue
// 			}
//
// 			if !ok {
// 				errors.ReportCtx(
// 					"startup.verify",
// 					fmt.Errorf("package integrity mismath"),
// 					map[string]string{
// 						"name": ip.Name,
// 						"type": ip.Type.TypeName,
// 					},
// 				)
// 			}
// 		}
// 	}
// }
//
// // VerifyPackageIntegrity verifies packages integrity by computing hash and comparing it to hash in *PkgID*.id.json
// func VerifyPackageIntegrity(m *manifest.Manifest, pkg packages.Pkg) (bool, error) {
// 	iPkg, found := m.InstalledPackages[pkg.Type][pkg.ID]
// 	if !found {
// 		return false, fmt.Errorf("package not installed: %s", pkg.Title)
// 	}
//
// 	enabled := m.EnabledPackages[pkg.Type] == pkg.ID
//
// 	var path string
// 	if enabled {
// 		path = iPkg.FullActivePath
// 	} else {
// 		path = iPkg.FullStoragePath
// 	}
//
// 	hash, err := filesystem.ComputeDirHash(path)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	pkgIDPath := filepath.Join(path, string(pkg.ID)+".id.json")
// 	pkgID, err := manifest.ReadPackageIDFile(pkgIDPath)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return hash == pkgID.Hash, nil
// }
