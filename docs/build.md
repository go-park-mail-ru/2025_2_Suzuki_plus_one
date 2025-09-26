# Deployment Instructions

```bash
# Install dependencies
go mod tidy
# Build the application
go build -o build/server cmd/main.go
# Run the application
./build/server
```