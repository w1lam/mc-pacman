// Package resolver holds resolver service
package resolver

import (
	"context"
	"fmt"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

type ResolvedFile struct {
	ID       modrinth.ID
	Version  string
	FileName string
	Size     int64
	URL      string
	Hash     string
	Algo     string
}

// Resolve resolves a remote package to a slice of downloader.FileRequest ready for download
func (r *Service) Resolve(
	ctx context.Context,
	pkg packages.RemotePackage,
) ([]ResolvedFile, error) {
	filter := modrinth.Filter{
		McVersion: pkg.McVersion,
		Loader:    pkg.Loader,
	}

	versions := modrinth.FetchBestVersions(pkg.Entries, filter)

	var files []ResolvedFile

	for id, version := range versions {
		if version == nil {
			return nil, fmt.Errorf("no matching version for: %s", id)
		}

		var primaryFile *modrinth.MRFile
		for _, f := range version.Files {
			if f.Primary {
				primaryFile = &f
				break
			}
		}

		if primaryFile == nil {
			return nil, fmt.Errorf("no primary files for %s", id)
		}

		files = append(files, ResolvedFile{
			ID:       id,
			Version:  version.VersionNumber,
			FileName: primaryFile.FileName,
			Size:     primaryFile.Size,
			URL:      primaryFile.URL,
			Hash:     primaryFile.Hashes.Sha512,
			Algo:     "sha512",
		})
	}

	return files, nil
}
