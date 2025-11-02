// This code is in charge of loading configuration from environment variables.

package config

import (
	"log"
	"os"
	"time"
)

// TODO: resturucture config if it grows too big.
// See https://github.com/evrone/go-clean-template/blob/master/config/config.go for reference.

// Config holds all environment-driven settings.
// Please keep fields as in .env file for easier reference.
// See .env.example for documentation.
type Config struct {
	SERVER_SERVE_STRING           string
	SERVER_SERVE_PREFIX           string
	SERVER_JWT_SECRET             string
	SERVER_JWT_ACCESS_EXPIRATION  time.Duration
	SERVER_JWT_REFRESH_EXPIRATION time.Duration
	SERVER_NAME                   string
	SERVER_FRONTEND_URL           string
	POPFILMS_ENVIRONMENT          string

	// Database
	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	// Redis
	REDIS_HOST string

	// Minio
	MINIO_HOST          string
	MINIO_ROOT_USER     string
	MINIO_ROOT_PASSWORD string
}

// Load loads config from env vars with defaults and validation.
func Load() Config {
	cfg := Config{
		SERVER_SERVE_STRING:           getEnv("SERVER_SERVE_STRING", ":8080"),
		SERVER_SERVE_PREFIX:           trimTrailingSlash(getEnv("SERVER_SERVE_PREFIX", "")),
		SERVER_JWT_SECRET:             mustEnv("SERVER_JWT_SECRET"),
		SERVER_JWT_ACCESS_EXPIRATION:  parseDuration(getEnv("SERVER_JWT_ACCESS_EXPIRATION", "15m")),
		SERVER_JWT_REFRESH_EXPIRATION: parseDuration(getEnv("SERVER_JWT_REFRESH_EXPIRATION", "1440m")),
		SERVER_NAME:                   getEnv("SERVER_NAME", "Localhost"),
		SERVER_FRONTEND_URL:           trimTrailingSlash(mustEnv("SERVER_FRONTEND_URL")),
		POPFILMS_ENVIRONMENT:          getEnv("POPFILMS_ENVIRONMENT", "development"),

		// Database
		POSTGRES_HOST:     mustEnv("POSTGRES_HOST"),
		POSTGRES_USER:     mustEnv("POSTGRES_USER"),
		POSTGRES_PASSWORD: mustEnv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       mustEnv("POSTGRES_DB"),

		// Redis
		REDIS_HOST: mustEnv("REDIS_HOST"),

		// Minio
		MINIO_HOST:          mustEnv("MINIO_HOST"),
		MINIO_ROOT_USER:     mustEnv("MINIO_ROOT_USER"),
		MINIO_ROOT_PASSWORD: mustEnv("MINIO_ROOT_PASSWORD"),
	}

	return cfg
}

// getEnv returns the value of an env var or fallback if not set.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s is not set or empty, using default: <%s>", key, fallback)
	return fallback
}

// mustEnv panics if the env var is missing.
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		//value = "123123"
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// Convert string into Duration using time.ParseDuration.
// Calls fatal if input is invalid.
//
// # Example:
//
//	"15s" -> 15 seconds
//	"10m" -> 10 minutes
//	"24h" -> 24 hours
func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration %q: %v", value, err)
	}
	return duration
}

// trimTrailingSlash removes a trailing slash from the input string,
// unless the string is empty or just "/".
func trimTrailingSlash(s string) string {
	if len(s) >= 1 && s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}
