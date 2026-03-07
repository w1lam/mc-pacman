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
	"github.com/w1lam/mc-pacman/internal/ux"
)

// Downloader is the downloader service
type Downloader struct {
	events.EmitterBase
	client      *http.Client
	maxParallel int
}

func New(view ux.View) *Downloader {
	d := Downloader{
		EmitterBase: events.EmitterBase{
			Scope: events.ScopeDownloader,
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxParallel: 5,
	}
	d.SetEmitter(view)
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
	defer d.EmitEnd(op)

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
	op := d.StartOp(parentOp, fmt.Sprintf("download_file_%s", file.FileName))
	d.EmitStart(op, fmt.Sprintf("starting download: %s", file.FileName))
	defer d.EmitEnd(op)

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
		return FileResult{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	outPath := filepath.Join(destDir, file.FileName)
	outFile, err := os.Create(outPath)
	if err != nil {
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
			d.Emit(events.Event{
				Type:       events.EventDownload,
				Op:         op,
				FileName:   file.FileName,
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
		d.EmitError(op, fmt.Errorf("hash mismath for %s", file.ID))
		return FileResult{}, fmt.Errorf("hash mismatch for %s", file.ID)
	}

	d.EmitComplete(op, fmt.Sprintf("download complete: %s", file.FileName))

	return FileResult{
		ID:       file.ID,
		FileName: file.FileName,
		Hash:     computed,
	}, nil
}
