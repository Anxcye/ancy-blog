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
	"github.com/anxcye/ancy-blog/backend/internal/worker"
)

type App struct {
	AuthHandler       *handler.AuthHandler
	PublicHandler     *handler.PublicHandler
	AdminHandler      *handler.AdminHandler
	UploadHandler     *handler.UploadHandler
	AuthService       *service.AuthService
	TranslationWorker *worker.TranslationWorker
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

	// Wire credential store for bcrypt-based persistent password
	authService.WithCredentialStore(pgRepo)

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
	articleService := service.NewArticleService(contentService)
	commentService := service.NewCommentService(contentService)
	linkService := service.NewLinkService(contentService)
	siteService := service.NewSiteService(contentService)
	integrationService := service.NewIntegrationService(contentService)
	translationService := service.NewTranslationService(contentService)
	timelineService := service.NewTimelineService(contentService)
	aiAssistService := service.NewAIAssistService(articleService, integrationService)
	var translationWorker *worker.TranslationWorker
	if cfg.Translation.WorkerEnabled {
		translationWorker = worker.NewTranslationWorker(
			logger,
			translationService,
			integrationService,
			time.Duration(cfg.Translation.PollIntervalMS)*time.Millisecond,
			time.Duration(cfg.Translation.BackoffBaseMS)*time.Millisecond,
			time.Duration(cfg.Translation.BackoffMaxMS)*time.Millisecond,
		)
	}

	return &App{
		AuthHandler:       handler.NewAuthHandler(authService),
		PublicHandler:     handler.NewPublicHandler(articleService, commentService, linkService, siteService, timelineService),
		AdminHandler:      handler.NewAdminHandler(articleService, commentService, linkService, siteService, integrationService, translationService, aiAssistService, authService),
		UploadHandler:     handler.NewUploadHandler(nil, integrationService),
		AuthService:       authService,
		TranslationWorker: translationWorker,
	}, nil
}
