# Branch Wrangler

## Project Summary

Branch Wrangler is a cross-platform, full-screen terminal application (TUI) built in Go that helps manage local Git branches by reconciling them with GitHub. The tool provides a comprehensive branch taxonomy system that classifies every branch into one of 16 possible states (detached HEAD, no upstream, in sync, ahead/behind, diverged, various PR states, etc.) and offers an interactive UI for cleaning up stale branches safely.

The application features GitHub API integration with OAuth device flow authentication, rate-limited API calls with caching, and a split-pane TUI interface built with the Bubble Tea framework. It supports both interactive mode for branch management and headless CLI mode for automation and CI/CD integration.

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

## Directory Structure

```
├── CLAUDE.md                           # Claude Code instructions and project guidance
├── README.md                           # This file - project documentation
├── branch-wrangler                     # Compiled binary executable
├── go.mod                              # Go module definition and dependencies
├── go.sum                              # Go module checksums for dependency verification
├── cmd/
│   └── branch-wrangler/
│       └── main.go                     # Application entry point and CLI setup
├── docs/
│   ├── requirements.md                 # Complete project requirements specification
│   └── issues.md                       # Implementation gap analysis report
└── internal/
    ├── config/
    │   └── config.go                   # Configuration file handling and YAML support
    ├── git/
    │   ├── classifier.go               # Branch state classification engine
    │   ├── operations.go               # Git operations and command execution
    │   └── types.go                    # Git-related type definitions and constants
    ├── github/
    │   ├── auth.go                     # GitHub authentication and OAuth flow
    │   └── client.go                   # GitHub API client with caching
    └── ui/
        ├── actions.go                  # TUI action handlers (delete, checkout, etc.)
        ├── confirm.go                  # Confirmation dialog components
        ├── filter.go                   # Branch filtering and search functionality
        ├── model.go                    # Bubble Tea model and main TUI logic
        └── views.go                    # TUI view rendering and layout
```

## File Descriptions

### Root Files

- **`CLAUDE.md`** - Project instructions and guidance for Claude Code, containing architecture decisions, development commands, and coding standards for the project.

- **`branch-wrangler`** - Compiled binary executable of the application.

- **`go.mod`** - Go module file defining the project dependencies including Bubble Tea TUI framework, GitHub API client, OAuth2, and Cobra CLI library.

- **`go.sum`** - Checksums for Go module dependencies to ensure reproducible builds.

### Command Line Interface

- **`cmd/branch-wrangler/main.go`** - Application entry point containing the Cobra CLI setup, command-line flag definitions (--version, --list, --json, --delete-stale, etc.), and the main TUI initialization logic that coordinates Git operations, GitHub client setup, and branch classification.

### Documentation

- **`docs/requirements.md`** - Comprehensive requirements specification including the 16-state branch taxonomy, functional requirements (FR-1 through FR-26), non-functional requirements, TUI specifications, and CLI command definitions.

- **`docs/issues.md`** - Implementation gap analysis report identifying 15 major issues between current codebase and requirements, organized by priority with estimated effort for full compliance.

### Configuration Management

- **`internal/config/config.go`** - Configuration file handling with support for YAML format, GitHub token management, saved filter sets, base branches configuration, and XDG config directory support. Currently contains stub implementations for config loading/saving.

### Git Operations

- **`internal/git/types.go`** - Core type definitions including the 16 BranchState constants, Branch struct with metadata (name, state, commits, PR info), GitStatus struct, and display name mappings for all branch states.

- **`internal/git/operations.go`** - Git command execution client providing branch listing, ahead/behind calculation, commit counting, remote existence checking, merge-base analysis, and GitHub repository URL parsing. Implements the core Git operations needed for branch classification.

- **`internal/git/classifier.go`** - Branch state classification engine that maps each branch to exactly one of the 16 states by analyzing Git status, GitHub PR information, and merge relationships. Coordinates between Git operations and GitHub API calls.

### GitHub Integration

- **`internal/github/auth.go`** - GitHub authentication system with support for environment variable tokens, config file tokens, and OAuth Device Flow. Currently has stub implementations for device flow and config file integration.

- **`internal/github/client.go`** - GitHub API client with 15-minute caching, pull request retrieval, branch existence checking, and rate limiting support. Includes both direct client and cached client implementations for performance optimization.

### User Interface

- **`internal/ui/model.go`** - Main Bubble Tea model implementing the TUI state management, keyboard event handling, window resizing, branch loading, and coordination between different UI components. Contains the core application loop and message processing.

- **`internal/ui/views.go`** - TUI view rendering functions for the split-pane layout including branch list view, details pane, header, help overlay, and filter interface. Handles the visual presentation and styling using Lipgloss.

- **`internal/ui/filter.go`** - Branch filtering and search functionality with support for state-based filtering, text search, custom filter sets, and predefined filters. Includes filter application logic and display name generation.

- **`internal/ui/actions.go`** - TUI action handlers for branch operations including checkout, delete (safe and force), PR opening in browser, and bulk operations. Implements the interactive commands triggered by keyboard shortcuts.

- **`internal/ui/confirm.go`** - Confirmation dialog components for dangerous operations like branch deletion, providing safety guards and user confirmation before executing destructive actions.

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
