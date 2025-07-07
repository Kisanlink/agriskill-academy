# Apply messages timestamp fix migration
# This script fixes the sent_at timestamp issue in the messages table

Write-Host "Applying messages timestamp fix migration..." -ForegroundColor Green

# Database connection parameters (update these as needed)
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_NAME = "agrijobs"
$DB_USER = "postgres"
$DB_PASSWORD = "your_password_here"

# Migration file path
$MIGRATION_FILE = "migrations/007_fix_messages_timestamp.sql"

# Check if migration file exists
if (-not (Test-Path $MIGRATION_FILE)) {
    Write-Host "Error: Migration file not found: $MIGRATION_FILE" -ForegroundColor Red
    exit 1
}

# Read the migration SQL
$SQL_CONTENT = Get-Content $MIGRATION_FILE -Raw

Write-Host "Migration SQL:" -ForegroundColor Yellow
Write-Host $SQL_CONTENT -ForegroundColor Gray

# Apply the migration using psql
try {
    $env:PGPASSWORD = $DB_PASSWORD
    $result = & psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c $SQL_CONTENT 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migration applied successfully!" -ForegroundColor Green
        Write-Host "Messages table has been recreated with proper timestamp handling." -ForegroundColor Green
    } else {
        Write-Host "Error applying migration:" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error executing migration: $_" -ForegroundColor Red
    exit 1
} finally {
    $env:PGPASSWORD = $null
}

Write-Host "Migration completed successfully!" -ForegroundColor Green
Write-Host "Please restart your backend server and test the messaging functionality." -ForegroundColor Yellow 