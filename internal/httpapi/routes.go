package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Travisaurus-Rex/go-file-indexer/internal/config"
)

type FileInfo struct {
	Path      string `json:"Path"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
	Modified  string `json:"modified"`
}

func RegisterRoutes(mux *http.ServeMux, cfg config.Config) {
	mux.HandleFunc("/healthz", handleHealthz)
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		handleConfig(w, r, cfg)
	})
	mux.HandleFunc("/files", handleFiles)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func handleConfig(w http.ResponseWriter, r *http.Request, cfg config.Config) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to encode config"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files := []FileInfo{
		{
			Path:      "./data/example1.txt",
			SizeBytes: 12345,
			MimeType:  "text/plain",
			Modified:  "2025-11-07T13:00:00z",
		},
		{
			Path:      "./data/image.jpg",
			SizeBytes: 987654,
			MimeType:  "image/jpeg",
			Modified:  "2025-11-07T12:45:00Z",
		},
		{
			Path:      "./data/archive.zip",
			SizeBytes: 5000000,
			MimeType:  "application/zip",
			Modified:  "2025-11-06T09:30:00Z",
		},
	}

	data, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to encode files"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
