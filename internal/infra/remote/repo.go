// Package remote holds remote repo that handles fetching of remote packages from github
package remote

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/w1lam/mc-pacman/internal/core/packages"
)

const (
	// BaseURL is the base url to github repo for api calls
	BaseURL = "https://api.github.com/repos/w1lam/mc-pacman/"
	// RawURL is the base url to github repo for raw file access
	RawURL = "https://raw.githubusercontent.com/w1lam/mc-pacman/refs/heads/main/"
)

// RemoteRepository handles fetching remote packages from github
type RemoteRepository struct {
	baseURL string
	rawURL  string
	client  *http.Client

	mu    sync.Mutex
	cache map[packages.PkgID]packages.RemotePackage
}

// New creates a new remote repository
func New() *RemoteRepository {
	return &RemoteRepository{
		baseURL: BaseURL,
		rawURL:  RawURL,
		client: &http.Client{
			Timeout: time.Second * 30,
		},

		cache: make(map[packages.PkgID]packages.RemotePackage),
	}
}

// GetAll gets all available packages from catalogue
func (r *RemoteRepository) GetAll(ctx context.Context) ([]packages.RemotePackage, error) {
	url := fmt.Sprintf("%scontents/packages", r.baseURL)

	var contents []githubContentResponse
	if err := r.decodeRemoteJSON(ctx, url, &contents); err != nil {
		return nil, err
	}

	out := make([]packages.RemotePackage, 0, len(contents))

	for _, item := range contents {
		if item.Type != "dir" {
			continue
		}

		pkgURL := fmt.Sprintf("%spackages/%s/pkg.json", r.rawURL, item.Name)

		var remotePkg packages.RemotePackage
		if err := r.decodeRemoteJSON(ctx, pkgURL, &remotePkg); err != nil {
			return nil, err
		}

		out = append(out, remotePkg)
	}

	return out, nil
}

// GetByID gets a remote package by id
func (r *RemoteRepository) GetByID(ctx context.Context, id packages.PkgID) (packages.RemotePackage, error) {
	url := fmt.Sprintf("%spackages/%s/pkg.json", r.rawURL, id)

	var pkg packages.RemotePackage
	if err := r.decodeRemoteJSON(ctx, url, &pkg); err != nil {
		return packages.RemotePackage{}, err
	}

	return pkg, nil
}
