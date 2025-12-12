package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ENVFILE = "../../.env"

// Reads [ENVFILE], sets env vars, loads config and checks correctness.
func TestLoadEnvExample(t *testing.T) {
	// Open the .env file
	file, err := os.Open(ENVFILE)
	require.NoError(t, err, "Failed to open %s", ENVFILE)
	defer file.Close()

	// Create a map to store environment variables
	envVars := make(map[string]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "export ")

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Cut suffix from # symbol (remove inline comments)
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		require.Len(t, parts, 2, "Invalid line format: %s", line)

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

		// Set the environment variable
		envVars[key] = value
		fmt.Println(line)
		require.NoError(t, os.Setenv(key, value), "Failed to set environment variable: %s", key)
	}

	require.NoError(t, scanner.Err(), "Error reading .env.example")

	// Load the configuration using the Load function
	config := Load()

	// Compare the loaded configuration with the environment variables
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_SERVESTRING"], config.SERVICE_HTTP_SERVESTRING)
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_METRICS_SERVESTRING"], config.SERVICE_HTTP_METRICS_SERVESTRING)
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_SERVE_PREFIX"], config.SERVICE_HTTP_SERVE_PREFIX)
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_JWT_SECRET"], config.SERVICE_HTTP_JWT_SECRET)
	require.Equal(t, parseDuration(envVars["POPFILMS_SERVICE_HTTP_JWT_ACCESS_EXPIRATION"]), config.SERVICE_HTTP_JWT_ACCESS_EXPIRATION)
	require.Equal(t, parseDuration(envVars["POPFILMS_SERVICE_HTTP_JWT_REFRESH_EXPIRATION"]), config.SERVICE_HTTP_JWT_REFRESH_EXPIRATION)
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_NAME"], config.SERVICE_HTTP_NAME)
	require.Equal(t, envVars["POPFILMS_SERVICE_HTTP_FRONTEND_URL"], config.SERVICE_HTTP_FRONTEND_URL)
	require.Equal(t, envVars["POPFILMS_ENVIRONMENT"], config.ENVIRONMENT)

	// Database
	require.Equal(t, envVars["POSTGRES_HOST"], config.POSTGRES_HOST)
	require.Equal(t, envVars["APP_DB_USER"], config.APP_DB_USER)
	require.Equal(t, envVars["APP_DB_PASSWORD"], config.APP_DB_PASSWORD)
	require.Equal(t, envVars["POSTGRES_DB"], config.POSTGRES_DB)

	// Redis
	require.Equal(t, envVars["REDIS_HOST"], config.REDIS_HOST)

	// AWS S3
	require.Equal(t, envVars["AWS_S3_EXTERNAL_HOST"], config.AWS_S3_PUBLIC_URL)
	require.Equal(t, envVars["AWS_ACCESS_KEY_ID"], config.AWS_ACCESS_KEY_ID)
	require.Equal(t, envVars["AWS_SECRET_ACCESS_KEY"], config.AWS_SECRET_ACCESS_KEY)
	require.Equal(t, envVars["AWS_DEFAULT_REGION"], config.AWS_REGION)
	require.Equal(t, envVars["AWS_ENDPOINT_URL"], config.AWS_S3_ENDPOINT)

	// Services

	// Auth
	require.Equal(t, envVars["POPFILMS_SERVICE_AUTH_SERVESTRING"], config.SERVICE_AUTH_SERVE_STRING)
	require.Equal(t, envVars["POPFILMS_SERVICE_AUTH_METRICS_SERVESTRING"], config.SERVICE_AUTH_METRICS_SERVE_STRING)

	// Search
	require.Equal(t, envVars["POPFILMS_SERVICE_SEARCH_SERVESTRING"], config.SERVICE_SEARCH_SERVE_STRING)
	require.Equal(t, envVars["POPFILMS_SERVICE_SEARCH_METRICS_SERVESTRING"], config.SERVICE_SEARCH_METRICS_SERVE_STRING)
}
