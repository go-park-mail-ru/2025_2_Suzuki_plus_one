# How to launch server

## Prerequisites

```bash
# Create .env file in the root directory according to .env.example
cp .env.example .env
# Load environment variables
source .env
```

## TODO : clean this up. Use makefile instead

## Launch server

```bash
# Install dependencies
go mod tidy
# Build or run the application
go build -o build/server cmd/main.go && ./build/server
# OR
go run ./cmd/main.go
```

## Access the server

```bash
# Access the server
curl -X GET http://localhost:8080
```

## Database

```bash
# Make sure Postgres settings are set in the environment variables
source .env
# Launch Postgres using Docker
docker compose up -d db

# OR
source .env
make up-db
```
