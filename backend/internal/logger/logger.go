// File: logger.go
// Purpose: Build a shared structured logger for all application layers.
// Module: backend/internal/logger, infrastructure logging layer.
// Related: cmd/server startup, internal/server runtime logs.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(env string) *slog.Logger {
	level := slog.LevelInfo
	if strings.ToLower(env) == "dev" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
