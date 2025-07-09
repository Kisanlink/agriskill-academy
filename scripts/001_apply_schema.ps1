# Database Schema Migration Script for ASA Backend
# This script applies the complete database schema (001_complete_database_schema.sql)
# Reads database configuration from .env file

# Load environment variables from .env file
if (Test-Path "../.env") {
    Get-Content "../.env" | ForEach-Object {
        if ($_ -match "^([^#][^=]+)=(.*)$") {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
    Write-Host "✓ Loaded configuration from .env file" -ForegroundColor Green
} else {
    Write-Host "⚠ .env file not found, using default values" -ForegroundColor Yellow
}

# Get database configuration from environment variables
$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_NAME = $env:DB_NAME
$DB_USER = $env:POSTGRESS_USER
$DB_PASSWORD = $env:POSTGRESS_PASS

Write-Host "=== ASA Backend Schema Migration Script ===" -ForegroundColor Green
Write-Host "Database: $DB_NAME" -ForegroundColor Cyan
Write-Host "Host: $DB_HOST:$DB_PORT" -ForegroundColor Cyan
Write-Host "User: $DB_USER" -ForegroundColor Cyan
Write-Host "Applying complete database schema..." -ForegroundColor Yellow

# Function to execute SQL file
function Execute-SqlFile {
    param(
        [string]$SqlFile,
        [string]$Description
    )
    
    Write-Host "Applying: $Description" -ForegroundColor Cyan
    
    if (Test-Path $SqlFile) {
        try {
            # Use psql to execute the SQL file
            $env:POSTGRESS_PASS = $DB_PASSWORD
            psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $SqlFile
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✓ Successfully applied: $Description" -ForegroundColor Green
            } else {
                Write-Host "✗ Failed to apply: $Description" -ForegroundColor Red
                exit 1
            }
        } catch {
            Write-Host "✗ Error executing $SqlFile : $_" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "✗ SQL file not found: $SqlFile" -ForegroundColor Red
        exit 1
    }
}

# Check if psql is available
try {
    $null = Get-Command psql -ErrorAction Stop
    Write-Host "✓ PostgreSQL client (psql) found" -ForegroundColor Green
} catch {
    Write-Host "✗ PostgreSQL client (psql) not found. Please install PostgreSQL client tools." -ForegroundColor Red
    Write-Host "Download from: https://www.postgresql.org/download/" -ForegroundColor Yellow
    exit 1
}

# Test database connection
Write-Host "Testing database connection..." -ForegroundColor Yellow
try {
    $env:POSTGRESS_PASS = $DB_PASSWORD
    $testResult = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Database connection successful" -ForegroundColor Green
    } else {
        Write-Host "✗ Database connection failed: $testResult" -ForegroundColor Red
        Write-Host "Please check your database credentials in .env file and ensure the database exists." -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host "✗ Database connection failed: $_" -ForegroundColor Red
    exit 1
}

# Apply the complete database schema
Write-Host "`nApplying complete database schema..." -ForegroundColor Yellow
Execute-SqlFile -SqlFile "../migrations/001_complete_database_schema.sql" -Description "Complete Database Schema"

Write-Host "`n=== Schema Migration Summary ===" -ForegroundColor Green
Write-Host "✓ Complete database schema applied successfully!" -ForegroundColor Green
Write-Host "Database: $DB_NAME" -ForegroundColor Cyan
Write-Host "Host: $DB_HOST:$DB_PORT" -ForegroundColor Cyan
Write-Host "User: $DB_USER" -ForegroundColor Cyan

Write-Host "`nNext steps:" -ForegroundColor Yellow
Write-Host "1. Verify the database schema: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt'" -ForegroundColor White
Write-Host "2. Check tables created: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\d user_profiles'" -ForegroundColor White
Write-Host "3. Run next migration if needed: ./scripts/007_fix_messages.ps1" -ForegroundColor White

Write-Host "`nSchema migration script completed!" -ForegroundColor Green 