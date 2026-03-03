// File: health_test.go
// Purpose: Verify health endpoint returns expected payload.
// Module: backend/internal/handler, health HTTP test layer.
// Related: health.go and response envelope.
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/gin-gonic/gin"
)

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", Healthz)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
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
