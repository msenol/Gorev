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
VERSION="${VERSION:-v0.11.0}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
DATA_DIR="${DATA_DIR:-$HOME/.gorev}"

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

# Download binary
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
echo -e "${YELLOW}Downloading from ${DOWNLOAD_URL}...${NC}"
curl -L -o "$TEMP_DIR/gorev" "$DOWNLOAD_URL" || {
    echo -e "${RED}Failed to download gorev${NC}"
    exit 1
}

# Make executable
chmod +x "$TEMP_DIR/gorev"

# Create data directory
echo -e "${YELLOW}Creating data directory at ${DATA_DIR}...${NC}"
mkdir -p "$DATA_DIR"

# Download and extract source for migrations
echo -e "${YELLOW}Downloading migrations and data files...${NC}"
SOURCE_URL="https://github.com/${REPO}/archive/refs/tags/${VERSION}.tar.gz"
curl -L -o "$TEMP_DIR/source.tar.gz" "$SOURCE_URL" || {
    echo -e "${RED}Failed to download source files${NC}"
    exit 1
}

# Extract source
tar -xzf "$TEMP_DIR/source.tar.gz" -C "$TEMP_DIR"

# Find the extracted directory (it will be named Gorev-VERSION without 'v' prefix)
VERSION_WITHOUT_V=${VERSION#v}
SOURCE_DIR="$TEMP_DIR/Gorev-${VERSION_WITHOUT_V}"

if [ ! -d "$SOURCE_DIR" ]; then
    # Try alternative naming
    SOURCE_DIR=$(find "$TEMP_DIR" -maxdepth 1 -type d -name "Gorev-*" | head -1)
fi

if [ ! -d "$SOURCE_DIR" ]; then
    echo -e "${RED}Failed to find source directory${NC}"
    exit 1
fi

# Copy migrations
echo -e "${YELLOW}Copying migration files...${NC}"
mkdir -p "$DATA_DIR/internal/veri"
cp -r "$SOURCE_DIR/gorev-mcpserver/internal/veri/migrations" "$DATA_DIR/internal/veri/" || {
    echo -e "${RED}Failed to copy migration files${NC}"
    exit 1
}

# Verify binary with new GOREV_ROOT
echo -e "${YELLOW}Verifying binary...${NC}"
export GOREV_ROOT="$DATA_DIR"
"$TEMP_DIR/gorev" version || {
    echo -e "${RED}Binary verification failed${NC}"
    exit 1
}

# Install binary
echo -e "${YELLOW}Installing binary to ${INSTALL_DIR}...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    cp "$TEMP_DIR/gorev" "$INSTALL_DIR/gorev.bin"
else
    sudo cp "$TEMP_DIR/gorev" "$INSTALL_DIR/gorev.bin" || {
        echo -e "${RED}Failed to install gorev. Try with sudo or change INSTALL_DIR${NC}"
        exit 1
    }
fi

# Create wrapper script
echo -e "${YELLOW}Creating wrapper script...${NC}"
WRAPPER_CONTENT="#!/bin/bash
# Gorev wrapper script
export GOREV_ROOT=\"$DATA_DIR\"
exec \"$INSTALL_DIR/gorev.bin\" \"\$@\"
"

# Install wrapper
if [ -w "$INSTALL_DIR" ]; then
    echo "$WRAPPER_CONTENT" > "$INSTALL_DIR/gorev"
    chmod +x "$INSTALL_DIR/gorev"
else
    echo "$WRAPPER_CONTENT" | sudo tee "$INSTALL_DIR/gorev" > /dev/null
    sudo chmod +x "$INSTALL_DIR/gorev"
fi

echo -e "${GREEN}✓ Gorev installed successfully!${NC}"
echo ""
echo "Installed components:"
echo "  Binary: $INSTALL_DIR/gorev"
echo "  Data files: $DATA_DIR"
echo ""
echo "Run 'gorev version' to verify installation"
echo "Run 'gorev help' to see available commands"
echo ""
echo "To use with Claude Desktop, VS Code, or other MCP clients,"
echo "see: https://github.com/${REPO}#mcp-editör-entegrasyonu"

# Add to shell profile if not already there
SHELL_PROFILE=""
if [ -n "$ZSH_VERSION" ]; then
    SHELL_PROFILE="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_PROFILE="$HOME/.bashrc"
fi

if [ -n "$SHELL_PROFILE" ] && [ -f "$SHELL_PROFILE" ]; then
    if ! grep -q "GOREV_ROOT" "$SHELL_PROFILE"; then
        echo ""
        echo -e "${YELLOW}Adding GOREV_ROOT to $SHELL_PROFILE...${NC}"
        echo "export GOREV_ROOT=\"$DATA_DIR\"" >> "$SHELL_PROFILE"
        echo -e "${GREEN}✓ Added to shell profile. Run 'source $SHELL_PROFILE' or restart your terminal.${NC}"
    fi
fi