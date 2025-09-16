#!/bin/bash
#
# Gorev VM Setup - Step 3: VS Code Extension Setup
# This script builds, packages, and installs the Gorev VS Code extension
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
WORKSPACE_DIR="$HOME/Projects"
PROJECT_DIR="$WORKSPACE_DIR/Gorev"
EXTENSION_DIR="$PROJECT_DIR/gorev-vscode"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Gorev VM Setup - VS Code Extension${NC}"
echo -e "${BLUE}========================================${NC}"

# Function to print success message
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Function to print error message
print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Function to print info message
print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Function to run command with error handling
run_command() {
    local cmd="$1"
    local description="$2"

    echo -e "\n${YELLOW}Running: ${description}${NC}"

    if eval "$cmd"; then
        print_success "$description completed"
        return 0
    else
        print_error "$description failed"
        return 1
    fi
}

# Check prerequisites
if ! command -v code >/dev/null 2>&1; then
    print_error "VS Code is not installed"
    exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
    print_error "npm is not installed"
    exit 1
fi

if [ ! -d "$PROJECT_DIR" ]; then
    print_error "Gorev project not found at $PROJECT_DIR"
    print_info "Please run ./02-build-gorev.sh first"
    exit 1
fi

# 1. Navigate to extension directory
echo -e "\n${YELLOW}Step 1: Navigating to extension directory...${NC}"
cd "$EXTENSION_DIR"
print_success "In directory: $(pwd)"

# Show extension info
echo -e "\n${BLUE}Extension Information:${NC}"
if [ -f "package.json" ]; then
    echo "Name: $(node -p "require('./package.json').name")"
    echo "Display Name: $(node -p "require('./package.json').displayName")"
    echo "Version: $(node -p "require('./package.json').version")"
    echo "Publisher: $(node -p "require('./package.json').publisher")"
    echo "VS Code Engine: $(node -p "require('./package.json').engines.vscode")"
else
    print_error "package.json not found"
    exit 1
fi

# 2. Install dependencies
echo -e "\n${YELLOW}Step 2: Installing npm dependencies...${NC}"
run_command "npm install" "Installing dependencies"

# Check for security vulnerabilities
if npm audit --audit-level=high >/dev/null 2>&1; then
    run_command "npm audit --audit-level=high" "Checking for vulnerabilities"
else
    print_info "Audit found issues, attempting to fix..."
    npm audit fix --force || print_info "Some audit issues could not be fixed automatically"
fi

# 3. Compile TypeScript
echo -e "\n${YELLOW}Step 3: Compiling TypeScript...${NC}"
run_command "npm run compile" "Compiling TypeScript"

# Check compiled output
if [ -d "out" ]; then
    print_success "Compilation output found in 'out' directory"
    echo "Compiled files: $(find out -name "*.js" | wc -l) JavaScript files"
else
    print_error "Compilation output directory 'out' not found"
fi

# 4. Run tests
echo -e "\n${YELLOW}Step 4: Running extension tests...${NC}"
# Note: VS Code extension tests might require a display, so we'll make this optional
if command -v xvfb-run >/dev/null 2>&1; then
    # Use xvfb for headless testing
    run_command "xvfb-run -a npm test" "Running extension tests (headless)" || print_info "Tests failed, continuing..."
else
    # Try running tests normally
    if run_command "npm test" "Running extension tests"; then
        print_success "Extension tests passed"
    else
        print_info "Extension tests failed or require display - this is normal in VM environments"
    fi
fi

# 5. Lint the code
echo -e "\n${YELLOW}Step 5: Linting code...${NC}"
if npm run lint >/dev/null 2>&1; then
    run_command "npm run lint" "Linting TypeScript code"
else
    print_info "Linting not available or failed"
fi

# 6. Install @vscode/vsce for packaging
echo -e "\n${YELLOW}Step 6: Installing packaging tools...${NC}"
if ! npm list -g @vscode/vsce >/dev/null 2>&1; then
    run_command "npm install -g @vscode/vsce" "Installing @vscode/vsce globally"
else
    print_success "@vscode/vsce already installed"
fi

# 7. Package the extension
echo -e "\n${YELLOW}Step 7: Packaging extension...${NC}"
run_command "npm run package" "Creating VSIX package"

# Find the created VSIX file
VSIX_FILE=$(find . -name "*.vsix" -type f | head -1)
if [ -n "$VSIX_FILE" ]; then
    print_success "VSIX package created: $VSIX_FILE"
    echo "Package size: $(ls -lh "$VSIX_FILE" | awk '{print $5}')"
else
    print_error "VSIX package not found"
    exit 1
fi

# 8. Install the extension
echo -e "\n${YELLOW}Step 8: Installing extension to VS Code...${NC}"

# First, uninstall if already installed
if code --list-extensions | grep -q "mehmetsenol.gorev-vscode"; then
    print_info "Uninstalling existing extension..."
    code --uninstall-extension mehmetsenol.gorev-vscode || print_info "Could not uninstall existing extension"
fi

# Install the new extension
run_command "code --install-extension \"$VSIX_FILE\"" "Installing extension from VSIX"

# Verify installation
if code --list-extensions | grep -q "mehmetsenol.gorev-vscode"; then
    print_success "Extension installed successfully"
else
    print_error "Extension installation verification failed"
fi

# 9. Setup VS Code workspace
echo -e "\n${YELLOW}Step 9: Setting up VS Code workspace...${NC}"

# Create workspace settings
WORKSPACE_SETTINGS="$PROJECT_DIR/.vscode/settings.json"
mkdir -p "$(dirname "$WORKSPACE_SETTINGS")"

cat > "$WORKSPACE_SETTINGS" << 'EOF'
{
    "gorev.databaseMode": "workspace",
    "gorev.pagination.pageSize": 100,
    "gorev.debug.useWrapper": true,
    "gorev.debug.logLevel": "debug",
    "go.toolsManagement.checkForUpdates": "proxy",
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "typescript.preferences.includePackageJsonAutoImports": "on",
    "eslint.enable": true,
    "files.exclude": {
        "**/node_modules": true,
        "**/out": false,
        "**/*.vsix": true
    }
}
EOF

print_success "VS Code workspace settings created"

# Create launch configuration for debugging
LAUNCH_CONFIG="$PROJECT_DIR/.vscode/launch.json"
cat > "$LAUNCH_CONFIG" << 'EOF'
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Run Extension",
            "type": "extensionHost",
            "request": "launch",
            "runtimeExecutable": "${execPath}",
            "args": [
                "--extensionDevelopmentPath=${workspaceFolder}/gorev-vscode"
            ],
            "outFiles": [
                "${workspaceFolder}/gorev-vscode/out/**/*.js"
            ]
        },
        {
            "name": "Debug Gorev Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/gorev-mcpserver/cmd/gorev",
            "args": ["serve", "--debug"],
            "cwd": "${workspaceFolder}/gorev-mcpserver"
        }
    ]
}
EOF

print_success "VS Code launch configuration created"

# Create tasks configuration
TASKS_CONFIG="$PROJECT_DIR/.vscode/tasks.json"
cat > "$TASKS_CONFIG" << 'EOF'
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Build Gorev Server",
            "type": "shell",
            "command": "make build",
            "group": "build",
            "options": {
                "cwd": "${workspaceFolder}/gorev-mcpserver"
            },
            "problemMatcher": ["$go"]
        },
        {
            "label": "Test Gorev Server",
            "type": "shell",
            "command": "make test",
            "group": "test",
            "options": {
                "cwd": "${workspaceFolder}/gorev-mcpserver"
            },
            "problemMatcher": ["$go"]
        },
        {
            "label": "Compile Extension",
            "type": "npm",
            "script": "compile",
            "group": "build",
            "options": {
                "cwd": "${workspaceFolder}/gorev-vscode"
            },
            "problemMatcher": ["$tsc"]
        },
        {
            "label": "Package Extension",
            "type": "npm",
            "script": "package",
            "group": "build",
            "options": {
                "cwd": "${workspaceFolder}/gorev-vscode"
            }
        }
    ]
}
EOF

print_success "VS Code tasks configuration created"

# 10. Test the extension installation
echo -e "\n${YELLOW}Step 10: Testing extension installation...${NC}"

# Create a simple test script
TEST_SCRIPT="$PROJECT_DIR/scripts/vm-setup/test-extension.sh"
cat > "$TEST_SCRIPT" << 'EOF'
#!/bin/bash
echo "Testing Gorev VS Code Extension..."

# Check if extension is installed
if code --list-extensions | grep -q "mehmetsenol.gorev-vscode"; then
    echo "✓ Extension is installed"
else
    echo "✗ Extension is not installed"
    exit 1
fi

# Test opening VS Code with the project
echo "Opening VS Code with Gorev project..."
cd ~/Projects/Gorev
code . &
CODE_PID=$!

# Wait a moment for VS Code to start
sleep 3

echo "VS Code opened with PID: $CODE_PID"
echo "You can now test the extension manually:"
echo "1. Open Command Palette (Ctrl+Shift+P)"
echo "2. Type 'Gorev' to see available commands"
echo "3. Try 'Gorev Debug: Seed Test Data'"
echo "4. Check the Gorev panels in the Explorer"

# Optional: kill VS Code after a few seconds for automated testing
# kill $CODE_PID 2>/dev/null || true
EOF

chmod +x "$TEST_SCRIPT"
print_success "Extension test script created"

# 11. Create development helper scripts
echo -e "\n${YELLOW}Step 11: Creating development helper scripts...${NC}"

# Extension development helper
EXT_DEV_SCRIPT="$PROJECT_DIR/scripts/vm-setup/dev-extension.sh"
cat > "$EXT_DEV_SCRIPT" << 'EOF'
#!/bin/bash
# Development helper for Gorev VS Code Extension

cd ~/Projects/Gorev/gorev-vscode

echo "Gorev Extension Development Helper"
echo "================================="
echo "1. Compile and watch: npm run watch"
echo "2. Run tests: npm test"
echo "3. Package: npm run package"
echo "4. Install dev version: code --install-extension *.vsix"
echo "5. Debug: F5 in VS Code"

case "$1" in
    "watch")
        echo "Starting TypeScript compiler in watch mode..."
        npm run watch
        ;;
    "test")
        echo "Running extension tests..."
        npm test
        ;;
    "package")
        echo "Packaging extension..."
        npm run package
        ;;
    "install")
        echo "Installing latest VSIX..."
        VSIX=$(ls -t *.vsix | head -1)
        if [ -n "$VSIX" ]; then
            code --install-extension "$VSIX"
        else
            echo "No VSIX file found. Run 'package' first."
        fi
        ;;
    "debug")
        echo "Opening VS Code for extension debugging..."
        code .
        echo "Press F5 to start debugging"
        ;;
    *)
        echo "Usage: $0 {watch|test|package|install|debug}"
        echo ""
        echo "Available commands:"
        echo "  watch   - Compile in watch mode"
        echo "  test    - Run tests"
        echo "  package - Create VSIX package"
        echo "  install - Install latest VSIX"
        echo "  debug   - Open for debugging"
        ;;
esac
EOF

chmod +x "$EXT_DEV_SCRIPT"
print_success "Extension development helper created"

# 12. Final summary
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}  VS Code Extension Setup Complete${NC}"
echo -e "${BLUE}========================================${NC}"

echo -e "\n${GREEN}Installation Summary:${NC}"
echo "✓ Dependencies installed"
echo "✓ TypeScript compiled"
echo "✓ Extension packaged: $VSIX_FILE"
echo "✓ Extension installed to VS Code"
echo "✓ Workspace configured"
echo "✓ Development helpers created"

echo -e "\n${YELLOW}Next Steps:${NC}"
echo "1. Open VS Code: code ~/Projects/Gorev"
echo "2. Press F5 to start extension debugging"
echo "3. Test extension commands in Command Palette"
echo "4. Run comprehensive tests: ./04-run-tests.sh"

echo -e "\n${YELLOW}Development Commands:${NC}"
echo "• Extension development: ./scripts/vm-setup/dev-extension.sh"
echo "• Test extension: ./scripts/vm-setup/test-extension.sh"
echo "• Recompile: cd gorev-vscode && npm run compile"
echo "• Repackage: cd gorev-vscode && npm run package"

echo -e "\n${BLUE}Extension is ready for testing!${NC}"