#!/bin/bash

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

echo "Applying migration: Add contact_requests table..."

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

# Apply the migration
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/014_add_contact_requests_table.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration applied successfully"
else
    echo "❌ Migration failed"
    exit 1
fi 