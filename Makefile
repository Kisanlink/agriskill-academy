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
	@if [ -f "scripts/001_apply_schema.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/001_apply_schema.ps1; \
	else \
		echo "[ERROR] Schema migration script not found (scripts/001_apply_schema.ps1)"; \
		exit 1; \
	fi
	@echo "2. Applying messages timestamp fix..."
	@if [ -f "scripts/007_fix_messages.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/007_fix_messages.ps1; \
	else \
		echo "[ERROR] Messages fix script not found (scripts/007_fix_messages.ps1)"; \
		exit 1; \
	fi
	@echo "3. Applying student profiles rename..."
	@if [ -f "scripts/008_rename_profiles.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/008_rename_profiles.ps1; \
	else \
		echo "[ERROR] Profiles rename script not found (scripts/008_rename_profiles.ps1)"; \
		exit 1; \
	fi
	@echo "✓ All migrations completed successfully!"

.PHONY: migrate-schema
migrate-schema: ## Apply only the database schema migration
	@echo "Applying database schema migration..."
	@if [ -f "scripts/001_apply_schema.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/001_apply_schema.ps1; \
	else \
		echo "[ERROR] Schema migration script not found (scripts/001_apply_schema.ps1)"; \
		exit 1; \
	fi

.PHONY: migrate-messages
migrate-messages: ## Apply only the messages timestamp fix
	@echo "Applying messages timestamp fix..."
	@if [ -f "scripts/007_fix_messages.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/007_fix_messages.ps1; \
	else \
		echo "[ERROR] Messages fix script not found (scripts/007_fix_messages.ps1)"; \
		exit 1; \
	fi

.PHONY: migrate-profiles
migrate-profiles: ## Apply only the student profiles rename
	@echo "Applying student profiles rename..."
	@if [ -f "scripts/008_rename_profiles.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/008_rename_profiles.ps1; \
	else \
		echo "[ERROR] Profiles rename script not found (scripts/008_rename_profiles.ps1)"; \
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

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

.PHONY: clean-all
clean-all: clean ## Clean everything including go mod cache
	@echo "Cleaning everything..."
	@go clean -modcache
	@rm -rf uploads/resumes/*
	@rm -rf uploads/certificates/*

# Dependencies targets
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

.PHONY: deps-check
deps-check: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	@go list -u -m all

# Docker targets (if using Docker)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p $(PORT):$(PORT) $(APP_NAME)

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	@docker stop $(APP_NAME) || true

# Utility targets
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully!"

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)
	@echo "All checks completed!"

.PHONY: prepare
prepare: deps check ## Prepare for development (deps + checks)
	@echo "Project prepared for development!"

# Production targets
.PHONY: prod-build
prod-build: build-linux ## Build for production (Linux)
	@echo "Production build complete!"

.PHONY: prod-run
prod-run: prod-build ## Build and run for production
	@echo "Starting production server..."
	@./$(BUILD_DIR)/$(BINARY_NAME)-linux

# Default make target
.PHONY: make
make: build ## Default make target (alias for build)
	@echo "Build complete!"

# Show current status
.PHONY: status
status: ## Show current project status
	@echo "=== ASA Backend Status ==="
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build directory: $(BUILD_DIR)"
	@echo "Main file: $(MAIN_PATH)"
	@echo "================================" 