package config

import (
	"log"
	"os"
)

// Config holds all environment-driven settings.
// Please keep fields as in .env file for easier reference.
// See .env.example for documentation.
type Config struct {
	SERVER_SERVE_STRING string
}

// Load loads config from env vars with defaults and validation.
func Load() Config {
	cfg := Config{
		SERVER_SERVE_STRING: getEnv("SERVER_SERVE_STRING", ":8080"),
	}

	return cfg
}

// getEnv returns the value of an env var or fallback if not set.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s is not set, using default: %s", key, fallback)
	return fallback
}

// mustEnv panics if the env var is missing.
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}
