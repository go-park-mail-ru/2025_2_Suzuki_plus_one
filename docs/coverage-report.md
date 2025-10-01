# Coverage Report

> 64.5%

## Commands

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
# Display coverage summary
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```