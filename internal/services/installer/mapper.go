package installer

import (
	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/services/downloader"
	"github.com/w1lam/mc-pacman/internal/services/resolver"
)

// buildInstalledPackage builds an installed package from a remote package
func buildInstalledPackage(
	remote packages.RemotePackage,
	resolved []resolver.ResolvedFile,
	results []downloader.FileResult,
	storagePath string,
	activePath string,
	fullHash string,
) packages.InstalledPackage {
	entries := make(map[modrinth.ID]packages.InstalledPackageEntry)

	resultMap := make(map[modrinth.ID]downloader.FileResult)
	for _, r := range results {
		resultMap[r.ID] = r
	}

	for _, r := range resolved {
		result, ok := resultMap[r.ID]
		if !ok {
			panic("missing downloader result for resolved file: " + string(r.ID))
		}

		entries[r.ID] = buildInstalledPackageEntry(r, result)
	}

	return packages.InstalledPackage{
		Name:             remote.Name,
		ID:               remote.ID,
		InstalledVersion: remote.ListVersion,
		McVersion:        remote.McVersion,
		Loader:           remote.Loader,
		Type:             remote.Type,
		ListSource:       remote.ListSource,
		Hash:             fullHash,
		Entries:          entries,
		FullActivePath:   activePath,
		FullStoragePath:  storagePath,
	}
}

// buildInstalledPackageEntry builds an installed package entry from downloader results
func buildInstalledPackageEntry(resolved resolver.ResolvedFile, result downloader.FileResult) packages.InstalledPackageEntry {
	return packages.InstalledPackageEntry{
		ID:               resolved.ID,
		InstalledVersion: resolved.Version,
		FileName:         result.FileName,
		Hash:             result.Hash,
		Algo:             "sha512",
	}
}

// buildFileRequests builds file requests from resolved files
func buildFileRequests(resolvedFiles []resolver.ResolvedFile) []downloader.FileRequest {
	files := make([]downloader.FileRequest, 0, len(resolvedFiles))

	for _, r := range resolvedFiles {
		files = append(files, downloader.FileRequest{
			ID:       r.ID,
			URL:      r.URL,
			FileName: r.FileName,
			Hash:     r.Hash,
			Algo:     r.Algo,
		})
	}

	return files
}
