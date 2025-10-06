#!/bin/bash
#
# Markdown Linter for Gorev Documentation
# Uses markdownlint-cli to enforce consistent markdown style
#
# Usage:
#   ./scripts/lint-docs.sh [--fix]
#
# Options:
#   --fix    Automatically fix linting issues where possible
#
# Installation:
#   npm install -g markdownlint-cli
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Options
FIX_MODE=false
if [ "$1" == "--fix" ]; then
    FIX_MODE=true
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Gorev Documentation Linter${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Project root: $PROJECT_ROOT"
echo "Fix mode: $FIX_MODE"
echo ""

# Check if markdownlint-cli is installed
if ! command -v markdownlint &> /dev/null; then
    echo -e "${RED}✗ markdownlint-cli is not installed${NC}"
    echo ""
    echo "Please install it globally with:"
    echo "  npm install -g markdownlint-cli"
    echo ""
    echo "Or run locally with npx:"
    echo "  npx markdownlint-cli '**/*.md'"
    echo ""
    exit 1
fi

echo -e "${BLUE}Linting markdown files...${NC}"
echo ""

# Find all markdown files to lint
MD_FILES=$(find "$PROJECT_ROOT" -name "*.md" -type f \
    -not -path "*/node_modules/*" \
    -not -path "*/.vscode-test/*" \
    -not -path "*/build/*" \
    -not -path "*/dist/*" \
    -not -path "*/gorev-npm/*" \
    | sort)

FILE_COUNT=$(echo "$MD_FILES" | wc -l)
echo "Found $FILE_COUNT markdown files to lint"
echo ""

# Run markdownlint
if [ "$FIX_MODE" = true ]; then
    echo -e "${YELLOW}Running in FIX mode - attempting automatic fixes...${NC}"
    echo ""

    if markdownlint --config .markdownlint.json --fix $MD_FILES; then
        echo ""
        echo -e "${GREEN}✓ All linting issues fixed successfully${NC}"
        exit 0
    else
        EXIT_CODE=$?
        echo ""
        echo -e "${RED}✗ Some linting issues could not be auto-fixed${NC}"
        echo ""
        echo "Please review the errors above and fix manually."
        echo "You can also check the documentation:"
        echo "  https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md"
        echo ""
        exit $EXIT_CODE
    fi
else
    # Lint without fixing
    if markdownlint --config .markdownlint.json $MD_FILES; then
        echo ""
        echo -e "${GREEN}✓ All markdown files are compliant${NC}"
        exit 0
    else
        EXIT_CODE=$?
        echo ""
        echo -e "${RED}✗ Linting failed${NC}"
        echo ""
        echo "To attempt automatic fixes, run:"
        echo "  ./scripts/lint-docs.sh --fix"
        echo ""
        echo "For rule details, see:"
        echo "  https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md"
        echo ""
        exit $EXIT_CODE
    fi
fi
