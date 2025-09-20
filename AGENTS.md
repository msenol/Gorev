# Gorev Project - Agent Development Guide

## Build/Lint/Test Commands

### Go Server (gorev-mcpserver)
```bash
make build              # Build for current platform
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make fmt                # Format code
make lint               # Run golangci-lint
make deps               # Download and tidy dependencies
go test -v ./...        # Run all tests with verbose output
go test -v -run TestName  # Run single test
```

### VS Code Extension (gorev-vscode)
```bash
npm install             # Install dependencies
npm run compile         # Compile TypeScript
npm test                # Run tests
npm run lint            # Run ESLint
npm run package         # Package extension
```

### Root Project
```bash
make build              # Build both modules
make test               # Run all tests
make fmt                # Format both modules
make lint               # Lint both modules
make pre-commit         # Run pre-commit checks
```

## Code Style Guidelines

### Go Code
- **Imports**: Group standard library, third-party, local packages
- **Formatting**: Use `make fmt` (go fmt)
- **Naming**: Turkish domain terms (gorev, proje), English technical terms
- **Error Handling**: Always return explicit errors, wrap with context
- **i18n**: Use `i18n.T("key", templateData)` for user-facing strings
- **Testing**: Use standardized `internal/testing/helpers.go` patterns

### TypeScript Code
- **Imports**: Use ES6 imports, organize by type
- **Formatting**: TypeScript strict mode, ESLint configuration
- **Naming**: camelCase for variables, PascalCase for classes/interfaces
- **Error Handling**: Try-catch with proper error messages
- **Testing**: Mocha with Sinon mocks, 100% coverage required

### Project Rules
- **Rule 15 Compliance**: NO workarounds, NO technical debt, NO quick fixes
- **Template Usage**: Mandatory for all task creation (v0.10.0+)
- **Database**: SQLite with migrations, workspace/global support
- **Localization**: Full TR/EN support with proper i18n keys