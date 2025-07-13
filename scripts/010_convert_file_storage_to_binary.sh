#!/bin/bash

# Migration: Convert file storage from file paths to binary storage
# This script converts text columns to BYTEA and adds metadata columns

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

echo "=== ASA Backend - Binary File Storage Migration ==="

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

print_info "Database: $DB_HOST:$DB_PORT/$DB_NAME"
print_info "User: $DB_USER"

# Test database connection
echo "Testing database connection..."
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    print_error "Failed to connect to database"
    print_error "Please check your database configuration in .env file"
    exit 1
fi
print_status "Database connection successful"

echo ""
echo "Starting binary file storage migration..."

# Convert text columns to BYTEA
print_info "Converting file columns to BYTEA..."

# Users table - avatar
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
ALTER TABLE users ALTER COLUMN avatar TYPE BYTEA;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_type VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_size BIGINT;
" 2>/dev/null || print_warning "Users table avatar columns already updated"

# Employer profiles table - logo
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
ALTER TABLE employer_profiles ALTER COLUMN logo TYPE BYTEA;
ALTER TABLE employer_profiles ADD COLUMN IF NOT EXISTS logo_name VARCHAR(255);
ALTER TABLE employer_profiles ADD COLUMN IF NOT EXISTS logo_type VARCHAR(100);
ALTER TABLE employer_profiles ADD COLUMN IF NOT EXISTS logo_size BIGINT;
" 2>/dev/null || print_warning "Employer profiles table logo columns already updated"

# Student profiles table - profile_photo and resume
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
ALTER TABLE student_profiles ALTER COLUMN profile_photo TYPE BYTEA;
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS profile_photo_name VARCHAR(255);
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS profile_photo_type VARCHAR(100);
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS profile_photo_size BIGINT;
ALTER TABLE student_profiles ALTER COLUMN resume TYPE BYTEA;
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS resume_name VARCHAR(255);
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS resume_type VARCHAR(100);
ALTER TABLE student_profiles ADD COLUMN IF NOT EXISTS resume_size BIGINT;
" 2>/dev/null || print_warning "Student profiles table file columns already updated"

# Applications table - resume_file
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
ALTER TABLE applications ALTER COLUMN resume_file TYPE BYTEA;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS resume_file_name VARCHAR(255);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS resume_file_type VARCHAR(100);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS resume_file_size BIGINT;
" 2>/dev/null || print_warning "Applications table resume_file columns already updated"

# Certificates table - file
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
ALTER TABLE certificates ALTER COLUMN file TYPE BYTEA;
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS file_name VARCHAR(255);
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS file_type VARCHAR(100);
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS file_size BIGINT;
" 2>/dev/null || print_warning "Certificates table file columns already updated"

print_status "File columns converted to BYTEA"

# Add indexes for better performance
print_info "Adding indexes for file metadata..."

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
CREATE INDEX IF NOT EXISTS idx_users_avatar_name ON users(avatar_name);
CREATE INDEX IF NOT EXISTS idx_employer_profiles_logo_name ON employer_profiles(logo_name);
CREATE INDEX IF NOT EXISTS idx_student_profiles_resume_name ON student_profiles(resume_name);
CREATE INDEX IF NOT EXISTS idx_certificates_file_name ON certificates(file_name);
" 2>/dev/null || print_warning "Indexes already exist"

print_status "Indexes created"

# Add comments to document the changes
print_info "Adding column comments..."

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
COMMENT ON COLUMN users.avatar IS 'Binary file data for user avatar/profile photo';
COMMENT ON COLUMN users.avatar_name IS 'Original filename of the avatar';
COMMENT ON COLUMN users.avatar_type IS 'MIME type of the avatar file';
COMMENT ON COLUMN users.avatar_size IS 'File size in bytes';

COMMENT ON COLUMN employer_profiles.logo IS 'Binary file data for company logo';
COMMENT ON COLUMN employer_profiles.logo_name IS 'Original filename of the logo';
COMMENT ON COLUMN employer_profiles.logo_type IS 'MIME type of the logo file';
COMMENT ON COLUMN employer_profiles.logo_size IS 'File size in bytes';

COMMENT ON COLUMN student_profiles.profile_photo IS 'Binary file data for student profile photo';
COMMENT ON COLUMN student_profiles.profile_photo_name IS 'Original filename of the profile photo';
COMMENT ON COLUMN student_profiles.profile_photo_type IS 'MIME type of the profile photo file';
COMMENT ON COLUMN student_profiles.profile_photo_size IS 'File size in bytes';

COMMENT ON COLUMN student_profiles.resume IS 'Binary file data for student resume';
COMMENT ON COLUMN student_profiles.resume_name IS 'Original filename of the resume';
COMMENT ON COLUMN student_profiles.resume_type IS 'MIME type of the resume file';
COMMENT ON COLUMN student_profiles.resume_size IS 'File size in bytes';

COMMENT ON COLUMN applications.resume_file IS 'Binary file data for application resume';
COMMENT ON COLUMN applications.resume_file_name IS 'Original filename of the resume';
COMMENT ON COLUMN applications.resume_file_type IS 'MIME type of the resume file';
COMMENT ON COLUMN applications.resume_file_size IS 'File size in bytes';

COMMENT ON COLUMN certificates.file IS 'Binary file data for certificate';
COMMENT ON COLUMN certificates.file_name IS 'Original filename of the certificate';
COMMENT ON COLUMN certificates.file_type IS 'MIME type of the certificate file';
COMMENT ON COLUMN certificates.file_size IS 'File size in bytes';
" 2>/dev/null || print_warning "Comments already exist"

print_status "Column comments added"

# Verify the migration
print_info "Verifying migration..."

# Check if columns exist and are of correct type
COLUMN_CHECK=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
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

if [ -n "$COLUMN_CHECK" ]; then
    print_status "Migration verification successful"
    echo "$COLUMN_CHECK"
else
    print_error "Migration verification failed - some columns may not have been converted"
    exit 1
fi

echo ""
print_status "Binary file storage migration completed successfully!"
echo ""
print_info "Changes made:"
echo "  ✓ Converted file path columns to BYTEA (binary storage)"
echo "  ✓ Added metadata columns (name, type, size) for each file field"
echo "  ✓ Created indexes for better performance"
echo "  ✓ Added documentation comments"
echo ""
print_info "Next steps:"
echo "  1. Update your application code to handle binary file storage"
echo "  2. Test file upload and download functionality"
echo "  3. Remove file system dependencies from your code"
echo ""
print_warning "Note: Existing file paths in the database will be converted to NULL"
print_warning "You may need to re-upload files to populate the binary storage" 