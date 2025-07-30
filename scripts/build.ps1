# Build script for HEIC-2-Go (Windows)
# Builds HEIC-2-Go for Windows with version information

param(
    [string]$Version = "0.1.0",
    [string]$Platform = "windows",
    [string]$Arch = "amd64"
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Set build variables
$APP_NAME = "heic2go"
$BUILD_DIR = "..\bin"
$OUTPUT_DIR = "$BUILD_DIR\$Platform-$Arch"
$OUTPUT_FILE = "$OUTPUT_DIR\${APP_NAME}.exe"
$LDFLAGS = "-X 'main.Version=$Version' -X 'main.BuildTime=$(Get-Date -Format "2006-01-02T15:04:05Z07:00")'"

Write-Host "Building $APP_NAME v$Version for $Platform/$Arch..." -ForegroundColor Cyan

# Create output directory if it doesn't exist
if (-not (Test-Path -Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

# Set environment variables for the build
$env:GOOS = $Platform
$env:GOARCH = $Arch

# Build the application
$buildCmd = "go build -ldflags='$LDFLAGS' -o '$OUTPUT_FILE' ..\cmd\heic2go"
Write-Host "Executing: $buildCmd" -ForegroundColor DarkGray
Invoke-Expression $buildCmd

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    $fileSize = (Get-Item $OUTPUT_FILE).Length / 1MB
    Write-Host " Build successful!" -ForegroundColor Green
    Write-Host "   Output: $OUTPUT_FILE"
    Write-Host ("   Size: {0:N2} MB" -f $fileSize)
    
    # Create a zip archive
    $zipPath = "$BUILD_DIR\${APP_NAME}-${Version}-${Platform}-${Arch}.zip"
    Compress-Archive -Path $OUTPUT_FILE -DestinationPath $zipPath -Force
    Write-Host "   Archive: $zipPath"
} else {
    Write-Host " Build failed!" -ForegroundColor Red
    exit 1
}
