# Markdown Linting Guide

**Version:** v0.16.3 | **Last Updated:** October 6, 2025

This document describes the markdown linting setup for Gorev project documentation.

---

## Overview

Gorev uses [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli) to enforce consistent markdown formatting across all documentation files.

**Configuration File:** `.markdownlint.json`
**Linting Script:** `scripts/lint-docs.sh`

---

## Installation

### Global Installation (Recommended)

```bash
npm install -g markdownlint-cli
```

### Local Usage (npx)

```bash
npx markdownlint-cli '**/*.md'
```

---

## Usage

### Check All Documentation

```bash
./scripts/lint-docs.sh
```

### Auto-Fix Issues

```bash
./scripts/lint-docs.sh --fix
```

---

## Configuration Rules

### Enabled Rules

| Rule | Description | Setting |
|------|-------------|---------|
| MD001 | Heading levels increment by one | ✅ Enabled |
| MD003 | Heading style | ATX (`#` style) |
| MD004 | Unordered list style | Dash (`-`) |
| MD007 | List indentation | 2 spaces |
| MD009 | Trailing spaces | Max 2 (for line breaks) |
| MD012 | Multiple blank lines | Max 2 consecutive |
| MD013 | Line length | 120 characters |
| MD022 | Blank lines around headings | 1 line above/below |
| MD024 | Duplicate headings | Siblings only |
| MD029 | Ordered list numbering | Sequential |
| MD032 | Blank lines around lists | Required |
| MD040 | Code block language | Required |
| MD046 | Code block style | Fenced (```) |
| MD047 | File end newline | Required |
| MD048 | Code fence style | Backtick |
| MD049 | Emphasis style | Asterisk |
| MD050 | Strong emphasis style | Asterisk |

### Disabled Rules

| Rule | Reason |
|------|--------|
| MD034 | Bare URLs allowed (useful for references) |
| MD036 | Emphasis as heading allowed (stylistic choice) |
| MD051 | Link fragments not validated (too strict) |

### Allowed HTML Elements

The following HTML elements are permitted in markdown files:

- `<br>` - Line breaks
- `<details>`, `<summary>` - Collapsible sections
- `<sup>`, `<sub>` - Superscript/subscript
- `<div>` - Container elements
- `<img>`, `<a>` - Images and links (when markdown syntax insufficient)
- `<b>` - Bold text (when markdown insufficient)
- `<kbd>` - Keyboard shortcuts

---

## Common Issues and Fixes

### MD040: Missing Code Block Language

**Problem:**
```markdown
```
code here
```
```

**Fix:**
```markdown
```bash
code here
```
```

### MD013: Line Too Long

**Problem:**
```markdown
This is a very long line that exceeds 120 characters and should be broken into multiple lines for better readability and maintainability.
```

**Fix:**
```markdown
This is a very long line that exceeds 120 characters and should be broken
into multiple lines for better readability and maintainability.
```

### MD022: Missing Blank Lines Around Headings

**Problem:**
```markdown
Some text here
## Heading
More text
```

**Fix:**
```markdown
Some text here

## Heading

More text
```

### MD047: Missing Final Newline

**Problem:**
```markdown
Last line without newline
```

**Fix:**
```markdown
Last line with newline

```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Lint Documentation

on: [push, pull_request]

jobs:
  lint-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install markdownlint-cli
        run: npm install -g markdownlint-cli

      - name: Lint markdown files
        run: ./scripts/lint-docs.sh
```

### Pre-commit Hook Example

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running markdown linter..."
./scripts/lint-docs.sh

if [ $? -ne 0 ]; then
    echo "Markdown linting failed. Run './scripts/lint-docs.sh --fix' to auto-fix."
    exit 1
fi
```

---

## Best Practices

### 1. Run Before Committing

Always lint your markdown changes before committing:

```bash
# Check for issues
./scripts/lint-docs.sh

# Auto-fix what can be fixed
./scripts/lint-docs.sh --fix

# Review remaining issues
git diff
```

### 2. Code Block Languages

Always specify language for code blocks:

```markdown
✅ Good:
```bash
npm install
```

❌ Bad:
```
npm install
```
```

### 3. Line Length for Code vs Prose

- **Prose text:** Break at natural sentence boundaries
- **Code blocks:** Line length rule disabled in code
- **URLs:** Allowed to exceed limit (MD013 ignores them)
- **Tables:** Line length disabled for tables

### 4. Headings Structure

Maintain proper heading hierarchy:

```markdown
✅ Good:
# Title
## Section
### Subsection

❌ Bad:
# Title
### Subsection (skipped level 2)
```

---

## Script Details

### `scripts/lint-docs.sh`

**Features:**
- ✅ Colorized output (errors in red, success in green)
- ✅ Auto-fix mode with `--fix` flag
- ✅ Excludes `node_modules/`, `dist/`, `build/`, `.vscode-test/`
- ✅ File count statistics
- ✅ Exit codes for CI/CD integration
- ✅ Helpful error messages with documentation links

**Exit Codes:**
- `0` - All files compliant
- `1` - Linting errors found

---

## Troubleshooting

### markdownlint-cli not found

```bash
# Install globally
npm install -g markdownlint-cli

# Or use npx (no installation)
npx markdownlint-cli '**/*.md'
```

### Too many errors

Start with auto-fix to resolve most issues:

```bash
./scripts/lint-docs.sh --fix
```

Then review remaining issues manually.

### Rule is too strict

Update `.markdownlint.json` to disable or configure the rule:

```json
{
  "MD013": false,  // Disable line length check
  "MD033": {       // Allow specific HTML elements
    "allowed_elements": ["div", "br"]
  }
}
```

---

## Rule Documentation

Full rule reference: [markdownlint Rules](https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md)

Quick lookup:
- `MD0XX` - Rule number
- Each rule has detailed explanation
- Examples of violations and fixes
- Configuration options

---

## See Also

- [Link Checker](../../scripts/check-links.sh) - Validate internal/external links
- [CLAUDE.md](../../CLAUDE.md) - Project coding standards
- [Contributing Guide](../../CONTRIBUTING.md) - Contribution guidelines

---

**Last Updated:** October 6, 2025 | **Version:** v0.16.3
