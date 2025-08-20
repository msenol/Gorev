# Değişkenler
BINARY_NAME=gorev
GO=go
GOFLAGS=-v

# Versiyon bilgisi
VERSION ?= 0.12.0
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test run install lint fmt

all: clean test build

build:
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gorev

build-all:
	# Linux with CGO enabled (native)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/gorev
	# macOS and Windows without CGO (cross-compile friendly, note: may have SQLite limitations)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/gorev
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/gorev

# Alternative: Build all with CGO (requires platform-specific toolchains)
build-all-cgo:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64-cgo ./cmd/gorev
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64-cgo ./cmd/gorev || echo "macOS CGO build failed, need osxcross"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64-cgo.exe ./cmd/gorev || echo "Windows CGO build failed, need mingw"

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