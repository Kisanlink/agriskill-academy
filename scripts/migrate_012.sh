#!/bin/bash

# Migration Script 012: Create job_hires table
# This script creates the job_hires table for tracking multiple hired candidates per job

set -e  # Exit immediately if a command exits with a non-zero status

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Check if required environment variables are set
if [ -z "$DB_HOST" ] || [ -z "$POSTGRES_USER" ] || [ -z "$DB_NAME" ]; then
    echo "Error: Required environment variables are not set"
    echo "Please ensure DB_HOST, POSTGRES_USER, and DB_NAME are set in .env file"
    exit 1
fi

echo "Starting Migration 012: Create job_hires table..."
echo "Database: $DB_NAME"
echo "Host: $DB_HOST"

# Set PGPASSWORD for non-interactive authentication
export PGPASSWORD="$POSTGRES_PASS"

# Run the migration
psql -h "$DB_HOST" -p "${DB_PORT:-5432}" -U "$POSTGRES_USER" -d "$DB_NAME" -f migrations/012_create_job_hires.sql

# Check if migration was successful
if [ $? -eq 0 ]; then
    echo "Migration 012 completed successfully!"
else
    echo "Migration 012 failed!"
    exit 1
fi

# Unset password for security
unset PGPASSWORD

echo "Migration 012 finished."
