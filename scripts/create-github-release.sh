#!/bin/bash

# GitHub Release Creation Script for Gorev v0.10.2
# This script creates a GitHub release with all attachments

set -e

VERSION="v0.10.2"
RELEASE_DIR="release-${VERSION}"
GITHUB_REPO="msenol/Gorev"

echo "📦 Creating GitHub Release for Gorev ${VERSION}..."

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Please install it first:"
    echo "   https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub. Run 'gh auth login' first."
    exit 1
fi

# Check if release directory exists
if [ ! -d "$RELEASE_DIR" ]; then
    echo "❌ Release directory not found. Run build-release.sh first."
    exit 1
fi

# Create release with release notes
echo "📝 Creating release..."
gh release create "${VERSION}" \
    --repo "${GITHUB_REPO}" \
    --title "Gorev ${VERSION} - Enhanced MCP Debug System & Pagination Fixes" \
    --notes-file "RELEASE_NOTES_v0.10.2.md" \
    --draft

# Upload binaries
echo "📤 Uploading binaries..."
gh release upload "${VERSION}" \
    --repo "${GITHUB_REPO}" \
    "${RELEASE_DIR}/gorev-linux-amd64" \
    "${RELEASE_DIR}/gorev-darwin-amd64" \
    "${RELEASE_DIR}/gorev-darwin-arm64" \
    "${RELEASE_DIR}/gorev-windows-amd64.exe"

# Upload archives
echo "📤 Uploading archives..."
gh release upload "${VERSION}" \
    --repo "${GITHUB_REPO}" \
    "${RELEASE_DIR}/gorev-${VERSION}-linux-amd64.tar.gz" \
    "${RELEASE_DIR}/gorev-${VERSION}-darwin-amd64.tar.gz" \
    "${RELEASE_DIR}/gorev-${VERSION}-darwin-arm64.tar.gz" \
    "${RELEASE_DIR}/gorev-${VERSION}-windows-amd64.zip"

# Upload checksums
echo "📤 Uploading checksums..."
gh release upload "${VERSION}" \
    --repo "${GITHUB_REPO}" \
    "${RELEASE_DIR}/checksums.txt"

echo "✅ Draft release created successfully!"
echo ""
echo "📌 Next steps:"
echo "1. Review the draft release at: https://github.com/${GITHUB_REPO}/releases"
echo "2. Edit release notes if needed"
echo "3. Publish the release when ready"
echo ""
echo "🎨 Don't forget to:"
echo "- Update VS Code extension in marketplace"
echo "- Update documentation website"
echo "- Announce on social media"