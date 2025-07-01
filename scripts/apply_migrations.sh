# File: scripts/apply_migrations.sh
#!/bin/bash

DB_URL="postgresql://postgres:12345678@localhost:5432/jb?sslmode=disable"

echo "Applying migrations..."
for file in $(ls ../migrations/*.sql | sort); do
  echo "Applying $file"
  psql "$DB_URL" -f "$file"
done

echo "All migrations applied."
