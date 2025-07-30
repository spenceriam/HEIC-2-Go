# Build script for Windows

# Set error action preference
$ErrorActionPreference = "Stop"

# Create bin directory if it doesn't exist
if (-not (Test-Path -Path ".\bin")) {
    New-Item -ItemType Directory -Path ".\bin" | Out-Null
}

# Build the application
Write-Host "Building HEIC-2-Go..." -ForegroundColor Cyan
go build -o ".\bin\heic2go.exe" .\cmd\heic2go

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Executable created at: $(Resolve-Path ".\bin\heic2go.exe")" -ForegroundColor Green
} else {
    Write-Host "Build failed with exit code $LASTEXITCODE" -ForegroundColor Red
    exit $LASTEXITCODE
}
