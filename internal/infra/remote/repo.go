// Package remote holds remote repo that handles fetching of remote packages from github
package remote

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/ux"
)

const (
	// BaseURL is the base url to github repo for api calls
	BaseURL = "https://api.github.com/repos/w1lam/mc-pacman/"
	// RawURL is the base url to github repo for raw file access
	RawURL = "https://raw.githubusercontent.com/w1lam/mc-pacman/refs/heads/main/"
)

// repo handles fetching remote packages from github
type repo struct {
	events.EmitterBase
	baseURL string
	rawURL  string
	client  *http.Client
}

// New creates a new remote repository
func New(view ux.View) *repo {
	r := repo{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeRemoteRepo,
		},
		baseURL: BaseURL,
		rawURL:  RawURL,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
	r.SetEmitter(view)

	return &r
}

// GetAll gets all available packages from github
func (r *repo) GetAll(ctx context.Context) ([]packages.RemotePackage, error) {
	parentOp, _ := events.OpFromCtx(ctx)
	op := r.StartOp(parentOp, "get_remote_packages")
	r.EmitStart(op, "")
	defer r.EmitEnd(op)

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
func (r *repo) GetByID(ctx context.Context, pkgID packages.PkgID) (packages.RemotePackage, error) {
	parentOp, _ := events.OpFromCtx(ctx)
	op := r.StartOp(parentOp, fmt.Sprintf("get_remote_%s", pkgID))
	r.EmitStart(op, "")
	defer r.EmitEnd(op)

	url := fmt.Sprintf("%spackages/%s/pkg.json", r.rawURL, pkgID)

	var pkg packages.RemotePackage
	if err := r.decodeRemoteJSON(ctx, url, &pkg); err != nil {
		return packages.RemotePackage{}, err
	}

	return pkg, nil
}
