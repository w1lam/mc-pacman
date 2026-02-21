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
	BaseURL = "https://api.github.com/repos/w1lam/mc-pacman/"
	RawURL  = "https://raw.githubusercontent.com/w1lam/mc-pacman/refs/heads/main/"
)

type RemoteRepository struct {
	baseURL string
	rawURL  string
	client  *http.Client

	mu    sync.Mutex
	cache *Catalogue
}

// NewCatalogue creates blank catalogue
func NewCatalogue() *Catalogue {
	return &Catalogue{
		SchemaVersion: 1,
		Packages:      []CataloguePackage{},
	}
}

// New creates a new remote repository
func New() *RemoteRepository {
	return &RemoteRepository{
		baseURL: BaseURL,
		rawURL:  RawURL,
		client: &http.Client{
			Timeout: time.Second * 30,
		},

		cache: NewCatalogue(),
	}
}

// GetAll gets all available packages from catalogue
func (r *RemoteRepository) GetAll(ctx context.Context) (packages.RemotePackageIndex, error) {
	url := fmt.Sprintf("%spacakges/catalogue.json", r.rawURL)

	var cat Catalogue
	err := r.githubJSONResp(ctx, url, &cat)
	if err != nil {
		return nil, err
	}

	out := packages.BlankRemotePackageIndex()

	for _, pkg := range cat.Packages {
		pType := packages.PackageTypeIndex[pkg.Type]

		rURL := fmt.Sprintf("%spackages/%s/%s/pkg.json", r.rawURL, pType.StorageDir, pkg.ID)

		var remotePkg packages.RemotePackage
		err := r.githubJSONResp(ctx, rURL, &remotePkg)
		if err != nil {
			return nil, err
		}

		out[pkg.Type][pkg.ID] = remotePkg
	}

	return out, nil
}

func (r *RemoteRepository) GetByID(ctx context.Context, id packages.PkgID) (packages.RemotePackage, error) {
	url := fmt.Sprintf("%spackages/catalogue.json", r.rawURL)

	var cat Catalogue
	err := r.githubJSONResp(ctx, url, &cat)
	if err != nil {
		return packages.RemotePackage{}, err
	}

	for _, p := range cat.Packages {
		if p.ID == id {
			rURL := fmt.Sprintf("%spackages/%s/%s/pkg.json", r.rawURL, p.Type, p.ID)

			var pkg packages.RemotePackage
			err := r.githubJSONResp(ctx, rURL, &pkg)
			if err != nil {
				return packages.RemotePackage{}, err
			}

			return pkg, nil
		}
	}

	return packages.RemotePackage{}, fmt.Errorf("failed to find specified package")
}
