// File: logger_test.go
// Purpose: Verify logger construction for different environments.
// Module: backend/internal/logger, logger unit test layer.
// Related: logger.go.
package logger

import "testing"

func TestNewLogger(t *testing.T) {
	if New("dev") == nil {
		t.Fatalf("expected logger for dev env")
	}
	if New("prod") == nil {
		t.Fatalf("expected logger for prod env")
	}
}
