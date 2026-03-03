// File: config.go
// Purpose: Define runtime configuration schema and load values from environment variables.
// Module: backend/internal/config, application configuration layer.
// Related: cmd/server bootstrap and internal/server startup.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	Auth AuthConfig
	DB   DBConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Host string
	Port int
}

type AuthConfig struct {
	AdminUsername          string
	AdminPassword          string
	AccessTokenTTLSeconds  int
	RefreshTokenTTLSeconds int
}

type DBConfig struct {
	Host         string
	Port         int
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

func Load() (*Config, error) {
	port, err := parseInt(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}
	accessTTL, err := parseInt(getEnv("AUTH_ACCESS_TOKEN_TTL_SECONDS", "3600"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_ACCESS_TOKEN_TTL_SECONDS: %w", err)
	}
	refreshTTL, err := parseInt(getEnv("AUTH_REFRESH_TOKEN_TTL_SECONDS", "604800"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_REFRESH_TOKEN_TTL_SECONDS: %w", err)
	}
	dbPort, err := parseInt(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	maxOpenConns, err := parseInt(getEnv("DB_MAX_OPEN_CONNS", "20"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}
	maxIdleConns, err := parseInt(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "ancy-blog-api"),
			Env:  strings.ToLower(getEnv("APP_ENV", "dev")),
		},
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: port,
		},
		Auth: AuthConfig{
			AdminUsername:          getEnv("AUTH_ADMIN_USERNAME", "admin"),
			AdminPassword:          getEnv("AUTH_ADMIN_PASSWORD", "123456"),
			AccessTokenTTLSeconds:  accessTTL,
			RefreshTokenTTLSeconds: refreshTTL,
		},
		DB: DBConfig{
			Host:         getEnv("DB_HOST", "127.0.0.1"),
			Port:         dbPort,
			Name:         getEnv("DB_NAME", "ancy_blog"),
			User:         getEnv("DB_USER", "ancy"),
			Password:     getEnv("DB_PASSWORD", "ancy_dev_password"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseInt(raw string) (int, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return v, nil
}
