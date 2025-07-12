# Branch Wrangler - Implementation Issues Report

This document outlines the gaps between the current implementation and the requirements specified in `docs/requirements.md`.

## Critical Issues (Must Fix)

### 1. Authentication System (FR-14 through FR-26)
**Status**: Partially implemented  
**Issue**: OAuth Device Flow and config file token management not implemented  
**Files**: `internal/github/auth.go`, `internal/config/config.go`  

- OAuth Device Flow returns "not implemented" error in `auth.go:48`
- Config file reading/writing not implemented (`config.go:60-64`)
- Token precedence logic incomplete (only checks env var)
- Missing login/logout CLI commands
- No token validation on startup

### 2. Git Operations Core Functionality  
**Status**: Missing implementation  
**Issue**: Git client operations not implemented  
**Files**: `internal/git/operations.go`  

- Missing branch listing implementation
- No ahead/behind calculation logic
- Missing remote existence checks
- No merge-base checking for `FULLY_MERGED_BASE` state
- Missing current branch detection

### 3. GitHub API Integration  
**Status**: Missing implementation  
**Issue**: GitHub client functionality not implemented  
**Files**: `internal/github/client.go`  

- Pull request retrieval not implemented
- Branch existence checking not implemented  
- No API rate limiting or caching
- Missing concurrent request throttling (≤5 concurrent)

### 4. Branch State Classification Logic Gaps
**Status**: Partially implemented  
**Issue**: Several branch states missing detection logic  
**Files**: `internal/git/classifier.go`  

Missing states:
- `UPSTREAM_CHANGED` - Force-push/rebase detection
- `REMOTE_RENAMED` - Remote branch rename detection  
- `UPSTREAM_GONE` - Explicit upstream deletion vs orphan

### 5. TUI Action Implementation
**Status**: Missing implementation  
**Issue**: Core TUI actions not implemented  
**Files**: `internal/ui/actions.go`  

- Branch deletion (safe and force) not implemented
- Checkout functionality missing
- Browser opening for PRs not implemented
- Undo functionality completely missing
- Bulk actions not implemented

## Major Issues (High Priority)

### 6. CLI Argument Handling
**Status**: Flags defined but not processed  
**Issue**: Command-line flags are defined but not handled  
**Files**: `cmd/branch-wrangler/main.go:29-40`  

- `--version`, `--list`, `--json` flags ignored in main logic
- `--delete-stale`, `--dry-run` headless mode not implemented
- `--config`, `--github-token-path` overrides not processed
- Shell completion generation missing

### 7. Configuration File System  
**Status**: Stub implementation  
**Issue**: YAML config loading/saving not implemented  
**Files**: `internal/config/config.go`  

- `Load()` always returns default config
- `Save()` is a no-op
- No YAML parsing/serialization
- Saved filter sets not loaded from config
- File permissions (600) not enforced

### 8. TUI Layout and Keyboard Handling
**Status**: Basic implementation missing features  
**Issue**: Several required TUI features missing  
**Files**: `internal/ui/model.go`, `internal/ui/views.go`  

Missing features:
- Adjustable split-pane layouts
- Help overlay (`?` key) partially implemented
- Color-coded state badges missing proper styling
- Missing keyboard shortcuts (u for undo, F for force, etc.)
- Confirmation dialogs incomplete

### 9. Performance Requirements  
**Status**: Not addressed  
**Issue**: No performance optimizations implemented  

- No concurrent GitHub API calls (requirement: ≤5 concurrent)
- No 15-minute cache TTL implementation
- No performance monitoring for 200 branches in <2s requirement

## Minor Issues (Medium Priority)

### 10. Branch State Display Names
**Status**: Implemented but inconsistent  
**Issue**: Some display names don't match requirements exactly  
**Files**: `internal/git/types.go:27-66`  

Discrepancies from requirements table:
- Requirements use different casing/formatting in some cases

### 11. Error Handling and Logging
**Status**: Basic error handling, no structured logging  
**Issue**: Missing observability features  

- No structured JSON logging (`--log-format json`)
- No verbose logging levels (`--log-level debug`)
- Basic error handling without detailed context

### 12. Cross-Platform Considerations
**Status**: Not addressed  
**Issue**: Windows-specific paths and behaviors not handled  

- Config path uses XDG_CONFIG_HOME but fallback to `%APPDATA%` on Windows not implemented
- File permissions (chmod 600) won't work on Windows

## Testing and Documentation Issues

### 13. Test Coverage
**Status**: No tests implemented  
**Issue**: Requirements specify 90% line coverage target  

- No unit tests
- No integration tests  
- No end-to-end TUI tests with expect harness
- No mocked GitHub API tests

### 14. Documentation
**Status**: Missing  
**Issue**: Requirements specify GoDoc for all exported functions  

- No GoDoc comments on exported functions
- No `make docs` target for running godoc
- README quick-start missing

## Dependency Issues

### 15. Missing Dependencies
**Status**: Core dependencies present, some missing  
**Issue**: Some required functionality needs additional dependencies  

Present:
- ✅ Bubble Tea framework
- ✅ Lipgloss styling  
- ✅ GitHub API client
- ✅ OAuth2 support
- ✅ Cobra CLI framework

Missing:
- YAML parsing library (for config files)
- Potentially missing Bubbles components library

## Implementation Priority Recommendations

### Phase 1 (Core Functionality)
1. Implement Git operations (`internal/git/operations.go`)
2. Implement GitHub API client (`internal/github/client.go`)
3. Complete authentication system including OAuth Device Flow
4. Implement config file loading/saving with YAML support

### Phase 2 (User Interface)
1. Complete branch classification logic for all 16 states
2. Implement TUI actions (delete, checkout, open PR)
3. Complete CLI argument processing
4. Add confirmation dialogs and safety guards

### Phase 3 (Polish and Performance)  
1. Implement caching and performance optimizations
2. Add comprehensive error handling and logging
3. Cross-platform compatibility fixes
4. Comprehensive testing suite

### Phase 4 (Advanced Features)
1. Undo functionality with reflog caching
2. Saved filter sets and advanced filtering
3. Shell completion generation
4. Documentation and help system

## Estimated Effort

- **Phase 1**: ~2-3 weeks (critical path)
- **Phase 2**: ~2 weeks  
- **Phase 3**: ~1-2 weeks
- **Phase 4**: ~1 week

**Total estimated effort**: 6-8 weeks for full compliance with requirements.

The current codebase provides a solid foundation with proper architecture and framework choices, but significant implementation work is needed to meet the specified requirements.