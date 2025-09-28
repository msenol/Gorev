# VS Code Extension Marketplace Publishing Guide

This guide covers the process of publishing the Gorev VS Code extension to the Visual Studio Code Marketplace.

## Prerequisites

- Node.js and npm installed
- VS Code Extension CLI (`vsce`) installed: `npm install -g @vscode/vsce`
- A Microsoft account
- An Azure DevOps account

## 1. Create Publisher Account

### Step 1: Sign in to Azure DevOps

1. Go to https://dev.azure.com/
2. Sign in with your Microsoft account

### Step 2: Create Personal Access Token (PAT)

1. Navigate to User Settings → Personal Access Tokens
2. Click "New Token"
3. Configure the token:
   - Name: `vsce-publish`
   - Organization: All accessible organizations
   - Expiration: 90 days (or longer)
   - Scopes: Custom defined → "Marketplace" → Check "Manage"
4. Copy and securely save the token!

### Step 3: Create Publisher

1. Go to https://marketplace.visualstudio.com/manage
2. Click "Create Publisher"
3. Fill in details:
   - Publisher ID: `mehmetsenol` (or your chosen ID)
   - Display Name: Your name or organization
   - Description: Brief description of your work

## 2. Configure VSCE

```bash
# Login with your publisher ID
vsce login mehmetsenol
# Paste your PAT token when prompted
```

## 3. Prepare Extension for Publishing

### Update package.json

Ensure all required fields are present:

```json
{
  "name": "gorev-vscode",
  "displayName": "Gorev",
  "description": "Powerful task management for VS Code powered by MCP protocol",
  "version": "0.3.0",
  "publisher": "mehmetsenol",
  "icon": "media/icon.png",
  "repository": {
    "type": "git",
    "url": "https://github.com/msenol/gorev"
  },
  "categories": ["Other"],
  "keywords": ["task", "todo", "project management", "mcp", "gorev"]
}
```

### Create/Update .vscodeignore

```
.vscode/**
.vscode-test/**
src/**
.gitignore
.yarnrc
tsconfig.json
webpack.config.js
*.vsix
test/**
.eslintrc.json
```

### Ensure Icon Requirements

- Format: PNG (SVG not supported)
- Size: 128x128 or 256x256 pixels
- Location: As specified in package.json

## 4. Package and Test

```bash
cd gorev-vscode

# Package the extension
vsce package

# This creates a .vsix file
# Test it locally:
code --install-extension gorev-vscode-0.3.0.vsix
```

## 5. Publish to Marketplace

```bash
# Publish current version
vsce publish

# Or publish with version bump
vsce publish minor  # 0.3.0 → 0.4.0
vsce publish patch  # 0.3.0 → 0.3.1

# Or specify exact version
vsce publish 0.3.0
```

## 6. Post-Publishing

### Marketplace URL

Your extension will be available at:

```
https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode
```

### Installation

Users can install via:

```bash
# Command line
code --install-extension mehmetsenol.gorev-vscode

# Or from VS Code
# Extensions panel → Search "gorev" → Install
```

## 7. Updates and Maintenance

### Publishing Updates

```bash
# Update version in package.json
npm version patch  # or minor/major

# Publish update
vsce publish
```

### Monitor Statistics

View download counts, ratings, and reviews at:
https://marketplace.visualstudio.com/manage/publishers/mehmetsenol

## Best Practices

1. **Version Management**
   - Follow semantic versioning (MAJOR.MINOR.PATCH)
   - Update CHANGELOG.md with each release
   - Tag releases in git

2. **Quality Checks**
   - Test extension thoroughly before publishing
   - Run linter and fix all issues
   - Ensure README is professional and complete
   - Include screenshots/GIFs in README

3. **Metadata**
   - Choose appropriate categories
   - Add relevant keywords for discoverability
   - Write clear, concise description
   - Keep display name short and memorable

## Troubleshooting

### "Personal Access Token verification failure"

- Ensure PAT has Marketplace > Manage scope
- Check PAT hasn't expired
- Try creating a new PAT

### "Missing publisher name"

- Add `"publisher": "your-id"` to package.json
- Ensure you're logged in: `vsce login your-id`

### "Icon not found"

- Check icon path in package.json
- Ensure icon is PNG format
- Verify icon file exists

## Checklist

Before publishing, ensure:

- [ ] Publisher account created
- [ ] PAT token obtained and saved
- [ ] package.json has all required fields
- [ ] Icon is PNG format (128x128 or 256x256)
- [ ] README is professional and complete
- [ ] .vscodeignore excludes unnecessary files
- [ ] Extension tested locally
- [ ] Version number updated
- [ ] CHANGELOG.md updated
- [ ] Git tag created for release

## Additional Resources

- [VS Code Publishing Extensions](https://code.visualstudio.com/api/working-with-extensions/publishing-extension)
- [Extension Manifest](https://code.visualstudio.com/api/references/extension-manifest)
- [VS Code Marketplace](https://marketplace.visualstudio.com/vscode)
