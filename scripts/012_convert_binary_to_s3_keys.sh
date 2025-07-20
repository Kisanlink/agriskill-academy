#!/bin/bash

# Migration: Convert binary storage to S3 key storage
# This script converts BYTEA columns to TEXT columns for S3 keys

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${POSTGRES_USER:-postgres}
DB_PASSWORD=${POSTGRES_PASS:-password}
DB_NAME=${DB_NAME:-asa_db}

print_status "=== ASA Backend - Binary to S3 Key Storage Migration ==="
print_info "Database: $DB_NAME"
print_info "Host: $DB_HOST:$DB_PORT"
print_info "User: $DB_USER"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    print_error "PostgreSQL client (psql) is not installed or not in PATH"
    exit 1
fi

# Test database connection
print_info "Testing database connection..."
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    print_error "Failed to connect to database. Please check your connection parameters."
    exit 1
fi
print_status "Database connection successful"

# Check if migration has already been applied
print_info "Checking if migration has already been applied..."
MIGRATION_CHECK=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.columns 
WHERE table_name = 'users' AND column_name = 'avatar_key';")

if [ "$MIGRATION_CHECK" -gt 0 ]; then
    print_warning "Migration appears to have already been applied (avatar_key column exists)"
    read -p "Do you want to continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Migration cancelled"
        exit 0
    fi
fi

print_status "Starting binary to S3 key storage migration..."

# Run the migration
print_info "Converting binary columns to S3 key columns..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/012_convert_binary_to_s3_keys.sql

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
WHERE column_name IN ('avatar_key', 'logo_key', 'profile_photo_key', 'resume_key', 'file_key')
ORDER BY table_name, column_name;")

if [ -n "$VERIFICATION" ]; then
    print_status "Migration verification successful - S3 key columns found:"
    echo "$VERIFICATION"
else
    print_error "Migration verification failed - no S3 key columns found"
    exit 1
fi

print_status "=== Migration Summary ==="
echo "  ✓ Converted BYTEA columns to TEXT columns for S3 key storage"
echo "  ✓ Updated users.avatar → users.avatar_key"
echo "  ✓ Updated employer_profiles.logo → employer_profiles.logo_key"
echo "  ✓ Updated applications.resume_file → applications.resume_key"
echo "  ✓ Updated certificates.file → certificates.file_key"
echo "  ✓ Removed binary data columns (profile_photo, resume from student_profiles)"
echo ""
print_warning "Next steps:"
echo "  1. Update your application code to use S3 storage"
echo "  2. Configure S3 bucket and credentials"
echo "  3. Test file upload and download functionality"
echo "  4. Re-upload any existing files to S3" 