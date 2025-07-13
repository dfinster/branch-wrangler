# Branch Wrangler

This is a work-in-progress, and currently not functional. Please standby.

## Project Summary

Branch Wrangler is a cross-platform, full-screen terminal application (TUI) built in Go that helps manage local Git branches by reconciling them with GitHub and offers an interactive UI for cleaning up stale branches safely.

The tool provides a comprehensive branch taxonomy system that classifies every branch into one of 16 possible states, such as detached HEAD, no upstream, in sync, ahead/behind, diverged, various PR states, etc.

The application features GitHub API integration with OAuth device flow authentication, rate-limited API calls with caching, and a split-pane TUI interface built with the Bubble Tea framework. It supports both interactive mode for branch management and headless CLI mode for automation and CI/CD integration.

## Installation

Branch Wrangler supports macOS and Linux. Multiple installation methods are available.

### Quick Install

**Homebrew (Recommended for macOS and Linux):**
```bash
brew install dfinster/tap/branch-wrangler
```

**Download Binary:**
Visit the [releases page](https://github.com/dfinster/branch-wrangler/releases) and download the appropriate binary for your platform.

### Verification

After installation:
```bash
branch-wrangler --version
branch-wrangler --help
```

### Comprehensive Installation Guide

For detailed installation instructions, including:
- Platform-specific requirements
- Binary verification
- Multiple installation methods
- Troubleshooting

See the **[Installation Guide](docs/user/installation.md)**

### Building from Source

Linux users and developers can build from source. See the **[Building from Source Guide](docs/user/building-from-source.md)** for:
- Complete Linux compilation instructions
- Build system documentation
- Cross-platform building
- Development setup

## Documentation

- **[Installation Guide](docs/user/installation.md)** - Comprehensive installation instructions for all platforms
- **[Building from Source](docs/user/building-from-source.md)** - Complete build instructions and Linux compilation guide
- **[Troubleshooting](docs/user/troubleshooting.md)** - Common issues and solutions
- **[User Guide](docs/user/README.md)** - Usage instructions and features
- **[Admin Documentation](docs/admin/README.md)** - Project status and development information

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

See [Project Status](docs/admin/project-status.md) for a detailed analysis of implementation gaps and estimated effort for full compliance with requirements.
