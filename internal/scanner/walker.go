package scanner

import (
	"fmt"
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
		walkErr = walkDirectory(root, paths)
	}()

	for i := range numWorkers {
		wg.Add(1)
		go worker(i, paths, &files, &mu, &wg)
	}

	wg.Wait()

	if walkErr != nil {
		return nil, walkErr
	}

	return files, nil
}

func walkDirectory(root string, paths chan<- string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		paths <- path
		return nil
	})
}

func worker(id int, paths <-chan string, files *[]FileInfo, mu *sync.Mutex, wg *sync.WaitGroup) {
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
		*files = append(*files, file)
		mu.Unlock()

		fmt.Printf("[worker %d] processed %s\n", id, filepath.Base(path))
	}

}
