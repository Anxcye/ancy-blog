// File: main.go
// Purpose: Bootstrap the API service process and wire application dependencies.
// Module: backend/cmd/server, executable entrypoint layer.
// Related: internal/config, internal/logger, internal/server.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/logger"
	"github.com/anxcye/ancy-blog/backend/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	lg := logger.New(cfg.App.Env)
	lg.Info("service booting", "env", cfg.App.Env, "port", cfg.HTTP.Port)

	srv := server.NewHTTPServer(cfg, lg)
	if err := srv.Start(ctx); err != nil {
		lg.Error("service exited with error", "error", err)
		log.Fatalf("service exited: %v", err)
	}
}
