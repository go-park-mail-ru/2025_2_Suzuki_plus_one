# Makefile for Go projects template
# https://www.mohitkhare.com/blog/go-makefile/

# Makefile settings
SHELL := /bin/bash # Use bash syntax

# Go settings
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)

# Default value for ENV_FILE if not set
ENV_FILE ?= .env

SERVICES = http auth search

# Default service to build/run/test
SERVICE ?= http


ENTRYPOINT = ./cmd/${SERVICE}/main.go
APP_EXECUTABLE = build/${SERVICE}

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
	set -a && source $(ENV_FILE) && set +a && $(APP_EXECUTABLE)

live: ## launch the application using air for live reloading
	@echo "Starting live reload for $(SERVICE)..."
	set -a && source $(ENV_FILE) && set +a && SERVICE=$(SERVICE) air

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

swagger:
	@echo "Launching swagger UI on http://localhost"
	docker run -p 80:8080 -e SWAGGER_JSON=/api/popfilms.yaml -v ./api:/api swaggerapi/swagger-ui 

## Image
image-create: ## creates a docker image for the application
	docker build -t ghcr.io/sorrtory/popfilms-backend:latest .

image-push: ## pushes the docker image to the registry
	docker push ghcr.io/sorrtory/popfilms-backend:latest

image-update: ## builds and pushes the docker image to the registry
	make image-create
	make image-push

## Ansible
ansible-api: ## runs the ansible playbook to deploy the OpenAPI spec server
	cd deployments && ansible-playbook update-api.yaml --vault-password-file=vault_password.sh
ansible-nginx: ## runs the ansible playbook to deploy the nginx server
	cd deployments && ansible-playbook update-nginx.yaml --vault-password-file=vault_password.sh
ansible-frontend: ## runs the ansible playbook to deploy the frontend server
	cd deployments && ansible-playbook update-frontend.yaml --vault-password-file=vault_password.sh
ansible-backend: ## runs the ansible playbook to deploy the backend server
	cd deployments && ansible-playbook update-backend.yaml --vault-password-file=vault_password.sh
ansible-bootstrap: ## runs the ansible playbook to bootstrap the backend server
	cd deployments && ansible-playbook update-backend.yaml --vault-password-file=vault_password.sh -e deploy_mode=bootstrap
update-deploy-branch: ## updates deploy branch by merging dev into deploy
	git switch deploy
	git pull origin deploy
	git merge dev
	git push origin deploy
	git switch dev
update-deploy-backend: ## updates backend deployment by building and pushing image and running ansible playbook
	make image-update
	make ansible-backend

## Database

# Postgres commands using docker-compose
db-start: ## starts the database using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) up -d db
db-stop: ## stops the database using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down --env-file $(ENV_FILE) db
db-wipe: ## wipes the database container and its volume
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down -v --env-file $(ENV_FILE) db
db-logs: ## views database logs
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) logs -f db
db-ps:   ## lists docker status of the database container
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) ps db
db-fill: ## fills the database with test data from testdata/postgres
	@echo "Filling Postgres with test data..."
	source $(ENV_FILE) && \
	sql_files=$$(ls -1 testdata/postgres/*.sql | sort); \
	for file in $$sql_files; do \
		echo "Executing $$file..."; \
		PGPASSWORD=$$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB -f $$file; \
	done
	@echo "Test data inserted"

# Psql shell command
db-shell: ## connects to the database shell with psql
	source $(ENV_FILE) && \
	PGPASSWORD=$$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB

db-check: ## checks database connection
	source $(ENV_FILE) && \
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
	source $(ENV_FILE) && \
	POSTGRESQL_URL="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@localhost:5432/$$POSTGRES_DB?sslmode=disable" && \
	migrate -path ./migrations -database $$POSTGRESQL_URL up

migrate-down: ## applies all down migrations
	source $(ENV_FILE) && \
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
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) up -d redis
redis-stop: ## stops the redis using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down --env-file $(ENV_FILE) redis
redis-wipe: ## wipes the redis container and its volume
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down -v --env-file $(ENV_FILE) redis

## [DEPRECATED] Minio

minio-pull: ## fetches data from vkedu minio to local minio
	@echo "Fetching data from vkedu drive to local minio..."
	rclone sync vkedu: testdata/minio/ -P

minio-start: ## starts the minio using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) up -d minio
minio-stop: ## stops the minio using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down minio
minio-wipe: ## wipes the minio container and its volume
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down -v minio
minio-create: ## creates required buckets in minio
	@echo "Creating buckets in Minio..."
	source $(ENV_FILE) && \
	docker compose --env-file $(ENV_FILE) exec minio /bin/sh -c " \
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
	source $(ENV_FILE) && \
	docker compose --env-file $(ENV_FILE) exec minio /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc ls local \
	"
minio-fill:
	@echo "Mirroring test data directly..."
	source $(ENV_FILE) && \
	docker compose --env-file $(ENV_FILE) exec minio /bin/sh -c " \
	mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD && \
	mc mb -p local/actors local/avatars local/medias local/posters local/trailers && \
	mc mirror --overwrite /testdata/actors/ local/actors/ && \
	mc mirror --overwrite /testdata/avatars/ local/avatars/ && \
	mc mirror --overwrite /testdata/medias/ local/medias/ && \
	mc mirror --overwrite /testdata/posters/ local/posters/ && \
	mc mirror --overwrite /testdata/trailers/ local/trailers/"
	@echo "Test data uploaded"

## All Services
all-start: ## starts all services using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) up -d

all-stop: ## stops all services using docker-compose
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down

all-wipe: ## wipes all services containers and their volumes
	source $(ENV_FILE) && docker compose --env-file $(ENV_FILE) down -v

all-fill: ## fills all services with test data
# 	make minio-fill
	make db-fill

all-migrate: ## migrates database
	make migrate-up
# 	make minio-create

all-prepare: ## prepare all database: starts, migrates, fills with test data
	docker compose --env-file $(ENV_FILE) down -v
# 	docker compose up --build -d db redis minio
	docker compose --env-file $(ENV_FILE) up --build -d db redis

	@echo "Waiting for services to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		ready=true; \
		source $(ENV_FILE) && pg_isready -h localhost -p 5432 -U $$POSTGRES_USER -d $$POSTGRES_DB || ready=false; \
		docker compose --env-file $(ENV_FILE) exec redis redis-cli ping | grep -q PONG || ready=false; \
		# docker compose --env-file $(ENV_FILE) exec minio mc alias set local http://localhost:9000 $$MINIO_ROOT_USER $$MINIO_ROOT_PASSWORD >/dev/null 2>&1 && \
		# docker compose --env-file $(ENV_FILE) exec minio mc ls local >/dev/null 2>&1 || ready=false; \
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

all-service-deploy: ## runs all microservices in docker-compose
	docker compose --env-file $(ENV_FILE) -f compose.yaml up --build -d backend-http backend-auth backend-search

all-deploy: ## deploy all services using docker-compose
	docker compose --env-file $(ENV_FILE) -f compose.yaml up --build -d

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
