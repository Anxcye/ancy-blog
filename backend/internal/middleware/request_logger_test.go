// File: request_logger_test.go
// Purpose: Verify request logger middleware executes downstream handlers.
// Module: backend/internal/middleware, middleware test layer.
// Related: request_logger.go.
package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggerPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	r := gin.New()
	r.Use(RequestLogger(logger))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "pong" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
