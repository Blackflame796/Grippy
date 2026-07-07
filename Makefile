include .env
export

# Variables
MIGRATIONS_PATH = ./migrations
DB_URL = "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

.PHONY: create_migration migrate_up migrate_down

create_migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext=sql -dir=$(MIGRATIONS_PATH) -seq $$name

migrate_up:
	@read -p "Enter migration number: " number; \
	migrate -path=$(MIGRATIONS_PATH) -database $(DB_URL) -verbose up $$number

migrate_down:
	@read -p "Enter migration number: " number; \
	migrate -path=$(MIGRATIONS_PATH) -database $(DB_URL) -verbose down $$number

force_migrate_down:
	@read -p "Enter version to force (usually the last successful version): " v; \
	migrate -path=$(MIGRATIONS_PATH) -database $(DB_URL) force $$v
