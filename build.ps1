# StegGo Build Script for Windows
# Created by pranaykumar2
# Updated: 2025-03-24 03:08:08 (UTC)

# Colors for output
$Green = [System.ConsoleColor]::Green
$Blue = [System.ConsoleColor]::Cyan
$Red = [System.ConsoleColor]::Red
$Yellow = [System.ConsoleColor]::Yellow

# Get command line args
param(
    [switch]$StartApi = $false
)

Write-Host "Building StegGo - Image Steganography Tool (CLI & API)" -ForegroundColor $Blue
Write-Host "===============================================" -ForegroundColor $Blue
Write-Host "Current Date: 2025-03-24 03:08:08 (UTC)" -ForegroundColor $Yellow
Write-Host "User: pranaykumar2" -ForegroundColor $Yellow
Write-Host ""

# Create necessary directories
Write-Host "Creating required directories..." -ForegroundColor $Blue
if (-not (Test-Path -Path "uploads")) { New-Item -ItemType Directory -Path "uploads" -Force | Out-Null }
if (-not (Test-Path -Path "temp")) { New-Item -ItemType Directory -Path "temp" -Force | Out-Null }
if (-not (Test-Path -Path "test-images")) { New-Item -ItemType Directory -Path "test-images" -Force | Out-Null }

# Update dependencies
Write-Host "Tidying modules..." -ForegroundColor $Blue
go mod tidy

# Build CLI application
Write-Host "Building CLI application..." -ForegroundColor $Blue
go build -v -o stego.exe .\cmd\stego

if ($LASTEXITCODE -eq 0) {
    Write-Host "CLI build successful!" -ForegroundColor $Green
} else {
    Write-Host "CLI build failed!" -ForegroundColor $Red
    exit 1
}

# Build API server
Write-Host "Building API server..." -ForegroundColor $Blue
go build -v -o steggo-api.exe .\cmd\api

if ($LASTEXITCODE -eq 0) {
    Write-Host "API server build successful!" -ForegroundColor $Green
} else {
    Write-Host "API server build failed!" -ForegroundColor $Red
    exit 1
}

Write-Host ""
Write-Host "=====================================================" -ForegroundColor $Green
Write-Host "Build completed successfully!" -ForegroundColor $Green
Write-Host "=====================================================" -ForegroundColor $Green
Write-Host ""
Write-Host "You can now run:" -ForegroundColor $Blue
Write-Host "  .\stego.exe               - For the CLI application" -ForegroundColor $Yellow
Write-Host "  .\steggo-api.exe          - For the API server" -ForegroundColor $Yellow
Write-Host ""
Write-Host "API Documentation:" -ForegroundColor $Blue
Write-Host "  Start the API server and visit: http://localhost:8080/swagger/index.html" -ForegroundColor $Yellow
Write-Host ""

# Check for test images
if (-not (Test-Path -Path ".\test-images\sample.png")) {
    Write-Host "Note: No test images found in .\test-images directory." -ForegroundColor $Yellow
    Write-Host "      Add test images for API testing." -ForegroundColor $Yellow
}

# Optionally start the API server
if ($StartApi) {
    Write-Host "Starting API server..." -ForegroundColor $Blue
    Write-Host "Press Ctrl+C to stop the server" -ForegroundColor $Blue
    Write-Host ""
    .\steggo-api.exe
}
