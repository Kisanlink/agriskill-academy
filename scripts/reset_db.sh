# File: scripts/reset_db.sh
#!/bin/bash

DB_URL="postgresql://postgres:12345678@localhost:5432/jb?sslmode=disable"

echo "Dropping and recreating public schema..."
psql "$DB_URL" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

echo "Applying migrations..."
for file in $(ls ../migrations/*.sql | sort); do
  echo "Applying $file"
  psql "$DB_URL" -f "$file"
done

echo "Database reset and migrations applied."
