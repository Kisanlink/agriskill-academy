#!/bin/bash

# Student Profiles Rename Migration Script for ASA Backend
# This script applies the student profiles rename (008_rename_user_profiles_to_student_profiles.sql)
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

echo "=== ASA Backend Student Profiles Rename Migration Script ==="

# Load environment variables from .env file
if [ -f ".env" ]; then
    export $(grep -v '^#' .env | xargs)
    print_status "Loaded configuration from .env file"
elif [ -f "../.env" ]; then
    export $(grep -v '^#' ../.env | xargs)
    print_status "Loaded configuration from .env file"
else
    print_warning ".env file not found, using default values"
fi

# Get database configuration from environment variables
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_NAME=${DB_NAME}
DB_USER=${POSTGRESS_USER}
DB_PASSWORD=${POSTGRESS_PASS}

# Check if all required environment variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ]; then
    print_error "Missing required database configuration in .env file"
    echo "Required variables: DB_HOST, DB_PORT, DB_NAME, POSTGRESS_USER, POSTGRESS_PASS"
    echo "Please create a .env file with all required database configuration."
    exit 1
fi

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo "Applying student profiles rename..."

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

# Apply the student profiles rename
echo ""
echo "Applying student profiles rename..."
execute_sql_file "../migrations/008_rename_user_profiles_to_student_profiles.sql" "Rename User Profiles to Student Profiles"

echo ""
echo "=== Student Profiles Rename Migration Summary ==="
print_status "Student profiles rename applied successfully!"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"

echo ""
echo "Next steps:"
echo "1. Verify the student_profiles table: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\d student_profiles'"
echo "2. Check that user_profiles table is renamed: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt' | grep profiles"
echo "3. Verify foreign key constraints: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c 'SELECT conname, conrelid::regclass, confrelid::regclass FROM pg_constraint WHERE conname LIKE '\''%student_profiles%'\'';'"

echo ""
print_status "Student profiles rename migration script completed!" 