# Branch Wrangler

## Project Summary

Branch Wrangler is a cross-platform, full-screen terminal application (TUI) built in Go that helps manage local Git branches by reconciling them with GitHub and offers an interactive UI for cleaning up stale branches safely.

The tool provides a comprehensive branch taxonomy system that classifies every branch into one of 16 possible states, such as detached HEAD, no upstream, in sync, ahead/behind, diverged, various PR states, etc.

The application features GitHub API integration with OAuth device flow authentication, rate-limited API calls with caching, and a split-pane TUI interface built with the Bubble Tea framework. It supports both interactive mode for branch management and headless CLI mode for automation and CI/CD integration.

## Installation

Install on macOS or Linux with Homebrew, download the pre-built binary from the release page, or build from source. Windows users can use WSL or build from source.

### Install with Homebrew

Homebrew installation is available for macOS (Intel and Apple Silicon) and Linux (Intel and Arm).

To install Branch Wrangler with Homebrew:

```bash
brew install dfinster/tap/branch-wrangler
```

Or alternatively:

```bash
# Add the tap first
brew tap dfinster/tap

# Then install
brew install branch-wrangler
```

After installation, you can run:

```bash
  branch-wrangler --version
  branch-wrangler --help
```

## Building on macOS

### Prerequisites

- **Go 1.19+**: Install via [Homebrew](https://brew.sh/) or [official installer](https://golang.org/dl/)

```bash
# Using Homebrew (recommended)
brew install go

# Verify installation
go version
```

- **Git**: Pre-installed on macOS or via Homebrew

```bash
brew install git
```

### Build Instructions

1. **Clone the repository**:

   ```bash
   git clone https://github.com/dfinster/branch-wrangler.git
   cd branch-wrangler
   ```

2. **Build using Make** (recommended):
   ```bash
   # Development build with race detection
   make build

   # Optimized release build
   make build-release

   # Build with checksums for distribution
   make dist
   ```

3. **Alternative: Direct Go build**:
   ```bash
   go build -o branch-wrangler ./cmd/branch-wrangler
   ```

4. **Install system-wide** (optional):
   ```bash
   make install  # Installs to /usr/local/bin (requires sudo)
   ```

### Make Targets

| Target               | Description                           |
|----------------------|---------------------------------------|
| `make help`          | Show all available targets            |
| `make build`         | Development build with race detection |
| `make build-release` | Optimized release build               |
| `make test`          | Run all tests                         |
| `make clean`         | Clean build artifacts                 |
| `make dist`          | Build with checksums for distribution |
| `make version`       | Show version information              |

## Current Implementation Status

The codebase provides a solid architectural foundation with proper Go structure and framework choices (Bubble Tea, GitHub API client, OAuth2). Key components are implemented but many features are incomplete:

- ✅ **Architecture & Dependencies**: Complete
- ✅ **Basic TUI Structure**: Functional
- ✅ **Git Operations**: Well implemented
- ✅ **GitHub API Integration**: Functional with caching
- ⚠️ **Authentication**: OAuth Device Flow not implemented
- ⚠️ **Configuration**: YAML loading/saving incomplete
- ⚠️ **CLI Commands**: Flags defined but not processed
- ❌ **Testing**: No tests implemented

See `docs/issues.md` for a detailed analysis of implementation gaps and estimated effort (6-8 weeks) for full compliance with requirements.
