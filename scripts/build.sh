#!/bin/bash

# Build script for Linux/macOS

# Exit on error
set -e

# Create bin directory if it doesn't exist
mkdir -p ./bin

# Build the application
echo "Building HEIC-2-Go..."
go build -o "./bin/heic2go" ./cmd/heic2go

# Make the binary executable
chmod +x "./bin/heic2go"

echo "Build successful!"
echo "Executable created at: $(pwd)/bin/heic2go"
