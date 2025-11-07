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
	@echo "Running tests excluding /mocks..."
	PACKAGES=$$(go list ./... | grep -v '/mocks') &&  \
	go test $$PACKAGES -coverprofile=coverage.out

coverage: ## generate test coverage report in html mode
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out


## Build
build: ## build the go application
	mkdir -p build/
	go build -o $(APP_EXECUTABLE) $(ENTRYPOINT)
	@echo "Build passed"

start: ## starts the application (requires built binary)
	@echo "Starting $(APP_EXECUTABLE)..."
	set -a && source .env && set +a && $(APP_EXECUTABLE)

live: ## launch the application using air for live reloading
	set -a && source .env && set +a && air

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


## Image
image-create: ## creates a docker image for the application
	docker build -t ghcr.io/sorrtory/popfilms-backend:latest .

image-push: ## pushes the docker image to the registry
	docker push ghcr.io/sorrtory/popfilms-backend:latest

## Ansible
ansible-api: ## runs the ansible playbook to deploy the OpenAPI spec server
	cd deployments && ansible-playbook update-api.yaml --vault-password-file=vault_password.sh
ansible-nginx: ## runs the ansible playbook to deploy the nginx server
	cd deployments && ansible-playbook update-nginx.yaml --vault-password-file=vault_password.sh
ansible-backend: ## runs the ansible playbook to deploy the backend server
	cd deployments && ansible-playbook update-backend.yaml --vault-password-file=vault_password.sh

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

## Database Migrations
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

remigrate: ## reapplies all migrations
	@echo "Restarting database"
	make db-wipe
	make db-start
	@echo "Waiting for database to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		pg_isready -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB && break; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	make migrate-up

## Redis
redis-start: ## starts the redis using docker-compose
	source .env && docker compose up -d redis
redis-stop: ## stops the redis using docker-compose
	source .env && docker compose down redis
redis-wipe: ## wipes the redis container and its volume
	source .env && docker compose down -v redis

## Minio

minio-pull: ## fetches data from vkedu minio to local minio
	@echo "Fetching data from vkedu drive to local minio..."
	rclone sync vkedu: testdata/minio/ -P

minio-start: ## starts the minio using docker-compose
	source .env && docker compose up -d minio
minio-stop: ## stops the minio using docker-compose
	source .env && docker compose down minio
minio-wipe: ## wipes the minio container and its volume
	source .env && docker compose down -v minio
minio-create: ## creates required buckets in minio
	@echo "Creating buckets in Minio..."
	source .env && \
	docker compose exec minio /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc mb local/actors || true && \
	mc mb local/avatars || true && \
	mc mb local/posters || true && \
	mc mb local/trailers || true && \
	mc mb local/medias || true && \
	mc anonymous set public local/actors || true && \
	mc anonymous set public local/avatars || true && \
	mc anonymous set public local/posters || true && \
	mc anonymous set public local/trailers || true && \
	mc anonymous set none local/medias || true && \
	echo 'Buckets created' \
	"
minio-list: ## lists buckets in minio
	@echo "Listing buckets in Minio..."
	source .env && \
	docker compose exec minio /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc ls local \
	"
minio-fill: ## fills minio with test data from testdata/minio using docker cp
	@echo "Filling Minio with test data using docker cp..."
	source .env && \
	docker cp testdata/minio/. $$(docker compose ps -q minio):/testdata/
	@echo "Test data uploaded"
	@echo "Filling Minio with test data using mc cp..."
	source .env && \
	docker compose exec minio /bin/sh -c " \
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

all-migrate: ## migrates database and minio services
	make migrate-up
	make minio-create

all-prepare: ## prepare all database: starts, migrates, fills with test data
	docker compose down -v
	docker compose up --build -d db redis minio

	@echo "Waiting for services to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		ready=true; \
		source .env && pg_isready -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB || ready=false; \
		docker compose exec redis redis-cli ping | grep -q PONG || ready=false; \
		docker compose exec minio mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD >/dev/null 2>&1 && \
		docker compose exec minio mc ls local >/dev/null 2>&1 || ready=false; \
		if [ "$$ready" = true ]; then \
			echo "All services are ready."; \
			break; \
		fi; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	if [ "$$ready" != true ]; then \
		echo "Timeout waiting for services to be ready."; \
		exit 1; \
	fi
	make all-migrate
	make all-fill
	@echo "All services prepared"

all-deploy: ## deploy all services using docker-compose
	docker compose -f compose.yaml up --build -d

all-bootstrap: ## bootstrap all: wipe, prepare, build and run
	@echo "Filling data"
	make all-prepare
	make all-stop
	@echo "Running application..."
	make all-deploy

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
