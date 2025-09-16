#!/bin/bash
#
# Gorev VM Setup - Step 1: Install Prerequisites
# This script installs all required dependencies for Gorev development
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GO_VERSION="1.23.5"
NODE_VERSION="20"  # LTS version

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Gorev VM Setup - Prerequisites${NC}"
echo -e "${BLUE}========================================${NC}"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

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

# 1. Update system packages
echo -e "\n${YELLOW}Step 1: Updating system packages...${NC}"
sudo apt update && sudo apt upgrade -y
print_success "System packages updated"

# 2. Install build essentials
echo -e "\n${YELLOW}Step 2: Installing build essentials...${NC}"
sudo apt install -y \
    build-essential \
    git \
    curl \
    wget \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    sqlite3 \
    jq
print_success "Build essentials installed"

# 3. Install Go
echo -e "\n${YELLOW}Step 3: Installing Go ${GO_VERSION}...${NC}"
if command_exists go; then
    CURRENT_GO=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go ${CURRENT_GO} already installed"
    read -p "Do you want to reinstall Go ${GO_VERSION}? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Skipping Go installation"
    else
        wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        rm "go${GO_VERSION}.linux-amd64.tar.gz"
        print_success "Go ${GO_VERSION} installed"
    fi
else
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
    rm "go${GO_VERSION}.linux-amd64.tar.gz"
    print_success "Go ${GO_VERSION} installed"
fi

# Add Go to PATH if not already there
if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
    print_success "Go added to PATH"
fi

# 4. Install Node.js and npm
echo -e "\n${YELLOW}Step 4: Installing Node.js and npm...${NC}"
if command_exists node; then
    CURRENT_NODE=$(node --version)
    print_info "Node.js ${CURRENT_NODE} already installed"
else
    curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
    sudo apt install -y nodejs
    print_success "Node.js $(node --version) installed"
fi

# 5. Install VS Code
echo -e "\n${YELLOW}Step 5: Installing Visual Studio Code...${NC}"
if command_exists code; then
    print_info "VS Code already installed: $(code --version | head -1)"
else
    # Method 1: Using snap (simpler but might not work in all VMs)
    if command_exists snap; then
        sudo snap install code --classic
        print_success "VS Code installed via snap"
    else
        # Method 2: Using apt repository
        wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
        sudo install -o root -g root -m 644 packages.microsoft.gpg /usr/share/keyrings/packages.microsoft.gpg
        sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/usr/share/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
        sudo apt update
        sudo apt install -y code
        print_success "VS Code installed via apt"
    fi
fi

# 6. Install additional development tools
echo -e "\n${YELLOW}Step 6: Installing additional development tools...${NC}"
sudo apt install -y \
    vim \
    tmux \
    htop \
    net-tools \
    ripgrep \
    fd-find
print_success "Development tools installed"

# 7. Install Go development tools
echo -e "\n${YELLOW}Step 7: Installing Go development tools...${NC}"
export PATH=$PATH:/usr/local/go/bin:~/go/bin
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
print_success "Go development tools installed"

# 8. Create workspace directories
echo -e "\n${YELLOW}Step 8: Creating workspace directories...${NC}"
mkdir -p ~/Projects/gorev-test
mkdir -p ~/.gorev
print_success "Workspace directories created"

# 9. Configure Git (optional)
echo -e "\n${YELLOW}Step 9: Git configuration...${NC}"
if ! git config --global user.name > /dev/null 2>&1; then
    read -p "Enter your Git username: " git_username
    git config --global user.name "$git_username"
fi
if ! git config --global user.email > /dev/null 2>&1; then
    read -p "Enter your Git email: " git_email
    git config --global user.email "$git_email"
fi
print_success "Git configured"

# 10. System information
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}  System Information${NC}"
echo -e "${BLUE}========================================${NC}"
echo "OS: $(lsb_release -ds)"
echo "Kernel: $(uname -r)"
echo "Architecture: $(uname -m)"
echo "Go version: $(go version 2>/dev/null || echo 'Not installed')"
echo "Node version: $(node --version 2>/dev/null || echo 'Not installed')"
echo "npm version: $(npm --version 2>/dev/null || echo 'Not installed')"
echo "VS Code version: $(code --version 2>/dev/null | head -1 || echo 'Not installed')"

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  Prerequisites installation complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo "1. Run: source ~/.bashrc"
echo "2. Clone the Gorev repository"
echo "3. Run the build script: ./02-build-gorev.sh"