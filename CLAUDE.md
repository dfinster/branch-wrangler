# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Branch Wrangler is a cross-platform, full-screen terminal application (TUI) built in Go that helps manage local Git branches by reconciling them with GitHub. The tool provides a comprehensive branch taxonomy system and interactive UI for cleaning up stale branches.

## Architecture

This is a greenfield Go project that implements:

- **Branch State Classification**: A comprehensive taxonomy of 16 different branch states (detached HEAD, no upstream, in sync, ahead/behind, diverged, various PR states, etc.)
- **GitHub API Integration**: OAuth device flow authentication, rate-limited API calls with caching
- **TUI Interface**: Built with Bubble Tea framework featuring split-pane layout (branch list + details)
- **Cross-platform Support**: No CGO dependencies, binaries for macOS/Linux/Windows

## Key Components (To Be Implemented)

The codebase structure should follow Go conventions:

- `cmd/` - CLI entry points and command definitions
- `internal/` - Private application code
  - `git/` - Git operations and branch discovery
  - `github/` - GitHub API client and authentication
  - `ui/` - TUI components and layouts
  - `config/` - Configuration file handling
- `pkg/` - Public packages (if any)

## Branch State Taxonomy

The core of the application is a 16-state classification system that maps every possible branch condition:

- Repository states: `DETACHED_HEAD`, `NO_UPSTREAM`, `ORPHAN_REMOTE_DELETED`
- Sync states: `IN_SYNC`, `UNPUSHED_AHEAD`, `BEHIND_REMOTE`, `DIVERGED`
- PR states: `DRAFT_PR`, `OPEN_PR`, `CLOSED_PR`, `MERGED_REMOTE_EXISTS`, `STALE_LOCAL`
- Special states: `FULLY_MERGED_BASE`, `NO_COMMITS`, `UPSTREAM_CHANGED`, `REMOTE_RENAMED`, `UPSTREAM_GONE`

## Configuration

- Location: `$XDG_CONFIG_HOME/branch-wrangler/config.yml` (or `%APPDATA%` on Windows)
- Format: YAML with support for GitHub token, saved filter sets, and UI preferences
- Authentication: Supports `GITHUB_TOKEN` env var, config file token, or OAuth device flow

## TUI Framework Choice

**Bubble Tea** is recommended over tview for this project due to:

- **Complex State Management**: 16 branch states with filtering/sorting require reactive architecture
- **Split Pane Layouts**: Bubble Tea explicitly supports adjustable split-screen layouts
- **Modern Architecture**: Elm-style Model → Update → View pattern ideal for interactive UIs
- **Rich Ecosystem**: Bubbles (components) + Lipgloss (styling) provide comprehensive tools

### Key Dependencies

- `github.com/charmbracelet/bubbletea` - Core TUI framework
- `github.com/charmbracelet/bubbles` - Reusable components (lists, tables, inputs)
- `github.com/charmbracelet/lipgloss` - Styling and layout capabilities

## Development Commands

Since this is a new project, standard Go commands will apply:

- `go build` - Build the application
- `go test ./...` - Run all tests
- `go run main.go` - Run the application
- `go mod tidy` - Clean up dependencies

## Security Requirements

- HTTPS-only GitHub API calls
- Never log authentication tokens
- Set restrictive file permissions on config files (600 on Unix)
- Support for Personal Access Tokens with appropriate scopes (`repo`, `workflow`)

## Performance Targets

- Scan and classify ≤200 branches in <2 seconds
- Maximum 5 concurrent GitHub API calls
- 15-minute default cache TTL for API responses

## Testing Requirements

- 90% line coverage target
- Unit tests with mocked GitHub API
- End-to-end TUI tests using expect harness
- Integration tests for Git operations