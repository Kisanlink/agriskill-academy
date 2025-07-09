#!/bin/bash

# Database Schema Migration Script for ASA Backend
# This script applies the complete database schema (001_complete_database_schema.sql)
# Reads database configuration from .env file

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

echo "=== ASA Backend Schema Migration Script ==="

# Load environment variables from .env file
if [ -f "../.env" ]; then
    export $(grep -v '^#' ../.env | xargs)
    print_status "Loaded configuration from .env file"
else
    print_warning ".env file not found, using default values"
fi

# Get database configuration from environment variables
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-asa}
DB_USER=${POSTGRESS_USER:-postgres}
DB_PASSWORD=${POSTGRESS_PASS:-password}

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo "Applying complete database schema..."

# Function to execute SQL file
execute_sql_file() {
    local sql_file="$1"
    local description="$2"
    
    echo "Applying: $description"
    
    if [ -f "$sql_file" ]; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$sql_file"
        
        if [ $? -eq 0 ]; then
            print_status "Successfully applied: $description"
        else
            print_error "Failed to apply: $description"
            return 1
        fi
    else
        print_error "SQL file not found: $sql_file"
        return 1
    fi
}

# Check if psql is available
if ! command -v psql &> /dev/null; then
    print_error "PostgreSQL client (psql) not found. Please install PostgreSQL client tools."
    echo "Download from: https://www.postgresql.org/download/"
    exit 1
fi
print_status "PostgreSQL client (psql) found"

# Test database connection
echo "Testing database connection..."
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" >/dev/null 2>&1; then
    print_status "Database connection successful"
else
    print_error "Database connection failed"
    echo "Please check your database credentials in .env file and ensure the database exists."
    exit 1
fi

# Apply the complete database schema
echo ""
echo "Applying complete database schema..."
execute_sql_file "../migrations/001_complete_database_schema.sql" "Complete Database Schema"

echo ""
echo "=== Schema Migration Summary ==="
print_status "Complete database schema applied successfully!"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"

echo ""
echo "Next steps:"
echo "1. Verify the database schema: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt'"
echo "2. Check tables created: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\d user_profiles'"
echo "3. Run next migration if needed: ./scripts/007_fix_messages.sh"

echo ""
print_status "Schema migration script completed!" 