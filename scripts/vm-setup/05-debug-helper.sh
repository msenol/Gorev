#!/bin/bash
#
# Gorev VM Setup - Step 5: Debug Helper
# This script provides debugging and troubleshooting tools
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
WORKSPACE_DIR="$HOME/Projects"
PROJECT_DIR="$WORKSPACE_DIR/Gorev"
SERVER_DIR="$PROJECT_DIR/gorev-mcpserver"
EXTENSION_DIR="$PROJECT_DIR/gorev-vscode"
DEBUG_LOG_DIR="$PROJECT_DIR/debug-logs"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Gorev Debug Helper & Troubleshooter${NC}"
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

# Function to print menu option
print_option() {
    echo -e "${MAGENTA}$1)${NC} $2"
}

# Create debug logs directory
mkdir -p "$DEBUG_LOG_DIR"

# Main menu function
show_menu() {
    echo -e "\n${YELLOW}Gorev Debug Helper Menu:${NC}"
    echo "========================="
    print_option "1" "System Information"
    print_option "2" "Project Health Check"
    print_option "3" "Database Inspection"
    print_option "4" "Server Debug Mode"
    print_option "5" "Extension Debug Mode"
    print_option "6" "Log Analysis"
    print_option "7" "Network & Process Check"
    print_option "8" "Performance Profiling"
    print_option "9" "Clean & Reset"
    print_option "10" "Generate Debug Report"
    print_option "0" "Exit"
    echo ""
}

# System information function
system_info() {
    print_section "System Information"

    echo -e "${YELLOW}Operating System:${NC}"
    echo "  Distribution: $(lsb_release -ds 2>/dev/null || echo 'Unknown')"
    echo "  Kernel: $(uname -r)"
    echo "  Architecture: $(uname -m)"
    echo "  Uptime: $(uptime | cut -d',' -f1)"

    echo -e "\n${YELLOW}Hardware:${NC}"
    echo "  CPU: $(lscpu | grep 'Model name' | cut -d':' -f2 | xargs)"
    echo "  Memory: $(free -h | grep Mem | awk '{print $2 " total, " $3 " used, " $7 " available"}')"
    echo "  Disk: $(df -h / | tail -1 | awk '{print $2 " total, " $3 " used, " $4 " available"}')"

    echo -e "\n${YELLOW}Development Tools:${NC}"
    echo "  Go: $(go version 2>/dev/null | awk '{print $3}' || echo 'Not installed')"
    echo "  Node.js: $(node --version 2>/dev/null || echo 'Not installed')"
    echo "  npm: $(npm --version 2>/dev/null || echo 'Not installed')"
    echo "  VS Code: $(code --version 2>/dev/null | head -1 || echo 'Not installed')"
    echo "  Git: $(git --version 2>/dev/null || echo 'Not installed')"
    echo "  SQLite: $(sqlite3 --version 2>/dev/null | cut -d' ' -f1 || echo 'Not installed')"

    echo -e "\n${YELLOW}Environment Variables:${NC}"
    echo "  PATH: ${PATH:0:100}..."
    echo "  GOPATH: ${GOPATH:-'Not set'}"
    echo "  GOROOT: ${GOROOT:-'Not set'}"
    echo "  HOME: $HOME"
}

# Project health check function
project_health() {
    print_section "Project Health Check"

    echo -e "${YELLOW}Directory Structure:${NC}"
    if [ -d "$PROJECT_DIR" ]; then
        print_success "Project directory: $PROJECT_DIR"
        echo "  Size: $(du -sh "$PROJECT_DIR" | cut -f1)"
        echo "  Files: $(find "$PROJECT_DIR" -type f | wc -l)"
    else
        print_error "Project directory not found: $PROJECT_DIR"
        return 1
    fi

    echo -e "\n${YELLOW}Server Component:${NC}"
    if [ -d "$SERVER_DIR" ]; then
        print_success "Server directory exists"
        if [ -f "$SERVER_DIR/gorev" ]; then
            print_success "Server binary exists"
            echo "  Binary size: $(ls -lh "$SERVER_DIR/gorev" | awk '{print $5}')"
            echo "  Last modified: $(ls -l "$SERVER_DIR/gorev" | awk '{print $6, $7, $8}')"
        else
            print_error "Server binary not found"
        fi

        if [ -f "$SERVER_DIR/go.mod" ]; then
            print_success "Go module file exists"
            echo "  Module: $(grep '^module' "$SERVER_DIR/go.mod" | cut -d' ' -f2)"
            echo "  Go version: $(grep '^go ' "$SERVER_DIR/go.mod" | cut -d' ' -f2)"
        else
            print_error "go.mod not found"
        fi
    else
        print_error "Server directory not found"
    fi

    echo -e "\n${YELLOW}Extension Component:${NC}"
    if [ -d "$EXTENSION_DIR" ]; then
        print_success "Extension directory exists"
        if [ -f "$EXTENSION_DIR/package.json" ]; then
            print_success "package.json exists"
            echo "  Name: $(node -p "require('$EXTENSION_DIR/package.json').name" 2>/dev/null || echo 'Unknown')"
            echo "  Version: $(node -p "require('$EXTENSION_DIR/package.json').version" 2>/dev/null || echo 'Unknown')"
        else
            print_error "package.json not found"
        fi

        if [ -d "$EXTENSION_DIR/out" ]; then
            print_success "Compiled output exists"
            echo "  JS files: $(find "$EXTENSION_DIR/out" -name "*.js" | wc -l)"
        else
            print_error "Compiled output not found"
        fi

        VSIX_COUNT=$(find "$EXTENSION_DIR" -name "*.vsix" | wc -l)
        if [ $VSIX_COUNT -gt 0 ]; then
            print_success "VSIX packages found: $VSIX_COUNT"
            find "$EXTENSION_DIR" -name "*.vsix" -exec ls -lh {} \;
        else
            print_info "No VSIX packages found"
        fi
    else
        print_error "Extension directory not found"
    fi

    echo -e "\n${YELLOW}Database Status:${NC}"
    if [ -f "$HOME/.gorev/gorev.db" ]; then
        print_success "Global database exists"
        echo "  Location: $HOME/.gorev/gorev.db"
        echo "  Size: $(ls -lh "$HOME/.gorev/gorev.db" | awk '{print $5}')"
        echo "  Tables: $(echo '.tables' | sqlite3 "$HOME/.gorev/gorev.db" | wc -l)"
    else
        print_error "Global database not found"
    fi

    if [ -f "$PROJECT_DIR/.gorev/gorev.db" ]; then
        print_success "Workspace database exists"
        echo "  Location: $PROJECT_DIR/.gorev/gorev.db"
        echo "  Size: $(ls -lh "$PROJECT_DIR/.gorev/gorev.db" | awk '{print $5}')"
    else
        print_info "Workspace database not found (optional)"
    fi
}

# Database inspection function
database_inspection() {
    print_section "Database Inspection"

    local db_path="$HOME/.gorev/gorev.db"
    if [ ! -f "$db_path" ]; then
        print_error "Database not found: $db_path"
        return 1
    fi

    echo -e "${YELLOW}Database: $db_path${NC}"

    echo -e "\n${YELLOW}Tables:${NC}"
    echo '.tables' | sqlite3 "$db_path" | tr '\t' '\n' | sort

    echo -e "\n${YELLOW}Schema Overview:${NC}"
    for table in $(echo '.tables' | sqlite3 "$db_path" | tr '\t' ' '); do
        echo -e "\n${CYAN}Table: $table${NC}"
        echo ".schema $table" | sqlite3 "$db_path"
        row_count=$(echo "SELECT COUNT(*) FROM $table;" | sqlite3 "$db_path" 2>/dev/null || echo "0")
        echo "Rows: $row_count"
    done

    echo -e "\n${YELLOW}Sample Data:${NC}"
    echo -e "\n${CYAN}Projects:${NC}"
    echo "SELECT id, ad, aciklama FROM projeler LIMIT 5;" | sqlite3 -header -column "$db_path" 2>/dev/null || echo "No projects found"

    echo -e "\n${CYAN}Tasks:${NC}"
    echo "SELECT id, baslik, durum, oncelik FROM gorevler LIMIT 5;" | sqlite3 -header -column "$db_path" 2>/dev/null || echo "No tasks found"

    echo -e "\n${CYAN}Templates:${NC}"
    echo "SELECT id, ad, aciklama FROM gorev_templateleri LIMIT 5;" | sqlite3 -header -column "$db_path" 2>/dev/null || echo "No templates found"

    echo -e "\n${YELLOW}Database Integrity Check:${NC}"
    integrity_result=$(echo "PRAGMA integrity_check;" | sqlite3 "$db_path")
    if [ "$integrity_result" = "ok" ]; then
        print_success "Database integrity: OK"
    else
        print_error "Database integrity issues found: $integrity_result"
    fi
}

# Server debug mode function
server_debug() {
    print_section "Server Debug Mode"

    if [ ! -f "$SERVER_DIR/gorev" ]; then
        print_error "Server binary not found. Run build script first."
        return 1
    fi

    cd "$SERVER_DIR"

    echo -e "${YELLOW}Server Debug Options:${NC}"
    print_option "1" "Start server with debug logging"
    print_option "2" "Test server commands"
    print_option "3" "MCP protocol test"
    print_option "4" "Server with memory profiling"
    print_option "5" "Server with CPU profiling"
    print_option "0" "Back to main menu"

    read -p "Choose option: " debug_option

    case $debug_option in
        1)
            print_info "Starting server with debug logging (Ctrl+C to stop)..."
            ./gorev serve --debug --lang=tr
            ;;
        2)
            echo -e "\n${CYAN}Testing server commands:${NC}"
            echo "Version: $(./gorev version)"
            echo "Templates: $(./gorev template aliases)"
            echo "Help: $(./gorev help | head -5)"
            ;;
        3)
            print_info "Starting MCP protocol test..."
            echo "Starting server in background..."
            ./gorev serve --debug > "$DEBUG_LOG_DIR/mcp_test.log" 2>&1 &
            SERVER_PID=$!
            sleep 3

            if ps -p $SERVER_PID > /dev/null; then
                print_success "Server started (PID: $SERVER_PID)"
                print_info "Check logs: tail -f $DEBUG_LOG_DIR/mcp_test.log"
                read -p "Press Enter to stop server..."
                kill $SERVER_PID
                wait $SERVER_PID 2>/dev/null || true
                print_success "Server stopped"
            else
                print_error "Server failed to start"
            fi
            ;;
        4)
            if command -v valgrind >/dev/null 2>&1; then
                print_info "Starting server with memory profiling..."
                valgrind --tool=memcheck --leak-check=full --log-file="$DEBUG_LOG_DIR/memory_profile.log" ./gorev version
                print_success "Memory profile saved to: $DEBUG_LOG_DIR/memory_profile.log"
            else
                print_error "valgrind not installed"
            fi
            ;;
        5)
            print_info "CPU profiling requires Go pprof integration"
            print_info "Starting server with GODEBUG for basic profiling..."
            GODEBUG=gctrace=1 ./gorev version
            ;;
        0)
            return 0
            ;;
        *)
            print_error "Invalid option"
            ;;
    esac
}

# Extension debug mode function
extension_debug() {
    print_section "Extension Debug Mode"

    if [ ! -d "$EXTENSION_DIR" ]; then
        print_error "Extension directory not found"
        return 1
    fi

    cd "$EXTENSION_DIR"

    echo -e "${YELLOW}Extension Debug Options:${NC}"
    print_option "1" "Check extension installation"
    print_option "2" "Recompile extension"
    print_option "3" "Run extension tests"
    print_option "4" "Package extension"
    print_option "5" "Open extension in VS Code for debugging"
    print_option "6" "Check extension logs"
    print_option "0" "Back to main menu"

    read -p "Choose option: " ext_option

    case $ext_option in
        1)
            echo -e "\n${CYAN}Extension Installation Check:${NC}"
            if code --list-extensions | grep -q "mehmetsenol.gorev-vscode"; then
                print_success "Extension is installed"
                echo "Version: $(code --list-extensions --show-versions | grep gorev-vscode || echo 'Unknown')"
            else
                print_error "Extension is not installed"
            fi

            echo -e "\n${CYAN}Available Extensions:${NC}"
            code --list-extensions | grep -i gorev || echo "No Gorev extensions found"
            ;;
        2)
            print_info "Recompiling extension..."
            npm run compile
            print_success "Extension compiled"
            ;;
        3)
            print_info "Running extension tests..."
            if command -v xvfb-run >/dev/null 2>&1; then
                xvfb-run -a npm test
            else
                npm test
            fi
            ;;
        4)
            print_info "Packaging extension..."
            npm run package
            VSIX_FILE=$(find . -name "*.vsix" -type f | head -1)
            if [ -n "$VSIX_FILE" ]; then
                print_success "Package created: $VSIX_FILE"
                echo "Installing package..."
                code --install-extension "$VSIX_FILE"
            else
                print_error "Package creation failed"
            fi
            ;;
        5)
            print_info "Opening extension in VS Code for debugging..."
            code .
            print_info "Press F5 in VS Code to start debugging"
            ;;
        6)
            print_info "Checking extension logs..."
            if [ -d "$HOME/.vscode/logs" ]; then
                echo "VS Code logs directory: $HOME/.vscode/logs"
                find "$HOME/.vscode/logs" -name "*gorev*" -type f 2>/dev/null || echo "No Gorev logs found"
            else
                print_info "VS Code logs directory not found"
            fi
            ;;
        0)
            return 0
            ;;
        *)
            print_error "Invalid option"
            ;;
    esac
}

# Log analysis function
log_analysis() {
    print_section "Log Analysis"

    echo -e "${YELLOW}Available Log Sources:${NC}"

    # System logs
    echo -e "\n${CYAN}System Logs:${NC}"
    if [ -f "/var/log/syslog" ]; then
        echo "Recent system errors:"
        tail -50 /var/log/syslog | grep -i error | tail -5 || echo "No recent errors"
    fi

    # Application logs
    echo -e "\n${CYAN}Application Logs:${NC}"
    if [ -d "$DEBUG_LOG_DIR" ]; then
        echo "Debug logs directory: $DEBUG_LOG_DIR"
        find "$DEBUG_LOG_DIR" -name "*.log" -type f | while read log_file; do
            echo "  $(basename "$log_file"): $(wc -l < "$log_file") lines"
        done
    else
        echo "No debug logs found"
    fi

    # VS Code logs
    echo -e "\n${CYAN}VS Code Logs:${NC}"
    if [ -d "$HOME/.vscode/logs" ]; then
        LATEST_LOG_DIR=$(find "$HOME/.vscode/logs" -type d -name "*" | sort | tail -1)
        if [ -n "$LATEST_LOG_DIR" ]; then
            echo "Latest VS Code log session: $LATEST_LOG_DIR"
            find "$LATEST_LOG_DIR" -name "*.log" -type f | head -5
        fi
    fi

    # Process information
    echo -e "\n${CYAN}Running Processes:${NC}"
    ps aux | grep -E "(gorev|code)" | grep -v grep || echo "No relevant processes running"

    # Network connections
    echo -e "\n${CYAN}Network Connections:${NC}"
    ss -tuln | grep -E ":(8080|3000|4000)" || echo "No development servers running"
}

# Network and process check function
network_process_check() {
    print_section "Network & Process Check"

    echo -e "${YELLOW}Process Information:${NC}"
    echo -e "\n${CYAN}Gorev Processes:${NC}"
    ps aux | grep gorev | grep -v grep || echo "No Gorev processes running"

    echo -e "\n${CYAN}VS Code Processes:${NC}"
    ps aux | grep code | grep -v grep | head -5 || echo "No VS Code processes running"

    echo -e "\n${CYAN}Memory Usage:${NC}"
    free -h

    echo -e "\n${CYAN}CPU Usage:${NC}"
    top -bn1 | head -10

    echo -e "\n${YELLOW}Network Information:${NC}"
    echo -e "\n${CYAN}Network Interfaces:${NC}"
    ip addr show | grep -E "(inet |inet6)" | head -5

    echo -e "\n${CYAN}Listening Ports:${NC}"
    ss -tuln | head -10

    echo -e "\n${CYAN}DNS Resolution:${NC}"
    nslookup github.com | head -5 2>/dev/null || echo "DNS resolution test failed"
}

# Performance profiling function
performance_profiling() {
    print_section "Performance Profiling"

    if [ ! -f "$SERVER_DIR/gorev" ]; then
        print_error "Server binary not found"
        return 1
    fi

    cd "$SERVER_DIR"

    echo -e "${YELLOW}Performance Tests:${NC}"

    echo -e "\n${CYAN}Command Performance:${NC}"
    echo "Testing 'gorev version' command:"
    time ./gorev version

    echo -e "\n${CYAN}Database Query Performance:${NC}"
    echo "Testing database query:"
    time echo "SELECT COUNT(*) FROM sqlite_master;" | sqlite3 "$HOME/.gorev/gorev.db"

    echo -e "\n${CYAN}Memory Usage:${NC}"
    echo "Server binary size: $(ls -lh ./gorev | awk '{print $5}')"

    if command -v valgrind >/dev/null 2>&1; then
        echo -e "\n${CYAN}Memory Leak Check:${NC}"
        print_info "Running memory leak check (this may take a moment)..."
        valgrind --tool=memcheck --leak-check=summary --show-leak-kinds=definite ./gorev version 2>&1 | grep -E "(definitely lost|ERROR SUMMARY)"
    else
        print_info "Install valgrind for memory leak checking: sudo apt install valgrind"
    fi

    echo -e "\n${CYAN}Disk Usage:${NC}"
    echo "Project size: $(du -sh "$PROJECT_DIR")"
    echo "Database sizes:"
    [ -f "$HOME/.gorev/gorev.db" ] && echo "  Global: $(ls -lh "$HOME/.gorev/gorev.db" | awk '{print $5}')"
    [ -f "$PROJECT_DIR/.gorev/gorev.db" ] && echo "  Workspace: $(ls -lh "$PROJECT_DIR/.gorev/gorev.db" | awk '{print $5}')"
}

# Clean and reset function
clean_reset() {
    print_section "Clean & Reset"

    echo -e "${YELLOW}Clean & Reset Options:${NC}"
    print_option "1" "Clean build artifacts"
    print_option "2" "Reset databases"
    print_option "3" "Clean extension build"
    print_option "4" "Clean all logs"
    print_option "5" "Full reset (keep source code)"
    print_option "0" "Back to main menu"

    read -p "Choose option: " clean_option

    case $clean_option in
        1)
            print_info "Cleaning build artifacts..."
            cd "$SERVER_DIR"
            make clean 2>/dev/null || rm -f ./gorev
            print_success "Build artifacts cleaned"
            ;;
        2)
            read -p "This will delete all databases. Are you sure? (y/N): " confirm
            if [[ $confirm =~ ^[Yy]$ ]]; then
                print_info "Resetting databases..."
                rm -f "$HOME/.gorev/gorev.db"
                rm -f "$PROJECT_DIR/.gorev/gorev.db"
                cd "$SERVER_DIR"
                ./gorev init --global 2>/dev/null || print_error "Failed to reinitialize global database"
                ./gorev init 2>/dev/null || print_error "Failed to reinitialize workspace database"
                print_success "Databases reset"
            fi
            ;;
        3)
            print_info "Cleaning extension build..."
            cd "$EXTENSION_DIR"
            rm -rf out node_modules *.vsix
            npm install
            npm run compile
            print_success "Extension build cleaned and recompiled"
            ;;
        4)
            print_info "Cleaning logs..."
            rm -rf "$DEBUG_LOG_DIR"
            mkdir -p "$DEBUG_LOG_DIR"
            print_success "Logs cleaned"
            ;;
        5)
            read -p "This will reset everything except source code. Are you sure? (y/N): " confirm
            if [[ $confirm =~ ^[Yy]$ ]]; then
                print_info "Performing full reset..."
                cd "$SERVER_DIR"
                make clean 2>/dev/null || rm -f ./gorev
                cd "$EXTENSION_DIR"
                rm -rf out node_modules *.vsix
                rm -f "$HOME/.gorev/gorev.db"
                rm -f "$PROJECT_DIR/.gorev/gorev.db"
                rm -rf "$DEBUG_LOG_DIR"
                print_success "Full reset completed. Run build scripts to restore."
            fi
            ;;
        0)
            return 0
            ;;
        *)
            print_error "Invalid option"
            ;;
    esac
}

# Generate debug report function
generate_debug_report() {
    print_section "Generate Debug Report"

    local report_file="$DEBUG_LOG_DIR/debug_report_$(date +%Y%m%d_%H%M%S).txt"

    print_info "Generating comprehensive debug report..."

    {
        echo "Gorev Debug Report"
        echo "=================="
        echo "Generated: $(date)"
        echo "Hostname: $(hostname)"
        echo ""

        echo "SYSTEM INFORMATION"
        echo "=================="
        echo "OS: $(lsb_release -ds 2>/dev/null || echo 'Unknown')"
        echo "Kernel: $(uname -r)"
        echo "Architecture: $(uname -m)"
        echo "Memory: $(free -h | grep Mem)"
        echo "Disk: $(df -h /)"
        echo ""

        echo "DEVELOPMENT TOOLS"
        echo "================="
        echo "Go: $(go version 2>/dev/null || echo 'Not installed')"
        echo "Node: $(node --version 2>/dev/null || echo 'Not installed')"
        echo "npm: $(npm --version 2>/dev/null || echo 'Not installed')"
        echo "VS Code: $(code --version 2>/dev/null | head -1 || echo 'Not installed')"
        echo ""

        echo "PROJECT STATUS"
        echo "=============="
        echo "Project dir: $PROJECT_DIR"
        echo "Server binary: $([ -f "$SERVER_DIR/gorev" ] && echo "EXISTS" || echo "MISSING")"
        echo "Extension compiled: $([ -d "$EXTENSION_DIR/out" ] && echo "YES" || echo "NO")"
        echo "Global DB: $([ -f "$HOME/.gorev/gorev.db" ] && echo "EXISTS" || echo "MISSING")"
        echo "Workspace DB: $([ -f "$PROJECT_DIR/.gorev/gorev.db" ] && echo "EXISTS" || echo "MISSING")"
        echo ""

        echo "EXTENSION STATUS"
        echo "================"
        if code --list-extensions | grep -q "mehmetsenol.gorev-vscode"; then
            echo "Extension installed: YES"
        else
            echo "Extension installed: NO"
        fi
        echo ""

        echo "RUNNING PROCESSES"
        echo "================="
        ps aux | grep -E "(gorev|code)" | grep -v grep || echo "No relevant processes"
        echo ""

        echo "RECENT ERRORS"
        echo "============="
        if [ -f "/var/log/syslog" ]; then
            echo "System errors (last 24h):"
            grep "$(date +'%b %d')" /var/log/syslog 2>/dev/null | grep -i error | tail -5 || echo "No errors found"
        fi
        echo ""

        echo "NETWORK STATUS"
        echo "=============="
        echo "Network interfaces:"
        ip addr show 2>/dev/null | grep -E "(inet |inet6)" | head -3
        echo ""
        echo "Listening ports:"
        ss -tuln | head -5
        echo ""

    } > "$report_file"

    print_success "Debug report generated: $report_file"
    echo -e "${YELLOW}Report contents:${NC}"
    echo "  - System information"
    echo "  - Development tools status"
    echo "  - Project health check"
    echo "  - Running processes"
    echo "  - Network status"
    echo "  - Recent errors"

    read -p "View the report now? (y/N): " view_report
    if [[ $view_report =~ ^[Yy]$ ]]; then
        less "$report_file"
    fi
}

# Main loop
main() {
    while true; do
        show_menu
        read -p "Choose an option (0-10): " choice

        case $choice in
            1) system_info ;;
            2) project_health ;;
            3) database_inspection ;;
            4) server_debug ;;
            5) extension_debug ;;
            6) log_analysis ;;
            7) network_process_check ;;
            8) performance_profiling ;;
            9) clean_reset ;;
            10) generate_debug_report ;;
            0)
                echo -e "\n${GREEN}Debug helper session ended.${NC}"
                exit 0
                ;;
            *)
                print_error "Invalid option. Please choose 0-10."
                ;;
        esac

        echo ""
        read -p "Press Enter to continue..."
    done
}

# Check if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi