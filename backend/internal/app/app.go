// File: app.go
// Purpose: Wire repository, services, and handlers into a single application container.
// Module: backend/internal/app, dependency assembly layer.
// Related: server bootstrap and route registration.
package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/cache"
	rediscache "github.com/anxcye/ancy-blog/backend/internal/cache/redis"
	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/handler"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
	"github.com/anxcye/ancy-blog/backend/internal/repository/postgres"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/anxcye/ancy-blog/backend/internal/storage"
)

type App struct {
	AuthHandler   *handler.AuthHandler
	PublicHandler *handler.PublicHandler
	AdminHandler  *handler.AdminHandler
	UploadHandler *handler.UploadHandler
	AuthService   *service.AuthService
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	authService := service.NewAuthService(
		cfg.Auth.AdminUsername,
		cfg.Auth.AdminPassword,
		time.Duration(cfg.Auth.AccessTokenTTLSeconds)*time.Second,
		time.Duration(cfg.Auth.RefreshTokenTTLSeconds)*time.Second,
	)

	pgRepo, err := postgres.New(context.Background(), cfg.DB)
	if err != nil {
		logger.Error("postgres repository init failed", "error", err)
		return nil, errors.New("postgres repository initialization failed")
	}
	logger.Info("postgres repository initialized", "host", cfg.DB.Host, "db", cfg.DB.Name)
	var repo repository.ContentRepository = pgRepo

	var cacheClient cache.Cache
	if cfg.Redis.Enabled {
		rc, err := rediscache.New(context.Background(), cfg.Redis)
		if err != nil {
			logger.Error("redis cache init failed, disabling cache", "error", err)
		} else {
			cacheClient = rc
			logger.Info("redis cache initialized", "addr", cfg.Redis.Addr, "db", cfg.Redis.DB)
		}
	}

	contentService := service.NewContentService(repo, cacheClient)

	var uploader storage.Uploader
	return &App{
		AuthHandler:   handler.NewAuthHandler(authService),
		PublicHandler: handler.NewPublicHandler(contentService),
		AdminHandler:  handler.NewAdminHandler(contentService),
		UploadHandler: handler.NewUploadHandler(uploader),
		AuthService:   authService,
	}, nil
}
