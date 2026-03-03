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
