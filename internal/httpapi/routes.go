package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Travisaurus-Rex/go-file-indexer/internal/config"
	"github.com/Travisaurus-Rex/go-file-indexer/internal/scanner"
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

	root := "./data"
	files, err := scanner.ScanDir(root)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to scan directory"}`))
		return
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
