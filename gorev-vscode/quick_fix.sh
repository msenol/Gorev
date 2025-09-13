#!/bin/bash

# Remove all duplicate imports of t
for file in src/commands/*.ts src/providers/*.ts src/ui/*.ts src/debug/*.ts; do
    if [[ -f "$file" ]]; then
        # Keep only the first occurrence of t import
        awk '!seen && /import { t } from/ {seen=1; print; next} !/import { t } from/ {print}' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
        echo "Fixed $file"
    fi
done

echo "All duplicate imports fixed!"