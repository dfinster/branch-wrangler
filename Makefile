# Branch Wrangler Build System
# Supports macOS Apple Silicon (darwin/arm64) as per FR-27

# Application configuration
APP_NAME := branch-wrangler
MAIN_PATH := ./cmd/branch-wrangler
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GO := /usr/local/go/bin/go
GO_VERSION := $(shell $(GO) version | cut -d' ' -f3)

# Build configuration
GOOS := darwin
GOARCH := arm64
BINARY_NAME := $(APP_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)
DIST_DIR := dist
BUILD_DIR := build

# Go build flags for version injection
LDFLAGS := -ldflags "\
	-X github.com/dfinster/branch-wrangler/internal/version.Version=$(VERSION) \
	-X github.com/dfinster/branch-wrangler/internal/version.BuildDate=$(BUILD_DATE) \
	-X github.com/dfinster/branch-wrangler/internal/version.CommitHash=$(COMMIT_HASH)"

# Build flags for optimized releases
RELEASE_FLAGS := -trimpath
DEBUG_FLAGS := -race

# Colors for output
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_RED := \033[31m

.PHONY: help build build-dev build-release dist clean test lint fmt vet check-deps version checksums install uninstall

# Default target
all: build

## help: Show this help message
help:
	@echo "$(COLOR_BLUE)Branch Wrangler Build System$(COLOR_RESET)"
	@echo "Usage: make [target]"
	@echo ""
	@echo "$(COLOR_GREEN)Development targets:$(COLOR_RESET)"
	@echo "  build         Build development binary (with race detection)"
	@echo "  build-dev     Alias for build"
	@echo "  test          Run all tests"
	@echo "  lint          Run golangci-lint"
	@echo "  fmt           Format Go code"
	@echo "  vet           Run go vet"
	@echo "  clean         Clean build artifacts"
	@echo ""
	@echo "$(COLOR_GREEN)Release targets:$(COLOR_RESET)"
	@echo "  build-release Build optimized release binary"
	@echo "  dist          Build release binary with checksums"
	@echo "  checksums     Generate checksums for existing binaries"
	@echo ""
	@echo "$(COLOR_GREEN)Utility targets:$(COLOR_RESET)"
	@echo "  version       Show version information"
	@echo "  check-deps    Check for required dependencies"
	@echo "  install       Install binary to /usr/local/bin"
	@echo "  uninstall     Remove binary from /usr/local/bin"
	@echo "  help          Show this help message"

## version: Show version information that will be embedded in build
version:
	@echo "$(COLOR_BLUE)Version Information:$(COLOR_RESET)"
	@echo "  Version:     $(VERSION)"
	@echo "  Build Date:  $(BUILD_DATE)"
	@echo "  Commit Hash: $(COMMIT_HASH)"
	@echo "  Go Version:  $(GO_VERSION)"
	@echo "  Target OS:   $(GOOS)"
	@echo "  Target Arch: $(GOARCH)"
	@echo "  Binary Name: $(BINARY_NAME)"

## check-deps: Check for required dependencies
check-deps:
	@echo "$(COLOR_BLUE)Checking dependencies...$(COLOR_RESET)"
	@$(GO) version >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: Go is not installed$(COLOR_RESET)"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: Git is not installed$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)✓ All dependencies satisfied$(COLOR_RESET)"

## build: Build development binary with race detection
build: build-dev

build-dev: check-deps
	@echo "$(COLOR_BLUE)Building development binary...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(DEBUG_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(COLOR_GREEN)✓ Development binary built: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)"

## build-release: Build optimized release binary
build-release: check-deps
	@echo "$(COLOR_BLUE)Building release binary...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(RELEASE_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(COLOR_GREEN)✓ Release binary built: $(DIST_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

## dist: Build release binary and generate checksums
dist: build-release checksums
	@echo "$(COLOR_GREEN)✓ Distribution package ready in $(DIST_DIR)/$(COLOR_RESET)"
	@ls -la $(DIST_DIR)/

## checksums: Generate SHA256 checksums for binaries
checksums:
	@echo "$(COLOR_BLUE)Generating checksums...$(COLOR_RESET)"
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "$(COLOR_RED)Error: $(DIST_DIR) directory not found. Run 'make build-release' first.$(COLOR_RESET)"; \
		exit 1; \
	fi
	@cd $(DIST_DIR) && \
	find . -name "$(APP_NAME)-*" -type f ! -name "*.txt" -exec shasum -a 256 {} \; > checksums.txt
	@echo "$(COLOR_GREEN)✓ Checksums generated: $(DIST_DIR)/checksums.txt$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Checksum contents:$(COLOR_RESET)"
	@cat $(DIST_DIR)/checksums.txt

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	$(GO) test -v ./...
	@echo "$(COLOR_GREEN)✓ All tests passed$(COLOR_RESET)"

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(COLOR_GREEN)✓ Linting completed$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)Warning: golangci-lint not installed, skipping...$(COLOR_RESET)"; \
	fi

## fmt: Format Go code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GO) vet ./...
	@echo "$(COLOR_GREEN)✓ Vet completed$(COLOR_RESET)"

## clean: Clean build artifacts
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)/
	rm -rf $(DIST_DIR)/
	rm -f $(APP_NAME)
	@echo "$(COLOR_GREEN)✓ Build artifacts cleaned$(COLOR_RESET)"

## install: Install binary to /usr/local/bin (requires sudo)
install: build-release
	@echo "$(COLOR_BLUE)Installing $(APP_NAME) to /usr/local/bin...$(COLOR_RESET)"
	@if [ ! -f "$(DIST_DIR)/$(BINARY_NAME)" ]; then \
		echo "$(COLOR_RED)Error: Binary not found. Run 'make build-release' first.$(COLOR_RESET)"; \
		exit 1; \
	fi
	sudo cp $(DIST_DIR)/$(BINARY_NAME) /usr/local/bin/$(APP_NAME)
	sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(COLOR_GREEN)✓ $(APP_NAME) installed successfully$(COLOR_RESET)"
	@echo "Run '$(APP_NAME) --version' to verify installation"

## uninstall: Remove binary from /usr/local/bin (requires sudo)
uninstall:
	@echo "$(COLOR_BLUE)Uninstalling $(APP_NAME) from /usr/local/bin...$(COLOR_RESET)"
	@if [ -f "/usr/local/bin/$(APP_NAME)" ]; then \
		sudo rm /usr/local/bin/$(APP_NAME); \
		echo "$(COLOR_GREEN)✓ $(APP_NAME) uninstalled successfully$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)$(APP_NAME) not found in /usr/local/bin$(COLOR_RESET)"; \
	fi

# Debug target to show all variables
debug-vars:
	@echo "$(COLOR_BLUE)Build Variables:$(COLOR_RESET)"
	@echo "APP_NAME: $(APP_NAME)"
	@echo "VERSION: $(VERSION)"
	@echo "BUILD_DATE: $(BUILD_DATE)"
	@echo "COMMIT_HASH: $(COMMIT_HASH)"
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "BINARY_NAME: $(BINARY_NAME)"
	@echo "LDFLAGS: $(LDFLAGS)"