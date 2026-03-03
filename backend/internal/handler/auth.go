// File: auth.go
// Purpose: Expose authentication endpoints for login, token refresh, and profile lookup.
// Module: backend/internal/handler, auth HTTP presentation layer.
// Related: service.AuthService and auth middleware.
package handler

import (
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	res, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.JSON(c, 401, response.Envelope{Code: "AUTH_INVALID_CREDENTIALS", Message: "username or password is incorrect"})
		return
	}
	response.JSON(c, 200, response.Envelope{Code: "OK", Message: "success", Data: res})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	res, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		response.JSON(c, 401, response.Envelope{Code: "AUTH_REFRESH_INVALID", Message: "refresh token is invalid"})
		return
	}
	response.JSON(c, 200, response.Envelope{Code: "OK", Message: "success", Data: res})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user := middleware.MustUser(c)
	response.JSON(c, 200, response.Envelope{Code: "OK", Message: "success", Data: user})
}
