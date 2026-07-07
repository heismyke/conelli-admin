APP_NAME := conelli-admin-backend
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)
GOCACHE ?= $(CURDIR)/.cache/go-build
MIGRATE_ARGS := $(filter-out migrate up,$(MAKECMDGOALS))
COMPOSE := docker compose --project-name conelli -f devops/docker-compose.yml
DEV_COMPOSE := docker compose -f devops/docker-compose.dev.yml
PROD_COMPOSE := docker compose -f devops/docker-compose.prod.yml

.PHONY: help run dev dev-db-up dev-up dev-down dev-logs dev-migrate seed-projects prod-up prod-down prod-logs build db-create migrate up docker-build docker-up docker-down docker-logs docker-migrate fmt tidy clean

help:
	@echo "Available targets:"
	@echo "  make run             - run the API server"
	@echo "  make dev             - alias for run"
	@echo "  make dev-db-up       - run Postgres in Docker"
	@echo "  make dev-up          - run API and Postgres in Docker"
	@echo "  make dev-migrate     - run migrations in Docker"
	@echo "  make seed-projects   - upsert bundled public projects into admin data"
	@echo "  make dev-down        - stop Docker services"
	@echo "  make dev-logs        - tail Docker API logs"
	@echo "  make prod-up         - run production API behind Nginx"
	@echo "  make prod-down       - stop production Docker services"
	@echo "  make prod-logs       - tail production Docker logs"
	@echo "  make build           - build the API binary"
	@echo "  make db-create       - create the configured database if missing"
	@echo "  make migrate up      - run all up migrations"
	@echo "  make migrate up N    - run up migrations through number N"
	@echo "  make docker-build    - build the Docker image"
	@echo "  make docker-up       - run API and Postgres with Docker Compose"
	@echo "  make docker-migrate  - run migrations inside Docker Compose"
	@echo "  make docker-down     - stop Docker Compose services"
	@echo "  make docker-logs     - tail API logs"
	@echo "  make fmt             - format Go files"
	@echo "  make tidy            - tidy Go modules"
	@echo "  make clean           - remove build output"

run:
	GOCACHE=$(GOCACHE) go run main.go

dev: run

dev-db-up:
	$(COMPOSE) up -d conelli-postgres

dev-up:
	$(DEV_COMPOSE) up -d --build

dev-migrate:
	$(COMPOSE) run --rm conelli-admin-backend conelli-admin-migrate up

dev-down:
	$(DEV_COMPOSE) down

dev-logs:
	$(DEV_COMPOSE) logs -f conelli-admin-backend

prod-up:
	$(PROD_COMPOSE) up -d --build

prod-down:
	$(PROD_COMPOSE) down

prod-logs:
	$(PROD_COMPOSE) logs -f

build:
	mkdir -p $(BIN_DIR)
	GOCACHE=$(GOCACHE) go build -o $(BIN) .

db-create:
	./scripts/ensure-db.sh

migrate: db-create
	GOCACHE=$(GOCACHE) go run ./cmd/migrate $(if $(MIGRATE_ARGS),$(MIGRATE_ARGS),up)

seed-projects: db-create
	GOCACHE=$(GOCACHE) go run ./cmd/seed-projects

up:
	@:

%:
	@:

docker-build:
	$(COMPOSE) build

docker-up:
	$(COMPOSE) up -d --build

docker-migrate:
	$(COMPOSE) run --rm conelli-admin-backend conelli-admin-migrate up

docker-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f api

fmt:
	gofmt -w .

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR) .cache
