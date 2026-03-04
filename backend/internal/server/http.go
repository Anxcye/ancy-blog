// File: http.go
// Purpose: Construct and run the HTTP server, including route registration and graceful shutdown.
// Module: backend/internal/server, transport runtime layer.
// Related: cmd/server bootstrap, internal/app container, and middleware chain.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/app"
	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/handler"
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	cfg               *config.Config
	logger            *slog.Logger
	server            *http.Server
	translationWorker interface{ Run(context.Context) }
}

func NewHTTPServer(cfg *config.Config, logger *slog.Logger) (*HTTPServer, error) {
	if cfg.App.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	container, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("app initialization failed", "error", err)
		return nil, err
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(logger))
	engine.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	engine.GET("/healthz", handler.Healthz)

	api := engine.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", container.AuthHandler.Login)
			auth.POST("/refresh", container.AuthHandler.Refresh)
			auth.GET("/me", middleware.AuthRequired(container.AuthService), container.AuthHandler.Me)
		}

		pub := api.Group("/public")
		{
			pub.GET("/articles", container.PublicHandler.Articles)
			pub.GET("/articles/:slug", container.PublicHandler.ArticleBySlug)
			pub.GET("/articles/by-category/:categorySlug", container.PublicHandler.ArticleByCategory)
			pub.GET("/moments", container.PublicHandler.Moments)
			pub.GET("/comments/article/:articleId", container.PublicHandler.CommentByArticle)
			pub.GET("/comments/:id/children", container.PublicHandler.CommentChildren)
			pub.GET("/comments/article/:articleId/total", container.PublicHandler.CommentArticleTotal)
			pub.POST("/comments", container.PublicHandler.AddComment)
			pub.POST("/links/submissions", container.PublicHandler.SubmitLink)
			pub.GET("/links", container.PublicHandler.ApprovedLinks)
			pub.GET("/categories", container.PublicHandler.Categories)
			pub.GET("/tags", container.PublicHandler.Tags)
			pub.GET("/site/settings", container.PublicHandler.SiteSettings)
			pub.GET("/site/footer", container.PublicHandler.SiteFooter)
			pub.GET("/site/social-links", container.PublicHandler.SiteSocialLinks)
			pub.GET("/site/nav", container.PublicHandler.SiteNav)
			pub.GET("/site/slots/:slotKey", container.PublicHandler.SiteSlotContent)
			pub.GET("/timeline", container.PublicHandler.Timeline)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(container.AuthService))
		{
			admin.GET("/articles", container.AdminHandler.ListArticles)
			admin.GET("/articles/:id", container.AdminHandler.ArticleDetail)
			admin.DELETE("/articles/:id", container.AdminHandler.DeleteArticle)
			admin.POST("/articles/batch-status", container.AdminHandler.BatchUpdateArticleStatus)
			admin.POST("/articles/batch-delete", container.AdminHandler.BatchDeleteArticle)
			admin.POST("/articles", container.AdminHandler.CreateArticle)
			admin.PUT("/articles/:id", container.AdminHandler.UpdateArticle)
			admin.GET("/moments", container.AdminHandler.ListMoments)
			admin.POST("/moments", container.AdminHandler.CreateMoment)
			admin.PUT("/moments/:id", container.AdminHandler.UpdateMoment)
			admin.DELETE("/moments/:id", container.AdminHandler.DeleteMoment)
			admin.POST("/moments/batch-status", container.AdminHandler.BatchUpdateMomentStatus)
			admin.POST("/moments/batch-delete", container.AdminHandler.BatchDeleteMoment)
			admin.GET("/comments", container.AdminHandler.CommentPage)
			admin.PUT("/comments/:id", container.AdminHandler.CommentUpdate)

			admin.GET("/links", container.AdminHandler.ListLinkSubmissions)
			admin.PATCH("/links/:id/review", container.AdminHandler.ReviewLink)

			admin.PUT("/site/settings", container.AdminHandler.UpdateSiteSettings)
			admin.POST("/site/footer-items", container.AdminHandler.CreateFooterItem)
			admin.PUT("/site/footer-items/:id", container.AdminHandler.UpdateFooterItem)
			admin.DELETE("/site/footer-items/:id", container.AdminHandler.DeleteFooterItem)

			admin.POST("/site/social-links", container.AdminHandler.CreateSocialLink)
			admin.PUT("/site/social-links/:id", container.AdminHandler.UpdateSocialLink)
			admin.DELETE("/site/social-links/:id", container.AdminHandler.DeleteSocialLink)

			admin.POST("/site/nav-items", container.AdminHandler.CreateNavItem)
			admin.PUT("/site/nav-items/:id", container.AdminHandler.UpdateNavItem)
			admin.DELETE("/site/nav-items/:id", container.AdminHandler.DeleteNavItem)

			admin.GET("/site/slots", container.AdminHandler.ListSlots)
			admin.POST("/site/slots", container.AdminHandler.CreateSlot)
			admin.GET("/site/slots/:slotKey/items", container.AdminHandler.ListSlotItems)
			admin.POST("/site/slots/:slotKey/items", container.AdminHandler.CreateSlotItem)
			admin.DELETE("/site/slots/:slotKey/items/:id", container.AdminHandler.DeleteSlotItem)

			admin.GET("/integrations", container.AdminHandler.ListIntegrations)
			admin.PUT("/integrations/:providerKey", container.AdminHandler.UpdateIntegration)
			admin.POST("/integrations/:providerKey/test", container.AdminHandler.TestIntegration)

			admin.POST("/translations/jobs", container.AdminHandler.CreateTranslationJob)
			admin.GET("/translations/jobs", container.AdminHandler.ListTranslationJobs)
			admin.GET("/translations/jobs/:id", container.AdminHandler.TranslationJobDetail)
			admin.POST("/translations/jobs/:id/retry", container.AdminHandler.RetryTranslationJob)
			admin.GET("/translations/contents", container.AdminHandler.ListTranslationContents)
			admin.GET("/translations/contents/:sourceType/:sourceId/:locale", container.AdminHandler.TranslationContentDetail)
			admin.PUT("/translations/contents", container.AdminHandler.UpsertTranslationContent)
			admin.POST("/ai/summary", container.AdminHandler.GenerateSummary)
			admin.POST("/ai/slug", container.AdminHandler.SuggestSlug)

			admin.POST("/upload/image", container.UploadHandler.UploadImage)
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &HTTPServer{
		cfg:               cfg,
		logger:            logger,
		server:            srv,
		translationWorker: container.TranslationWorker,
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	if s.translationWorker != nil {
		go s.translationWorker.Run(ctx)
	}

	go func() {
		s.logger.Info("http server started", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
