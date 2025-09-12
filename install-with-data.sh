#!/bin/bash
#
# Gorev Installation Script with Data Files
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
VERSION="${VERSION:-v0.14.0}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
DATA_DIR="${DATA_DIR:-/opt/gorev}"

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
echo "Install directory: $INSTALL_DIR"
echo "Data directory: $DATA_DIR"

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download archive
ARCHIVE_NAME="gorev-${VERSION}-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

echo -e "${YELLOW}Downloading from ${DOWNLOAD_URL}...${NC}"
curl -L -o "$TEMP_DIR/gorev.tar.gz" "$DOWNLOAD_URL" || {
    echo -e "${YELLOW}Archive not found, trying binary only...${NC}"
    
    # Fallback to binary only
    BINARY_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
    curl -L -o "$TEMP_DIR/gorev" "$BINARY_URL" || {
        echo -e "${RED}Failed to download gorev${NC}"
        exit 1
    }
    
    # Binary only installation
    chmod +x "$TEMP_DIR/gorev"
    
    # Verify binary
    echo -e "${YELLOW}Verifying binary...${NC}"
    "$TEMP_DIR/gorev" version || {
        echo -e "${RED}Binary verification failed${NC}"
        exit 1
    }
    
    # Install binary
    echo -e "${YELLOW}Installing to ${INSTALL_DIR}...${NC}"
    sudo mv "$TEMP_DIR/gorev" "$INSTALL_DIR/gorev" || {
        echo -e "${RED}Failed to install gorev. Try with sudo or change INSTALL_DIR${NC}"
        exit 1
    }
    
    echo -e "${GREEN}✓ Gorev binary installed successfully!${NC}"
    echo ""
    echo -e "${YELLOW}NOTE: This is a binary-only installation.${NC}"
    echo -e "${YELLOW}You need to set GOREV_ROOT environment variable:${NC}"
    echo ""
    echo "  export GOREV_ROOT=/path/to/gorev-mcpserver"
    echo ""
    echo "Add this to your ~/.bashrc or ~/.zshrc file."
    echo ""
    echo "Run 'gorev version' to verify installation"
    echo "Run 'gorev help' to see available commands"
    exit 0
}

# Extract archive
echo -e "${YELLOW}Extracting archive...${NC}"
tar -xzf "$TEMP_DIR/gorev.tar.gz" -C "$TEMP_DIR" || {
    echo -e "${RED}Failed to extract archive${NC}"
    exit 1
}

# Find extracted directory
EXTRACT_DIR=$(find "$TEMP_DIR" -name "gorev-*" -type d | head -1)
if [ -z "$EXTRACT_DIR" ]; then
    echo -e "${RED}Failed to find extracted directory${NC}"
    exit 1
fi

# Verify binary
echo -e "${YELLOW}Verifying binary...${NC}"
chmod +x "$EXTRACT_DIR/gorev"
"$EXTRACT_DIR/gorev" version || {
    echo -e "${RED}Binary verification failed${NC}"
    exit 1
}

# Create data directory
echo -e "${YELLOW}Creating data directory at ${DATA_DIR}...${NC}"
sudo mkdir -p "$DATA_DIR" || {
    echo -e "${RED}Failed to create data directory${NC}"
    exit 1
}

# Copy files
echo -e "${YELLOW}Copying files...${NC}"
sudo cp -r "$EXTRACT_DIR/internal" "$DATA_DIR/" || {
    echo -e "${RED}Failed to copy data files${NC}"
    exit 1
}

# Install binary
echo -e "${YELLOW}Installing binary to ${INSTALL_DIR}...${NC}"
sudo cp "$EXTRACT_DIR/gorev" "$INSTALL_DIR/gorev.bin" || {
    echo -e "${RED}Failed to install binary${NC}"
    exit 1
}

# Create wrapper script
echo -e "${YELLOW}Creating wrapper script...${NC}"
WRAPPER_CONTENT="#!/bin/bash
# Gorev wrapper script
export GOREV_ROOT=\"$DATA_DIR\"
exec \"$INSTALL_DIR/gorev.bin\" \"\$@\"
"

echo "$WRAPPER_CONTENT" | sudo tee "$INSTALL_DIR/gorev" > /dev/null
sudo chmod +x "$INSTALL_DIR/gorev"

echo -e "${GREEN}✓ Gorev installed successfully!${NC}"
echo ""
echo "Installed files:"
echo "  Binary: $INSTALL_DIR/gorev"
echo "  Data: $DATA_DIR"
echo ""
echo "Run 'gorev version' to verify installation"
echo "Run 'gorev help' to see available commands"
echo ""
echo "To use with Claude Desktop, VS Code, or other MCP clients,"
echo "see: https://github.com/${REPO}#mcp-editör-entegrasyonu"