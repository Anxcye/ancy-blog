// File: response.go
// Purpose: Provide a common HTTP response envelope for API handlers.
// Module: backend/internal/response, transport utility layer.
// Related: internal/handler implementations and API contract definitions.
package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(c *gin.Context, status int, payload Envelope) {
	c.JSON(status, payload)
}
