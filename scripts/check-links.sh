#!/bin/bash
#
# Link Checker for Gorev Documentation
# Validates all internal and external links in markdown files
#
# Usage:
#   ./scripts/check-links.sh [--fix]
#
# Options:
#   --fix    Attempt to fix broken internal links automatically
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_LINKS=0
BROKEN_LINKS=0
FIXED_LINKS=0
EXTERNAL_LINKS=0

# Options
FIX_MODE=false
if [ "$1" == "--fix" ]; then
    FIX_MODE=true
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Gorev Documentation Link Checker${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Project root: $PROJECT_ROOT"
echo "Fix mode: $FIX_MODE"
echo ""

# Function to extract links from markdown
extract_links() {
    local file="$1"
    # Extract markdown links: [text](url)
    grep -oP '\[.*?\]\(\K[^)]+' "$file" 2>/dev/null || true
}

# Function to check if file exists (relative to project root)
check_internal_link() {
    local link="$1"
    local source_file="$2"
    local source_dir="$(dirname "$source_file")"

    # Remove anchor
    local file_path="${link%%#*}"

    # Skip if empty (anchor only)
    [ -z "$file_path" ] && return 0

    # Check if absolute path from project root
    if [[ "$file_path" == /* ]]; then
        if [ -f "$PROJECT_ROOT$file_path" ]; then
            return 0
        else
            return 1
        fi
    fi

    # Check relative path
    local full_path="$source_dir/$file_path"
    if [ -f "$full_path" ]; then
        return 0
    fi

    # Check from project root
    if [ -f "$PROJECT_ROOT/$file_path" ]; then
        return 0
    fi

    return 1
}

# Function to attempt fixing a broken link
fix_broken_link() {
    local link="$1"
    local source_file="$2"

    # Extract filename from link
    local filename="$(basename "$link" | cut -d'#' -f1)"

    # Search for file in project
    local found_files=$(find "$PROJECT_ROOT" -name "$filename" -type f \
        -not -path "*/node_modules/*" \
        -not -path "*/.vscode-test/*" \
        -not -path "*/build/*" 2>/dev/null)

    local count=$(echo "$found_files" | grep -c '.' || true)

    if [ "$count" -eq 1 ]; then
        # Found exactly one match - can fix confidently
        local found_path="$found_files"
        local source_dir="$(dirname "$source_file")"
        local relative_path=$(realpath --relative-to="$source_dir" "$found_path")

        echo -e "    ${GREEN}✓${NC} Found at: $relative_path"

        if [ "$FIX_MODE" = true ]; then
            # Replace the link in file
            sed -i "s|$link|$relative_path|g" "$source_file"
            echo -e "    ${GREEN}✓${NC} Fixed in file"
            FIXED_LINKS=$((FIXED_LINKS + 1))
            return 0
        else
            echo -e "    ${YELLOW}!${NC} Run with --fix to update automatically"
            return 1
        fi
    elif [ "$count" -gt 1 ]; then
        echo -e "    ${YELLOW}!${NC} Multiple matches found:"
        echo "$found_files" | while read -r match; do
            echo "      - $match"
        done
        return 1
    else
        echo -e "    ${RED}✗${NC} File not found in project"
        return 1
    fi
}

# Function to check external link
check_external_link() {
    local url="$1"

    # Use curl with timeout
    if curl -s -f -L --max-time 5 --head "$url" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

echo -e "${BLUE}Scanning markdown files...${NC}"
echo ""

# Find all markdown files
MD_FILES=$(find "$PROJECT_ROOT" -name "*.md" -type f \
    -not -path "*/node_modules/*" \
    -not -path "*/.vscode-test/*" \
    -not -path "*/build/*" \
    -not -path "*/dist/*" \
    | sort)

FILE_COUNT=$(echo "$MD_FILES" | wc -l)
echo "Found $FILE_COUNT markdown files to check"
echo ""

# Check each file
while IFS= read -r file; do
    relative_file="${file#$PROJECT_ROOT/}"

    # Extract links from file
    links=$(extract_links "$file")

    if [ -z "$links" ]; then
        continue
    fi

    has_broken=false

    while IFS= read -r link; do
        [ -z "$link" ] && continue

        TOTAL_LINKS=$((TOTAL_LINKS + 1))

        # Categorize link
        if [[ "$link" =~ ^https?:// ]]; then
            # External link
            EXTERNAL_LINKS=$((EXTERNAL_LINKS + 1))
            # Skip external link checking by default (too slow)
            continue
        elif [[ "$link" =~ ^mailto: ]] || [[ "$link" =~ ^# ]]; then
            # Skip mailto and anchor-only links
            continue
        else
            # Internal link
            if ! check_internal_link "$link" "$file"; then
                if [ "$has_broken" = false ]; then
                    echo -e "${YELLOW}$relative_file${NC}"
                    has_broken=true
                fi

                echo -e "  ${RED}✗${NC} Broken: $link"
                BROKEN_LINKS=$((BROKEN_LINKS + 1))

                # Attempt fix
                fix_broken_link "$link" "$file"
            fi
        fi
    done <<< "$links"

done <<< "$MD_FILES"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Total links checked: $TOTAL_LINKS"
echo "External links (skipped): $EXTERNAL_LINKS"
echo -e "${RED}Broken links: $BROKEN_LINKS${NC}"

if [ "$FIX_MODE" = true ] && [ "$FIXED_LINKS" -gt 0 ]; then
    echo -e "${GREEN}Fixed links: $FIXED_LINKS${NC}"
fi

echo ""

if [ "$BROKEN_LINKS" -gt 0 ]; then
    echo -e "${RED}✗ Link check failed${NC}"
    echo ""
    echo "To attempt automatic fixes, run:"
    echo "  ./scripts/check-links.sh --fix"
    echo ""
    exit 1
else
    echo -e "${GREEN}✓ All links are valid${NC}"
    exit 0
fi
