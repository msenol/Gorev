# 🚀 Gorev v0.9.0 Release Notes

**Release Date:** July 9, 2025

## 🎯 Highlights

This release introduces **AI Context Management & Automation**, a major feature that revolutionizes how AI assistants interact with the Gorev task management system. With 6 new AI-optimized MCP tools, automatic state transitions, and natural language query support, AI assistants can now work more efficiently and intuitively with tasks.

## ✨ What's New

### 🤖 AI Context Management System

#### New MCP Tools (6 tools added - Total: 25 tools)

1. **`gorev_set_active`** - Set and track active task
   - Automatically transitions tasks from "beklemede" to "devam_ediyor"
   - Maintains context across AI interactions
   - Tracks last 10 recent tasks

2. **`gorev_get_active`** - Get current active task
   - Returns detailed markdown-formatted task information
   - Shows dependencies, tags, and progress

3. **`gorev_recent`** - List recently interacted tasks
   - Configurable limit (default: 5)
   - Shows task status and priority

4. **`gorev_context_summary`** - AI-optimized session overview
   - Active task information
   - Session statistics (created, updated, completed)
   - Priority tasks and blockers
   - Working project context

5. **`gorev_batch_update`** - Bulk update multiple tasks
   - Efficient bulk operations for AI workflows
   - Supports status updates (more fields coming soon)
   - Detailed success/failure reporting

6. **`gorev_nlp_query`** - Natural language task search
   - Turkish language support
   - Query patterns: "bugün", "yüksek öncelikli", "tamamlanmamış"
   - Tag search: "etiket:bug" or "tag:frontend"
   - Smart text matching in titles and descriptions

### 🔄 Automatic State Management

- Tasks automatically transition from "beklemede" to "devam_ediyor" when viewed by AI
- Silent updates without disrupting workflow
- Interaction tracking for better context awareness

### 🗄️ Database Enhancements

- New tables: `ai_interactions`, `ai_context`
- New columns in tasks: `last_ai_interaction`, `estimated_hours`, `actual_hours`
- Migration 000006 adds full AI context support

### 📊 Performance Improvements

- Pagination support in `gorev_listele` and `proje_gorevleri`
- Response size optimization (60% reduction)
- Token limit prevention with 20K character safety limit
- Compact formatting for better AI consumption

## 🐛 Bug Fixes

- Fixed token limit errors in MCP tools
- Fixed task count display issues
- Improved error handling in batch operations

### VS Code Extension v0.3.7 Fixes
- Fixed task list not showing due to parser not recognizing 🔄 emoji
- Fixed subtask hierarchy preservation in TreeView
- Fixed multiline task description parsing
- Enhanced parser for MCP v0.8.1+ compact format compatibility

## 📚 Documentation

- Comprehensive AI tools documentation in `docs/mcp-araclari-ai.md`
- Updated all documentation to reflect 25 total MCP tools
- Added AI usage patterns and examples

## 🔧 Technical Details

### Dependencies
- Added `github.com/adrg/strutil` v0.3.1 for NLP support

### Breaking Changes
- None - All changes are backward compatible

### Migration Notes
- Run migrations automatically on first start
- Existing tasks will work seamlessly with new AI features

## 📦 Installation

### Quick Install (Linux/macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | VERSION=v0.9.0 bash
```

### Quick Install (Windows PowerShell)
```powershell
$env:VERSION="v0.9.0"; irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

### Manual Download
Download binaries for your platform:
- [gorev-linux-amd64](https://github.com/msenol/Gorev/releases/download/v0.9.0/gorev-linux-amd64)
- [gorev-darwin-amd64](https://github.com/msenol/Gorev/releases/download/v0.9.0/gorev-darwin-amd64)
- [gorev-darwin-arm64](https://github.com/msenol/Gorev/releases/download/v0.9.0/gorev-darwin-arm64)
- [gorev-windows-amd64.exe](https://github.com/msenol/Gorev/releases/download/v0.9.0/gorev-windows-amd64.exe)

## 🎨 VS Code Extension

The VS Code extension (v0.3.7) has been updated with:
- Pagination support for large task lists
- Improved performance with token limit prevention
- Better integration with AI tools
- **CRITICAL FIX**: Parser now correctly displays all tasks with MCP v0.8.1+ format
- Full support for task status emojis (⏳, 🚀, ✅, ✓, 🔄)
- Preserved subtask hierarchy in TreeView

Install from [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode) or download [gorev-vscode-0.3.7.vsix](https://github.com/msenol/Gorev/releases/download/v0.9.0/gorev-vscode-0.3.7.vsix)

## 🙏 Acknowledgments

Special thanks to all contributors and users who provided feedback for this release.

---

**Full Changelog**: https://github.com/msenol/Gorev/compare/v0.8.1...v0.9.0