# Makefile for Go projects template
# https://www.mohitkhare.com/blog/go-makefile/

# Makefile settings
SHELL := /bin/bash # Use bash syntax

# Go settings
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
APP=server
APP_EXECUTABLE="./build/$(APP)"
ENTRYPOINT=./cmd/main.go

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor


## Quality
check-quality: ## runs code quality checks
	make lint
	make fmt
	make vet

# Append || true below if blocking local developement
lint: ## go linting. Update and use specific lint tool and options
	golangci-lint run --enable-all

vet: ## go vet
	go vet ./...

fmt: ## runs go formatter
	go fmt ./...

tidy: ## runs tidy to fix go.mod dependencies
	go mod tidy

## Test
test: ## runs tests and create generates coverage report
	go tool cover -func=coverage.out
	go test -v -timeout 10m ./... -coverprofile=coverage.out

coverage: ## displays test coverage report in html mode
	make test
	go tool cover -html=coverage.out -o coverage.html


## Build
build: ## build the go application
	mkdir -p build/
	go build -o $(APP_EXECUTABLE) $(ENTRYPOINT)
	@echo "Build passed"

start: ## starts the application (requires built binary)
	@echo "Starting $(APP_EXECUTABLE)..."
	source .env && $(APP_EXECUTABLE)

live: ## launch the application using air for live reloading
	source .env && air

run: ## builds and starts the application
	make build
	chmod +x $(APP_EXECUTABLE)
	make start

clean: ## cleans binary and other generated files
	go clean
	rm -rf build/
	rm -f coverage*.out

vendor: ## all packages required to support builds and tests in the /vendor directory
	go mod vendor

## Database

# Postgres commands using docker-compose
db-start: ## starts the database using docker-compose
	source .env && docker compose up -d db
db-stop: ## stops the database using docker-compose
	source .env && docker compose down db
db-wipe: ## wipes the database container and its volume
	source .env && docker compose down -v db
db-logs: ## views database logs
	source .env && docker compose logs -f db
db-ps:   ## lists docker status of the database container
	source .env && docker compose ps db

# Psql shell command
db-shell: ## connects to the database shell with psql
	source .env && \
	PGPASSWORD=$$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB

db-check: ## checks database connection
	source .env && \
	pg_isready -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB

## Migrations
migrate-create: ## creates a new migration file. Usage: make migrate-create NAME=<name>
	@if [ -z "$(NAME)" ]; then \
		echo "NAME is undefined. Usage: make migrate-create NAME=<name>"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(NAME)"
	migrate create -ext sql -dir ./migrations -seq $(NAME)

migrate-up: ## applies all up migrations
	source .env && \
	POSTGRESQL_URL="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@localhost:5432/$$POSTGRES_DB?sslmode=disable" && \
	migrate -path ./migrations -database $$POSTGRESQL_URL up

migrate-down: ## applies all down migrations
	source .env && \
	POSTGRESQL_URL="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@localhost:5432/$$POSTGRES_DB?sslmode=disable" && \
	migrate -path ./migrations -database $$POSTGRESQL_URL down


.PHONY: help
## Help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
