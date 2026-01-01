package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	SupabaseURL      string
	SupabaseAnonKey  string
	SupabaseJWTSecret string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		SupabaseURL:      getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:  getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseJWTSecret: getEnv("SUPABASE_JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
