// Package activation handles package enabling/disabling
package activation

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
	"github.com/w1lam/mc-pacman/internal/usecases"
)

// TODO: IMPLEMENT ACTIVATOR IE ENABLE/DISABLE

type Activator struct {
	usecases.Base

	paths *paths.Paths
	iRepo packages.InstalledRepo
	sRepo state.Repo
}

func New(base usecases.Base, iRepo packages.InstalledRepo, sRepo state.Repo, paths *paths.Paths) *Activator {
	return &Activator{
		Base:  base,
		paths: paths,
		iRepo: iRepo,
		sRepo: sRepo,
	}
}

// Enable enables the specified package
func (a *Activator) Enable(ctx context.Context, pkgID packages.PkgID) error {
	if err := a.paths.Validate(); err != nil {
		return err
	}

	pOp, _ := events.OpFromCtx(ctx)
	op := a.StartOp(pOp, fmt.Sprintf("enable_package_%s", pkgID))
	a.EmitStart(op, fmt.Sprintf("enabling %s", pkgID))
	defer a.EmitEnd(op)

	pkg, err := a.iRepo.GetByID(ctx, pkgID)
	if err != nil {
		return fmt.Errorf("package not found: %w", err)
	}

	st, err := a.sRepo.Load()
	if err != nil {
		return err
	}

	if currentID, ok := st.EnabledPackageIDs[pkg.Type]; ok && currentID != pkgID {
		if err := a.Disable(ctx, pkg.Type); err != nil {
			return fmt.Errorf("failed to disable current package: %w", err)
		}
	}

	entriesDir := filepath.Join(a.paths.PackagesDir(), string(pkgID), "entries")
	if err := a.copyEntries(op, pkg, entriesDir); err != nil {
		return err
	}

	return a.sRepo.Update(func(s *state.State) error {
		s.EnabledPackageIDs[pkg.Type] = pkgID
		return nil
	})
}

// Disable disables a package of specified package type.
func (a *Activator) Disable(ctx context.Context, pkgType packages.PkgTypeID) error {
	pOp, _ := events.OpFromCtx(ctx)
	op := a.StartOp(pOp, fmt.Sprintf("disable_package_type_%s", pkgType))
	a.EmitStart(op, fmt.Sprintf("disabling %s", pkgType))
	defer a.EmitEnd(op)

	st, err := a.sRepo.Load()
	if err != nil {
		return err
	}

	pkgID, ok := st.EnabledPackageIDs[pkgType]
	if !ok {
		return nil
	}

	pkg, err := a.iRepo.GetByID(ctx, pkgID)
	if err != nil {
		return err
	}

	if err := a.removeEntries(op, pkg); err != nil {
		return err
	}

	return a.sRepo.Update(func(s *state.State) error {
		delete(s.EnabledPackageIDs, pkgType)
		return nil
	})
}
