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
db-fill: ## fills the database with test data from testdata/postgres
	@echo "Filling Postgres with test data..."
	source .env && \
	sql_files=$$(ls -1 testdata/postgres/*.sql | sort); \
	for file in $$sql_files; do \
		echo "Executing $$file..."; \
		PGPASSWORD=$$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB -f $$file; \
	done
	@echo "Test data inserted"

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

## Redis
redis-start: ## starts the redis using docker-compose
	source .env && docker compose up -d redis
redis-stop: ## stops the redis using docker-compose
	source .env && docker compose down redis
redis-wipe: ## wipes the redis container and its volume
	source .env && docker compose down -v redis

## Minio
minio-start: ## starts the minio using docker-compose
	source .env && docker compose up -d s3
minio-stop: ## stops the minio using docker-compose
	source .env && docker compose down s3
minio-wipe: ## wipes the minio container and its volume
	source .env && docker compose down -v s3
minio-create: ## creates required buckets in minio
	@echo "Creating buckets in Minio..."
	source .env && \
	docker compose exec s3 /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc mb local/avatars || true && \
	mc mb local/posters || true && \
	mc mb local/trailers || true && \
	mc mb local/medias || true && \
	mc anonymous set public local/avatars || true && \
	mc anonymous set public local/posters || true && \
	mc anonymous set public local/trailers || true && \
	mc anonymous set none local/medias || true && \
	echo 'Buckets created' \
	"
minio-list: ## lists buckets in minio
	@echo "Listing buckets in Minio..."
	source .env && \
	docker compose exec s3 /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc ls local \
	"
minio-fill: ## fills minio with test data from testdata/minio using docker cp
	@echo "Filling Minio with test data using docker cp..."
	source .env && \
	docker cp testdata/minio/. $$(docker compose ps -q s3):/testdata/
	@echo "Test data uploaded"
	@echo "Filling Minio with test data using mc cp..."
	source .env && \
	docker compose exec s3 /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc cp --recursive /testdata/ local "
	@echo "Test data uploaded"

## All Services
all-start: ## starts all services using docker-compose
	source .env && docker compose up -d

all-stop: ## stops all services using docker-compose
	source .env && docker compose down

all-wipe: ## wipes all services containers and their volumes
	source .env && docker compose down -v

all-fill: ## fills all services with test data
	make minio-fill
	make db-fill

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
