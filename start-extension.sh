#!/bin/bash

echo "🚀 Starting Gorev VS Code Extension Development..."

# Check if we're in the right directory
if [ ! -d "gorev-vscode" ]; then
    echo "❌ Error: gorev-vscode directory not found!"
    echo "Please run this script from the Gorev project root directory."
    exit 1
fi

# Check if gorev server exists
if [ ! -f "gorev-mcpserver/gorev" ]; then
    echo "⚠️  Warning: Gorev server not found. Building..."
    cd gorev-mcpserver
    make build
    cd ..
fi

# Make sure server is executable
chmod +x gorev-mcpserver/gorev

# Navigate to extension directory
cd gorev-vscode

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Compile TypeScript
echo "🔨 Compiling TypeScript..."
npm run compile

echo "✅ Ready! Now:"
echo "1. Open VS Code in the gorev-vscode directory"
echo "2. Press F5 to launch the Extension Development Host"
echo "3. The extension will automatically connect to the Gorev server"
echo ""
echo "Server path configured as: /mnt/f/Development/Projects/Gorev/gorev-mcpserver/gorev"