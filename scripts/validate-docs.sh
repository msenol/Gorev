#!/bin/bash

# Documentation Quality Validation Script
# Usage: ./scripts/validate-docs.sh [--fix] [--verbose]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
FIX_MODE=false
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --fix)
            FIX_MODE=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [--fix] [--verbose]"
            echo "  --fix     Automatically fix issues where possible"
            echo "  --verbose Show detailed output"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_verbose() {
    if [[ "$VERBOSE" == true ]]; then
        echo -e "${BLUE}[VERBOSE]${NC} $1"
    fi
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()

    if ! command -v markdownlint &> /dev/null; then
        missing_deps+=("markdownlint-cli")
    fi

    if ! command -v markdown-link-check &> /dev/null; then
        missing_deps+=("markdown-link-check")
    fi

    if ! command -v cspell &> /dev/null; then
        missing_deps+=("cspell")
    fi

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Install them with: npm install -g ${missing_deps[*]}"
        return 1
    fi

    log_success "All dependencies are installed"
    return 0
}

# Validate markdown files
validate_markdown() {
    log_info "Validating markdown files..."

    cd "$PROJECT_ROOT"

    local markdownlint_args=(
        "**/*.md"
        "--ignore" "node_modules"
        "--ignore" ".git"
        "--config" ".markdownlint.json"
    )

    if [[ "$FIX_MODE" == true ]]; then
        markdownlint_args+=("--fix")
        log_info "Fixing markdown issues..."
    fi

    if [[ "$VERBOSE" == true ]]; then
        markdownlint "${markdownlint_args[@]}"
    else
        markdownlint "${markdownlint_args[@]}" 2>/dev/null || {
            log_error "Markdown linting failed"
            return 1
        }
    fi

    log_success "Markdown validation completed"
    return 0
}

# Check for broken links
check_links() {
    log_info "Checking for broken links..."

    cd "$PROJECT_ROOT"

    local link_check_config='{
        "ignorePatterns": [
            {
                "pattern": "^https://github.com/msenol/Gorev/(issues|discussions|releases|blob|tree)"
            },
            {
                "pattern": "^https://marketplace.visualstudio.com"
            },
            {
                "pattern": "^http://localhost"
            },
            {
                "pattern": "^https://registry.modelcontextprotocol.io"
            }
        ],
        "timeout": "20s",
        "retryOn429": true,
        "retryCount": 3,
        "fallbackRetryDelay": "30s"
    }'

    # Create temporary config file
    echo "$link_check_config" > /tmp/link-check-config.json

    local failed_files=()

    # Check links in key documentation files
    for file in "README.md" "docs/"**"/*.md"; do
        if [[ -f "$file" ]]; then
            log_verbose "Checking links in $file"
            if ! markdown-link-check "$file" --config /tmp/link-check-config.json --quiet; then
                failed_files+=("$file")
            fi
        fi
    done

    rm -f /tmp/link-check-config.json

    if [[ ${#failed_files[@]} -gt 0 ]]; then
        log_error "Link check failed for: ${failed_files[*]}"
        return 1
    fi

    log_success "Link check completed"
    return 0
}

# Check spelling
check_spelling() {
    log_info "Checking spelling..."

    cd "$PROJECT_ROOT"

    if cspell "**/*.md" --config cspell.json --no-progress --quiet; then
        log_success "Spelling check completed"
        return 0
    else
        log_error "Spelling check failed"
        if [[ "$VERBOSE" == true ]]; then
            log_info "Running detailed spelling check..."
            cspell "**/*.md" --config cspell.json --no-progress
        fi
        return 1
    fi
}

# Check version consistency
check_version_consistency() {
    log_info "Checking version consistency..."

    cd "$PROJECT_ROOT"

    # Extract version from different files
    local makefile_version=""
    local package_json_version=""
    local server_json_version=""
    local changelog_version=""

    if [[ -f "gorev-mcpserver/Makefile" ]]; then
        makefile_version=$(grep "^VERSION" gorev-mcpserver/Makefile | cut -d'=' -f2 | tr -d ' ' | tr -d 'v')
    fi

    if [[ -f "gorev-npm/package.json" ]]; then
        package_json_version=$(grep '"version"' gorev-npm/package.json | cut -d'"' -f4)
    fi

    if [[ -f "server.json" ]]; then
        server_json_version=$(grep '"version"' server.json | cut -d'"' -f4)
    fi

    if [[ -f "CHANGELOG.md" ]]; then
        changelog_version=$(grep -E "^## \[v[0-9]+\.[0-9]+\.[0-9]+\]" CHANGELOG.md | head -1 | sed -E 's/^## \[v([0-9]+\.[0-9]+\.[0-9]+)\].*/\1/')
    fi

    log_verbose "Found versions:"
    log_verbose "  Makefile: $makefile_version"
    log_verbose "  package.json: $package_json_version"
    log_verbose "  server.json: $server_json_version"
    log_verbose "  CHANGELOG.md: $changelog_version"

    # Check if all versions match
    local versions=("$makefile_version" "$package_json_version" "$server_json_version" "$changelog_version")
    local first_version="${versions[0]}"
    local inconsistent=false

    for version in "${versions[@]}"; do
        if [[ "$version" != "$first_version" && -n "$version" ]]; then
            inconsistent=true
            break
        fi
    done

    if [[ "$inconsistent" == true ]]; then
        log_error "Version inconsistency detected"
        log_error "  Makefile: $makefile_version"
        log_error "  package.json: $package_json_version"
        log_error "  server.json: $server_json_version"
        log_error "  CHANGELOG.md: $changelog_version"
        return 1
    fi

    log_success "Version consistency check completed (v$first_version)"
    return 0
}

# Check file structure
check_file_structure() {
    log_info "Checking documentation file structure..."

    cd "$PROJECT_ROOT"

    local required_files=(
        "README.md"
        "CHANGELOG.md"
        "LICENSE"
        "docs/tr/kurulum.md"
        "docs/tr/kullanim.md"
        "docs/tr/mcp-araclari.md"
        "docs/api/api-referans.md"
        "docs/development/gelistirme.md"
        ".markdownlint.json"
        "cspell.json"
    )

    local missing_files=()

    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            missing_files+=("$file")
        fi
    done

    if [[ ${#missing_files[@]} -gt 0 ]]; then
        log_error "Missing required files: ${missing_files[*]}"
        return 1
    fi

    log_success "File structure check completed"
    return 0
}

# Main execution
main() {
    log_info "Starting documentation quality validation..."
    log_info "Project root: $PROJECT_ROOT"
    log_info "Fix mode: $FIX_MODE"
    log_info "Verbose mode: $VERBOSE"
    echo

    local exit_code=0
    local checks_passed=0
    local total_checks=6

    # Run checks
    if check_dependencies; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    if check_file_structure; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    if validate_markdown; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    if check_version_consistency; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    if check_spelling; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    if check_links; then
        ((checks_passed++))
    else
        exit_code=1
    fi

    echo
    log_info "Documentation validation completed"
    log_info "Checks passed: $checks_passed/$total_checks"

    if [[ $exit_code -eq 0 ]]; then
        log_success "All documentation quality checks passed!"
    else
        log_error "Some documentation quality checks failed"
        if [[ "$FIX_MODE" == false ]]; then
            log_info "Run with --fix to automatically fix issues where possible"
        fi
    fi

    return $exit_code
}

# Execute main function
main "$@"