# Gorev Release Process

## Pre-release Checks

- Run linting: `make lint`
- Run tests: `make test`
- Run test coverage: `make test-coverage`
- Test VS Code extension: `cd gorev-vscode && npm test`
- Test NPM package: `cd gorev-npm && npm test`

## Version Management

- Check current version: `cat VERSION`
- Verify version consistency:
  - `cd gorev-vscode && grep -r "version" package.json`
  - `cd gorev-npm && grep -r "version" package.json`
- Update version numbers in all package.json files and VERSION file

## Build All Components

### Go Server (gorev-mcpserver)

- Build: `cd gorev-mcpserver && make build`
- Build all platforms: `cd gorev-mcpserver && make build-all-platforms`

### VS Code Extension (gorev-vscode)

- Install dependencies: `cd gorev-vscode && npm install`
- Compile: `cd gorev-vscode && npm run compile`
- Package: `cd gorev-vscode && npm run package`
- Lint: `cd gorev-vscode && npm run lint`

### NPM Package (gorev-npm)

- Install dependencies: `cd gorev-npm && npm install`
- Build: `cd gorev-npm && npm run build`
- Lint: `cd gorev-npm && npm run lint`

### Docker Images

- Build Docker image: `cd docker && docker build -f Dockerfile.release -t gorev:latest .`
- Build with compose: `cd docker && docker-compose -f docker-compose.release.yml build`

## Create Release Artifacts

- Build release: `./scripts/build-release.sh`
- Verify binaries: `ls -la release-v*/`
- Generate checksums: `cd release-v* && sha256sum * > checksums.txt`
- Test binary functionality locally

## GitHub Release Preparation

- Check git status: `git status`
- Create and push tag: `git tag -a v$(cat VERSION) -m "Release v$(cat VERSION)"` && `git push origin v$(cat VERSION)`
- Generate changelog: `git log --oneline $(git describe --tags --abbrev=0)..HEAD`

## NPM Package Release

- Login to npm: `cd gorev-npm && npm login`
- Publish package: `cd gorev-npm && npm publish --access public`
- Verify publication: `npm view gorev-mcp`

## VS Code Extension Release

- Install vsce: `npm install -g @vscode/vsce`
- Publish extension: `cd gorev-vscode && vsce publish`
- Verify marketplace publication

## GitHub Release Creation

Create GitHub release with:

- All binary artifacts (gorev-linux-amd64, gorev-darwin-amd64, gorev-windows-amd64.zip)
- Checksums file
- Comprehensive release notes
- Installation instructions
- Changelog
- Breaking changes documentation

## Docker Image Publishing

- Tag images: `docker tag gorev:latest gorev:$(cat VERSION)`
- Push images: `docker push gorev:latest` && `docker push gorev:$(cat VERSION)`

## Quality Assurance

- Test npm installation: `npm install -g @mehmetsenol/gorev-mcp-server`
- Test VS Code extension installation
- Verify binary downloads from GitHub release
- Test Docker image functionality
- Validate documentation links

## Documentation Updates

Update version references in:

- README.md files
- Documentation files
- CHANGELOG.md
- API documentation
- Installation guides

## Post-release Tasks

- Create release announcement
- Update project roadmap
- Archive release artifacts
- Clean up temporary files
- Monitor for post-release issues
- Update version to next development version

## Verification Checklist

- [ ] All tests pass
- [ ] Linting successful
- [ ] Version numbers consistent
- [ ] All binaries built successfully
- [ ] Checksums generated
- [ ] Git tag created and pushed
- [ ] NPM package published
- [ ] VS Code extension published
- [ ] GitHub release created
- [ ] Docker images pushed
- [ ] Documentation updated
- [ ] Installation tested
- [ ] Release notes published
