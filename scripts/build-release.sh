#!/bin/bash

# Build script for Gorev v0.12.0 release
# This script builds binaries for all supported platforms

set -e

VERSION="v0.12.0"
BUILD_DIR="release-${VERSION}"

echo "🚀 Building Gorev ${VERSION} for all platforms..."

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Change to mcpserver directory
cd gorev-mcpserver

# Ensure dependencies are up to date
echo "📦 Updating dependencies..."
go mod download
go mod tidy

# Run tests first
echo "🧪 Running tests..."
go test -v -cover ./...

# Build for all platforms
echo "🔨 Building binaries..."

# Linux AMD64
echo "  → Linux AMD64"
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +\"%Y-%m-%dT%H:%M:%SZ\") -X main.gitCommit=$(git rev-parse --short HEAD)" -o ../${BUILD_DIR}/gorev-linux-amd64 ./cmd/gorev

# Darwin AMD64 (Intel Mac)
echo "  → Darwin AMD64"
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +\"%Y-%m-%dT%H:%M:%SZ\") -X main.gitCommit=$(git rev-parse --short HEAD)" -o ../${BUILD_DIR}/gorev-darwin-amd64 ./cmd/gorev

# Darwin ARM64 (M1/M2 Mac)
echo "  → Darwin ARM64"
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +\"%Y-%m-%dT%H:%M:%SZ\") -X main.gitCommit=$(git rev-parse --short HEAD)" -o ../${BUILD_DIR}/gorev-darwin-arm64 ./cmd/gorev

# Windows AMD64
echo "  → Windows AMD64"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +\"%Y-%m-%dT%H:%M:%SZ\") -X main.gitCommit=$(git rev-parse --short HEAD)" -o ../${BUILD_DIR}/gorev-windows-amd64.exe ./cmd/gorev

cd ..

# Create checksums
echo "🔒 Creating checksums..."
cd ${BUILD_DIR}
sha256sum gorev-* > checksums.txt
cd ..

# Create archives
echo "📦 Creating archives..."
cd ${BUILD_DIR}

# Create tar.gz for Unix systems
tar -czf gorev-${VERSION}-linux-amd64.tar.gz gorev-linux-amd64
tar -czf gorev-${VERSION}-darwin-amd64.tar.gz gorev-darwin-amd64
tar -czf gorev-${VERSION}-darwin-arm64.tar.gz gorev-darwin-arm64

# Create zip for Windows
zip gorev-${VERSION}-windows-amd64.zip gorev-windows-amd64.exe

cd ..

# Copy release notes
echo "📝 Copying release notes..."
cp CHANGELOG.md ${BUILD_DIR}/

echo "✅ Build complete! Release artifacts in ${BUILD_DIR}/"
echo ""
echo "📁 Release contents:"
ls -la ${BUILD_DIR}/