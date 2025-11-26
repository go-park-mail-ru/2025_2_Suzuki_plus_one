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
	SERVICE_HTTP_SERVESTRING            string
	SERVICE_HTTP_METRICS_SERVESTRING    string
	SERVICE_HTTP_SERVE_PREFIX           string
	SERVICE_HTTP_JWT_SECRET             string
	SERVICE_HTTP_JWT_ACCESS_EXPIRATION  time.Duration
	SERVICE_HTTP_JWT_REFRESH_EXPIRATION time.Duration
	SERVICE_HTTP_NAME                   string
	SERVICE_HTTP_FRONTEND_URL           string
	ENVIRONMENT                         string

	// Database
	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	// Redis
	REDIS_HOST string

	// Minio
	MINIO_INTERNAL_HOST string
	MINIO_EXTERNAL_HOST string
	MINIO_ROOT_USER     string
	MINIO_ROOT_PASSWORD string

	// Services

	// Auth
	SERVICE_AUTH_SERVE_STRING         string
	SERVICE_AUTH_METRICS_SERVE_STRING string

	// Search
	SERVICE_SEARCH_SERVE_STRING         string
	SERVICE_SEARCH_METRICS_SERVE_STRING string
}

// Load loads config from env vars with defaults and validation.
func Load() Config {
	cfg := Config{
		SERVICE_HTTP_SERVESTRING:            getEnv("POPFILMS_SERVICE_HTTP_SERVESTRING", ":8080"),
		SERVICE_HTTP_METRICS_SERVESTRING:    getEnv("POPFILMS_SERVICE_HTTP_METRICS_SERVESTRING", ":8880"),
		SERVICE_HTTP_SERVE_PREFIX:           trimTrailingSlash(getEnv("POPFILMS_SERVICE_HTTP_SERVE_PREFIX", "")),
		SERVICE_HTTP_JWT_SECRET:             mustEnv("POPFILMS_SERVICE_HTTP_JWT_SECRET"),
		SERVICE_HTTP_JWT_ACCESS_EXPIRATION:  parseDuration(getEnv("POPFILMS_SERVICE_HTTP_JWT_ACCESS_EXPIRATION", "15m")),
		SERVICE_HTTP_JWT_REFRESH_EXPIRATION: parseDuration(getEnv("POPFILMS_SERVICE_HTTP_JWT_REFRESH_EXPIRATION", "1440m")),
		SERVICE_HTTP_NAME:                   getEnv("POPFILMS_SERVICE_HTTP_NAME", "Localhost"),
		SERVICE_HTTP_FRONTEND_URL:           trimTrailingSlash(mustEnv("POPFILMS_SERVICE_HTTP_FRONTEND_URL")),
		ENVIRONMENT:                         getEnv("POPFILMS_ENVIRONMENT", "development"),

		// Database
		POSTGRES_HOST:     mustEnv("POSTGRES_HOST"),
		POSTGRES_USER:     mustEnv("POSTGRES_USER"),
		POSTGRES_PASSWORD: mustEnv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       mustEnv("POSTGRES_DB"),

		// Redis
		REDIS_HOST: mustEnv("REDIS_HOST"),

		// Minio
		MINIO_INTERNAL_HOST: mustEnv("MINIO_INTERNAL_HOST"),
		MINIO_EXTERNAL_HOST: mustEnv("MINIO_EXTERNAL_HOST"),
		MINIO_ROOT_USER:     mustEnv("MINIO_ROOT_USER"),
		MINIO_ROOT_PASSWORD: mustEnv("MINIO_ROOT_PASSWORD"),

		// Services

		// Auth
		SERVICE_AUTH_SERVE_STRING:         mustEnv("POPFILMS_SERVICE_AUTH_SERVESTRING"),
		SERVICE_AUTH_METRICS_SERVE_STRING: mustEnv("POPFILMS_SERVICE_AUTH_METRICS_SERVESTRING"),

		// Search
		SERVICE_SEARCH_SERVE_STRING:         mustEnv("POPFILMS_SERVICE_SEARCH_SERVESTRING"),
		SERVICE_SEARCH_METRICS_SERVE_STRING: mustEnv("POPFILMS_SERVICE_SEARCH_METRICS_SERVESTRING"),
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
