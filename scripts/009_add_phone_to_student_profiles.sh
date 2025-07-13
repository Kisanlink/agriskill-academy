#!/bin/bash

# Load environment variables
source .env

echo "Applying migration: Add phone_number to student_profiles..."

# Apply the migration
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/009_add_phone_to_student_profiles.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration applied successfully"
else
    echo "❌ Migration failed"
    exit 1
fi 