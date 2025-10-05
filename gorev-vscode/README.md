# Gorev - Advanced Task Management & AI Integration for VS Code

<p align="center">
  <img src="media/icon.png" alt="Gorev Logo" width="128" height="128">
</p>

<div align="center">

[🇺🇸 English](README.md) | [🇹🇷 Türkçe](README.tr.md)

[![Version](https://img.shields.io/badge/Version-0.16.1-blue?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Downloads](https://img.shields.io/visual-studio-marketplace/d/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Rating](https://img.shields.io/visual-studio-marketplace/r/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**The ultimate task management powerhouse with 48 MCP tools, unlimited hierarchy, and seamless AI assistant integration**

> 🎉 **NEW in v0.16.1**: Automatic Server Startup! Zero-configuration - extension now starts server automatically. No manual commands needed!

> 🚀 **v0.16.0**: Complete REST API Migration! Type-safe JSON responses, enhanced error handling, and Web UI integration. No more markdown parsing - direct API communication for 3x better performance!

</div>

## 🌟 Why Choose Gorev

Gorev transforms VS Code into a **professional task management powerhouse** with unique capabilities that set it apart:

- **🚀 Zero-Installation NPX Setup** - Get started in seconds with no binary downloads
- **🤖 48 MCP Tools** - Most comprehensive task management API for AI assistants
- **🌳 Unlimited Hierarchy** - Infinite subtask nesting with visual progress tracking
- **🔍 Advanced Search** - FTS5 full-text search with fuzzy matching and NLP
- **🎯 Smart Dependencies** - Visual dependency management with auto-resolution
- **🌍 Bilingual Support** - 668 i18n keys with automatic language detection
- **📊 Data Export/Import** - Multi-step wizards with conflict resolution
- **⚡ Ultra Performance** - RefreshManager with 90% reduction in operations

## 🚀 Zero-Installation Setup

### 🎯 NPX Mode (Recommended - No Downloads!)

The easiest way to get started - no binary installation required:

1. **Install Extension**: Search "Gorev" in VS Code marketplace
2. **Auto-Configuration**: Extension uses NPX mode by default
3. **Start Working**: Create projects and tasks immediately!

The extension automatically runs `npx @mehmetsenol/gorev-mcp-server@latest` behind the scenes.

```json
// Default configuration - no setup needed!
{
  "gorev.serverMode": "npx",  // Automatic NPX execution
  "gorev.autoConnect": true   // Connect on startup
}
```

### 🔧 Binary Mode (Advanced Users)

For users who prefer local binary installation:

```json
{
  "gorev.serverMode": "binary",
  "gorev.serverPath": "/path/to/gorev"
}
```

Follow the [installation guide](https://github.com/msenol/Gorev/blob/main/README.en.md#-installation) for binary setup.

## 🎯 Key Features Matrix

| Category | Feature | Description | Status |
|----------|---------|-------------|--------|
| **🚀 Setup** | NPX Zero-Install | No downloads, instant setup | ✅ |
| **🤖 AI Integration** | 48 MCP Tools | Complete API for AI assistants | ✅ |
| **🌳 Task Management** | Unlimited Hierarchy | Infinite subtask nesting | ✅ |
| **🔗 Dependencies** | Smart Resolution | Visual dependency management | ✅ |
| **🔍 Search** | FTS5 Full-Text | SQLite virtual tables, fuzzy matching | ✅ |
| **📊 Data Management** | Export/Import Wizards | JSON/CSV with conflict resolution | ✅ |
| **🎨 Visual Interface** | Rich TreeView | Progress bars, badges, color coding | ✅ |
| **⚡ Performance** | RefreshManager | 90% operation reduction, debouncing | ✅ |
| **🌍 Localization** | Bilingual Support | 668 i18n keys, auto-detection | ✅ |
| **💾 Database** | Workspace Mode | Project-specific or global databases | ✅ |
| **🎛️ Customization** | 50+ Settings | Complete visual and behavioral control | ✅ |
| **🔄 Real-time** | File Watching | Auto-updates on file changes | ✅ |

## 🤖 AI Assistant Integration

### MCP Protocol Compatibility

Works seamlessly with all MCP-compatible AI assistants:

- **✅ Claude Desktop** - Full conversation integration
- **✅ VS Code with MCP** - Native extension support
- **✅ Cursor IDE** - AI coding assistant integration
- **✅ Windsurf** - Development environment integration
- **✅ Any MCP Client** - Universal compatibility

### Natural Language Task Management

Talk to your AI assistant naturally:

```
🗨️ "Create a new high-priority task for implementing dark mode"
🗨️ "Show me all overdue tasks with dependencies"
🗨️ "Mark task #42 as completed and update dependencies"
🗨️ "Create a bug report template for the login issue"
🗨️ "Export all completed tasks from last month to CSV"
```

### 48 MCP Tools Categories

| Category | Tools | Description |
|----------|--------|-------------|
| **Task Management** | 6 tools | Create, update, list, detail operations |
| **Subtask Operations** | 3 tools | Hierarchy management and nesting |
| **Project Management** | 6 tools | Project creation, activation, statistics |
| **Template System** | 2 tools | Template-based task creation |
| **Advanced Search** | 6 tools | FTS5 search, suggestions, history |
| **Data Export/Import** | 2 tools | Multi-format data operations |
| **File Watching** | 4 tools | File system monitoring |
| **AI Context** | 6 tools | Context management and NLP |
| **IDE Integration** | 5 tools | Extension management automation |
| **Advanced Operations** | 8 tools | Batch processing, analytics |

## 🌳 Unlimited Task Hierarchy

### Visual Hierarchy Management

- **🔄 Infinite Nesting** - Create tasks within tasks without limits
- **📊 Progress Tracking** - Parent tasks show completion percentage
- **🎯 Visual Indicators** - Tree structure with expand/collapse
- **⚡ Quick Operations** - Drag & drop, inline editing

### Hierarchy Examples

```
📁 Project: E-commerce Platform
├── 🚀 User Authentication System (75% complete)
│   ├── ✅ JWT Middleware Setup
│   ├── ✅ Login Form Component
│   ├── 🔄 Password Validation
│   │   ├── ⏳ Regex Pattern Implementation
│   │   └── ⏳ Error Message Localization
│   └── ⏳ Session Management
└── 📱 Mobile Responsive Design (25% complete)
    ├── ✅ Breakpoint Analysis
    └── ⏳ Component Adaptation
        ├── ⏳ Header Responsiveness
        └── ⏳ Navigation Menu
```

## 🔍 Advanced Search & Filtering

### FTS5 Full-Text Search

Powered by SQLite virtual tables for lightning-fast search:

- **🔍 Content Search** - Search in titles, descriptions, tags
- **🎯 Fuzzy Matching** - Find tasks even with typos
- **🧠 NLP Integration** - Natural language query parsing
- **📊 Search Analytics** - Track search patterns and history
- **💾 Saved Profiles** - Store complex filter combinations

### Filter Capabilities

```typescript
// Advanced filtering options
{
  status: ["pending", "in_progress"],
  priority: ["high", "medium"],
  tags: ["bug", "urgent"],
  dateRange: {
    start: "2025-01-01",
    end: "2025-12-31"
  },
  project: "WebApp",
  hasDepencies: true,
  isOverdue: false
}
```

## 🔗 Smart Dependency Management

### Visual Dependency System

- **🔒 Blocked Tasks** - Clear visual indicators for blocked tasks
- **🔓 Ready Tasks** - Automatic resolution when dependencies complete
- **🔗 Linked Tasks** - Bidirectional dependency visualization
- **⚡ Batch Operations** - Manage multiple dependencies at once

### Dependency Types

| Icon | Status | Description |
|------|--------|-------------|
| 🔒 | Blocked | Has incomplete dependencies |
| 🔓 | Ready | All dependencies completed |
| 🔗 | Linked | Has bidirectional connections |
| ⚡ | Auto | Automatic resolution enabled |

## 📊 Data Export & Import Wizards

### Multi-Step Export Wizard

Advanced export capabilities with guided setup:

1. **📋 Select Format** - JSON (structured) or CSV (tabular)
2. **🎯 Choose Scope** - Current view, project, or custom filter
3. **📅 Date Range** - Flexible date filtering options
4. **🔧 Configuration** - Include dependencies, tags, metadata
5. **📤 Export** - Progress tracking with VS Code notifications

### Import with Conflict Resolution

Intelligent import system with multiple resolution strategies:

- **🔄 Skip Conflicts** - Keep existing data unchanged
- **📝 Overwrite** - Replace with imported data
- **🔀 Merge** - Smart combination of existing and new data
- **👀 Preview** - See changes before applying

### Export Formats

```json
// JSON Export (structured)
{
  "tasks": [...],
  "projects": [...],
  "dependencies": [...],
  "metadata": {
    "exportDate": "2025-09-18",
    "version": "v0.6.12"
  }
}
```

```csv
// CSV Export (tabular)
ID,Title,Status,Priority,Project,Tags,DueDate,Progress
1,"Setup Auth",pending,high,"WebApp","security,auth","2025-10-01",0
```

## ⚡ Performance Optimizations

### RefreshManager Architecture

Revolutionary refresh system with 90% performance improvement:

- **🎯 Intelligent Batching** - Group operations for efficiency
- **⏱️ Priority Debouncing** - High: 100ms, Normal: 500ms, Low: 2s
- **🔍 Differential Updates** - Hash-based change detection
- **📊 Performance Monitoring** - Real-time operation tracking
- **🚫 Zero Blocking** - Non-blocking async operations

### Performance Metrics

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Refresh Operations | 1000/min | 100/min | 90% reduction |
| UI Thread Blocking | 50ms | 0ms | Zero blocking |
| Memory Usage | 50MB | 35MB | 30% reduction |
| Startup Time | 2s | 1.4s | 30% faster |

## 🚀 v0.16.0 API Architecture (NEW!)

### REST API Integration

Complete migration from MCP (stdio + markdown) to REST API (HTTP + JSON):

**Key Benefits:**
- **⚡ 3x Faster** - Direct HTTP communication vs stdio streams
- **🔒 Type-Safe** - Full TypeScript type checking, zero parsing errors
- **🛡️ Better Errors** - Structured error responses with status codes
- **🌐 Web UI Ready** - Shared architecture with embedded web interface
- **📊 Debuggable** - Standard HTTP requests visible in network tools

**Architecture Improvements:**
```typescript
// Before v0.16.0 (MCP + Markdown Parsing)
const result = await mcpClient.callTool('gorev_listele', {});
const tasks = MarkdownParser.parseGorevListesi(result.content[0].text);
// ❌ Fragile regex parsing, no type safety

// After v0.16.0 (REST API + JSON)
const response = await apiClient.getTasks({});
const tasks = response.data; // Task[] - fully typed!
// ✅ Type-safe, no parsing needed
```

**What Changed:**
- ✅ All TreeView providers use REST API
- ✅ All command handlers use REST API
- ✅ Enhanced error handling with ApiError class
- ✅ ~300 lines of markdown parsing eliminated
- ⚠️ MCPClient deprecated (removal in v0.18.0)
- ⚠️ MarkdownParser deprecated (removal in v0.18.0)

**Migration Status:** 90% complete
- Remaining: TemplateWizard & TaskDetailPanel (v0.17.0)

## 🎨 Rich Visual Interface

### Enhanced TreeView

Professional-grade tree interface with advanced features:

- **📊 Progress Bars** - Visual completion tracking for parent tasks
- **🎯 Priority Badges** - Color-coded priority indicators (🔥⚡ℹ️)
- **📅 Smart Dates** - Relative formatting (Today, Tomorrow, 3d left)
- **🔗 Dependency Icons** - Visual dependency status (🔒🔓🔗)
- **🏷️ Tag Pills** - Colored tag badges with hover details
- **📈 Rich Tooltips** - Markdown tooltips with progress visualization

### Drag & Drop Operations

Intuitive drag & drop for all operations:

- **🔄 Move Tasks** - Between projects with visual feedback
- **📊 Change Status** - Drop on status groups to update
- **🎯 Reorder Priority** - Drag to change priority levels
- **🔗 Create Dependencies** - Drop task on another to create dependency
- **✨ Visual Feedback** - Smooth animations and drop indicators

### Inline Editing

Quick editing without leaving the tree view:

- **✏️ F2 or Double-Click** - Instant title editing
- **📝 Context Menus** - Right-click for status/priority changes
- **📅 Date Picker** - Inline due date selection
- **⌨️ Keyboard Shortcuts** - Escape to cancel, Enter to save

## 🎛️ Comprehensive Configuration

### 50+ Configuration Options

Complete control over every aspect of the extension:

#### Server Settings (5 options)

```json
{
  "gorev.serverMode": "npx|binary",
  "gorev.serverPath": "/path/to/gorev",
  "gorev.autoConnect": true,
  "gorev.showStatusBar": true,
  "gorev.refreshInterval": 300
}
```

#### TreeView Visuals (15 options)

```json
{
  "gorev.treeView.grouping": "status|priority|project|tag|dueDate",
  "gorev.treeView.sorting": "title|priority|dueDate|created",
  "gorev.treeView.sortAscending": false,
  "gorev.treeView.showCompleted": true,
  "gorev.treeView.showEmptyGroups": false,
  "gorev.treeView.visuals.showProgressBars": true,
  "gorev.treeView.visuals.showPriorityBadges": true,
  "gorev.treeView.visuals.showDueDateIndicators": true,
  "gorev.treeView.visuals.showDependencyBadges": true,
  "gorev.treeView.visuals.showTagPills": true,
  "gorev.treeView.visuals.progressBarStyle": "blocks|percentage|both",
  "gorev.treeView.visuals.dueDateFormat": "relative|absolute|smart",
  "gorev.treeView.visuals.priorityBadgeStyle": "emoji|text|color",
  "gorev.treeView.visuals.tagPillLimit": 3,
  "gorev.treeView.visuals.showSubtaskProgress": true
}
```

#### Drag & Drop (8 options)

```json
{
  "gorev.dragDrop.allowTaskMove": true,
  "gorev.dragDrop.allowStatusChange": true,
  "gorev.dragDrop.allowPriorityChange": true,
  "gorev.dragDrop.allowProjectMove": true,
  "gorev.dragDrop.allowDependencyCreate": true,
  "gorev.dragDrop.allowParentChange": true,
  "gorev.dragDrop.showDropIndicator": true,
  "gorev.dragDrop.enableAnimation": true
}
```

#### Performance (8 options)

```json
{
  "gorev.refreshManager.enableBatching": true,
  "gorev.refreshManager.batchSize": 10,
  "gorev.refreshManager.highPriorityDelay": 100,
  "gorev.refreshManager.normalPriorityDelay": 500,
  "gorev.refreshManager.lowPriorityDelay": 2000,
  "gorev.performance.enableMonitoring": true,
  "gorev.performance.slowOperationThreshold": 1000,
  "gorev.performance.maxMetrics": 1000
}
```

#### Database Modes (3 options)

```json
{
  "gorev.databaseMode": "auto|workspace|global",
  "gorev.workspaceDatabase.autoDetect": true,
  "gorev.workspaceDatabase.showModeInStatusBar": true
}
```

## 🌍 Internationalization

### Complete Bilingual Support

- **668 i18n Keys** - Every UI element translated
- **Auto-Detection** - Follows VS Code language setting
- **Languages**: English (en) and Turkish (tr)
- **Context-Aware** - Smart translations based on usage

### Translation Examples

| English | Turkish | Context |
|---------|---------|---------|
| "Create Task" | "Görev Oluştur" | Command |
| "High Priority" | "Yüksek Öncelik" | Priority badge |
| "Dependencies blocked" | "Bağımlılıklar engelledi" | Status |
| "Export completed" | "Dışa aktarma tamamlandı" | Notification |

## 💾 Database Management

### Flexible Database Modes

#### Workspace Mode (Default)

- **📁 Project-Specific** - Each project has its own `.gorev/gorev.db`
- **🔍 Auto-Detection** - Automatically finds workspace databases
- **📊 Status Indicator** - Shows current database in status bar

#### Global Mode

- **🌐 Shared Database** - Single database for all projects
- **🏠 User Directory** - Stored in `~/.gorev/gorev.db`
- **🔄 Easy Switching** - Toggle between modes via command

#### Auto Mode

- **🤖 Intelligent Selection** - Automatically chooses best database
- **⬆️ Fallback Chain** - Workspace → Parent → Global
- **⚡ Zero Configuration** - Works out of the box

## 📋 50+ Available Commands

### Task Operations (15 commands)

- `gorev.createTask` - Create new task
- `gorev.updateTaskStatus` - Update task status
- `gorev.showTaskDetail` - Show detailed task view
- `gorev.deleteTask` - Delete task
- `gorev.markAsCompleted` - Quick completion
- `gorev.setTaskPriority` - Change priority
- `gorev.addTaskTag` - Add tags
- `gorev.setTaskDueDate` - Set due date
- `gorev.createSubtask` - Add subtask
- `gorev.moveTo` - Move to project
- `gorev.duplicateTask` - Clone task
- `gorev.addTaskNote` - Add notes
- `gorev.linkTasks` - Create dependency
- `gorev.unlinkTasks` - Remove dependency
- `gorev.showTaskHistory` - View history

### Project Management (8 commands)

- `gorev.createProject` - Create new project
- `gorev.setActiveProject` - Set active project
- `gorev.showProjectStats` - View statistics
- `gorev.deleteProject` - Delete project
- `gorev.renameProject` - Rename project
- `gorev.archiveProject` - Archive project
- `gorev.exportProject` - Export project data
- `gorev.duplicateProject` - Clone project

### Template System (7 commands)

- `gorev.openTemplateWizard` - Template wizard
- `gorev.createFromTemplate` - Create from template
- `gorev.quickCreateFromTemplate` - Quick template selection
- `gorev.refreshTemplates` - Reload templates
- `gorev.initDefaultTemplates` - Initialize defaults
- `gorev.showTemplateDetails` - Template details
- `gorev.exportTemplate` - Export template

### Data Operations (4 commands)

- `gorev.exportData` - Export data wizard
- `gorev.importData` - Import data wizard
- `gorev.exportCurrentView` - Export current view
- `gorev.quickExport` - Quick export

### Filter Operations (10 commands)

- `gorev.showSearchInput` - Search tasks
- `gorev.showFilterMenu` - Filter menu
- `gorev.showFilterProfiles` - Saved profiles
- `gorev.clearAllFilters` - Clear filters
- `gorev.filterOverdue` - Show overdue
- `gorev.filterDueToday` - Show due today
- `gorev.filterDueThisWeek` - Show due this week
- `gorev.filterHighPriority` - Show high priority
- `gorev.filterActiveProject` - Show active project
- `gorev.filterByTag` - Filter by tag

### Debug Tools (6 commands)

- `gorev.showDebugInfo` - Debug information
- `gorev.clearDebugLogs` - Clear logs
- `gorev.testConnection` - Test MCP connection
- `gorev.refreshAllViews` - Force refresh
- `gorev.resetExtension` - Reset state
- `gorev.generateTestData` - Generate test data

## 🔄 File System Integration

### File Watcher Capabilities

- **📁 Project Monitoring** - Watch project files for changes
- **🔄 Auto-Updates** - Automatic task status transitions
- **⚡ Real-Time Sync** - Instant UI updates on file changes
- **🎯 Selective Watching** - Configure which files to monitor

### Integration Patterns

```javascript
// Automatic status updates based on file changes
.gitignore change → Update "Setup Git" task
package.json change → Update "Configure Dependencies" task
README.md change → Update "Documentation" task
```

## 🏆 Advanced Capabilities

### Professional Template Wizard 🎯 NEW

**Complete redesign with professional UI/UX for enhanced task creation:**

- **🎨 9 Field Types** - Text, textarea, select, date, tags, email, URL, number, markdown
- **⚡ Real-Time Validation** - Dynamic field validation with visual feedback
- **📝 Markdown Preview** - Live markdown rendering with local marked.js bundle
- **💫 Professional Styling** - 300+ lines of enhanced CSS with animations
- **⭐ Favorites System** - Save and manage favorite templates with localStorage
- **🔄 Form States** - Loading states, error handling, and validation feedback
- **🛡️ Security Enhanced** - Local asset bundling, no CDN dependencies

### Template System with Aliases

Pre-built templates for common task types:

- **🐛 Bug Report** (`bug`) - Structured bug documentation with required fields
- **✨ Feature Request** (`feature`) - New feature specifications with validation
- **🔬 Research** (`research`) - Investigation and learning tasks with time tracking
- **⚡ Spike** (`spike`) - Time-boxed exploration with scope definition
- **🔒 Security** (`security`) - Security-related tasks with impact assessment
- **🚀 Performance** (`performance`) - Optimization tasks with metrics tracking
- **🔧 Refactoring** (`refactor`) - Code improvement tasks with before/after
- **💳 Technical Debt** (`debt`) - Code debt tracking with priority scoring

### Batch Operations

Efficient bulk operations for productivity:

- **✅ Multi-Select** - Ctrl/Cmd+Click for multiple selection
- **📊 Batch Status Update** - Change status for multiple tasks
- **🗑️ Bulk Delete** - Delete multiple tasks at once
- **🏷️ Tag Management** - Add/remove tags in bulk
- **📁 Project Migration** - Move multiple tasks between projects

### Analytics & Reporting

Built-in analytics for project insights:

- **📊 Progress Tracking** - Visual progress charts
- **⏱️ Time Analysis** - Task completion patterns
- **🎯 Priority Distribution** - Priority level analysis
- **📅 Due Date Insights** - Deadline compliance tracking
- **👥 Dependency Analysis** - Dependency complexity metrics

## 🎮 Usage Examples

### Getting Started Workflow

```
1. 📦 Install Extension → Search "Gorev" in VS Code marketplace
2. 🚀 Auto-Setup → Extension automatically configures NPX mode
3. 📁 Create Project → "Web Application Development"
4. 🎯 Add Tasks → Use template wizard for structured tasks
5. 🌳 Build Hierarchy → Create subtasks with unlimited nesting
6. 🔗 Set Dependencies → Link related tasks for workflow
7. 📊 Track Progress → Watch visual progress indicators
8. 🔍 Use Search → Find tasks quickly with FTS5 search
9. 📤 Export Data → Share progress with team via CSV/JSON
```

### AI Assistant Workflow

```
🤖 "Hey Claude, I need help organizing my project tasks"
🗨️ "Create a bug report task for the login form validation issue"
   → Creates structured bug report with severity, steps, environment
🗨️ "Show me all high-priority tasks that are overdue"
   → Filters and displays urgent tasks needing attention
🗨️ "Mark the JWT middleware task as completed"
   → Updates status and auto-resolves dependent tasks
🗨️ "Export all completed tasks from this sprint to CSV"
   → Generates report for sprint review meeting
```

### Advanced Search Examples

```
🔍 "authentication bug high"     → Fuzzy search across titles/descriptions
🔍 "status:pending priority:high" → Structured filter query
🔍 "project:WebApp overdue"       → Project-specific overdue tasks
🔍 "tags:security,urgent"        → Multi-tag intersection search
🔍 "created:last-week"           → Date-relative search
```

## 🛠️ Installation Methods

### Method 1: VS Code Marketplace (Recommended)

```
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Gorev"
4. Click Install
5. Start using immediately!
```

### Method 2: Command Line

```bash
code --install-extension mehmetsenol.gorev-vscode
```

### Method 3: VSIX File

Download from [GitHub Releases](https://github.com/msenol/Gorev/releases) and install manually.

## 🔧 Troubleshooting

### Common Issues

**NPX Mode Not Working?**

```bash
# Check Node.js version (requires 14+)
node --version

# Test NPX directly
npx @mehmetsenol/gorev-mcp-server@latest --version
```

**Binary Mode Connection Issues?**

```bash
# Verify binary installation
gorev version

# Check binary path in settings
"gorev.serverPath": "/usr/local/bin/gorev"
```

**Extension Not Loading?**

1. Check VS Code Output → Gorev channel
2. Restart VS Code
3. Check for conflicting extensions
4. Reset extension settings

### Debug Mode

Enable debug logging for troubleshooting:

```json
{
  "gorev.debug.enabled": true,
  "gorev.debug.logLevel": "debug",
  "gorev.debug.showInOutput": true
}
```

## 📈 Performance & Statistics

### Extension Metrics

- **📊 Test Coverage**: 100% (VS Code extension)
- **🎯 MCP Tools**: 48 tools available
- **🌍 Languages**: English + Turkish support
- **💻 Platforms**: Windows, macOS, Linux
- **⚡ Performance**: 90% operation reduction
- **🔧 Configuration**: 50+ customizable settings
- **📱 Commands**: 50+ available commands
- **🎨 Visual**: 15+ customization options

### Architecture Highlights

- **🏗️ TypeScript**: Strict mode with full type safety
- **🔒 Thread Safety**: Race-condition free operations
- **⚡ Async Operations**: Non-blocking UI interactions
- **📊 Memory Efficient**: Smart caching and cleanup
- **🔄 Reactive Updates**: Event-driven architecture
- **🎯 Modular Design**: Clean separation of concerns

## 🤝 Integration Points

### MCP Clients Compatibility

| Client | Status | Features |
|--------|--------|----------|
| **Claude Desktop** | ✅ Full | All 48 MCP tools, conversation integration |
| **VS Code MCP** | ✅ Full | Native extension, direct integration |
| **Cursor IDE** | ✅ Full | AI coding assistant, context awareness |
| **Windsurf** | ✅ Full | Development environment integration |
| **Zed Editor** | 🔄 Planned | Future MCP support integration |

### Development Tools

- **Git Integration** - Track file changes and task updates
- **Project Templates** - Scaffold new projects with task templates
- **CI/CD Hooks** - Integrate with build and deployment pipelines
- **Documentation** - Auto-generate documentation from task structure

## 📚 Resources & Support

### Documentation

- 📖 [Main Repository](https://github.com/msenol/Gorev) - Complete source code and docs
- 🔧 [MCP Tools Reference](https://github.com/msenol/Gorev/blob/main/docs/mcp-araclari.md) - All 48 tools documented
- 📋 [Installation Guide](https://github.com/msenol/Gorev/blob/main/README.en.md#-installation) - Binary setup instructions
- 🎯 [VS Code Extension Guide](https://github.com/msenol/Gorev/blob/main/docs/user-guide/vscode-extension.md) - Advanced usage

### Community & Support

- 🐛 [Issue Tracker](https://github.com/msenol/Gorev/issues) - Bug reports and feature requests
- 💬 [Discussions](https://github.com/msenol/Gorev/discussions) - Community discussions
- ❓ [FAQ](https://github.com/msenol/Gorev/wiki/FAQ) - Frequently asked questions
- 📧 [Contact](mailto:me@mehmetsenol.dev) - Direct developer contact

### Contribution

1. 🍴 Fork the repository
2. 🌿 Create a feature branch
3. ✨ Make your changes
4. 🧪 Add tests if applicable
5. 📝 Submit a pull request

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **MCP Protocol** - For enabling seamless AI integration
- **SQLite FTS5** - For powerful full-text search capabilities
- **VS Code API** - For extensible editor integration
- **Community** - For feedback, bug reports, and feature requests

---

<div align="center">

**Made with ❤️ for productive developers**

[⬆ Back to Top](#gorev---advanced-task-management--ai-integration-for-vs-code)

**Try it now:** [Install from VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

</div>
