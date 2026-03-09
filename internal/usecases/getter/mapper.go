package getter

import (
	"fmt"

	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
)

type DownloadedPackage struct {
	Resolved resolver.ResolvedPackage
	Files    []downloader.FileResult
}

func buildDownloadedPackage(resolved resolver.ResolvedPackage, files []downloader.FileResult) DownloadedPackage {
	return DownloadedPackage{
		Resolved: resolved,
		Files:    files,
	}
}

// buildInstalledPackage builds an installed package from a remote package
func buildInstalledPackage(downloaded DownloadedPackage, hash string) (packages.InstalledPackage, error) {
	entries := make(map[packages.EntryID]packages.InstalledPackageEntry)

	resultMap := make(map[packages.EntryID]downloader.FileResult)
	for _, r := range downloaded.Files {
		resultMap[packages.EntryID(r.ID)] = r
	}

	for _, r := range downloaded.Resolved.Files {
		result, ok := resultMap[packages.EntryID(r.ID)]
		if !ok {
			return packages.InstalledPackage{}, fmt.Errorf("missing downloader result for resolved file: %s", string(r.ID))
		}

		entries[r.ID] = buildInstalledPackageEntry(r, result)
	}

	return packages.InstalledPackage{
		PackageBase: downloaded.Resolved.Remote.PackageBase,
		Hash:        hash,
		Entries:     entries,
	}, nil
}

// buildInstalledPackageEntry builds an installed package entry from downloader results
func buildInstalledPackageEntry(resolved resolver.ResolvedFile, result downloader.FileResult) packages.InstalledPackageEntry {
	return packages.InstalledPackageEntry{
		ID:       resolved.ID,
		Type:     resolved.Type,
		Version:  resolved.Version,
		FileName: result.FileName,
		Hash:     result.Hash,
		Algo:     "sha512",
	}
}

// buildFileRequests builds file requests from resolved files
func buildFileRequests(resolvedPackage resolver.ResolvedPackage) []downloader.FileRequest {
	files := make([]downloader.FileRequest, 0, len(resolvedPackage.Files))

	for _, r := range resolvedPackage.Files {
		files = append(files, downloader.FileRequest{
			ID:       string(r.ID),
			URL:      r.URL,
			FileName: r.FileName,
			Size:     r.Size,
			Hash:     r.Hash,
		})
	}

	return files
}
