package dev

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )
//
// type githubContentResponse struct {
// 	Name   string `json:"name"`
// 	Path   string `json:"path"`
// 	Sha    string `json:"sha"`
// 	Size   int    `json:"size"`
// 	RawURL string `json:"download_url"`
// 	Type   string `json:"type"`
// }
//
// func GetPackagesFromFolder(folder string) ([]packages.RemotePackage, error) {
// 	var items []githubContentResponse
//
// 	url := fmt.Sprintf("%scontents/packages/%s", netcfg.GithubRepo, folder)
// 	if err := GithubGetJSON(url, &items); err != nil {
// 		return nil, err
// 	}
//
// 	var result []packages.RemotePackage
// 	for _, it := range items {
// 		if it.Type != "file" {
// 			continue
// 		}
//
// 		pkg, err := ResolvePackageJSON(it.RawURL)
// 		if err != nil {
// 			return nil, fmt.Errorf("folder %s: %w", folder, err)
// 		}
//
// 		result = append(result, pkg)
// 	}
// 	return result, nil
// }
//
// func (r *RemoteRepo) GithubGetJSON(url string, out any) error {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	req.Header.Set("Accept", "application/vnd.github+json")
// 	req.Header.Set("Cache-Control", "no-cache")
// 	req.Header.Set("User-Agent", "MyModInstaller/1.0")
//
// 	resp, err := r.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("github request failed: %s", resp.Status)
// 	}
//
// 	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func ScanPackagesFolder() ([]string, error) {
// 	var items []githubContentResponse
//
// 	url := netcfg.GithubRepo + "contents/packages"
// 	if err := GithubGetJSON(url, &items); err != nil {
// 		return nil, err
// 	}
//
// 	var folders []string
// 	for _, it := range items {
// 		if it.Type == "dir" {
// 			folders = append(folders, it.Name)
// 		}
// 	}
//
// 	if len(folders) == 0 {
// 		return nil, fmt.Errorf("no package folder found")
// 	}
//
// 	return folders, nil
// }
//
// func ResolvePackageJSON(url string) (packages.RemotePackage, error) {
// 	var pkg packages.RemotePackage
//
// 	if err := GithubGetJSON(url, &pkg); err != nil {
// 		return pkg, err
// 	}
//
// 	if pkg.Name == "" {
// 		return pkg, fmt.Errorf("package json missing name")
// 	}
//
// 	if pkg.ListSource == "" {
// 		pkg.ListSource = url
// 	}
//
// 	return pkg, nil
// }
