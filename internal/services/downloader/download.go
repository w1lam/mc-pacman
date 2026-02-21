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
	"github.com/w1lam/mc-pacman/internal/events"
	"github.com/w1lam/mc-pacman/internal/services"
)

// Service is the downloader service
type Downloader struct {
	services.Base
	client      *http.Client
	maxParallel int
}

func New() *Downloader {
	return &Downloader{
		Base: services.Base{
			Scope:   events.ScopeDownloader,
			Emitter: nil,
		},
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

func (s *Downloader) DownloadBatch(
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

	op := events.NewOperation(s.Scope, destDir)

	s.Emit(events.Event{
		Type:    events.EventStart,
		Op:      op,
		Message: fmt.Sprintf("[%s] starting download", s.Scope),
	})

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

			s.Emit(events.Event{
				Type:       events.EventUpdate,
				Op:         op,
				FileName:   file.FileName,
				Percentage: 0,
			})

			res, err := s.downloadOne(ctx, op, destDir, file)
			if err != nil {
				s.Emit(events.Event{
					Type:     events.EventFailure,
					Op:       op,
					FileName: file.FileName,
					Error:    err,
				})

				select {
				case errCh <- err:
				default:
				}
				return
			}

			mu.Lock()
			results = append(results, res)
			mu.Unlock()

			s.Emit(events.Event{
				Type:       events.EventSuccess,
				Op:         op,
				FileName:   file.FileName,
				Percentage: 100,
			})
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return results, err
	default:
	}

	s.Emit(events.Event{
		Type:       events.EventComplete,
		Op:         op,
		FileName:   destDir,
		Percentage: 100,
	})

	return results, nil
}

func (s *Downloader) downloadOne(
	ctx context.Context,
	op events.Operation,
	destDir string,
	file FileRequest,
) (FileResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.URL, nil)
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
			s.Emit(events.Event{
				Type:       events.EventUpdate,
				Op:         op,
				FileName:   file.FileName,
				Percentage: float64(downloaded) / float64(file.Size) * 100,
				Message:    "downloading",
			})
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
		s.Emit(events.Event{
			Type:    events.EventFailure,
			Message: fmt.Sprintf("[%s] hash mismath for %s", s.Scope, file.ID),
			Error:   fmt.Errorf("hash mismatch for %s", file.ID),
		})
		return FileResult{}, fmt.Errorf("hash mismatch for %s", file.ID)
	}

	s.Emit(events.Event{
		Type: events.EventComplete,
		Op:   op,
	})

	return FileResult{
		ID:       file.ID,
		FileName: file.FileName,
		Hash:     computed,
	}, nil
}
