# 📋 Release Checklist for Gorev v0.9.0

## Pre-Release Tasks

### Code & Testing

- [ ] All tests passing (`make test`)
- [ ] Code formatted (`make fmt`)
- [ ] No linting errors (`make lint`)
- [ ] Coverage > 80% (`make test-coverage`)
- [ ] Manual testing of new AI features

### Version Updates

- [x] Update Makefile version to 0.9.0
- [x] Update install.sh default version to v0.9.0
- [x] Update install.ps1 default version to v0.9.0
- [x] VS Code extension already at 0.3.5
- [x] Update CHANGELOG.md
- [x] Update README.md version badge

### Documentation

- [x] Release notes created (RELEASE_NOTES_v0.9.0.md)
- [x] AI tools documentation (docs/mcp-araclari-ai.md)
- [x] Updated mcp-tools.md for 25 tools
- [x] Updated CLAUDE.md with v0.9.0 features
- [ ] Review all documentation for accuracy

## Build & Package

### MCP Server Binaries

- [ ] Run `./scripts/build-release.sh`
- [ ] Verify all binaries created:
  - [ ] gorev-linux-amd64
  - [ ] gorev-darwin-amd64
  - [ ] gorev-darwin-arm64
  - [ ] gorev-windows-amd64.exe
- [ ] Verify checksums.txt generated

### VS Code Extension

- [ ] Run `./scripts/package-vscode-extension.sh`
- [ ] Verify gorev-vscode-0.3.5.vsix created

## GitHub Release

### Create Release

- [ ] Commit all changes
- [ ] Create and push tag: `git tag v0.9.0 && git push origin v0.9.0`
- [ ] Run `./scripts/create-github-release.sh`
- [ ] Review draft release on GitHub

### Release Artifacts

- [ ] Binary files uploaded
- [ ] Archive files uploaded (.tar.gz, .zip)
- [ ] Checksums file uploaded
- [ ] VS Code extension VSIX uploaded

### Publish Release

- [ ] Review release notes one final time
- [ ] Publish release (remove draft status)
- [ ] Verify download links work

## Post-Release

### VS Code Marketplace

- [ ] Login to VS Code marketplace
- [ ] Upload new extension version
- [ ] Update extension description if needed
- [ ] Verify extension installable

### Documentation & Communication

- [ ] Update project website (if exists)
- [ ] Create announcement blog post
- [ ] Post on social media
- [ ] Update Discord/Slack channels
- [ ] Email major users/contributors

### Verification

- [ ] Test installation script on fresh Linux VM
- [ ] Test installation script on fresh Windows VM
- [ ] Test VS Code extension installation
- [ ] Test MCP integration with Claude Desktop

### Monitoring

- [ ] Monitor GitHub issues for problems
- [ ] Check download statistics
- [ ] Respond to user feedback

## Rollback Plan

If critical issues found:

1. Mark release as pre-release on GitHub
2. Fix issues in hotfix branch
3. Create v0.9.1 patch release
4. Update installation scripts

## Notes

- AI Context Management is a major feature - monitor closely
- Token limit fixes are critical for large projects
- Natural language queries need real-world testing

---

**Release Manager**: _________________
**Date**: July 9, 2025
**Status**: [ ] Ready for Release
