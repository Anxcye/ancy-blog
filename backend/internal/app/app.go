// File: app.go
// Purpose: Wire repository, services, and handlers into a single application container.
// Module: backend/internal/app, dependency assembly layer.
// Related: server bootstrap and route registration.
package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/handler"
	"github.com/anxcye/ancy-blog/backend/internal/repository/memory"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/anxcye/ancy-blog/backend/internal/storage"
	"github.com/anxcye/ancy-blog/backend/internal/storage/r2"
)

type App struct {
	AuthHandler   *handler.AuthHandler
	PublicHandler *handler.PublicHandler
	AdminHandler  *handler.AdminHandler
	UploadHandler *handler.UploadHandler
	AuthService   *service.AuthService
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	repo := memory.NewRepository()
	authService := service.NewAuthService(
		cfg.Auth.AdminUsername,
		cfg.Auth.AdminPassword,
		time.Duration(cfg.Auth.AccessTokenTTLSeconds)*time.Second,
		time.Duration(cfg.Auth.RefreshTokenTTLSeconds)*time.Second,
	)
	contentService := service.NewContentService(repo)

	var uploader storage.Uploader
	if cfg.R2.Enabled {
		r2Uploader, err := r2.New(context.Background(), cfg.R2)
		if err != nil {
			logger.Error("r2 uploader init failed", "error", err)
		} else {
			uploader = r2Uploader
			logger.Info("r2 uploader initialized", "bucket", cfg.R2.Bucket)
		}
	}

	return &App{
		AuthHandler:   handler.NewAuthHandler(authService),
		PublicHandler: handler.NewPublicHandler(contentService),
		AdminHandler:  handler.NewAdminHandler(contentService),
		UploadHandler: handler.NewUploadHandler(uploader),
		AuthService:   authService,
	}
}
