package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/angith/issueboard/internal/logger"
)

type Config struct {
	DatabaseURL       string       `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:54322/postgres"`
	SupabaseURL       string       `env:"SUPABASE_URL,required"`
	SupabaseAnonKey   string       `env:"SUPABASE_ANON_KEY,required"`
	SupabaseJWTSecret string       `env:"SUPABASE_JWT_SECRET,required"`
	GitHubToken       string       `env:"GITHUB_TOKEN"`
	Port              string       `env:"PORT" envDefault:"8080"`
	LogLevel          logger.Level // populated after env parsing — see Load()
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, reading from environment variables")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Error("Failed to parse config")
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Convert the LOG_LEVEL string from the environment into the typed Level.
	cfg.LogLevel = logger.ParseLevel(os.Getenv("LOG_LEVEL"))

	return &cfg, nil
}
