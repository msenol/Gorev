#!/bin/bash

# List of files that need t import added
files=(
    "src/providers/inlineEditProvider.ts"
    "src/providers/enhancedGorevTreeProvider.ts"
    "src/providers/projeTreeProvider.ts"
    "src/providers/dragDropController.ts"
    "src/ui/taskDetailPanel.ts"
    "src/ui/templateWizard.ts"
    "src/ui/exportDialog.ts"
    "src/ui/importWizard.ts"
    "src/commands/enhancedGorevCommands.ts"
    "src/commands/dataCommands.ts"
    "src/commands/projeCommands.ts"
    "src/commands/mcpDebugCommands.ts"
    "src/commands/debugCommands.ts"
    "src/commands/inlineEditCommands.ts"
    "src/commands/filterCommands.ts"
    "src/commands/templateCommands.ts"
    "src/commands/gorevCommands.ts"
    "src/debug/testDataSeeder.ts"
)

for file in "${files[@]}"; do
    if [[ -f "$file" ]]; then
        # Check if the file already has the import
        if ! grep -q "import { t } from '../utils/l10n'" "$file" && ! grep -q "import { t } from '../../utils/l10n'" "$file"; then
            # Determine the correct relative path
            if [[ "$file" == src/commands/* ]] || [[ "$file" == src/debug/* ]]; then
                import_statement="import { t } from '../utils/l10n';"
            elif [[ "$file" == src/providers/* ]] || [[ "$file" == src/ui/* ]]; then
                import_statement="import { t } from '../utils/l10n';"
            else
                import_statement="import { t } from '../utils/l10n';"
            fi

            # Find the last import statement and add our import after it
            sed -i "/^import.*from.*$/a\\$import_statement" "$file"
            echo "Added import to $file"
        else
            echo "Import already exists in $file"
        fi
    else
        echo "File not found: $file"
    fi
done

echo "Import fixing completed!"