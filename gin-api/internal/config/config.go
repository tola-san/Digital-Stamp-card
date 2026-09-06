package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
}

func Load() (Config, error) {

	cfg := Config{
		Port:        valueOrDefault("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		FrontendURL: valueOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
