# Gorev VS Code Extension Diagnostics

## Issue: Tasks Not Showing in TreeView

Based on the debugging, here are the findings and solutions:

### 1. **Root Cause**
The VS Code extension is not displaying tasks even though the MCP server returns them. The issue appears to be related to:
- The new compact format parser in v0.8.1 not assigning project IDs to tasks
- The executable path configured in VS Code settings may not exist

### 2. **Changes Made**
I've updated the following files to fix the issue:

#### a) Enhanced TreeView Provider (`src/providers/enhancedGorevTreeProvider.ts`)
- Added code to fetch the active project ID when loading tasks
- Tasks without project_id now get assigned the active project ID automatically
- Added detailed logging for debugging

#### b) Markdown Parser (`src/utils/markdownParser.ts`)
- Enhanced the compact format parser with better error handling
- Added support for tasks with empty descriptions
- Added detailed logging to help debug parsing issues

### 3. **Steps to Fix**

1. **Build the MCP Server**
   ```bash
   cd /mnt/f/Development/Projects/Gorev/gorev-mcpserver
   go build -o gorev cmd/gorev/main.go
   ```

2. **Verify the executable exists**
   ```bash
   ls -la /mnt/f/Development/Projects/Gorev/gorev-mcpserver/gorev
   ```

3. **In VS Code:**
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
   - Run "Developer: Reload Window"
   - Open the Output panel (`View > Output`)
   - Select "Gorev" from the dropdown
   - Try refreshing tasks (`Ctrl+Alt+R` or click refresh button)

4. **Check the logs**
   Look for these log entries in the Output panel:
   - `[EnhancedGorevTreeProvider] Active project ID: ...`
   - `[EnhancedGorevTreeProvider] Parsed tasks count: ...`
   - `[MarkdownParser] Total tasks parsed from compact format: ...`

### 4. **If Tasks Still Don't Show**

Run these commands in VS Code's Developer Console (Help > Toggle Developer Tools):

```javascript
// Check if tasks are being parsed
console.log('[Debug] Checking parsed tasks...');
const logs = document.querySelector('.monaco-workbench .output-editor').innerText;
console.log(logs.split('\n').filter(line => line.includes('Parsed tasks') || line.includes('MarkdownParser')));
```

### 5. **Additional Debugging**

If the issue persists, check:

1. **Server Path Configuration**
   - Open VS Code settings (`Ctrl+,`)
   - Search for "gorev.serverPath"
   - Ensure it points to the correct executable

2. **Active Project**
   - Make sure you have an active project set
   - Use the command "Gorev: Show Active Project" to verify

3. **Enable Debug Logging**
   - The extension now has debug logging enabled by default
   - Check the Output panel for detailed logs

### 6. **Expected Behavior**
After these fixes, you should see:
- Tasks appearing in the TreeView under their respective groups (by status)
- Progress bars for parent tasks with subtasks
- Priority badges and due date indicators
- Dependency counts with lock/unlock icons

### 7. **Notes on v0.8.1 Changes**
The MCP server now uses a compact format to prevent token limit errors:
- Status icons: [⏳] Beklemede, [🚀] DevamEdiyor, [✅] Tamamlandi
- Priority: Y (Yüksek), O (Orta), D (Düşük)
- Format: `[Icon] Title (Priority)` followed by details line

The VS Code extension has been updated to parse this new format correctly.