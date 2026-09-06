package config

import "testing"

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected missing DATABASE_URL to return an error")
	}
}

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("PORT", "")
	t.Setenv("FRONTEND_URL", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if cfg.Port != "8080" || cfg.FrontendURL != "http://localhost:3000" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}
