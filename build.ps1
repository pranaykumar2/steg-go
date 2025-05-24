# StegGo Build Script for Windows
# Created by pranaykumar2
# Updated: 2025-03-31 14:07:35 (UTC)

param(
    [switch]$StartServer = $false
)

$Green = [System.ConsoleColor]::Green
$Blue = [System.ConsoleColor]::Cyan
$Red = [System.ConsoleColor]::Red
$Yellow = [System.ConsoleColor]::Yellow

Write-Host "Building StegGo - Image Steganography Tool (CLI & API/Web Server)" -ForegroundColor $Blue
Write-Host "===============================================" -ForegroundColor $Blue
Write-Host "Last Updated On: 2025-03-31 14:07:35 (UTC)" -ForegroundColor $Yellow
Write-Host "Created By: pranaykumar2" -ForegroundColor $Yellow
Write-Host ""

Write-Host "Creating required directories..." -ForegroundColor $Blue
if (-not (Test-Path -Path "uploads")) { New-Item -ItemType Directory -Path "uploads" -Force | Out-Null }
if (-not (Test-Path -Path "temp")) { New-Item -ItemType Directory -Path "temp" -Force | Out-Null }
if (-not (Test-Path -Path "test-images")) { New-Item -ItemType Directory -Path "test-images" -Force | Out-Null }
if (-not (Test-Path -Path "web\static\img")) { New-Item -ItemType Directory -Path "web\static\img" -Force | Out-Null }

Write-Host "Tidying modules..." -ForegroundColor $Blue
go mod tidy

Write-Host "Building CLI application..." -ForegroundColor $Blue
go build -v -o stego.exe .\cmd\stego

if ($LASTEXITCODE -eq 0) {
    Write-Host "CLI build successful!" -ForegroundColor $Green
} else {
    Write-Host "CLI build failed!" -ForegroundColor $Red
    exit 1
}

Write-Host "Building API/Web server..." -ForegroundColor $Blue
go build -v -o steggo-server.exe .\cmd\api

if ($LASTEXITCODE -eq 0) {
    Write-Host "API/Web server build successful!" -ForegroundColor $Green
} else {
    Write-Host "API/Web server build failed!" -ForegroundColor $Red
    exit 1
}

Write-Host ""
Write-Host "=====================================================" -ForegroundColor $Green
Write-Host "Build completed successfully!" -ForegroundColor $Green
Write-Host "=====================================================" -ForegroundColor $Green
Write-Host ""
Write-Host "You can now run:" -ForegroundColor $Blue
Write-Host "  .\stego.exe               - For the CLI application" -ForegroundColor $Yellow
Write-Host "  .\steggo-server.exe       - For the API/Web server" -ForegroundColor $Yellow
Write-Host ""
Write-Host "Web Interface:" -ForegroundColor $Blue
Write-Host "  Start the server and visit: http://localhost:8080" -ForegroundColor $Yellow
Write-Host ""
Write-Host "API Documentation:" -ForegroundColor $Blue
Write-Host "  Start the server and visit: http://localhost:8080/swagger/index.html" -ForegroundColor $Yellow
Write-Host ""

if (-not (Test-Path -Path ".\test-images\sample.jpg")) {
    Write-Host "Note: No test images found in .\test-images directory." -ForegroundColor $Yellow
    Write-Host "      Add test images for testing." -ForegroundColor $Yellow
}

if ($StartServer) {
    Write-Host "Starting API/Web server..." -ForegroundColor $Blue
    Write-Host "Press Ctrl+C to stop the server" -ForegroundColor $Blue
    Write-Host ""
    .\steggo-server.exe
}
