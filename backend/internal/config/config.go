package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	APIPort     string
}

func Load() Config {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		APIPort:     port,
	}
}
