#!/bin/bash
#
# Gorev VM Setup - Step 2: Build and Setup Gorev
# This script clones, builds, and configures the Gorev project
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/msenol/Gorev.git"
WORKSPACE_DIR="$HOME/Projects"
PROJECT_DIR="$WORKSPACE_DIR/Gorev"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Gorev VM Setup - Build & Configure${NC}"
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
    echo "Command: $cmd"

    if eval "$cmd"; then
        print_success "$description completed"
    else
        print_error "$description failed"
        exit 1
    fi
}

# Check if Go is available
if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not installed or not in PATH"
    print_info "Please run: source ~/.bashrc"
    print_info "Then run this script again"
    exit 1
fi

# Check if Node.js is available
if ! command -v node >/dev/null 2>&1; then
    print_error "Node.js is not installed"
    print_info "Please run the prerequisites script first"
    exit 1
fi

# 1. Create workspace directory
echo -e "\n${YELLOW}Step 1: Creating workspace directory...${NC}"
mkdir -p "$WORKSPACE_DIR"
cd "$WORKSPACE_DIR"
print_success "Workspace directory created: $WORKSPACE_DIR"

# 2. Clone the repository
echo -e "\n${YELLOW}Step 2: Cloning Gorev repository...${NC}"
if [ -d "$PROJECT_DIR" ]; then
    print_info "Repository already exists. Updating..."
    cd "$PROJECT_DIR"
    git pull origin main
    print_success "Repository updated"
else
    git clone "$REPO_URL"
    cd "$PROJECT_DIR"
    print_success "Repository cloned"
fi

# Show current branch and latest commits
echo -e "\n${BLUE}Repository Information:${NC}"
echo "Branch: $(git branch --show-current)"
echo "Latest commit: $(git log -1 --oneline)"
echo "Total commits: $(git rev-list --count HEAD)"

# 3. Build MCP Server
echo -e "\n${YELLOW}Step 3: Building MCP Server...${NC}"
cd "$PROJECT_DIR/gorev-mcpserver"

# Check Go version compatibility
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_info "Using Go version: $GO_VERSION"

# Download dependencies
run_command "make deps" "Downloading Go dependencies"

# Run Go mod tidy
run_command "go mod tidy" "Tidying Go modules"

# Format code
run_command "make fmt" "Formatting Go code"

# Run static analysis
if command -v golangci-lint >/dev/null 2>&1; then
    run_command "golangci-lint run --timeout=5m" "Running Go linter"
else
    print_info "golangci-lint not available, skipping"
fi

# Build the project
run_command "make build" "Building MCP server"

# Verify binary
if [ -f "./gorev" ]; then
    print_success "Binary created successfully"
    echo "Binary size: $(ls -lh ./gorev | awk '{print $5}')"
    echo "Binary path: $(pwd)/gorev"
else
    print_error "Binary not found after build"
    exit 1
fi

# 4. Initialize database
echo -e "\n${YELLOW}Step 4: Initializing database...${NC}"

# Initialize global database
run_command "./gorev init --global" "Initializing global database"

# Initialize workspace database
mkdir -p "$PROJECT_DIR/.gorev"
cd "$PROJECT_DIR"
run_command "./gorev-mcpserver/gorev init" "Initializing workspace database"

# Verify databases
echo -e "\n${BLUE}Database Information:${NC}"
if [ -f "$HOME/.gorev/gorev.db" ]; then
    echo "Global DB: $HOME/.gorev/gorev.db ($(ls -lh $HOME/.gorev/gorev.db | awk '{print $5}'))"
else
    print_error "Global database not created"
fi

if [ -f "$PROJECT_DIR/.gorev/gorev.db" ]; then
    echo "Workspace DB: $PROJECT_DIR/.gorev/gorev.db ($(ls -lh $PROJECT_DIR/.gorev/gorev.db | awk '{print $5}'))"
else
    print_error "Workspace database not created"
fi

# 5. Test MCP Server
echo -e "\n${YELLOW}Step 5: Testing MCP Server...${NC}"
cd "$PROJECT_DIR/gorev-mcpserver"

# Test version command
run_command "./gorev version" "Testing version command"

# Test help command
run_command "./gorev help" "Testing help command"

# Test template command
run_command "./gorev template aliases" "Testing template aliases"

# Quick database test
run_command "echo 'SELECT COUNT(*) as table_count FROM sqlite_master WHERE type=\"table\";' | sqlite3 ~/.gorev/gorev.db" "Testing database structure"

# 6. Run unit tests
echo -e "\n${YELLOW}Step 6: Running unit tests...${NC}"
run_command "make test" "Running unit tests"

# Generate test coverage if possible
if command -v make >/dev/null 2>&1; then
    run_command "make test-coverage" "Generating test coverage"
    if [ -f "coverage.html" ]; then
        print_success "Coverage report generated: $(pwd)/coverage.html"
    fi
fi

# 7. Build VS Code Extension dependencies
echo -e "\n${YELLOW}Step 7: Preparing VS Code Extension...${NC}"
cd "$PROJECT_DIR/gorev-vscode"

# Check package.json
if [ -f "package.json" ]; then
    print_success "package.json found"
    echo "Extension version: $(node -p "require('./package.json').version")"
    echo "Extension name: $(node -p "require('./package.json').displayName || require('./package.json').name")"
else
    print_error "package.json not found in VS Code extension directory"
    exit 1
fi

# Install dependencies
run_command "npm install" "Installing npm dependencies"

# Compile TypeScript
run_command "npm run compile" "Compiling TypeScript"

# Run extension tests
if npm run test >/dev/null 2>&1; then
    run_command "npm run test" "Running extension tests"
else
    print_info "Extension tests not available or failed"
fi

# 8. Create useful aliases and scripts
echo -e "\n${YELLOW}Step 8: Creating useful aliases...${NC}"

# Create alias file
ALIAS_FILE="$PROJECT_DIR/scripts/vm-setup/gorev-aliases.sh"
mkdir -p "$(dirname "$ALIAS_FILE")"

cat > "$ALIAS_FILE" << 'EOF'
#!/bin/bash
# Gorev Development Aliases

# Navigation aliases
alias gorev-cd='cd ~/Projects/Gorev'
alias gorev-server='cd ~/Projects/Gorev/gorev-mcpserver'
alias gorev-ext='cd ~/Projects/Gorev/gorev-vscode'

# Server aliases
alias gorev-serve='cd ~/Projects/Gorev/gorev-mcpserver && ./gorev serve --debug'
alias gorev-build='cd ~/Projects/Gorev/gorev-mcpserver && make build'
alias gorev-test='cd ~/Projects/Gorev/gorev-mcpserver && make test'
alias gorev-clean='cd ~/Projects/Gorev/gorev-mcpserver && make clean'

# Extension aliases
alias gorev-ext-compile='cd ~/Projects/Gorev/gorev-vscode && npm run compile'
alias gorev-ext-test='cd ~/Projects/Gorev/gorev-vscode && npm run test'
alias gorev-ext-package='cd ~/Projects/Gorev/gorev-vscode && npm run package'

# Database aliases
alias gorev-db-global='sqlite3 ~/.gorev/gorev.db'
alias gorev-db-workspace='sqlite3 ~/Projects/Gorev/.gorev/gorev.db'

# Development aliases
alias gorev-logs='cd ~/Projects/Gorev && find . -name "*.log" -type f'
alias gorev-status='cd ~/Projects/Gorev && git status'
alias gorev-pull='cd ~/Projects/Gorev && git pull'

echo "Gorev development aliases loaded!"
echo "Available commands:"
echo "  gorev-cd, gorev-server, gorev-ext"
echo "  gorev-serve, gorev-build, gorev-test"
echo "  gorev-ext-compile, gorev-ext-test, gorev-ext-package"
echo "  gorev-db-global, gorev-db-workspace"
echo "  gorev-logs, gorev-status, gorev-pull"
EOF

chmod +x "$ALIAS_FILE"

# Add aliases to bashrc if not already there
if ! grep -q "gorev-aliases.sh" ~/.bashrc; then
    echo "" >> ~/.bashrc
    echo "# Gorev development aliases" >> ~/.bashrc
    echo "source \"$ALIAS_FILE\"" >> ~/.bashrc
    print_success "Aliases added to ~/.bashrc"
fi

# 9. Create development environment info
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}  Development Environment Summary${NC}"
echo -e "${BLUE}========================================${NC}"

echo "Project Location: $PROJECT_DIR"
echo "Server Binary: $PROJECT_DIR/gorev-mcpserver/gorev"
echo "Global Database: $HOME/.gorev/gorev.db"
echo "Workspace Database: $PROJECT_DIR/.gorev/gorev.db"
echo ""
echo "Key Commands:"
echo "  cd ~/workspace/Gorev/gorev-mcpserver && ./gorev serve --debug"
echo "  cd ~/workspace/Gorev/gorev-vscode && code ."
echo ""
echo "Database Inspection:"
echo "  sqlite3 ~/.gorev/gorev.db"
echo "  .tables"
echo "  SELECT * FROM gorevler LIMIT 5;"
echo ""
echo "Test Commands:"
echo "  cd ~/workspace/Gorev/gorev-mcpserver && make test"
echo "  cd ~/workspace/Gorev/gorev-vscode && npm test"

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  Build completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo "1. Run: source ~/.bashrc  (to load aliases)"
echo "2. Test server: cd ~/workspace/Gorev/gorev-mcpserver && ./gorev serve --debug"
echo "3. Setup VS Code: Run ./03-setup-vscode.sh"
echo "4. Run comprehensive tests: ./04-run-tests.sh"