# How to launch server

```bash
# Install dependencies
go mod tidy
# Build or run the application
go build -o build/server cmd/main.go && ./build/server
# OR
go run ./cmd/main.go 
```