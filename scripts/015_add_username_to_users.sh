#!/bin/bash

# Migration script: Add username column to users table
# This script adds a username field to separate username from email

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-asa_db}
DB_USER=${POSTGRES_USER:-postgres}
DB_PASS=${POSTGRES_PASS:-password}

echo "🔧 Applying migration: Add username column to users table"
echo "📊 Database: $DB_NAME"
echo "🏠 Host: $DB_HOST:$DB_PORT"
echo "👤 User: $DB_USER"

# Check if migration has already been applied
MIGRATION_CHECK=$(PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' 
        AND column_name = 'username'
    );
" 2>/dev/null | xargs)

if [ "$MIGRATION_CHECK" = "t" ]; then
    echo "✅ Migration already applied - username column exists"
    exit 0
fi

echo "🔄 Applying migration..."

# Apply the migration
PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << EOF
-- Migration: Add username column to users table
-- This migration adds a username field to separate username from email

-- Add username column to users table
ALTER TABLE public.users 
ADD COLUMN IF NOT EXISTS username text COLLATE pg_catalog."default";

-- Add unique constraint on username
ALTER TABLE public.users 
ADD CONSTRAINT IF NOT EXISTS users_username_key UNIQUE (username);

-- Add index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users (username);

-- Update existing users to have a username based on their email (temporary)
-- This will be handled by the application logic during login/signup
UPDATE public.users 
SET username = email 
WHERE username IS NULL;

-- Make username NOT NULL after setting default values
ALTER TABLE public.users 
ALTER COLUMN username SET NOT NULL;
EOF

if [ $? -eq 0 ]; then
    echo "✅ Migration applied successfully!"
    echo "📋 Summary:"
    echo "   - Added username column to users table"
    echo "   - Added unique constraint on username"
    echo "   - Added index for faster lookups"
    echo "   - Set existing users' username to their email"
    echo "   - Made username NOT NULL"
else
    echo "❌ Migration failed!"
    exit 1
fi
