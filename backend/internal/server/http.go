// File: http.go
// Purpose: Construct and run the HTTP server, including route registration and graceful shutdown.
// Module: backend/internal/server, transport runtime layer.
// Related: cmd/server bootstrap, internal/handler endpoints, internal/middleware.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/handler"
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	cfg    *config.Config
	logger *slog.Logger
	server *http.Server
}

func NewHTTPServer(cfg *config.Config, logger *slog.Logger) *HTTPServer {
	if cfg.App.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(logger))
	engine.GET("/healthz", handler.Healthz)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &HTTPServer{
		cfg:    cfg,
		logger: logger,
		server: srv,
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

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
