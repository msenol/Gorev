# Gorev Project Root Makefile
# Manages both gorev-mcpserver (Go) and gorev-vscode (TypeScript) modules

.PHONY: all build clean test run install lint fmt deps server-build server-test server-run extension-build extension-test help

# Default target
all: deps build test

help:
	@echo "Available targets:"
	@echo "  all           - Download deps, build both modules, run tests"
	@echo "  build         - Build both server and extension"
	@echo "  test          - Run all tests (server + extension)"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and install dependencies"
	@echo "  fmt           - Format code (Go + TypeScript)"
	@echo "  lint          - Run linters on both modules"
	@echo ""
	@echo "Server specific:"
	@echo "  server-build  - Build Go server"
	@echo "  server-test   - Run Go tests"
	@echo "  server-run    - Build and run server"
	@echo ""
	@echo "Extension specific:"
	@echo "  extension-build   - Compile TypeScript extension"
	@echo "  extension-test    - Run extension tests"
	@echo "  extension-package - Package extension for marketplace"

# Combined targets
build: server-build extension-build

test: server-test extension-test

clean: server-clean extension-clean

deps: server-deps extension-deps

fmt: server-fmt extension-fmt

lint: server-lint extension-lint

# Server targets (Go)
server-build:
	cd gorev-mcpserver && make build

server-test:
	cd gorev-mcpserver && make test

server-run:
	cd gorev-mcpserver && make run

server-clean:
	cd gorev-mcpserver && make clean

server-deps:
	cd gorev-mcpserver && make deps

server-fmt:
	cd gorev-mcpserver && make fmt

server-lint:
	cd gorev-mcpserver && make lint

server-coverage:
	cd gorev-mcpserver && make test-coverage

# Extension targets (TypeScript)
extension-deps:
	cd gorev-vscode && npm install --no-bin-links

extension-build:
	cd gorev-vscode && npm run compile

extension-test:
	cd gorev-vscode && npm test

extension-clean:
	cd gorev-vscode && rm -rf out/ && rm -rf node_modules/

extension-fmt:
	cd gorev-vscode && npm run format || echo "Format script not available"

extension-lint:
	cd gorev-vscode && npm run lint || echo "Lint script not available"

extension-package:
	cd gorev-vscode && npm run package || echo "Package script not available"

# Release targets
release-build:
	cd gorev-mcpserver && make build-all

install: server-build
	cd gorev-mcpserver && make install

# Development convenience
dev-server:
	cd gorev-mcpserver && ./gorev serve --debug

dev-extension:
	@echo "To develop extension: Open gorev-vscode in VS Code and press F5"

# Quality checks before commit
pre-commit: fmt lint test
	@echo "✅ Pre-commit checks passed!"