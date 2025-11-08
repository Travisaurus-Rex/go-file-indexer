package config

import "os"

type Config struct {
	Port     string
	ScanPath string
	LogLevel string
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

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	return Config{
		Port:     port,
		ScanPath: path,
		LogLevel: level,
	}
}
