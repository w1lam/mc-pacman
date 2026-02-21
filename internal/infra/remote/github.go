package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type githubContentResponse struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sha    string `json:"sha"`
	Size   int    `json:"size"`
	RawURL string `json:"download_url"`
	Type   string `json:"type"`
}

func (r *RemoteRepository) githubJSONResp(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "MyModInstaller/1.0")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github request failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
