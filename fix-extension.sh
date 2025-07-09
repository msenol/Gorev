#!/bin/bash

echo "=== Gorev Extension Fix Script ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MCP_SERVER_DIR="$PROJECT_ROOT/gorev-mcpserver"
VSCODE_EXT_DIR="$PROJECT_ROOT/gorev-vscode"

echo "Project root: $PROJECT_ROOT"
echo "MCP Server dir: $MCP_SERVER_DIR"
echo "VS Code ext dir: $VSCODE_EXT_DIR"
echo

# Step 1: Build the MCP server
echo -e "${YELLOW}Step 1: Building MCP server...${NC}"
cd "$MCP_SERVER_DIR"

if [ -f "go.mod" ]; then
    go build -o gorev cmd/gorev/main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ MCP server built successfully${NC}"
        ls -la gorev
    else
        echo -e "${RED}✗ Failed to build MCP server${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ go.mod not found in $MCP_SERVER_DIR${NC}"
    exit 1
fi

# Step 2: Check if database exists
echo
echo -e "${YELLOW}Step 2: Checking database...${NC}"
if [ -f "gorev.db" ]; then
    echo -e "${GREEN}✓ Database exists${NC}"
    echo "Database size: $(du -h gorev.db | cut -f1)"
else
    echo -e "${YELLOW}! Database doesn't exist, will be created on first run${NC}"
fi

# Step 3: Build VS Code extension
echo
echo -e "${YELLOW}Step 3: Building VS Code extension...${NC}"
cd "$VSCODE_EXT_DIR"

if [ -f "package.json" ]; then
    npm install
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Dependencies installed${NC}"
    else
        echo -e "${RED}✗ Failed to install dependencies${NC}"
        exit 1
    fi
    
    npm run compile
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Extension compiled successfully${NC}"
    else
        echo -e "${RED}✗ Failed to compile extension${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ package.json not found in $VSCODE_EXT_DIR${NC}"
    exit 1
fi

# Step 4: Instructions
echo
echo -e "${GREEN}=== Setup Complete ===${NC}"
echo
echo "Next steps:"
echo "1. Open VS Code"
echo "2. Press Ctrl+Shift+P (Cmd+Shift+P on macOS)"
echo "3. Run 'Developer: Reload Window'"
echo "4. Open the Gorev view in the sidebar"
echo "5. Click the refresh button or press Ctrl+Alt+R"
echo
echo "If you still have issues:"
echo "- Check VS Code Output panel (View > Output > Gorev)"
echo "- Verify server path in settings: gorev.serverPath"
echo "- Make sure you have an active project selected"
echo
echo -e "${YELLOW}Server executable path:${NC}"
echo "$MCP_SERVER_DIR/gorev"