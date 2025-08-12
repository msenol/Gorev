# VS Code Extension Guide

Complete guide for installing, configuring, and using the Gorev VS Code Extension across different VS Code-compatible editors.

## 🎯 Overview

The Gorev VS Code Extension provides a rich visual interface for task management that integrates seamlessly with the Gorev MCP Server. It works with VS Code and other VS Code-compatible editors like Cursor, Windsurf, VSCodium, and more.

### ✨ Key Features

- **🌍 Bilingual Support** - Automatic English/Turkish detection
- **🌳 Unlimited Subtask Hierarchy** - Infinite task nesting
- **🔗 Advanced Dependencies** - Visual dependency management
- **🤖 AI Assistant Integration** - MCP protocol support
- **📊 Rich Visual Interface** - Enhanced TreeView with progress tracking
- **🎛️ Advanced Task Management** - Bulk operations, drag & drop, inline editing

## 📦 Installation

### VS Code (Recommended)

**Method 1: VS Code Marketplace**
1. Open VS Code
2. Press `Ctrl+Shift+X` (or `Cmd+Shift+X` on macOS) to open Extensions
3. Search for "Gorev"
4. Click "Install" on the Gorev extension by mehmetsenol

**Method 2: Command Line**
```bash
code --install-extension mehmetsenol.gorev-vscode
```

**Method 3: Direct Download**
- Visit [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
- Click "Download Extension" 
- Install the downloaded `.vsix` file via VS Code

### Alternative Editors

For VS Code-compatible editors, you'll need to download and install the VSIX file manually.

#### Download VSIX File

**Option 1: GitHub Releases**
1. Go to [Gorev Releases](https://github.com/msenol/gorev/releases)
2. Find the latest release
3. Download `gorev-vscode-x.x.x.vsix` from Assets

**Option 2: Direct Download**
```bash
# Terminal/PowerShell
curl -L -o gorev-vscode.vsix https://marketplace.visualstudio.com/_apis/public/gallery/publishers/mehmetsenol/vsextensions/gorev-vscode/latest/vspackage
```

#### Install in Different Editors

**Cursor**
1. Open Command Palette: `Cmd/Ctrl + Shift + P`
2. Type "Extensions: Install from VSIX..."
3. Select the downloaded `.vsix` file
4. Restart Cursor

**Windsurf**
1. Command Palette: `Cmd/Ctrl + Shift + P`
2. Search for "Install from VSIX"
3. Select the VSIX file
4. Restart the IDE

**VSCodium**
```bash
# Command line
codium --install-extension gorev-vscode.vsix

# Or via GUI: Extensions panel → ⋯ (three dots) → Install from VSIX...
```

**Code-Server (Browser-based VS Code)**
```bash
# On server
code-server --install-extension gorev-vscode.vsix

# Or upload via web interface:
# Settings → Extensions → Install from VSIX → Upload
```

**Theia**
```bash
# Terminal
theia extension:install gorev-vscode.vsix
```

#### Manual Installation

If automatic installation fails:

1. **Locate extension directory:**
   - **Windows**: `%USERPROFILE%\.vscode\extensions\` or `%APPDATA%\[EDITOR-NAME]\extensions\`
   - **macOS**: `~/.vscode/extensions/` or `~/Library/Application Support/[EDITOR-NAME]/extensions/`
   - **Linux**: `~/.vscode/extensions/` or `~/.config/[EDITOR-NAME]/extensions/`

2. **Extract VSIX (it's a zip file):**
   ```bash
   unzip gorev-vscode.vsix -d gorev-temp
   ```

3. **Copy to extensions directory:**
   ```bash
   # Linux/macOS
   cp -r gorev-temp/extension ~/.vscode/extensions/mehmetsenol.gorev-vscode-0.5.0

   # Windows PowerShell
   Copy-Item -Recurse gorev-temp\extension "$env:USERPROFILE\.vscode\extensions\mehmetsenol.gorev-vscode-0.5.0"
   ```

4. **Restart your editor**

## ⚙️ Configuration

### Basic Setup

After installation, configure the extension to connect to your Gorev MCP Server:

1. **Open Settings**: `Ctrl/Cmd + ,`
2. **Search for "gorev"**
3. **Configure the server path**:

**Settings.json configuration:**
```json
{
  "gorev.serverPath": "gorev",           // If gorev is in PATH
  "gorev.autoConnect": true,             // Auto-connect on startup
  "gorev.showStatusBar": true,           // Show status in status bar
  "gorev.refreshInterval": 30            // Refresh interval in seconds
}
```

**Platform-specific paths:**
- **Linux/macOS**: `/usr/local/bin/gorev` or `gorev` (if in PATH)
- **Windows**: `C:\Users\[USERNAME]\AppData\Local\Programs\gorev\gorev.exe` or `gorev` (if in PATH)

### Advanced Configuration

**TreeView Preferences:**
```json
{
  "gorev.treeView.grouping": "status",                    // Group by: status, priority, project, tag, due_date
  "gorev.treeView.sorting": "priority",                   // Sort by: title, priority, due_date, created_at
  "gorev.treeView.visuals.showProgressBars": true,       // Show progress bars for parent tasks
  "gorev.treeView.visuals.showPriorityBadges": true,     // Show priority badges
  "gorev.treeView.visuals.colorCoding": "priority",      // Color code by priority
  "gorev.dragDrop.allowTaskMove": true,                  // Enable drag & drop
  "gorev.dragDrop.allowStatusChange": true              // Allow status change via drag & drop
}
```

**Filter Settings:**
```json
{
  "gorev.filters.autoSave": true,                       // Auto-save filter profiles
  "gorev.filters.defaultProfile": "active_tasks",       // Default filter profile
  "gorev.search.caseSensitive": false,                  // Case sensitive search
  "gorev.search.regex": false                           // Enable regex in search
}
```

## 🎮 Usage Guide

### Getting Started

1. **First Launch**
   - Click the Gorev icon in the Activity Bar (left sidebar)
   - The extension will attempt to connect to the MCP server
   - If no tasks exist, it will offer to create sample data

2. **Create Your First Project**
   - Right-click in the TreeView
   - Select "Create Project"
   - Enter project name and description

3. **Add Tasks**
   - Right-click on a project or in empty space
   - Select "Create Task" 
   - Fill in task details using the template system

### Core Features

#### 🌳 Unlimited Subtask Hierarchy

Create nested tasks with infinite levels:

```
📁 Web Development Project
├── 🎯 User Authentication (High Priority)
│   ├── ⚙️ Set up JWT middleware
│   ├── 🖼️ Create login form
│   │   ├── 📱 Mobile responsive design
│   │   └── ✅ Form validation
│   └── 🔒 Add password validation
├── 🎨 Frontend Development
└── 🔧 Backend API
```

**Controls:**
- **Expand/Collapse**: Click arrow icons or use keyboard shortcuts
- **Create Subtask**: Right-click on any task
- **Progress Tracking**: Parent tasks show completion percentage
- **Visual Nesting**: Indentation and lines show hierarchy

#### 🔗 Advanced Dependencies

Manage task dependencies with visual indicators:

**Creating Dependencies:**
1. **Drag & Drop**: Drag task A onto task B to make A depend on B
2. **Context Menu**: Right-click → "Add Dependency"
3. **Detail Panel**: Use the dependency section

**Dependency States:**
- 🔒 **Blocked**: Cannot start (dependencies not completed)
- 🔓 **Ready**: All dependencies completed
- 🔗 **Linked**: Has dependencies (some may be pending)

**Smart Validation:**
- Prevents circular dependencies
- Shows dependency chains
- Validates before allowing changes

#### 📊 Rich Visual Interface

**Enhanced TreeView:**
- **Grouping Options**: Status, Priority, Project, Tag, Due Date
- **Color Coding**: Priority-based visual distinction
- **Progress Bars**: Parent task completion tracking
- **Smart Badges**: Task counts, due dates, dependency indicators

**TreeView Controls:**
```
🎯 High Priority Tasks (3)          ← Group header with count
├── [🔴] Fix login bug               ← Priority color + checkbox
│   📅 Due: Today                    ← Due date badge
├── [⚪] Implement dark mode          ← Status indicator
│   🔗 Depends on: Theme system      ← Dependency indicator
└── [🟡] Update documentation        ← Progress indicator
    ▓▓▓░░░ 60%                      ← Progress bar
```

#### 🎛️ Advanced Task Management

**Multi-Select Operations:**
1. Hold `Ctrl/Cmd` and click multiple tasks
2. Right-click for bulk operations:
   - Change status
   - Update priority
   - Add tags
   - Delete tasks

**Drag & Drop Features:**
- **Move Tasks**: Drag between projects
- **Change Status**: Drop on status groups
- **Reorder Priority**: Drop in priority order
- **Create Dependencies**: Drop task on another

**Inline Editing:**
- **F2** or **double-click** to edit task title
- **Context menu** for quick status/priority changes
- **Escape** to cancel, **Enter** to save
- **Tab navigation** between fields

#### 🔍 Advanced Filtering

**Real-time Search:**
- Search across title, description, tags, and projects
- Highlighting of matching text
- Live filtering as you type

**Filter Profiles:**
```json
{
  "active_tasks": {
    "status": ["pending", "in_progress"],
    "priority": ["high", "medium"]
  },
  "overdue": {
    "due_date": "past",
    "status": ["pending", "in_progress"]
  },
  "my_bugs": {
    "tags": ["bug"],
    "assigned_to": "me"
  }
}
```

**Quick Filter Shortcuts:**
- `Ctrl+1`: Show all tasks
- `Ctrl+2`: Show active tasks only
- `Ctrl+3`: Show completed tasks
- `Ctrl+4`: Show high priority tasks

### 🤖 AI Assistant Integration

The extension works seamlessly with MCP-compatible AI assistants:

**Supported AI Tools:**
- Claude Desktop
- VS Code Copilot
- Windsurf AI
- Cursor AI
- Any MCP-compatible assistant

**Natural Language Commands:**
```
"Create a new bug report for the login issue"
→ Opens template wizard with bug report template

"Show me all overdue high priority tasks"
→ Applies filters and displays matching tasks

"Mark the authentication task as completed"  
→ Updates task status and refreshes TreeView

"Create subtasks for the API development task"
→ Creates nested tasks under the specified parent
```

**Template Integration:**
AI assistants can use predefined templates:
- **Bug Report**: Structured issue documentation
- **Feature Request**: New feature specifications
- **Technical Debt**: Code improvement tasks  
- **Research**: Investigation and learning tasks

### 📋 Template System

Templates ensure consistent task creation:

**Available Templates:**
1. **Bug Report**
   - Title pattern: `🐛 [Module] Issue Description`
   - Required fields: module, description, steps, expected, actual
   - Priority default: high

2. **Feature Request**
   - Title pattern: `✨ [Component] Feature Name`
   - Fields: description, requirements, acceptance_criteria
   - Priority default: medium

3. **Technical Debt**
   - Title pattern: `🔧 [Area] Improvement Description`
   - Fields: current_state, desired_state, impact
   - Priority default: low

**Using Templates:**
1. **Via Extension**: Right-click → "Create Task" → Select template
2. **Via AI**: "Create a bug report task" → AI uses appropriate template
3. **Via MCP Tools**: Use `templateden_gorev_olustur` command

## 🔧 Troubleshooting

### Common Issues

#### Extension Not Loading
**Symptoms**: Gorev icon doesn't appear, no TreeView
**Solutions**:
1. Check VS Code version compatibility (requires VS Code 1.74.0+)
2. Restart VS Code completely
3. Check extension installation: Extensions → Installed → Gorev
4. Check for extension errors in Developer Tools (`Ctrl+Shift+I`)

#### MCP Server Connection Failed
**Symptoms**: "Server not running" message, empty TreeView
**Solutions**:
1. **Verify server installation**:
   ```bash
   gorev version  # Should show version info
   ```

2. **Check server path in settings**:
   - Open Settings (`Ctrl+,`)
   - Search "gorev.serverPath"
   - Use full path if `gorev` command not in PATH

3. **Start server manually**:
   ```bash
   gorev serve --debug
   # Leave running, then restart VS Code
   ```

4. **Check Output panel**:
   - View → Output
   - Select "Gorev" from dropdown
   - Look for connection errors

#### Tasks Not Loading
**Symptoms**: Extension loads but no tasks appear
**Solutions**:
1. **Refresh TreeView**: Click refresh button or `Ctrl+R`
2. **Check data directory**: Ensure `~/.gorev/data/` exists
3. **Reset database** (if needed):
   ```bash
   mv ~/.gorev/data/gorev.db ~/.gorev/data/gorev.db.backup
   gorev serve  # Will create new database
   ```

#### Performance Issues
**Symptoms**: Slow loading, unresponsive interface
**Solutions**:
1. **Reduce refresh interval**:
   ```json
   {
     "gorev.refreshInterval": 60  // Increase from 30 to 60 seconds
   }
   ```

2. **Disable visual effects**:
   ```json
   {
     "gorev.treeView.visuals.showProgressBars": false,
     "gorev.treeView.visuals.animations": false
   }
   ```

3. **Archive old tasks**: Clean up completed tasks older than 3 months

### Debug Mode

Enable debug logging for troubleshooting:

1. **Enable in settings**:
   ```json
   {
     "gorev.debug": true,
     "gorev.verboseLogging": true
   }
   ```

2. **Check debug output**:
   - View → Output → Select "Gorev Debug"
   - Look for detailed connection and operation logs

3. **Extension Development Host** (for advanced debugging):
   - Open Gorev extension source code
   - Press `F5` to launch Extension Development Host
   - Test functionality in clean environment

### Platform-Specific Issues

#### Windows
- **Path issues**: Use forward slashes or escape backslashes in settings
- **PowerShell execution policy**: May need to enable script execution
- **Windows Defender**: May block gorev.exe, add exception if needed

#### macOS  
- **Gatekeeper**: May need to allow gorev binary in System Preferences → Security
- **Path issues**: Gorev installed via Homebrew goes to `/opt/homebrew/bin/`

#### Linux
- **Permissions**: Ensure gorev binary has execute permissions (`chmod +x gorev`)
- **PATH issues**: Check if `/usr/local/bin` is in PATH
- **AppImage**: If using AppImage version, use full path to AppImage file

## 📊 Supported Editors

### ✅ Fully Compatible
- **VS Code** - Complete feature support
- **Cursor** - All features work
- **Windsurf** - Full compatibility
- **VSCodium** - Open-source VS Code, fully supported
- **Code-Server** - Browser-based VS Code

### ⚠️ Partially Compatible  
- **Theia** - Core features work, some UI elements may differ
- **Eclipse Che** - Basic functionality, limited visual features
- **Gitpod** - Works in web environment with limitations

### ❌ Not Compatible
- **Sublime Text** - No VS Code API
- **Atom** - Different extension system (discontinued)
- **IntelliJ IDEA** - Different platform entirely
- **Vim/Neovim** - Text-based, no TreeView support

## 🔗 Integration Examples

### Development Workflow

**Daily Development Routine:**
1. **Morning**: Review active tasks in TreeView
2. **Work**: Update task status as you progress
3. **AI Assistance**: "Create task for fixing the authentication bug I found"
4. **Evening**: Mark completed tasks, plan next day

**Sprint Planning:**
1. Create sprint project
2. Import user stories as tasks
3. Break down into subtasks
4. Set dependencies between tasks
5. Track progress via TreeView grouping

### Team Collaboration

**Shared Project Setup:**
1. Each team member installs extension
2. Connect to shared Gorev MCP server
3. Use consistent templates for task creation
4. Regular status updates via AI commands

**Code Review Integration:**
```
"Create task for code review of PR #123"
→ Creates structured review task with checklist
→ Link to PR, assign reviewer, set due date
→ Track review progress in TreeView
```

## 📚 Additional Resources

- **[Main Documentation](../../../README.en.md)** - Gorev project overview
- **[Installation Guide](../getting-started/installation.md)** - Server installation
- **[MCP Tools Reference](mcp-tools.md)** - Complete command reference
- **[API Reference](../api/reference.md)** - Technical documentation
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions

## 🤝 Support

- **GitHub Issues**: [Report bugs and feature requests](https://github.com/msenol/Gorev/issues)
- **Discussions**: [Community support and questions](https://github.com/msenol/Gorev/discussions)
- **VS Code Marketplace**: [Extension reviews and ratings](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

---

*✨ This comprehensive VS Code Extension guide helps you maximize your productivity with Gorev's visual task management interface across all supported editors and platforms.*