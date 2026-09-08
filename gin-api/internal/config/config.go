package config

import (
	"errors"
	"os"
)

type Config struct {
	Port               string
	DatabaseURL        string
	DatabaseCACertFile string
	FrontendURL        string
	SeedStaffName      string
	SeedStaffEmail     string
	SeedStaffPassword  string
}

func Load() (Config, error) {

	cfg := Config{
		Port:               valueOrDefault("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		DatabaseCACertFile: os.Getenv("DATABASE_CA_CERT_FILE"),
		FrontendURL:        valueOrDefault("FRONTEND_URL", "http://localhost:3000"),
		SeedStaffName:      os.Getenv("SEED_STAFF_NAME"),
		SeedStaffEmail:     os.Getenv("SEED_STAFF_EMAIL"),
		SeedStaffPassword:  os.Getenv("SEED_STAFF_PASSWORD"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	seedValues := []string{cfg.SeedStaffName, cfg.SeedStaffEmail, cfg.SeedStaffPassword}
	provided := 0
	for _, value := range seedValues {
		if value != "" {
			provided++
		}
	}
	if provided != 0 && provided != len(seedValues) {
		return Config{}, errors.New("SEED_STAFF_NAME, SEED_STAFF_EMAIL, and SEED_STAFF_PASSWORD must be set together")
	}
	if cfg.SeedStaffPassword != "" && len(cfg.SeedStaffPassword) < 12 {
		return Config{}, errors.New("SEED_STAFF_PASSWORD must contain at least 12 characters")
	}
	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
