package scanner

import (
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Path      string    `json:"path"`
	SizeBytes int64     `json:"size_bytes"`
	Modified  time.Time `json:"modified"`
}

func ScanDir(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		files = append(files, FileInfo{
			Path:      path,
			SizeBytes: info.Size(),
			Modified:  info.ModTime(),
		})

		return nil
	})

	return files, err
}
