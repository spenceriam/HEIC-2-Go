#!/bin/bash

# Build script for HEIC-2-Go (Linux/macOS)
# Builds HEIC-2-Go for multiple platforms with version information

# Default values
VERSION=${1:-0.1.0}
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")
APP_NAME="heic2go"
BUILD_DIR="../bin"

# Clean and create build directory
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

echo "🚀 Building ${APP_NAME} v${VERSION}"
echo "========================================"

# Function to build for a specific platform
build_for_platform() {
    local os=$1
    local arch=$2
    local output_name="${APP_NAME}"
    
    # Add .exe extension for Windows
    if [[ "${os}" == "windows" ]]; then
        output_name+=".exe"
    fi
    
    local output_dir="${BUILD_DIR}/${os}-${arch}"
    local output_file="${output_dir}/${output_name}"
    
    # Set environment variables for the build
    export GOOS="${os}"
    export GOARCH="${arch}"
    
    # Create platform-specific output directory
    mkdir -p "${output_dir}"
    
    # Build the application with version information
    echo "🔨 Building for ${os}/${arch}..."
    
    go build \
        -ldflags "-X 'main.Version=${VERSION}' -X 'main.BuildTime=$(date -u +"%Y-%m-%dT%H:%M:%SZ")'" \
        -o "${output_file}" \
        "../cmd/heic2go"
    
    if [ $? -ne 0 ]; then
        echo "❌ Error building for ${os}/${arch}" >&2
        return 1
    fi
    
    # Make the binary executable (non-Windows)
    if [[ "${os}" != "windows" ]]; then
        chmod +x "${output_file}"
    fi
    
    # Create a tarball/zip archive
    echo "📦 Creating archive for ${os}/${arch}..."
    
    pushd "${output_dir}" > /dev/null || return 1
    
    if [[ "${os}" == "windows" ]]; then
        # Create zip for Windows
        zip -r "../${APP_NAME}-${VERSION}-${os}-${arch}.zip" "."
    else
        # Create tarball for Unix-like systems
        tar -czf "../${APP_NAME}-${VERSION}-${os}-${arch}.tar.gz" "."
    fi
    
    popd > /dev/null || return 1
    
    echo "✅ Successfully built for ${os}/${arch}"
    echo "   Binary: ${output_file}"
    
    if [[ "${os}" == "windows" ]]; then
        echo "   Archive: ${BUILD_DIR}/${APP_NAME}-${VERSION}-${os}-${arch}.zip"
    else
        echo "   Archive: ${BUILD_DIR}/${APP_NAME}-${VERSION}-${os}-${arch}.tar.gz"
    fi
    
    echo ""
}

# Build for all platforms
for platform in "${PLATFORMS[@]}"; do
    # Split platform into OS and architecture
    IFS='/' read -r -a parts <<< "${platform}"
    os="${parts[0]}"
    arch="${parts[1]}"
    
    # Build for this platform
    if ! build_for_platform "${os}" "${arch}"; then
        echo "❌ Build failed for ${os}/${arch}" >&2
        exit 1
    fi
done

echo "✨ All builds completed successfully!"
echo "📦 Output directory: ${BUILD_DIR}/"

# Create a checksums file
echo "🔍 Generating checksums..."
pushd "${BUILD_DIR}" > /dev/null || exit 1
find . -type f \( -name "*.tar.gz" -o -name "*.zip" \) -exec shasum -a 256 {} \; > "${APP_NAME}-${VERSION}-checksums.txt"
popd > /dev/null || exit 1

echo "✅ Build process completed successfully!"
exit 0
