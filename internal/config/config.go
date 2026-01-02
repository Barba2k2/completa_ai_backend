package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	FirebaseProjectID string
	FirebaseBucket    string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseBucket:    getEnv("FIREBASE_STORAGE_BUCKET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
