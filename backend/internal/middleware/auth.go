// File: auth.go
// Purpose: Enforce bearer-token authentication for protected API groups.
// Module: backend/internal/middleware, HTTP cross-cutting layer.
// Related: auth service and admin route groups.
package middleware

import (
	"strings"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

const ContextUserKey = "auth_user"

func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.JSON(c, 401, response.Envelope{Code: "AUTH_UNAUTHORIZED", Message: "missing bearer token"})
			c.Abort()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		user, err := authService.ResolveUser(token)
		if err != nil {
			response.JSON(c, 401, response.Envelope{Code: "AUTH_UNAUTHORIZED", Message: "unauthorized"})
			c.Abort()
			return
		}
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

func MustUser(c *gin.Context) domain.User {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return domain.User{}
	}
	user, ok := v.(domain.User)
	if !ok {
		return domain.User{}
	}
	return user
}
