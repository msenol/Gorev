#!/bin/bash

echo "🎨 Generating PNG icons from SVG..."

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "⚠️  ImageMagick is not installed. Installing..."
    sudo apt-get update && sudo apt-get install -y imagemagick
fi

# Generate main icon PNG from gorev-icon.svg
if [ -f "media/gorev-icon.svg" ]; then
    echo "Converting main icon..."
    convert -background transparent -resize 128x128 media/gorev-icon.svg media/icon.png
    echo "✅ Generated icon.png (128x128)"
else
    echo "❌ gorev-icon.svg not found!"
fi

# Generate smaller icons if needed
echo "🔍 Icons generated in media/ directory"
ls -la media/*.png 2>/dev/null || echo "No PNG files found yet."

echo "
📌 To manually create PNG icons:
1. Open the SVG files in a vector editor (Inkscape, Illustrator)
2. Export as PNG at desired resolutions
3. Save to media/ directory
"