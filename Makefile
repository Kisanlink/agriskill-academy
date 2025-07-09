# Makefile for ASA Backend
# Author: Karthikeya Akhandam
# Description: Build, run, and manage the ASA backend application

# Variables
APP_NAME=asa
MAIN_PATH=cmd/server/main.go
BUILD_DIR=build
BINARY_NAME=$(APP_NAME)
PORT=3000

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "ASA Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build targets
.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux"

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building $(APP_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME).exe"

.PHONY: build-mac
build-mac: ## Build for macOS
	@echo "Building $(APP_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_PATH)
	@echo "macOS build complete: $(BUILD_DIR)/$(BINARY_NAME)-mac"

# Run targets
.PHONY: run
run: ## Run the application
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	@go run $(MAIN_PATH)

.PHONY: run-build
run-build: build ## Build and run the application
	@echo "Starting $(APP_NAME) from build..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Development targets
.PHONY: air
air: ## Run with Air for hot reloading (requires air to be installed)
	@echo "Starting $(APP_NAME) with Air for hot reloading..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

.PHONY: dev
dev: ## Run in development mode (alias for air)
	@$(MAKE) air

# Database targets
.PHONY: migrate
migrate: ## Run all database migrations in order
	@echo "Running all database migrations..."
	@echo "1. Applying complete database schema..."
	@if [ -f "scripts/001_apply_schema.sh" ]; then \
		bash scripts/001_apply_schema.sh; \
	else \
		echo "[ERROR] Schema migration script not found (scripts/001_apply_schema.sh)"; \
		exit 1; \
	fi
	@echo "2. Applying messages timestamp fix..."
	@if [ -f "scripts/007_fix_messages.sh" ]; then \
		bash scripts/007_fix_messages.sh; \
	else \
		echo "[ERROR] Messages fix script not found (scripts/007_fix_messages.sh)"; \
		exit 1; \
	fi
	@echo "3. Applying student profiles rename..."
	@if [ -f "scripts/008_rename_profiles.sh" ]; then \
		bash scripts/008_rename_profiles.sh; \
	else \
		echo "[ERROR] Profiles rename script not found (scripts/008_rename_profiles.sh)"; \
		exit 1; \
	fi
	@echo "✅ All migrations completed successfully!"

.PHONY: migrate-schema
migrate-schema: ## Apply only the database schema migration
	@echo "Applying database schema migration..."
	@if [ -f "scripts/001_apply_schema.sh" ]; then \
		bash scripts/001_apply_schema.sh; \
	else \
		echo "[ERROR] Schema migration script not found (scripts/001_apply_schema.sh)"; \
		exit 1; \
	fi

.PHONY: migrate-messages
migrate-messages: ## Apply only the messages timestamp fix
	@echo "Applying messages timestamp fix..."
	@if [ -f "scripts/007_fix_messages.sh" ]; then \
		bash scripts/007_fix_messages.sh; \
	else \
		echo "[ERROR] Messages fix script not found (scripts/007_fix_messages.sh)"; \
		exit 1; \
	fi

.PHONY: migrate-profiles
migrate-profiles: ## Apply only the student profiles rename
	@echo "Applying student profiles rename..."
	@if [ -f "scripts/008_rename_profiles.sh" ]; then \
		bash scripts/008_rename_profiles.sh; \
	else \
		echo "[ERROR] Profiles rename script not found (scripts/008_rename_profiles.sh)"; \
		exit 1; \
	fi

.PHONY: migrate-reset
migrate-reset: ## Reset database and run all migrations
	@echo "Resetting database and running all migrations..."
	@if [ -f "scripts/reset_db.sh" ]; then \
		bash scripts/reset_db.sh; \
	else \
		echo "[ERROR] Reset script not found (scripts/reset_db.sh)"; \
		exit 1; \
	fi

# Testing targets
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@go test ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Code quality targets
.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

.PHONY: tidy
tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	@go mod tidy

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

.PHONY: clean-uploads
clean-uploads: ## Clean uploaded files
	@echo "Cleaning uploaded files..."
	@rm -rf uploads/resumes/*
	@rm -rf uploads/certificates/*

# Dependencies
.PHONY: deps
deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download

.PHONY: install-air
install-air: ## Install Air for hot reloading
	@echo "Installing Air for hot reloading..."
	@go install github.com/cosmtrek/air@latest

# Quick start
.PHONY: setup
setup: deps install-air ## Setup development environment
	@echo "Development environment setup complete!"

.PHONY: start
start: migrate run ## Setup database and start the application
	@echo "Application started successfully!" 