package scanner

import (
	"mime"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileInfo struct {
	Path      string    `json:"path"`
	SizeBytes int64     `json:"size_bytes"`
	MimeType  string    `json:"mime_type"`
	Modified  time.Time `json:"modified"`
}

func ScanDir(root string) ([]FileInfo, error) {
	const numWorkers = 4
	var paths = make(chan string, 100)
	var files []FileInfo
	var wg sync.WaitGroup
	var mu sync.Mutex
	var walkErr error

	go func() {
		defer close(paths)
		walkErr = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				return nil
			}

			paths <- path
			return nil
		})

	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for path := range paths {

				info, err := os.Stat(path)
				if err != nil {
					continue
				}

				ext := filepath.Ext(path)
				mimeType := mime.TypeByExtension(ext)

				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				file := FileInfo{
					Path:      path,
					SizeBytes: info.Size(),
					MimeType:  mimeType,
					Modified:  info.ModTime(),
				}
				mu.Lock()
				files = append(files, file)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if walkErr != nil {
		return nil, walkErr
	}

	return files, nil
}
