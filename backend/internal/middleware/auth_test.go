// File: auth_test.go
// Purpose: Verify bearer auth middleware behavior for unauthorized and authorized requests.
// Module: backend/internal/middleware, middleware unit test layer.
// Related: middleware/auth.go and service/auth.go.
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func TestAuthRequiredMissingBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := service.NewAuthService("admin", "secret", time.Hour, 24*time.Hour)
	engine := gin.New()
	engine.GET("/protected", AuthRequired(authService), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	var body response.Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if body.Code != "AUTH_UNAUTHORIZED" {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

func TestAuthRequiredValidBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := service.NewAuthService("admin", "secret", time.Hour, 24*time.Hour)
	loginResult, err := authService.Login("admin", "secret")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}

	engine := gin.New()
	engine.GET("/protected", AuthRequired(authService), func(c *gin.Context) {
		user := MustUser(c)
		c.JSON(http.StatusOK, gin.H{"username": user.Username})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+loginResult.AccessToken)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if payload["username"] != "admin" {
		t.Fatalf("unexpected user in context: %#v", payload)
	}
}
