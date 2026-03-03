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
	R2   R2Config
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

type R2Config struct {
	Enabled         bool
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicBaseURL   string
	Region          string
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
		R2: R2Config{
			Enabled:         parseBool(getEnv("R2_ENABLED", "false")),
			AccountID:       getEnv("R2_ACCOUNT_ID", ""),
			AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
			Bucket:          getEnv("R2_BUCKET", ""),
			PublicBaseURL:   strings.TrimRight(getEnv("R2_PUBLIC_BASE_URL", ""), "/"),
			Region:          getEnv("R2_REGION", "auto"),
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

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
