// File: main.go
// Purpose: Execute database migrations in a controlled, versioned workflow.
// Module: backend/cmd/migrate, migration CLI entrypoint layer.
// Related: backend/migrations and runtime DB configuration.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		action  string
		steps   int
		version int
		path    string
		verbose bool
	)

	flag.StringVar(&action, "action", "up", "migration action: up | down | steps | version | force")
	flag.IntVar(&steps, "steps", 1, "number of steps for action=steps (negative rolls back)")
	flag.IntVar(&version, "version", 0, "target version for action=force")
	flag.StringVar(&path, "path", "file://migrations", "migration source path")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)
	m, err := migrate.New(path, dsn)
	if err != nil {
		log.Fatalf("init migrate failed: %v", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if verbose {
		log.Printf("migration action=%s source=%s db=%s:%d/%s", action, path, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	}

	switch action {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up failed: %v", err)
		}
	case "down":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate down failed: %v", err)
		}
	case "steps":
		err = m.Steps(steps)
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate steps failed: %v", err)
		}
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("version: nil")
				return
			}
			log.Fatalf("get version failed: %v", err)
		}
		fmt.Printf("version: %d dirty: %v\n", v, dirty)
		return
	case "force":
		if version <= 0 {
			log.Fatalf("force action requires -version > 0")
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("migrate force failed: %v", err)
		}
	default:
		log.Fatalf("unsupported action: %s", action)
	}

	v, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("migration complete; version: nil")
			return
		}
		log.Fatalf("get version failed: %v", err)
	}
	fmt.Printf("migration complete; version: %d dirty: %v\n", v, dirty)
}
