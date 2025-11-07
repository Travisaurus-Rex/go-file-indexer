package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Travisaurus-Rex/go-file-indexer/internal/config"
)

type FileInfo struct {
	Path      string `json:"Path"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
	Modified  string `json:"modified"`
}

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("configuration loaded", "port", cfg.Port, "scan_path", cfg.ScanPath)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to encode config"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
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
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		logger.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	stop()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Error("shutdown failed", "err", err)
	} else {
		logger.Info("server stopped cleanly")
	}

	fmt.Println("Selamat jalan, mbokne amput!")
}
