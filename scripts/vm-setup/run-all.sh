#!/bin/bash
#
# Gorev VM Setup - Master Script
# This script runs all setup scripts in sequence
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$SCRIPT_DIR/setup.log"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Gorev VirtualBox VM Setup Master${NC}"
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

# Function to print section header
print_section() {
    echo -e "\n${CYAN}========================================${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}========================================${NC}"
}

# Function to run script with logging
run_script() {
    local script="$1"
    local description="$2"
    local script_path="$SCRIPT_DIR/$script"

    print_section "$description"

    if [ ! -f "$script_path" ]; then
        print_error "Script not found: $script_path"
        return 1
    fi

    if [ ! -x "$script_path" ]; then
        chmod +x "$script_path"
        print_info "Made script executable: $script"
    fi

    echo "Running: $script" | tee -a "$LOG_FILE"
    echo "Started at: $(date)" | tee -a "$LOG_FILE"

    if "$script_path" 2>&1 | tee -a "$LOG_FILE"; then
        print_success "$description completed successfully"
        echo "Completed at: $(date)" | tee -a "$LOG_FILE"
        return 0
    else
        print_error "$description failed"
        echo "Failed at: $(date)" | tee -a "$LOG_FILE"
        return 1
    fi
}

# Initialize log file
echo "Gorev VM Setup Master Script" > "$LOG_FILE"
echo "Started at: $(date)" >> "$LOG_FILE"
echo "========================================" >> "$LOG_FILE"

# Check if we're in the right directory
if [ ! -f "$SCRIPT_DIR/01-install-prerequisites.sh" ]; then
    print_error "Setup scripts not found in current directory"
    print_info "Please run this script from the vm-setup directory"
    exit 1
fi

# Show setup plan
echo -e "\n${YELLOW}Setup Plan:${NC}"
echo "1. Install Prerequisites (Go, Node.js, VS Code, etc.)"
echo "2. Build Gorev Project (Clone, compile, initialize)"
echo "3. Setup VS Code Extension (Compile, package, install)"
echo "4. Run Comprehensive Tests (Unit, integration, coverage)"
echo "5. Setup Debug Helper (Optional)"

echo -e "\n${YELLOW}Estimated time: 15-30 minutes${NC}"
echo -e "${YELLOW}Log file: $LOG_FILE${NC}"

# Confirmation
read -p "Do you want to proceed with the full setup? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Setup cancelled by user"
    exit 0
fi

# Start setup process
START_TIME=$(date +%s)

# Step 1: Prerequisites
if ! run_script "01-install-prerequisites.sh" "Step 1: Install Prerequisites"; then
    print_error "Prerequisites installation failed. Check log: $LOG_FILE"
    exit 1
fi

# Source bashrc to get new PATH
print_info "Sourcing ~/.bashrc to update environment..."
if [ -f ~/.bashrc ]; then
    # We can't source bashrc in a script effectively, so we'll ask the user
    echo -e "\n${YELLOW}IMPORTANT:${NC} Please run the following command before continuing:"
    echo -e "${CYAN}source ~/.bashrc${NC}"
    echo ""
    read -p "Press Enter after running 'source ~/.bashrc'..."
fi

# Check if Go is available
if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not available in PATH. Please:"
    echo "1. Run: source ~/.bashrc"
    echo "2. Verify: go version"
    echo "3. Re-run this script"
    exit 1
fi

# Step 2: Build Project
if ! run_script "02-build-gorev.sh" "Step 2: Build Gorev Project"; then
    print_error "Project build failed. Check log: $LOG_FILE"
    echo -e "\n${YELLOW}Common solutions:${NC}"
    echo "- Ensure Go is in PATH: go version"
    echo "- Check internet connection"
    echo "- Try running the script manually: ./02-build-gorev.sh"
    exit 1
fi

# Step 3: VS Code Extension
if ! run_script "03-setup-vscode.sh" "Step 3: Setup VS Code Extension"; then
    print_error "VS Code extension setup failed. Check log: $LOG_FILE"
    echo -e "\n${YELLOW}Common solutions:${NC}"
    echo "- Ensure VS Code is installed: code --version"
    echo "- Check if extension compiled: ls gorev-vscode/out/"
    echo "- Try running manually: ./03-setup-vscode.sh"

    # Ask if user wants to continue without extension
    read -p "Continue without VS Code extension? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Step 4: Run Tests
print_info "Running comprehensive tests..."
if ! run_script "04-run-tests.sh" "Step 4: Run Comprehensive Tests"; then
    print_error "Some tests failed. Check log: $LOG_FILE"
    echo -e "\n${YELLOW}This is often normal in VM environments.${NC}"
    echo "Key components should still work correctly."

    # Ask if user wants to continue
    read -p "Continue with setup completion? (Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        exit 1
    fi
fi

# Step 5: Setup Debug Helper (make executable)
if [ -f "$SCRIPT_DIR/05-debug-helper.sh" ]; then
    chmod +x "$SCRIPT_DIR/05-debug-helper.sh"
    print_success "Debug helper script is ready: ./05-debug-helper.sh"
fi

# Calculate total time
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))
MINUTES=$((TOTAL_TIME / 60))
SECONDS=$((TOTAL_TIME % 60))

# Final summary
print_section "Setup Complete!"

echo -e "${GREEN}🎉 Gorev VM setup completed successfully!${NC}"
echo -e "\n${YELLOW}Setup Summary:${NC}"
echo "Total time: ${MINUTES}m ${SECONDS}s"
echo "Log file: $LOG_FILE"

echo -e "\n${YELLOW}What was installed:${NC}"
echo "✓ Go $(go version 2>/dev/null | awk '{print $3}' || echo 'Unknown')"
echo "✓ Node.js $(node --version 2>/dev/null || echo 'Unknown')"
echo "✓ VS Code $(code --version 2>/dev/null | head -1 || echo 'Unknown')"
echo "✓ Gorev MCP Server"
echo "✓ Gorev VS Code Extension"
echo "✓ Development tools and aliases"

echo -e "\n${YELLOW}Quick Start Commands:${NC}"
echo "1. Start Gorev server:"
echo "   cd ~/Projects/Gorev/gorev-mcpserver"
echo "   ./gorev serve --debug"
echo ""
echo "2. Open VS Code for extension development:"
echo "   cd ~/Projects/Gorev"
echo "   code ."
echo "   Press F5 to debug extension"
echo ""
echo "3. Test basic functionality:"
echo "   cd ~/Projects/Gorev/gorev-mcpserver"
echo "   ./gorev version"
echo "   ./gorev template aliases"
echo ""
echo "4. Database inspection:"
echo "   sqlite3 ~/.gorev/gorev.db"
echo "   .tables"
echo ""
echo "5. Debug and troubleshooting:"
echo "   ./scripts/vm-setup/05-debug-helper.sh"

echo -e "\n${YELLOW}Development Aliases (after source ~/.bashrc):${NC}"
echo "• gorev-serve       - Start server with debug"
echo "• gorev-test        - Run server tests"
echo "• gorev-ext-compile - Compile extension"
echo "• gorev-db-global   - Open global database"

echo -e "\n${YELLOW}Project Locations:${NC}"
echo "• Project: ~/Projects/Gorev"
echo "• Server: ~/Projects/Gorev/gorev-mcpserver"
echo "• Extension: ~/Projects/Gorev/gorev-vscode"
echo "• Global DB: ~/.gorev/gorev.db"
echo "• Workspace DB: ~/Projects/Gorev/.gorev/gorev.db"

echo -e "\n${YELLOW}VS Code Extension:${NC}"
if code --list-extensions 2>/dev/null | grep -q "mehmetsenol.gorev-vscode"; then
    echo "✓ Extension installed and ready"
    echo "  Open Command Palette (Ctrl+Shift+P) and type 'Gorev'"
else
    echo "⚠ Extension installation may have failed"
    echo "  Try: cd ~/Projects/Gorev/gorev-vscode && npm run package"
fi

echo -e "\n${YELLOW}Next Steps:${NC}"
echo "1. Test the server: cd ~/Projects/Gorev/gorev-mcpserver && ./gorev serve --debug"
echo "2. Open VS Code: code ~/Projects/Gorev"
echo "3. Test extension: F5 in VS Code, then Ctrl+Shift+P -> Gorev commands"
echo "4. Run debug helper if needed: ./scripts/vm-setup/05-debug-helper.sh"

echo -e "\n${GREEN}Happy coding with Gorev! 🚀${NC}"

# Write final summary to log
{
    echo ""
    echo "========================================="
    echo "Setup completed successfully at: $(date)"
    echo "Total time: ${MINUTES}m ${SECONDS}s"
    echo "========================================="
} >> "$LOG_FILE"