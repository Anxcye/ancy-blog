// File: config_test.go
// Purpose: Verify environment configuration loading and parsing behavior.
// Module: backend/internal/config, config unit test layer.
// Related: config.go.
package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("AUTH_ACCESS_TOKEN_TTL_SECONDS", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected load success, got error: %v", err)
	}
	if cfg.HTTP.Port != 8080 {
		t.Fatalf("unexpected default port: %d", cfg.HTTP.Port)
	}
	if cfg.Auth.AdminUsername == "" {
		t.Fatalf("expected default admin username")
	}
	if len(cfg.CORS.AllowedOrigins) == 0 {
		t.Fatalf("expected default CORS allowed origins")
	}
}

func TestLoadInvalidHTTPPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "bad")
	if _, err := Load(); err == nil {
		t.Fatalf("expected invalid HTTP_PORT error")
	}
}

func TestParseBool(t *testing.T) {
	if !parseBool("yes") {
		t.Fatalf("expected yes to be true")
	}
	if parseBool("no") {
		t.Fatalf("expected no to be false")
	}
}

func TestParseCSV(t *testing.T) {
	values := parseCSV(" http://a.com, ,http://b.com ")
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}
	if values[0] != "http://a.com" || values[1] != "http://b.com" {
		t.Fatalf("unexpected parse result: %#v", values)
	}
}
