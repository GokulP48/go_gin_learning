# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Migration parameters
MIGRATE_CMD=cmd/migrate/main.go

.PHONY: help migrate-create migrate-up migrate-down migrate-version migrate-to run build clean test

# Default target
help:
	@echo "Available commands:"
	@echo "  make migrate-create name=migration_name  - Create new migration files"
	@echo "  make migrate-up                          - Run all pending migrations"
	@echo "  make migrate-down                        - Rollback one migration"
	@echo "  make migrate-version                     - Show current migration version"
	@echo "  make migrate-to version=N                - Migrate to specific version N"
	@echo "  make run                                 - Run the application (with auto-migration)"
	@echo "  make build                               - Build the application"
	@echo "  make test                                - Run tests"
	@echo "  make clean                               - Clean build files"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-create name=create_users_table"
	@echo "  make migrate-create name=add_phone_to_users"
	@echo "  make migrate-up"
	@echo "  make migrate-to version=3"

# Create a new migration
migrate-create:
ifndef name
	@echo "Error: name parameter is required"
	@echo "Usage: make migrate-create name=create_users_table"
	@exit 1
endif
	@$(GORUN) $(MIGRATE_CMD) -create $(name)

# Run all up migrations
migrate-up:
	@echo "Running migrations..."
	@$(GORUN) $(MIGRATE_CMD) -up

# Rollback one migration
migrate-down:
	@echo "Rolling back one migration..."
	@$(GORUN) $(MIGRATE_CMD) -down

# Show migration version
migrate-version:
	@$(GORUN) $(MIGRATE_CMD) -version

# Migrate to specific version
migrate-to:
ifndef version
	@echo "Error: version parameter is required"
	@echo "Usage: make migrate-to version=5"
	@exit 1
endif
	@echo "Migrating to version $(version)..."
	@$(GORUN) $(MIGRATE_CMD) -migrate-to $(version)

# Run the application (with auto-migration)
run:
	@echo "Starting application with auto-migration..."
	@$(GORUN) main.go

# Build the application
build:
	@echo "Building application..."
	@$(GOBUILD) -o bin/app main.go
	@echo "Built binary: bin/app"

# Build migration tool
build-migrate:
	@echo "Building migration tool..."
	@$(GOBUILD) -o bin/migrate $(MIGRATE_CMD)
	@echo "Built migration tool: bin/migrate"

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GOGET) github.com/golang-migrate/migrate/v4
	@$(GOGET) github.com/golang-migrate/migrate/v4/database/postgres
	@$(GOGET) github.com/lib/pq
	@$(GOMOD) tidy

# Development workflow shortcuts
dev-setup: deps migrate-up
	@echo "Development environment ready!"

dev-reset: 
	@echo "Resetting database to version 0..."
	@$(GORUN) $(MIGRATE_CMD) -migrate-to 0
	@$(GORUN) $(MIGRATE_CMD) -up
	@echo "Database reset complete!"