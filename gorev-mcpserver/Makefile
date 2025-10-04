# Değişkenler
BINARY_NAME=gorev
GO=go
GOFLAGS=-v

# Versiyon bilgisi
VERSION ?= 0.16.0
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: all build build-all build-web clean test run install lint fmt

all: clean test build

# Build web UI first, then Go binary
build: build-web
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gorev

# Build web UI production bundle
build-web:
	@echo "Building Web UI..."
	@cd ../gorev-web && npm run build
	@echo "Web UI built successfully"

build-all: build-web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/gorev
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/gorev
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/gorev

run: build
	./$(BINARY_NAME) serve

test:
	$(GO) test -v -cover ./...

test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f *.db
	rm -rf dist/
	rm -f coverage.out coverage.html

install:
	$(GO) install $(LDFLAGS) ./cmd/gorev

lint:
	golangci-lint run

fmt:
	$(GO) fmt ./...

deps:
	$(GO) mod download
	$(GO) mod tidy

docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	docker run --rm -it $(BINARY_NAME):$(VERSION)