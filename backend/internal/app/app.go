// File: app.go
// Purpose: Wire repository, services, and handlers into a single application container.
// Module: backend/internal/app, dependency assembly layer.
// Related: server bootstrap and route registration.
package app

import (
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/handler"
	"github.com/anxcye/ancy-blog/backend/internal/repository/memory"
	"github.com/anxcye/ancy-blog/backend/internal/service"
)

type App struct {
	AuthHandler   *handler.AuthHandler
	PublicHandler *handler.PublicHandler
	AdminHandler  *handler.AdminHandler
	AuthService   *service.AuthService
}

func New(cfg *config.Config) *App {
	repo := memory.NewRepository()
	authService := service.NewAuthService(
		cfg.Auth.AdminUsername,
		cfg.Auth.AdminPassword,
		time.Duration(cfg.Auth.AccessTokenTTLSeconds)*time.Second,
		time.Duration(cfg.Auth.RefreshTokenTTLSeconds)*time.Second,
	)
	contentService := service.NewContentService(repo)

	return &App{
		AuthHandler:   handler.NewAuthHandler(authService),
		PublicHandler: handler.NewPublicHandler(contentService),
		AdminHandler:  handler.NewAdminHandler(contentService),
		AuthService:   authService,
	}
}
