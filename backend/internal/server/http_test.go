// File: http_test.go
// Purpose: Verify server construction fails fast when app dependencies cannot initialize.
// Module: backend/internal/server, server bootstrap test layer.
// Related: http.go and app.New.
package server

import (
	"io"
	"log/slog"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/config"
)

func TestNewHTTPServerFailsWhenAppInitFails(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	cfg := &config.Config{
		App:  config.AppConfig{Name: "test", Env: "test"},
		HTTP: config.HTTPConfig{Host: "127.0.0.1", Port: 8080},
		Auth: config.AuthConfig{AdminUsername: "admin", AdminPassword: "secret", AccessTokenTTLSeconds: 3600, RefreshTokenTTLSeconds: 7200},
		DB: config.DBConfig{
			Host:         "127.0.0.1",
			Port:         1,
			Name:         "no_db",
			User:         "no_user",
			Password:     "no_pass",
			SSLMode:      "disable",
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		},
	}

	if _, err := NewHTTPServer(cfg, logger); err == nil {
		t.Fatalf("expected server initialization error when app init fails")
	}
}
