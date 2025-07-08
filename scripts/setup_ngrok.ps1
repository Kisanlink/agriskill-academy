# Setup script for ngrok environment
param(
    [Parameter(Mandatory=$true)]
    [string]$NgrokUrl
)

Write-Host "Setting up ngrok environment..." -ForegroundColor Green

# Set the BASE_URL environment variable
$env:BASE_URL = $NgrokUrl

Write-Host "Environment variables set:" -ForegroundColor Yellow
Write-Host "BASE_URL = $env:BASE_URL" -ForegroundColor Cyan

Write-Host "`nTo start your server with these settings:" -ForegroundColor Green
Write-Host "1. Run: go run cmd/server/main.go" -ForegroundColor White
Write-Host "2. Or set the environment variable permanently:" -ForegroundColor White
Write-Host "   [Environment]::SetEnvironmentVariable('BASE_URL', '$NgrokUrl', 'User')" -ForegroundColor Gray

Write-Host "`nTesting the setup..." -ForegroundColor Yellow
try {
    $testResponse = Invoke-WebRequest -Uri "$NgrokUrl/api/auth/login" -Method OPTIONS -UseBasicParsing
    Write-Host "CORS test successful!" -ForegroundColor Green
    Write-Host "Access-Control-Allow-Origin: $($testResponse.Headers['Access-Control-Allow-Origin'])" -ForegroundColor Cyan
} catch {
    Write-Host "CORS test failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Make sure your server is running and accessible via ngrok" -ForegroundColor Yellow
} 