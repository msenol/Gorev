# Gorev VS Code Extension

<p align="center">
  <img src="media/icon.png" alt="Gorev Logo" width="128" height="128">
</p>

<div align="center">

[🇺🇸 English](README.md) | [🇹🇷 Türkçe](README.tr.md)

[![Version](https://img.shields.io/badge/Version-0.5.1-blue?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Downloads](https://img.shields.io/visual-studio-marketplace/d/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Rating](https://img.shields.io/visual-studio-marketplace/r/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**Powerful task management extension for VS Code with unlimited subtask hierarchy, dependency visualization, and AI assistant integration**

> 🌐 **NEW in v0.5.0**: Full bilingual support! The extension now automatically displays in English or Turkish based on your VS Code language settings.

</div>

## ✨ Key Features

### 🌍 **Bilingual Support** (v0.5.0+)
- Automatic language detection based on VS Code settings
- Complete English and Turkish translations for all UI elements
- No configuration needed - seamless language switching

### 🌳 **Unlimited Subtask Hierarchy**
- Create tasks within tasks with infinite nesting levels
- Visual tree structure with progress tracking
- Parent task completion based on subtask status
- Intuitive expand/collapse navigation

### 🔗 **Advanced Dependencies**
- Create task dependencies with visual indicators
- Smart validation to prevent circular dependencies
- Dependency status tracking (🔒 blocked, 🔓 ready, 🔗 linked)
- Bulk dependency management

### 🤖 **AI Assistant Integration**
- **MCP Protocol**: Works with Claude, Windsurf, Cursor, and other AI assistants
- **Natural Language**: Create and manage tasks through conversation
- **Context Awareness**: AI understands your project structure
- **Template System**: Structured task creation with predefined formats

### 📊 **Rich Visual Interface**
- **Enhanced TreeView**: Group by status, priority, project, tags, or due dates
- **Progress Tracking**: Visual progress bars for parent tasks
- **Color Coding**: Priority-based visual distinction
- **Smart Badges**: Task counts, due dates, and dependency indicators

### 🎛️ **Advanced Task Management**
- **Multi-select Operations**: Bulk status updates and deletions
- **Drag & Drop**: Move tasks between projects, change status/priority
- **Inline Editing**: F2 or double-click for quick edits
- **Smart Filtering**: Real-time search with saved filter profiles
- **Template Wizard**: Create tasks from predefined templates

## 🚀 Getting Started

### Installation

**Option 1: VS Code Marketplace (Recommended)**
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Gorev"
4. Click Install

**Option 2: Command Line**
```bash
code --install-extension mehmetsenol.gorev-vscode
```

### Setup

1. **Install Gorev MCP Server**: Follow the [installation guide](https://github.com/msenol/Gorev/blob/main/README.en.md#-installation)

2. **Configure Extension**:
   - Open VS Code Settings (Ctrl+,)
   - Search for "Gorev"
   - Set the path to your Gorev executable:
   ```json
   {
     "gorev.serverPath": "/path/to/gorev"
   }
   ```

3. **Start Using**:
   - Click the Gorev icon in the Activity Bar
   - Create your first project
   - Start managing tasks!

## 🎮 Usage Examples

### Basic Task Management
```
1. Create a project: "Web Development"
2. Add task: "Implement user authentication"
3. Create subtasks:
   └─ Set up JWT middleware
   └─ Create login form
   └─ Add password validation
4. Set dependencies and priorities
5. Track progress automatically
```

### AI Assistant Integration
Talk to Claude, Cursor, or Windsurf:
```
"Create a new task for implementing dark mode with high priority"
"Show me all overdue tasks"
"Mark task #5 as completed"
"Create a bug report template task for the login issue"
```

### Template Usage
Use predefined templates for consistent task creation:
- **Bug Report**: Structured bug documentation
- **Feature Request**: New feature specifications  
- **Technical Debt**: Code improvement tasks
- **Research**: Investigation and learning tasks

## 📋 Features Overview

### Enhanced TreeView
- **Grouping**: Status, priority, project, tag, or due date grouping
- **Multi-Select**: Ctrl/Cmd+Click for bulk operations
- **Sorting**: Title, priority, due date, creation date
- **Color Coding**: Priority-based visual indicators
- **Quick Actions**: One-click completion with checkboxes

### Drag & Drop Support
- 🔄 Move tasks between projects
- 📊 Change status by dropping on status groups
- 🎯 Reorder priorities
- 🔗 Create dependencies (drop task on another)
- ✨ Visual feedback and animations

### Inline Editing
- ✏️ F2 or double-click for quick editing
- 📝 Context menu for status/priority changes
- 📅 Inline date picker
- ❌ Escape to cancel, ✅ Enter to save

### Advanced Filtering
- 🔍 Real-time search across all fields
- 🎛️ Advanced filters (status, priority, tags, dates)
- 💾 Saved filter profiles
- 📊 Status bar integration
- ⚡ Quick filter shortcuts

### Task Detail Panel
- 📝 Rich markdown editor for descriptions
- 🔗 Dependency visualization
- 📊 Progress tracking with charts
- 🏷️ Tag management
- 📅 Due date management

## 🔧 Configuration

### Server Settings
```json
{
  "gorev.serverPath": "/path/to/gorev",
  "gorev.autoConnect": true,
  "gorev.showStatusBar": true,
  "gorev.refreshInterval": 30
}
```

### Visual Preferences
```json
{
  "gorev.treeView.grouping": "status",
  "gorev.treeView.sorting": "priority", 
  "gorev.treeView.visuals.showProgressBars": true,
  "gorev.treeView.visuals.showPriorityBadges": true,
  "gorev.dragDrop.allowTaskMove": true
}
```

## 🎯 Use Cases

- **Software Development**: Track features, bugs, and technical debt
- **Project Management**: Organize complex projects with subtasks
- **Team Collaboration**: Share project status and dependencies
- **Personal Productivity**: Manage daily tasks and goals
- **Research Projects**: Track investigation and learning tasks

## 🔗 Related Links

- **Main Repository**: [GitHub](https://github.com/msenol/Gorev)
- **MCP Server Documentation**: [README.en.md](https://github.com/msenol/Gorev/blob/main/README.en.md)
- **Issue Tracker**: [GitHub Issues](https://github.com/msenol/Gorev/issues)
- **Discussions**: [GitHub Discussions](https://github.com/msenol/Gorev/discussions)

## 📊 Metrics

- **Test Coverage**: 100% (VS Code Extension)
- **MCP Tools**: 25+ tools available
- **Languages**: English + Turkish support
- **Platforms**: Windows, macOS, Linux

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ for productive developers**

[⬆ Back to Top](#gorev-vs-code-extension)

</div>