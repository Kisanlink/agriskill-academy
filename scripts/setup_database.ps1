# PowerShell script to setup database
$env:PGPASSWORD = "12345678"

Write-Host "Checking if agrijobs database exists..."
$dbExists = psql -h localhost -U postgres -t -c "SELECT 1 FROM pg_database WHERE datname='agrijobs';"

if ($dbExists -eq $null -or $dbExists.Trim() -eq "") {
    Write-Host "Creating agrijobs database..."
    psql -h localhost -U postgres -c "CREATE DATABASE agrijobs;"
    Write-Host "Database created successfully!"
} else {
    Write-Host "Database agrijobs already exists!"
}

$env:PGPASSWORD = "" 