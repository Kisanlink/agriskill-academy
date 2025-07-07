# PowerShell script to apply all migrations
$env:PGPASSWORD = "12345678"

Write-Host "Applying initial migrations..."

# Apply 001_create_tables.sql
Write-Host "Applying 001_create_tables.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/001_create_tables.sql

# Apply 002_alter_employer_profile.sql
Write-Host "Applying 002_alter_employer_profile.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/002_alter_employer_profile.sql

# Apply 003_update_job_posts_salary_structure.sql
Write-Host "Applying 003_update_job_posts_salary_structure.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/003_update_job_posts_salary_structure.sql

# Apply 004_add_updated_at_to_applications.sql
Write-Host "Applying 004_add_updated_at_to_applications.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/004_add_updated_at_to_applications.sql

# Apply 005_create_notification_preferences.sql
Write-Host "Applying 005_create_notification_preferences.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/005_create_notification_preferences.sql

# Apply 006_create_job_alerts.sql
Write-Host "Applying 006_create_job_alerts.sql..."
psql -h localhost -U postgres -d agrijobs -f migrations/006_create_job_alerts.sql

Write-Host "All migrations applied successfully!"

$env:PGPASSWORD = "" 