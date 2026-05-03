package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:54322/postgres"`
	SupabaseURL       string `env:"SUPABASE_URL,required"`
	SupabaseAnonKey   string `env:"SUPABASE_ANON_KEY,required"`
	SupabaseJWTSecret string `env:"SUPABASE_JWT_SECRET,required"`
	Port              string `env:"PORT" envDefault:"8080"`
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg
}
