// File: common.go
// Purpose: Provide shared handler helpers for pagination and consistent responses.
// Module: backend/internal/handler, HTTP presentation utility layer.
// Related: resource handlers and response envelope utilities.
package handler

import (
	"strconv"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/gin-gonic/gin"
)

type pageResult[T any] struct {
	Total int `json:"total"`
	Rows  []T `json:"rows"`
}

func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	raw := c.Query(key)
	if raw == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	return v
}

func badRequest(c *gin.Context, code, message string) {
	response.JSON(c, 400, response.Envelope{Code: code, Message: message})
}
