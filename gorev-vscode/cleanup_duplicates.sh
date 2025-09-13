#!/bin/bash

# Find all TypeScript files and remove duplicate import statements
find src -name "*.ts" -exec awk '!seen[$0]++' {} > {}.tmp \; -exec mv {}.tmp {} \;

echo "Duplicate import cleanup completed!"