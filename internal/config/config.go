package config

import "os"

type Config struct {
	Port     string
	ScanPath string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	path := os.Getenv("SCAN_PATH")
	if path == "" {
		path = "./data"
	}

	return Config{
		Port:     port,
		ScanPath: path,
	}
}
