// Package downloader holds donwload functionality
package downloader

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/w1lam/mc-pacman/internal/core/events"
	"github.com/w1lam/mc-pacman/internal/usecases"
)

// Downloader is the downloader service
type Downloader struct {
	usecases.Base

	client      *http.Client
	maxParallel int
}

func New(base usecases.Base) *Downloader {
	d := Downloader{
		Base: base,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxParallel: 5,
	}

	return &d
}

type FileRequest struct {
	ID       string
	URL      string
	FileName string
	Size     int64
	Hash     string
}

type FileResult struct {
	ID       string
	FileName string
	Hash     string
}

func (d *Downloader) Download(
	ctx context.Context,
	destDir string,
	files []FileRequest,
) ([]FileResult, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []FileResult
		sem     = make(chan struct{}, d.maxParallel)
		errCh   = make(chan error, 1)
	)

	pOp, _ := events.OpFromCtx(ctx)
	op := d.StartOp(pOp, fmt.Sprintf("download_to_%s", filepath.Base(destDir)))
	d.EmitStart(op, fmt.Sprintf("starting download to %s", destDir))

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

			res, err := d.downloadOne(ctx, op, destDir, file)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return results, err
	default:
	}

	d.EmitComplete(op, fmt.Sprintf("download to %s complete", destDir))

	return results, nil
}

func (d *Downloader) downloadOne(
	ctx context.Context,
	parentOp events.Operation,
	destDir string,
	file FileRequest,
) (FileResult, error) {
	op := d.StartOp(parentOp, fmt.Sprintf("download_file:%s", file.FileName))
	d.EmitStart(op, fmt.Sprintf("starting download: %s", file.FileName))

	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		result, err := d.attemptDownload(ctx, op, destDir, file)

		if err == nil {
			d.EmitComplete(op, fmt.Sprintf("download complete: %s", file.FileName))
			return result, nil
		}

		lastErr = err

		if attempt < maxRetries-1 {
			d.EmitWarn(op, err, fmt.Sprintf("attempt %d/%d failed, retrying...", attempt, maxRetries))
			time.Sleep(time.Second * time.Duration(1<<attempt))
			continue
		}
	}

	d.EmitError(op, lastErr, fmt.Sprintf("download failed after %d attempts", maxRetries))
	return FileResult{}, lastErr
}

func (d *Downloader) attemptDownload(
	ctx context.Context,
	op events.Operation,
	destDir string,
	file FileRequest,
) (FileResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.URL, nil)
	if err != nil {
		return FileResult{}, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return FileResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FileResult{}, fmt.Errorf("bad_status: %s", resp.Status)
	}

	outPath := filepath.Join(destDir, file.FileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		d.EmitError(op, err, "failed to create outFile:"+outFile.Name())
		return FileResult{}, err
	}
	defer outFile.Close()

	hasher := sha512.New()
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
			d.EmitProgress(op, events.Progress{
				Label:      file.FileName,
				Current:    downloaded,
				Total:      file.Size,
				Percentage: float64(downloaded) / float64(file.Size) * 100,
			})
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return FileResult{}, err
		}
	}

	computed := hex.EncodeToString(hasher.Sum(nil))

	if file.Hash != "" && computed != file.Hash {
		os.Remove(outPath)

		return FileResult{}, fmt.Errorf("hash_mismatch")
	}

	return FileResult{
		ID:       file.ID,
		FileName: file.FileName,
		Hash:     computed,
	}, nil
}
