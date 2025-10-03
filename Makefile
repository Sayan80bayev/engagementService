POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB_NAME ?= engagement_db
DB_SSLMODE ?= disable

DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_DIR = migrations

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-status"
	@echo "  make migrate-create NAME=name"
	@echo "  make migrate-force VERSION=version"

.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

.PHONY: migrate-status
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please specify a migration name with NAME=<migration_name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

.PHONY: migrate-force
migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: Please specify a version with VERSION=<version_number>"; \
		exit 1; \
	fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)