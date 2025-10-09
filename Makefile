# ASA Backend Makefile
# A comprehensive build and development tool for the ASA job portal backend

.PHONY: help build run air test clean migrate migrate-reset setup

# Default target
help:
	@echo "ASA Backend - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run          - Build and run the application"
	@echo "  make air          - Run with hot reload (requires air)"
	@echo "  make build        - Build the application"
	@echo ""
	@echo "Database:"
	@echo "  make migrate      - Apply all database migrations"
	@echo "  make migrate-reset- Reset database and apply all migrations"
	@echo "  make debug-migration - Debug migration issues"
	@echo ""
	@echo "Utilities:"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make setup        - Initial setup (install air, setup uploads)"
	@echo "  make help         - Show this help message"

# Build the application
build:
	@echo "Building ASA Backend..."
	go build -o bin/asa cmd/server/main.go
	@echo "Build complete! Binary: bin/asa"

# Run the application
run: build
	@echo "Starting ASA Backend..."
	./bin/asa

# Run with hot reload (requires air)
air:
	@echo "Starting ASA Backend with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete!"

# Database migrations
migrate:
	@echo "Applying database migrations..."
	@if [ -f "scripts/001_apply_schema.sh" ]; then \
		chmod +x scripts/001_apply_schema.sh; \
		./scripts/001_apply_schema.sh; \
	fi
	@if [ -f "scripts/007_fix_messages.sh" ]; then \
		chmod +x scripts/007_fix_messages.sh; \
		./scripts/007_fix_messages.sh; \
	fi
	@if [ -f "scripts/008_rename_profiles.sh" ]; then \
		chmod +x scripts/008_rename_profiles.sh; \
		./scripts/008_rename_profiles.sh; \
	fi
	@if [ -f "scripts/009_add_phone_to_student_profiles.sh" ]; then \
		chmod +x scripts/009_add_phone_to_student_profiles.sh; \
		./scripts/009_add_phone_to_student_profiles.sh; \
	fi
	@if [ -f "scripts/010_convert_file_storage_to_binary.sh" ]; then \
		chmod +x scripts/010_convert_file_storage_to_binary.sh; \
		./scripts/010_convert_file_storage_to_binary.sh; \
	fi
	@if [ -f "scripts/012_convert_binary_to_s3_keys.sh" ]; then \
		chmod +x scripts/012_convert_binary_to_s3_keys.sh; \
		./scripts/012_convert_binary_to_s3_keys.sh; \
	fi
	@if [ -f "scripts/014_add_contact_requests.sh" ]; then \
		chmod +x scripts/014_add_contact_requests.sh; \
		./scripts/014_add_contact_requests.sh; \
	fi
	@if [ -f "scripts/015_add_username_to_users.sh" ]; then \
		chmod +x scripts/015_add_username_to_users.sh; \
		./scripts/015_add_username_to_users.sh; \
	fi
	@echo "Migrations applied successfully!"

# Debug migration issues
debug-migration:
	@echo "Running migration debug script..."
	@if [ -f "scripts/debug_migration.sh" ]; then \
		chmod +x scripts/debug_migration.sh; \
		./scripts/debug_migration.sh; \
	fi

# Reset database and apply all migrations
migrate-reset:
	@echo "Resetting database and applying all migrations..."
	@if [ -f "scripts/reset_db.sh" ]; then \
		chmod +x scripts/reset_db.sh; \
		./scripts/reset_db.sh; \
	fi
	@echo "Database reset and migrations applied successfully!"

# Initial setup
setup:
	@echo "Setting up ASA Backend development environment..."
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Installing air for hot reload..."
	go install github.com/cosmtrek/air@latest
	@echo "Setting up uploads directory..."
	@if [ -f "scripts/setup_uploads.sh" ]; then \
		chmod +x scripts/setup_uploads.sh; \
		./scripts/setup_uploads.sh; \
	fi
	@echo "Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Create a .env file with your database configuration"
	@echo "2. Run 'make migrate' to set up the database"
	@echo "3. Run 'make run' to start the application"
	@echo "4. Or run 'make air' for development with hot reload"

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully!"

# Build with dependencies
build: deps
	@echo "Building ASA Backend..."
	go build -o bin/asa cmd/server/main.go
	@echo "Build complete! Binary: bin/asa" 
