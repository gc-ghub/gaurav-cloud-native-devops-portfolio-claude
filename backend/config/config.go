package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	Environment    string
	FrontendURL    string
	GitHubUsername string
	GitHubAPIURL   string
	GitHubToken    string
	CredlyUserID   string
	CredlyAPIURL   string
	CacheTTL       time.Duration
	MaxProjects    int
}

func Load() (*Config, error) {
	// Ignore a missing .env file — fine when real env vars are set directly (e.g. in prod).
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:5173"),
		GitHubUsername: getEnv("GITHUB_USERNAME", ""),
		GitHubAPIURL:   getEnv("GITHUB_API_URL", "https://api.github.com"),
		GitHubToken:    getEnv("GITHUB_TOKEN", ""),
		CredlyUserID:   getEnv("CREDLY_USER_ID", ""),
		CredlyAPIURL:   getEnv("CREDLY_API_URL", "https://www.credly.com"),
		CacheTTL:       time.Duration(getEnvInt("CACHE_TTL_MINUTES", 10)) * time.Minute,
		MaxProjects:    getEnvInt("MAX_PROJECTS", 6),
	}

	if cfg.GitHubUsername == "" {
		return nil, fmt.Errorf("GITHUB_USERNAME is required")
	}
	if cfg.CredlyUserID == "" {
		return nil, fmt.Errorf("CREDLY_USER_ID is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
