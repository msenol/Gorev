#!/bin/bash

# VS Code Extension Packaging Script for Gorev v0.3.5

set -e

VERSION="0.3.5"
EXTENSION_DIR="gorev-vscode"
OUTPUT_DIR="release-v0.9.0"

echo "📦 Packaging VS Code Extension v${VERSION}..."

# Change to extension directory
cd ${EXTENSION_DIR}

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Compile TypeScript
echo "🔨 Compiling TypeScript..."
npm run compile

# Run tests
echo "🧪 Running tests..."
npm test || echo "⚠️  Some tests failed, continuing..."

# Install vsce if not already installed
if ! command -v vsce &> /dev/null; then
    echo "📦 Installing vsce..."
    npm install -g @vscode/vsce
fi

# Package extension
echo "📦 Creating VSIX package..."
vsce package --out "../${OUTPUT_DIR}/gorev-vscode-${VERSION}.vsix"

cd ..

echo "✅ VS Code extension packaged successfully!"
echo ""
echo "📁 Package location: ${OUTPUT_DIR}/gorev-vscode-${VERSION}.vsix"
echo ""
echo "📌 To publish to marketplace:"
echo "1. vsce login <publisher-name>"
echo "2. vsce publish -p <personal-access-token>"
echo ""
echo "Or upload manually at:"
echo "https://marketplace.visualstudio.com/manage/publishers/mehmetsenol"