// Package resolver holds resolver service
package resolver

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/ux"
)

// Resolver handles resolving remote packages
type Resolver struct {
	events.EmitterBase
	modClient *modrinth.Client
}

// New creates a new resolver
func New(view ux.View, agent string) *Resolver {
	cfg := modrinth.Config{
		BaseURL: "",
		Agent:   agent,
		HTTP:    nil,
	}

	c := modrinth.NewClient(cfg)

	r := Resolver{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeResolver,
		},
		modClient: c,
	}

	r.SetEmitter(view)
	return &r
}

type ResolutionFailure struct {
	EntryID packages.EntryID
	Err     error
}

type ResolvedPackage struct {
	Remote   packages.RemotePackage
	Files    []ResolvedFile
	Failures []ResolutionFailure
}

// ResolvedFile is a file that has been resolved and is ready for download
type ResolvedFile struct {
	ID       packages.EntryID
	Version  string
	FileName string
	Size     int64
	URL      string
	Hash     string
	Algo     string
}

// Resolve resolves a remote package to a slice of downloader.FileRequest ready for download
func (r *Resolver) Resolve(
	ctx context.Context,
	pkg packages.RemotePackage,
) (ResolvedPackage, error) {
	pOp, _ := events.OpFromCtx(ctx)
	op := r.StartOp(pOp, fmt.Sprintf("resolve_package_%s", pkg.ID))
	r.EmitStart(op, fmt.Sprintf("resolving package: %s", pkg.ID))
	defer r.EmitEnd(op)

	filter := modrinth.VersionFilter{
		GameVersion: pkg.McVersion,
		Loader:      pkg.Loader,
	}

	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, 5)

	files := make([]ResolvedFile, 0, len(pkg.Entries))
	failures := make([]ResolutionFailure, 0, len(pkg.Entries))
	var mu sync.Mutex

	total := len(pkg.Entries)
	_ = total
	var completed int32

	for _, entry := range pkg.Entries {
		entry := entry

		g.Go(func() error {
			cOp := r.StartOp(op, fmt.Sprintf("resolving_%s_%s", entry.Type, entry.ID))
			r.EmitStart(cOp, "")
			defer r.EmitEnd(cOp)

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				r.EmitError(cOp, fmt.Errorf("cancelled: %w", ctx.Err()))
				return ctx.Err()
			}
			defer func() { <-sem }()

			version, err := r.modClient.ResolveBestVersion(
				ctx,
				string(entry.ID),
				entry.PinnedVer,
				filter,
			)
			if err != nil {
				r.EmitError(cOp, err)
				return err
			}

			if version == nil {
				fallbackFilter := modrinth.VersionFilter{Loader: filter.Loader}
				version, err = r.modClient.ResolveBestVersion(ctx, string(entry.ID), entry.PinnedVer, fallbackFilter)
				if err != nil {
					r.EmitError(cOp, err)
					return err
				}

				if version == nil {
					mu.Lock()
					failures = append(failures, ResolutionFailure{
						EntryID: entry.ID,
						Err:     fmt.Errorf("no_matching_version_for_%s", entry.ID),
					})
					mu.Unlock()
					r.EmitError(cOp, failures[len(failures)-1].Err)
					return nil
				}
				r.EmitInfo(cOp, fmt.Sprintf("fallback_version_%s_used_for_%s", version.VersionNumber, entry.ID))
			}

			var primary *modrinth.File
			for i := range version.Files {
				if version.Files[i].Primary {
					primary = &version.Files[i]
					break
				}
			}

			resolved := ResolvedFile{
				ID:       entry.ID,
				Version:  version.VersionNumber,
				FileName: primary.FileName,
				Size:     primary.Size,
				URL:      primary.URL,
				Hash:     primary.Hashes.Sha512,
				Algo:     "shaa512",
			}

			mu.Lock()
			files = append(files, resolved)
			mu.Unlock()

			newCompleted := atomic.AddInt32(&completed, 1)

			_ = newCompleted

			r.EmitComplete(cOp, fmt.Sprintf("version_%s_found_for_%s", version.VersionNumber, entry.ID))

			return nil
		})

	}
	if err := g.Wait(); err != nil {
		return ResolvedPackage{}, err
	}

	r.EmitComplete(op, fmt.Sprintf("version_found_for_%d/%d_entries", completed, len(pkg.Entries)))

	return ResolvedPackage{
		Remote:   pkg,
		Files:    files,
		Failures: failures,
	}, nil
}
