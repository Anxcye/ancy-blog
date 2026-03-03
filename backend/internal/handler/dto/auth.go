// File: auth.go
// Purpose: Define auth HTTP request DTOs to isolate transport payloads from domain models.
// Module: backend/internal/handler/dto, auth transport DTO layer.
// Related: handler/auth.go and service/auth.go.
package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
