// File: auth_test.go
// Purpose: Verify auth handler request validation and success/error responses.
// Module: backend/internal/handler, auth HTTP test layer.
// Related: auth.go and service.AuthService.
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func TestAuthHandlerLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(service.NewAuthService("admin", "secret", time.Hour, 24*time.Hour))
	r := gin.New()
	r.POST("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"admin","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if env.Code != "OK" {
		t.Fatalf("unexpected code: %s", env.Code)
	}
}

func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(service.NewAuthService("admin", "secret", time.Hour, 24*time.Hour))
	r := gin.New()
	r.POST("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"admin","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandlerRefreshInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(service.NewAuthService("admin", "secret", time.Hour, 24*time.Hour))
	r := gin.New()
	r.POST("/refresh", h.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
