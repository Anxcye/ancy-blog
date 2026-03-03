// File: health.go
// Purpose: Serve health-check endpoints used by local development and deployment probes.
// Module: backend/internal/handler, HTTP presentation layer.
// Related: internal/server route registration and deploy health checks.
package handler

import (
	"net/http"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{
		Code:    "OK",
		Message: "service is healthy",
		Data: map[string]any{
			"status": "up",
		},
	})
}
