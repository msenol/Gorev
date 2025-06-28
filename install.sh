#!/bin/bash
#
# Gorev Installation Script
# https://github.com/msenol/Gorev
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
REPO="msenol/Gorev"
VERSION="${VERSION:-v0.7.0-dev}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        echo -e "${RED}ARM architecture not supported yet${NC}"
        exit 1
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Set binary name
BINARY_NAME="gorev-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

echo -e "${GREEN}Installing Gorev ${VERSION}...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"
echo "Binary: $BINARY_NAME"

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download binary
echo -e "${YELLOW}Downloading from ${DOWNLOAD_URL}...${NC}"
curl -L -o "$TEMP_DIR/gorev" "$DOWNLOAD_URL" || {
    echo -e "${RED}Failed to download gorev${NC}"
    exit 1
}

# Make executable
chmod +x "$TEMP_DIR/gorev"

# Verify binary
echo -e "${YELLOW}Verifying binary...${NC}"
"$TEMP_DIR/gorev" version || {
    echo -e "${RED}Binary verification failed${NC}"
    exit 1
}

# Install
echo -e "${YELLOW}Installing to ${INSTALL_DIR}...${NC}"
sudo mv "$TEMP_DIR/gorev" "$INSTALL_DIR/gorev" || {
    echo -e "${RED}Failed to install gorev. Try with sudo or change INSTALL_DIR${NC}"
    exit 1
}

echo -e "${GREEN}✓ Gorev installed successfully!${NC}"
echo ""
echo "Run 'gorev version' to verify installation"
echo "Run 'gorev help' to see available commands"
echo ""
echo "To use with Claude Desktop, VS Code, or other MCP clients,"
echo "see: https://github.com/${REPO}#mcp-editör-entegrasyonu"