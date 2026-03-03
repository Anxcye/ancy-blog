// File: app_test.go
// Purpose: Verify app container initialization failure behavior when mandatory dependencies are unavailable.
// Module: backend/internal/app, dependency assembly test layer.
// Related: app.go and postgres repository bootstrap.
package app

import (
	"io"
	"log/slog"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/config"
)

func TestNewAppFailsWhenPostgresUnavailable(t *testing.T) {
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

	if _, err := New(cfg, logger); err == nil {
		t.Fatalf("expected app initialization error when postgres is unavailable")
	}
}
