// Package downloader holds donwload functionality
package downloader

import (
	"context"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/mc-pacman/internal/ux/progress"
)

// Service is the downloader service
type Service struct {
	client      *http.Client
	maxParallel int
	emitter     progress.ProgressEmitter
}

func New() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxParallel: 5,
	}
}

type FileRequest struct {
	ID       modrinth.ID
	URL      string
	FileName string
	Size     int64
	Hash     string
	Algo     string
}

type FileResult struct {
	ID       modrinth.ID
	FileName string
	Hash     string
}

func (s *Service) DownloadBatch(
	ctx context.Context,
	destDir string,
	files []FileRequest,
) ([]FileResult, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []FileResult
		sem     = make(chan struct{}, s.maxParallel)
		errCh   = make(chan error, 1)
	)

	for _, file := range files {
		file := file
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			if s.emitter != nil {
				s.emitter.Emit(progress.ProgressEvent{
					Type:       progress.ProgressStart,
					PackageID:  string(file.ID),
					FileName:   file.FileName,
					Percentage: 0,
				})
			}

			res, err := s.downloadOne(ctx, destDir, file)
			if err != nil {
				if s.emitter != nil {
					s.emitter.Emit(progress.ProgressEvent{
						Type:      progress.ProgressFailure,
						PackageID: string(file.ID),
						FileName:  file.FileName,
						Error:     err,
					})
				}

				select {
				case errCh <- err:
				default:
				}
				return
			}

			mu.Lock()
			results = append(results, res)
			mu.Unlock()

			if s.emitter != nil {
				s.emitter.Emit(progress.ProgressEvent{
					Type:       progress.ProgressSuccess,
					PackageID:  string(file.ID),
					FileName:   file.FileName,
					Percentage: 100,
				})
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return results, err
	default:
	}

	if s.emitter != nil {
		s.emitter.Emit(progress.ProgressEvent{
			Type: progress.ProgressComplete,
		})
	}

	return results, nil
}

func (s *Service) downloadOne(
	ctx context.Context,
	destDir string,
	file FileRequest,
) (FileResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", file.URL, nil)
	if err != nil {
		return FileResult{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return FileResult{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FileResult{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	outPath := filepath.Join(destDir, file.FileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return FileResult{}, nil
	}
	defer outFile.Close()

	var hasher hash.Hash

	switch file.Algo {
	case "sha512":
		hasher = sha512.New()
	case "sha1":
		hasher = sha1.New()
	default:
		return FileResult{}, fmt.Errorf("unsupported hash algo: %s", file.Algo)
	}

	writer := io.MultiWriter(outFile, hasher)

	buf := make([]byte, 32*1024)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := writer.Write(buf[:n]); werr != nil {
				return FileResult{}, werr
			}

			downloaded += int64(n)
			if s.emitter != nil && file.Size > 0 {
				s.emitter.Emit(progress.ProgressEvent{
					Type:       progress.ProgressUpdate,
					PackageID:  string(file.ID),
					FileName:   file.FileName,
					Percentage: float64(downloaded) / float64(file.Size) * 100,
				})
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return FileResult{}, err
		}
	}

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return FileResult{}, err
	}

	computed := hex.EncodeToString(hasher.Sum(nil))

	if file.Hash != "" && computed != file.Hash {
		return FileResult{}, fmt.Errorf("hash mismatch for %s", file.ID)
	}

	return FileResult{
		ID:       file.ID,
		FileName: file.FileName,
		Hash:     computed,
	}, nil
}

func (s *Service) SetEmitter(e progress.ProgressEmitter) {
	s.emitter = e
}
