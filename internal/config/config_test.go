package config

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ENVFILE = "../../.env.example"

// Reads [ENVFILE], sets env vars, loads config and checks correctness.
func TestLoadEnvExample(t *testing.T) {
	// Open the .env.example file
	file, err := os.Open(ENVFILE)
	require.NoError(t, err, "Failed to open .env.example")
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
		require.NoError(t, os.Setenv(key, value), "Failed to set environment variable: %s", key)
	}

	require.NoError(t, scanner.Err(), "Error reading .env.example")

	// Load the configuration using the Load function
	config := Load()

	// Compare the loaded configuration with the environment variables
	require.Equal(t, envVars["SERVER_SERVE_STRING"], config.SERVER_SERVE_STRING)
	require.Equal(t, envVars["SERVER_SERVE_PREFIX"], config.SERVER_SERVE_PREFIX)
	require.Equal(t, envVars["SERVER_JWT_SECRET"], config.SERVER_JWT_SECRET)
	require.Equal(t, parseDuration(envVars["SERVER_JWT_ACCESS_EXPIRATION"]), config.SERVER_JWT_ACCESS_EXPIRATION)
	require.Equal(t, parseDuration(envVars["SERVER_JWT_REFRESH_EXPIRATION"]), config.SERVER_JWT_REFRESH_EXPIRATION)
	require.Equal(t, envVars["SERVER_NAME"], config.SERVER_NAME)
	require.Equal(t, envVars["SERVER_FRONTEND_URL"], config.SERVER_FRONTEND_URL)
	require.Equal(t, envVars["POPFILMS_ENVIRONMENT"], config.POPFILMS_ENVIRONMENT)

	// Database
	require.Equal(t, envVars["POSTGRES_HOST"], config.POSTGRES_HOST)
	require.Equal(t, envVars["POSTGRES_USER"], config.POSTGRES_USER)
	require.Equal(t, envVars["POSTGRES_PASSWORD"], config.POSTGRES_PASSWORD)
	require.Equal(t, envVars["POSTGRES_DB"], config.POSTGRES_DB)

	// Redis
	require.Equal(t, envVars["REDIS_HOST"], config.REDIS_HOST)

	// Minio
	require.Equal(t, envVars["MINIO_HOST"], config.MINIO_HOST)
	require.Equal(t, envVars["MINIO_ROOT_USER"], config.MINIO_ROOT_USER)
	require.Equal(t, envVars["MINIO_ROOT_PASSWORD"], config.MINIO_ROOT_PASSWORD)
}
