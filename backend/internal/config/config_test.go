// File: config_test.go
// Purpose: Verify environment configuration loading and parsing behavior.
// Module: backend/internal/config, config unit test layer.
// Related: config.go.
package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("AUTH_ACCESS_TOKEN_TTL_SECONDS", "")
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
