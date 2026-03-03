// File: response_test.go
// Purpose: Verify response envelope helper writes JSON responses.
// Module: backend/internal/response, transport utility test layer.
// Related: response.go.
package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", func(c *gin.Context) {
		JSON(c, http.StatusCreated, Envelope{Code: "OK", Message: "created", Data: map[string]string{"id": "1"}})
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
