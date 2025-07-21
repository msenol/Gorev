# Release Notes - Gorev VS Code Extension v0.5.0

## 🌍 Bilingual Support Release

We're excited to announce **Gorev VS Code Extension v0.5.0**, featuring complete bilingual support for English and Turkish users!

### ✨ What's New

#### 🌐 Complete Bilingual Interface
- **Automatic Language Detection**: The extension now automatically displays in your preferred language based on your VS Code language setting
- **No Configuration Needed**: Simply use VS Code in English or Turkish, and Gorev adapts automatically
- **500+ Translations**: Every UI element, command, notification, and error message is now available in both languages

#### 📋 What's Localized
- ✅ All 21 VS Code commands with titles and descriptions
- ✅ TreeView panels (tasks, projects, templates)
- ✅ Filter toolbar with search and advanced filtering
- ✅ Status bar messages and tooltips
- ✅ Task detail panel with markdown editor
- ✅ Template wizard for guided task creation
- ✅ Drag & drop operation feedback
- ✅ Inline editing validation messages
- ✅ Debug tools and test data seeders
- ✅ All error messages and notifications

### 🛠️ Technical Details

#### Localization Implementation
- Uses VS Code's modern `l10n` API with `vscode.l10n.t()` for dynamic translations
- Bundle-based localization structure for efficient loading
- Preserves icon codes and formatting in all translations
- Supports dynamic placeholders for values like task counts and names

#### New Files
- `l10n/bundle.l10n.json` - English runtime strings
- `l10n/bundle.l10n.tr.json` - Turkish translations
- `package.nls.json` - English VS Code marketplace metadata
- `package.nls.tr.json` - Turkish VS Code marketplace metadata
- `README.tr.md` - Turkish README for Turkish users

### 🚀 How to Use

The extension automatically detects your VS Code language:
- If VS Code is in Turkish → Gorev displays in Turkish
- If VS Code is in any other language → Gorev displays in English

To change your VS Code language:
1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
2. Type "Configure Display Language"
3. Select your preferred language
4. Restart VS Code

### 🙏 Acknowledgments

This release represents a significant milestone in making Gorev accessible to a broader international audience. We're committed to providing an excellent user experience regardless of language preference.

### 📝 Coming Next

We're exploring support for additional languages based on user demand. If you'd like to see Gorev in your language, please let us know by opening an issue on GitHub!

---

**Full Changelog**: https://github.com/msenol/Gorev/compare/v0.4.6...v0.5.0