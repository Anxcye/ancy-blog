// File: repository.go
// Purpose: Define PostgreSQL repository core runtime construction and lifecycle.
// Module: backend/internal/repository/postgres, persistence implementation layer.
// Related: domain-specific repository files in this package and schema_v1.sql.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func New(ctx context.Context, cfg config.DBConfig) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() error { return r.db.Close() }
