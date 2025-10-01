# How to launch server

## Prerequisites
```bash
# Create .env file in the root directory according to .env.example
cp .env.example .env
# Load environment variables
source .env
```

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