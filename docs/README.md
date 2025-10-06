# Gorev Documentation

**Version:** v0.16.0
**Last Updated:** October 5, 2025
**Primary Language:** English
**Secondary Language:** Turkish (see [Legacy Documentation](#legacy-documentation))

---

## 📖 Overview

Welcome to the Gorev documentation! Gorev is a modern task management system designed for AI assistants (Claude, VS Code, Windsurf, Cursor) with MCP (Model Context Protocol) integration.

**Key Features:**

- 🌐 Embedded Web UI (React + TypeScript)
- 🗂️ Multi-Workspace Support (isolated databases per project)
- 🤖 41 MCP Tools for AI integration
- 📋 Template System with human-readable aliases
- 🔌 REST API (23 Fiber endpoints)
- 💻 VS Code Extension (optional)

---

## 🚀 Quick Navigation

### New to Gorev

Start here → **[Quick Start Guide](guides/getting-started/quick-start.md)** (10 minutes)

### Installing Gorev

See → **[Installation Guide](guides/getting-started/installation.md)** (platform-specific)

### Upgrading from v0.15.x

Read → **[Migration Guide](migration/v0.15-to-v0.16.md)** (15-30 minutes)

### Having Issues

Check → **[Troubleshooting Guide](guides/getting-started/troubleshooting.md)** (comprehensive solutions)

---

## 📚 Documentation Index

### Getting Started

| Guide | Description | Est. Time | Audience |
|-------|-------------|-----------|----------|
| [Quick Start](guides/getting-started/quick-start.md) | Get up and running with Gorev | 10 min | All users |
| [Installation Guide](guides/getting-started/installation.md) | Platform-specific installation instructions | 5 min | All users |
| [Troubleshooting](guides/getting-started/troubleshooting.md) | Common issues and solutions | As needed | All users |

### Core Features

| Guide | Description | Est. Time | Audience |
|-------|-------------|-----------|----------|
| [Web UI Guide](guides/features/web-ui.md) | Embedded React interface documentation | 20 min | Web UI users |
| [Multi-Workspace Support](guides/features/multi-workspace.md) | Managing multiple isolated workspaces | 15 min | Advanced users |
| [Template System](guides/features/template-system.md) | Task templates and aliases | 15 min | All users |
| [AI Context Management](guides/features/ai-context-management.md) | AI assistant integration | 15 min | AI users |

### Configuration & Setup

| Guide | Description | Est. Time | Audience |
|-------|-------------|-----------|----------|
| [MCP Configuration Examples](guides/mcp-config-examples.md) | IDE setup guides (Claude, VS Code, Cursor, Windsurf) | 10 min | AI users |
| [VS Code Extension](guides/user/vscode-extension.md) | Extension features and usage | 15 min | VS Code users |
| [VS Code Export/Import](guides/user/vscode-data-export-import.md) | Data migration guide | 10 min | VS Code users |
| [Usage Guide](guides/user/usage.md) | Detailed usage examples | 20 min | All users |

### Reference Documentation

| Reference | Description | Audience |
|-----------|-------------|----------|
| [MCP Tools Reference](legacy/tr/mcp-araclari.md) | Complete reference for 41 MCP tools (Turkish) | Developers |
| [MCP Tools Reference (API)](api/MCP_TOOLS_REFERENCE.md) | API documentation for MCP tools | Developers |

### Development

| Guide | Description | Audience |
|-------|-------------|----------|
| [System Architecture](architecture/architecture-v2.md) | Technical architecture details | Developers |
| [Contributing Guide](development/contributing.md) | How to contribute to Gorev | Contributors |
| [Development History](development/TASKS.md) | Complete project history | Developers |
| [Roadmap](../ROADMAP.md) | Development roadmap and future plans | All |

### Migration & Upgrades

| Guide | Description | Est. Time | Audience |
|-------|-------------|-----------|----------|
| [v0.15 → v0.16 Migration](migration/v0.15-to-v0.16.md) | Upgrade from v0.15.x to v0.16.0 | 15-30 min | Existing users |

### Release Information (v0.16.0)

| Document | Description | Audience |
|----------|-------------|----------|
| [Bug Fixes Summary](releases/v0.16.0_bug_fixes_summary.md) | Critical bug fixes and improvements | All users |
| [Testing Guide](guides/user/bug_fixes_testing_guide_v0.16.0.md) | Bug fix testing procedures | Testers |
| [Documentation Update Report](development/documentation_update_v0.16.0.md) | Documentation changes | Developers |
| [Release Notes](releases/RELEASE_NOTES_v0.16.0.md) | Full release documentation | All users |
| [Changelog](../CHANGELOG.md) | Complete version history | All users |

---

## 🌍 Language Support

### Primary Language: English

All new documentation (v0.16.0+) is written in English as the primary language. This includes:

- Getting Started guides
- Feature documentation
- Reference documentation
- Migration guides

### Secondary Language: Turkish

Legacy Turkish documentation has been preserved in the `legacy/tr/` directory:

- [MCP Araçları Referansı](legacy/tr/mcp-araclari.md) - Comprehensive MCP tools reference in Turkish
- [Kullanım Örnekleri](legacy/tr/ornekler.md) - Usage examples in Turkish
- And other Turkish legacy docs

**Main README files:**

- [README.md](../README.md) - English
- [README.tr.md](../README.tr.md) - Turkish

**AI Assistant Instructions:**

- [CLAUDE.en.md](../CLAUDE.en.md) - English
- [CLAUDE.md](../CLAUDE.md) - Turkish

---

## 📦 Documentation Structure

```
docs/
├── README.md                          # This file - Documentation index
├── guides/                            # User guides (English)
│   ├── getting-started/              # Getting started guides
│   │   ├── quick-start.md           # 10-minute quick start
│   │   ├── installation.md          # Installation guide
│   │   └── troubleshooting.md       # Troubleshooting guide
│   ├── features/                     # Feature documentation
│   │   ├── web-ui.md                # Web UI guide
│   │   ├── multi-workspace.md       # Multi-workspace guide
│   │   ├── template-system.md       # Template system guide
│   │   └── ai-context-management.md # AI context guide
│   ├── user/                         # User guides
│   │   ├── usage.md                 # Usage guide
│   │   ├── vscode-extension.md      # VS Code extension
│   │   └── vscode-data-export-import.md # Export/import
│   └── mcp-config-examples.md        # MCP configuration
├── legacy/                           # Legacy documentation
│   └── tr/                           # Turkish documentation (legacy)
│       ├── mcp-araclari.md          # MCP tools reference (TR)
│       ├── ornekler.md              # Usage examples (TR)
│       └── ... (other Turkish docs)
├── migration/                        # Migration guides
│   └── v0.15-to-v0.16.md            # v0.15 → v0.16 migration
├── architecture/                     # Architecture documentation
│   └── architecture-v2.md           # System architecture
├── development/                      # Development documentation
│   ├── TASKS.md                     # Development history
│   └── contributing.md              # Contributing guide
├── api/                             # API documentation
│   └── MCP_TOOLS_REFERENCE.md       # MCP API reference
└── releases/                        # Release documentation
    ├── v0.16.0_bug_fixes_summary.md
    └── RELEASE_NOTES_v0.16.0.md
```

---

## 🔍 Finding What You Need

### By Role

**End Users (Task Management)**

1. Start with [Quick Start Guide](guides/getting-started/quick-start.md)
2. Configure your AI assistant: [MCP Configuration](guides/mcp-config-examples.md)
3. Learn templates: [Template System](guides/features/template-system.md)
4. Explore Web UI: [Web UI Guide](guides/features/web-ui.md)

**VS Code Users**

1. Install extension: [VS Code Extension Guide](guides/user/vscode-extension.md)
2. Set up workspace: [Multi-Workspace Guide](guides/features/multi-workspace.md)
3. Export/import data: [VS Code Export/Import](guides/user/vscode-data-export-import.md)

**Developers (Contributing)**

1. Read architecture: [System Architecture](architecture/architecture-v2.md)
2. Review contributing guide: [Contributing Guide](development/contributing.md)
3. Understand MCP tools: [MCP Tools Reference](legacy/tr/mcp-araclari.md)
4. Check development history: [Development History](development/TASKS.md)

**AI Assistant Users (Claude, Copilot, etc.)**

1. Configure MCP: [MCP Configuration Examples](guides/mcp-config-examples.md)
2. Understand AI context: [AI Context Management](guides/features/ai-context-management.md)
3. Learn MCP tools: [MCP Tools Reference](legacy/tr/mcp-araclari.md)

### By Task

**Setting Up Gorev**
→ [Installation Guide](guides/getting-started/installation.md) → [Quick Start](guides/getting-started/quick-start.md)

**Managing Multiple Projects**
→ [Multi-Workspace Guide](guides/features/multi-workspace.md)

**Creating Structured Tasks**
→ [Template System Guide](guides/features/template-system.md)

**Integrating with AI**
→ [AI Context Management](guides/features/ai-context-management.md) → [MCP Configuration](guides/mcp-config-examples.md)

**Troubleshooting Issues**
→ [Troubleshooting Guide](guides/getting-started/troubleshooting.md)

**Upgrading Gorev**
→ [Migration Guide](migration/v0.15-to-v0.16.md)

---

## 📊 Documentation Stats

| Category | Files | Total Words | Status |
|----------|-------|-------------|--------|
| Getting Started | 3 | ~25,000 | ✅ Complete |
| Features | 4 | ~47,000 | ✅ Complete |
| User Guides | 3 | ~15,000 | ✅ Complete |
| Migration | 1 | ~8,000 | ✅ Complete |
| Reference | 2 | ~30,000 | ✅ Complete |
| Development | 3 | ~20,000 | ✅ Complete |
| **Total** | **16+** | **~145,000** | **✅ v0.16.0** |

---

## 🆘 Getting Help

### Documentation Issues

- **Broken links?** → [Open an issue](https://github.com/msenol/gorev/issues)
- **Unclear documentation?** → [Open an issue](https://github.com/msenol/gorev/issues)
- **Missing information?** → [Open an issue](https://github.com/msenol/gorev/issues)

### Technical Support

- **Bug reports** → [GitHub Issues](https://github.com/msenol/gorev/issues)
- **Feature requests** → [GitHub Discussions](https://github.com/msenol/gorev/discussions)
- **Questions** → [GitHub Discussions](https://github.com/msenol/gorev/discussions)

### Community

- **GitHub Repository**: https://github.com/msenol/gorev
- **VS Code Marketplace**: https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode
- **Issue Tracker**: https://github.com/msenol/gorev/issues
- **Discussions**: https://github.com/msenol/gorev/discussions
- **Wiki**: https://github.com/msenol/gorev/wiki

---

## 🔄 Version History

### v0.16.0 (October 4, 2025) - Current

- Embedded Web UI (React + TypeScript)
- Multi-workspace support
- REST API (23 endpoints)
- Template aliases (bug, feature, research, etc.)
- VS Code extension REST API migration
- **Documentation overhaul**: 60,000+ words of new English documentation

### v0.15.x (September 2025)

- Advanced search & filtering (FTS5, fuzzy matching)
- Filter profiles
- Performance improvements

### v0.14.x (August 2025)

- Data export/import (JSON/CSV)
- Enhanced error handling

### v0.13.x (July 2025)

- IDE extension management
- Multi-IDE support

See [ROADMAP.md](../ROADMAP.md) for future plans.

---

## 📝 Contributing to Documentation

We welcome documentation contributions! Please see:

- [Contributing Guide](development/contributing.md) for general guidelines
- Documentation follows [Markdown best practices](https://www.markdownguide.org/basic-syntax/)
- Primary language: English
- All guides should include: version info, estimated reading time, last updated date

### Documentation Checklist

- [ ] Clear, concise writing
- [ ] Code examples tested
- [ ] Screenshots up-to-date
- [ ] Links verified
- [ ] Version info included
- [ ] Last updated date current

---

## 📄 License

All documentation is released under the same [MIT License](../LICENSE) as the Gorev project.

---

<div align="center">

**[⬆ Back to Top](#gorev-documentation)**

Made with ❤️ by the [Gorev contributors](https://github.com/msenol/gorev/graphs/contributors)

*Documentation enhanced by Claude (Anthropic) - Your AI pair programming assistant*

</div>
