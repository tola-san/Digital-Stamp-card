package config

import "testing"

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected missing DATABASE_URL to return an error")
	}
}

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "mysql://example")
	t.Setenv("PORT", "")
	t.Setenv("FRONTEND_URL", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if cfg.Port != "8080" || cfg.FrontendURL != "http://localhost:3000" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.CookieSecure {
		t.Fatal("expected local cookies to default to insecure HTTP mode")
	}
}

func TestLoadReadsSecureCookieSetting(t *testing.T) {
	t.Setenv("DATABASE_URL", "mysql://example")
	t.Setenv("COOKIE_SECURE", "true")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if !cfg.CookieSecure {
		t.Fatal("expected COOKIE_SECURE=true to enable secure cookies")
	}
}
