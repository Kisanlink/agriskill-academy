#!/bin/bash

# Script: Remove all binary storage and ensure S3 key storage only
# This script runs the migration to completely remove binary storage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_status() {
    echo -e "${GREEN}[STATUS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo ""
echo "=========================================="
echo "  Remove Binary Storage Migration Script"
echo "=========================================="
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    print_error ".env file not found"
    print_error "Please create a .env file with your database configuration"
    exit 1
fi

# Load environment variables
print_info "Loading environment variables..."
source .env

# Check required environment variables
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    print_error "Missing required database environment variables"
    print_error "Please check your .env file for: DB_HOST, DB_PORT, DB_USER, DB_NAME"
    exit 1
fi

# Set default password if not provided
if [ -z "$DB_PASSWORD" ]; then
    print_warning "DB_PASSWORD not set, attempting connection without password"
    DB_PASSWORD=""
fi

print_info "Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# Test database connection
print_info "Testing database connection..."
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    print_error "Failed to connect to database"
    print_error "Please check your database configuration in .env file"
    exit 1
fi
print_status "Database connection successful"

echo ""
echo "Starting binary storage removal migration..."

# Check current state of binary columns
print_info "Checking current state of file storage columns..."
BINARY_CHECK=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT 
    table_name,
    column_name,
    data_type
FROM information_schema.columns 
WHERE table_schema = 'public' 
AND column_name IN ('avatar', 'logo', 'profile_photo', 'resume', 'resume_file', 'file')
AND data_type = 'bytea'
ORDER BY table_name, column_name;
")

if [ -n "$BINARY_CHECK" ]; then
    print_warning "Found binary columns that will be removed:"
    echo "$BINARY_CHECK"
    echo ""
    read -p "Do you want to continue with removing binary storage? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Migration cancelled"
        exit 0
    fi
else
    print_status "No binary columns found - database already using S3 keys"
fi

# Run the migration
print_info "Removing binary storage and ensuring S3 key storage..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/013_remove_all_binary_storage.sql

if [ $? -eq 0 ]; then
    print_status "Migration completed successfully!"
else
    print_error "Migration failed!"
    exit 1
fi

# Verify the migration
print_info "Verifying migration..."
VERIFICATION=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT 
    table_name,
    column_name,
    data_type
FROM information_schema.columns 
WHERE table_schema = 'public' 
AND column_name IN ('avatar_key', 'logo_key', 'profile_photo_key', 'resume_key', 'file_key')
ORDER BY table_name, column_name;")

if [ -n "$VERIFICATION" ]; then
    print_status "Migration verification successful - S3 key columns found:"
    echo "$VERIFICATION"
else
    print_error "Migration verification failed - no S3 key columns found"
    exit 1
fi

# Check that binary columns are gone
BINARY_REMAINING=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT 
    table_name,
    column_name,
    data_type
FROM information_schema.columns 
WHERE table_schema = 'public' 
AND column_name IN ('avatar', 'logo', 'profile_photo', 'resume', 'resume_file', 'file')
AND data_type = 'bytea'
ORDER BY table_name, column_name;")

if [ -n "$BINARY_REMAINING" ]; then
    print_warning "Some binary columns still exist:"
    echo "$BINARY_REMAINING"
    print_warning "You may need to manually remove these columns"
else
    print_status "All binary columns successfully removed"
fi

print_status "=== Migration Summary ==="
echo "  ✓ Removed all BYTEA columns"
echo "  ✓ Ensured S3 key columns exist"
echo "  ✓ Created indexes for S3 key columns"
echo "  ✓ Added documentation comments"
echo ""
print_warning "Next steps:"
echo "  1. Update your application code to use S3 storage only"
echo "  2. Remove any binary data handling from handlers and services"
echo "  3. Test file upload and download functionality"
echo "  4. Update any frontend code to handle S3 URLs"
echo ""
print_info "Migration completed successfully!" 