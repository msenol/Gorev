#!/bin/bash

# Documentation Version Update Script
# Usage: ./scripts/update-docs-version.sh NEW_VERSION [--dry-run]

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
DRY_RUN=false

# Usage function
usage() {
    echo "Usage: $0 NEW_VERSION [--dry-run]"
    echo "  NEW_VERSION  Version to update to (e.g., 0.15.25)"
    echo "  --dry-run    Show what would be changed without making changes"
    echo ""
    echo "Examples:"
    echo "  $0 0.15.25"
    echo "  $0 0.15.25 --dry-run"
    exit 1
}

# Parse arguments
if [[ $# -lt 1 ]]; then
    usage
fi

NEW_VERSION="$1"
shift

while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Validate version format
if ! [[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Version must be in format X.Y.Z (e.g., 0.15.25)${NC}"
    exit 1
fi

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

log_dry_run() {
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${YELLOW}[DRY RUN]${NC} $1"
    fi
}

# Update version in file
update_version_in_file() {
    local file="$1"
    local pattern="$2"
    local replacement="$3"
    local description="$4"

    if [[ ! -f "$file" ]]; then
        log_warning "File not found: $file"
        return 0
    fi

    if grep -q "$pattern" "$file"; then
        if [[ "$DRY_RUN" == true ]]; then
            log_dry_run "Would update $description in $file"
            echo "  $(grep "$pattern" "$file" | head -1)"
            echo "  → $(echo "$(grep "$pattern" "$file" | head -1)" | sed "$replacement")"
        else
            sed -i "$replacement" "$file"
            log_success "Updated $description in $file"
        fi
    else
        log_warning "Pattern not found in $file: $pattern"
    fi
}

# Main update function
update_documentation_versions() {
    log_info "Updating documentation versions to v$NEW_VERSION"

    cd "$PROJECT_ROOT"

    # Files to update with their patterns
    declare -A FILES_TO_UPDATE=(
        ["docs/tr/kullanim.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|Version validation"
        ["docs/tr/kurulum.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|Version validation"
        ["docs/api/api-referans.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|API version"
        ["docs/development/gelistirme.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g|Development guide version"
        ["docs/development/testing-strategy.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g|Testing strategy version"
        ["docs/development/contributing.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|Contributing guide version"
        ["docs/guides/getting-started/installation.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|Installation guide version"
        ["docs/api/reference.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\++/v$NEW_VERSION+/g|API reference version"
        ["docs/README.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g|Documentation README version"
        ["docs/tr/README.md"]="s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g|Turkish README version"
        ["package.json"]="s/\"version\":\s*\"[0-9]\+\.[0-9]\+\.[0-9]\+\"/\"version\": \"$NEW_VERSION\"/g|Package.json version"
    )

    local updated_files=()

    for file in "${!FILES_TO_UPDATE[@]}"; do
        IFS='|' read -r sed_pattern description <<< "${FILES_TO_UPDATE[$file]}"

        if [[ -f "$file" ]]; then
            if grep -qE "v?[0-9]+\.[0-9]+\.[0-9]+\+?" "$file"; then
                update_version_in_file "$file" "v[0-9]\+\.[0-9]\+\.[0-9]\+" "s/$sed_pattern" "$description"
                updated_files+=("$file")
            fi
        fi
    done

    # Update installation commands with specific version
    local install_files=(
        "docs/tr/kurulum.md"
        "docs/guides/getting-started/installation.md"
    )

    for file in "${install_files[@]}"; do
        if [[ -f "$file" ]]; then
            # Update version in download URLs
            local url_pattern="s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g"
            update_version_in_file "$file" "v[0-9]\+\.[0-9]\+\.[0-9]\+" "$url_pattern" "download URLs"
        fi
    done

    log_info "Version update summary:"
    echo "  Target version: v$NEW_VERSION"
    echo "  Files processed: ${#updated_files[@]}"

    if [[ ${#updated_files[@]} -gt 0 ]]; then
        echo "  Updated files:"
        printf '    %s\n' "${updated_files[@]}"
    fi
}

# Validate changes
validate_changes() {
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Dry run completed - no changes made"
        return 0
    fi

    log_info "Validating changes..."

    cd "$PROJECT_ROOT"

    # Check if version was actually updated
    local found_versions=()

    while IFS= read -r -d '' file; do
        if grep -q "v$NEW_VERSION" "$file"; then
            found_versions+=("$file")
        fi
    done < <(find docs -name "*.md" -print0)

    if [[ ${#found_versions[@]} -gt 0 ]]; then
        log_success "Version v$NEW_VERSION found in ${#found_versions[@]} files"
    else
        log_error "Version v$NEW_VERSION not found in any documentation files"
        return 1
    fi

    # Check for remaining old version references (basic check)
    local old_version_pattern="v0\.15\.(0[0-9]|1[0-9]|2[0-3])[^0-9]"
    local files_with_old_versions=()

    while IFS= read -r -d '' file; do
        if grep -qE "$old_version_pattern" "$file"; then
            files_with_old_versions+=("$file")
        fi
    done < <(find docs -name "*.md" -print0)

    if [[ ${#files_with_old_versions[@]} -gt 0 ]]; then
        log_warning "Found potential old version references in:"
        printf '  %s\n' "${files_with_old_versions[@]}"
        log_info "Review these files manually to ensure all versions are updated correctly"
    fi

    log_success "Validation completed"
}

# Generate commit message
generate_commit_message() {
    if [[ "$DRY_RUN" == true ]]; then
        return 0
    fi

    local commit_msg="docs: update documentation versions to v$NEW_VERSION

- Updated version references in documentation files
- Synchronized version numbers across guides and references
- Updated installation commands with new version
- Maintained backward compatibility notes where appropriate

Automated update via update-docs-version.sh script"

    echo "$commit_msg" > /tmp/docs-version-update-commit.txt

    log_info "Suggested commit message saved to /tmp/docs-version-update-commit.txt"
    log_info "Use: git commit -F /tmp/docs-version-update-commit.txt"
}

# Main execution
main() {
    log_info "Starting documentation version update..."
    log_info "Project root: $PROJECT_ROOT"
    log_info "Target version: v$NEW_VERSION"
    log_info "Dry run mode: $DRY_RUN"
    echo

    update_documentation_versions
    echo

    validate_changes
    echo

    generate_commit_message

    if [[ "$DRY_RUN" == true ]]; then
        log_info "Dry run completed. Run without --dry-run to make changes."
    else
        log_success "Documentation version update completed successfully!"
        log_info "Next steps:"
        echo "  1. Review the changes: git diff"
        echo "  2. Commit the changes: git commit -F /tmp/docs-version-update-commit.txt"
        echo "  3. Validate with: ./scripts/validate-docs.sh"
    fi
}

# Execute main function
main "$@"