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
	"github.com/w1lam/mc-pacman/internal/usecases"
)

// Resolver handles resolving remote packages
type Resolver struct {
	usecases.Base

	modClient *modrinth.Client
}

// New creates a new resolver
func New(base usecases.Base, agent string) *Resolver {
	cfg := modrinth.Config{
		BaseURL: "",
		Agent:   agent,
		HTTP:    nil,
	}

	c := modrinth.NewClient(cfg)

	r := Resolver{
		Base:      base,
		modClient: c,
	}

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
	Type     packages.EntryTypeID
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

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				r.EmitError(cOp, fmt.Errorf("cancelled_%w", ctx.Err()), fmt.Sprintf("operation canelled: %v", ctx.Err()))
				return ctx.Err()
			}
			defer func() { <-sem }()

			resolved, err := r.resolveOne(ctx, op, entry, filter)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failures = append(failures, ResolutionFailure{
					EntryID: entry.ID,
					Err:     err,
				})
			} else {
				files = append(files, resolved)
				atomic.AddInt32(&completed, 1)
			}

			return nil
		})

	}
	if err := g.Wait(); err != nil {
		r.EmitError(op, err, "")
		return ResolvedPackage{}, err
	}

	r.EmitComplete(op, fmt.Sprintf("resolved %d/%d entries", completed, len(pkg.Entries)))

	return ResolvedPackage{
		Remote:   pkg,
		Files:    files,
		Failures: failures,
	}, nil
}

func (r *Resolver) resolveOne(ctx context.Context, parentOp events.Operation, entry packages.RemoteEntry, filter modrinth.VersionFilter) (ResolvedFile, error) {
	op := r.StartOp(parentOp, fmt.Sprintf("resolve_%s", entry.ID))
	r.EmitStart(op, "")

	version, err := r.modClient.ResolveBestVersion(
		ctx,
		string(entry.ID),
		entry.PinnedVer,
		filter,
	)
	if err != nil {
		r.EmitError(op, err, "")
		return ResolvedFile{}, err
	}

	if version == nil {
		r.EmitInfo(op, "no version found, trying fallback filter")

		fallbackFilter := modrinth.VersionFilter{Loader: filter.Loader}
		version, err = r.modClient.ResolveBestVersion(ctx, string(entry.ID), entry.PinnedVer, fallbackFilter)
		if err != nil {
			r.EmitError(op, err, "")
			return ResolvedFile{}, err
		}

		if version == nil {
			err := fmt.Errorf("no_matching_version_for_%s", entry.ID)
			r.EmitError(op, err, "")
			return ResolvedFile{}, err
		}

		r.EmitInfo(op, fmt.Sprintf("using fallback version for: %s", version.VersionNumber))
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
		Type:     entry.Type,
		Version:  version.VersionNumber,
		FileName: primary.FileName,
		Size:     primary.Size,
		URL:      primary.URL,
		Hash:     primary.Hashes.Sha512,
		Algo:     "shaa512",
	}

	r.EmitComplete(op, fmt.Sprintf("resolved to version %s", version.VersionNumber))
	return resolved, nil
}
