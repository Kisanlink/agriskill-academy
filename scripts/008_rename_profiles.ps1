# Student Profiles Rename Migration Script for ASA Backend
# This script applies the user_profiles to student_profiles rename (008_rename_user_profiles_to_student_profiles.sql)
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

# Set defaults if not provided
if (-not $DB_HOST) { $DB_HOST = "localhost" }
if (-not $DB_PORT) { $DB_PORT = "5432" }
if (-not $DB_NAME) { $DB_NAME = "asa" }
if (-not $DB_USER) { $DB_USER = "postgres" }
if (-not $DB_PASSWORD) { $DB_PASSWORD = "password" }

Write-Host "=== ASA Backend Student Profiles Rename Migration Script ===" -ForegroundColor Green
Write-Host "Database: $DB_NAME" -ForegroundColor Cyan
Write-Host "Host: $DB_HOST:$DB_PORT" -ForegroundColor Cyan
Write-Host "User: $DB_USER" -ForegroundColor Cyan
Write-Host "Renaming user_profiles to student_profiles..." -ForegroundColor Yellow

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

# Check if user_profiles table exists before renaming
Write-Host "Checking if user_profiles table exists..." -ForegroundColor Yellow
try {
    $env:POSTGRESS_PASS = $DB_PASSWORD
    $tableCheck = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_profiles');" 2>&1
    
    if ($tableCheck -match "t") {
        Write-Host "✓ user_profiles table found, proceeding with rename" -ForegroundColor Green
    } else {
        Write-Host "⚠ user_profiles table not found, checking if student_profiles already exists" -ForegroundColor Yellow
        
        $studentTableCheck = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'student_profiles');" 2>&1
        
        if ($studentTableCheck -match "t") {
            Write-Host "✓ student_profiles table already exists, rename already completed" -ForegroundColor Green
            Write-Host "Migration already applied successfully!" -ForegroundColor Green
            exit 0
        } else {
            Write-Host "✗ Neither user_profiles nor student_profiles table found" -ForegroundColor Red
            Write-Host "Please run the schema migration first: ./scripts/001_apply_schema.ps1" -ForegroundColor Yellow
            exit 1
        }
    }
} catch {
    Write-Host "✗ Error checking table existence: $_" -ForegroundColor Red
    exit 1
}

# Apply the student profiles rename
Write-Host "`nRenaming user_profiles to student_profiles..." -ForegroundColor Yellow
Execute-SqlFile -SqlFile "../migrations/008_rename_user_profiles_to_student_profiles.sql" -Description "Rename User Profiles to Student Profiles"

Write-Host "`n=== Student Profiles Rename Migration Summary ===" -ForegroundColor Green
Write-Host "✓ Student profiles rename applied successfully!" -ForegroundColor Green
Write-Host "Database: $DB_NAME" -ForegroundColor Cyan
Write-Host "Host: $DB_HOST:$DB_PORT" -ForegroundColor Cyan
Write-Host "User: $DB_USER" -ForegroundColor Cyan

Write-Host "`nNext steps:" -ForegroundColor Yellow
Write-Host "1. Verify the renamed table: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\d student_profiles'" -ForegroundColor White
Write-Host "2. Check certificates foreign key: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\d certificates'" -ForegroundColor White
Write-Host "3. Verify no user_profiles table remains: psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt user_profiles'" -ForegroundColor White

Write-Host "`nStudent profiles rename migration script completed!" -ForegroundColor Green 