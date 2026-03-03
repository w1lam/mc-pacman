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

type ResolvedPackage struct {
	Remote packages.RemotePackage
	Files  []ResolvedFile
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
	parentOp events.Operation,
	pkg packages.RemotePackage,
) (ResolvedPackage, error) {
	op := r.StartOp(parentOp, fmt.Sprintf("resolve %s", pkg.ID))
	r.EmitStart(op, fmt.Sprintf("resolving package: %s", pkg.ID))
	defer r.EmitEnd(op)

	filter := modrinth.VersionFilter{
		GameVersion: pkg.McVersion,
		Loader:      pkg.Loader,
	}

	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, 5)

	files := make([]ResolvedFile, 0, len(pkg.Entries))
	var mu sync.Mutex

	total := len(pkg.Entries)
	_ = total
	var completed int32

	for _, entry := range pkg.Entries {
		entry := entry

		g.Go(func() error {
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return ctx.Err()
			}
			defer func() { <-sem }()

			r.Emit(events.Event{})

			version, err := r.modClient.ResolveBestVersion(
				ctx,
				string(entry.ID),
				entry.PinnedVer,
				filter,
			)
			if err != nil {
				return err
			}

			if version == nil {
				err := fmt.Errorf("no matching version for %s", entry.ID)
				return err
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
			r.Emit(events.Event{})

			_ = newCompleted

			return nil
		})

	}
	if err := g.Wait(); err != nil {
		return ResolvedPackage{}, err
	}

	return ResolvedPackage{
		Remote: pkg,
		Files:  files,
	}, nil
}
