# Contributing to Branch Wrangler

Thank you for your interest in contributing to Branch Wrangler! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct:

- **Be respectful** and inclusive in all interactions
- **Be constructive** when providing feedback
- **Be patient** with new contributors and maintainers
- **Be collaborative** and help others learn
- **Report issues** responsibly and professionally

## Getting Started

### Prerequisites

- **Go 1.24.2+** ([installation guide](https://golang.org/doc/install))
- **Git** for version control
- **Make** for build automation (optional but recommended)
- **GitHub account** for contributing

### First-Time Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/branch-wrangler.git
   cd branch-wrangler
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/dfinster/branch-wrangler.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Verify setup**
   ```bash
   make build
   make test
   ./build/branch-wrangler --version
   ```

## Development Setup

### Build System

Branch Wrangler uses a Makefile for common development tasks:

```bash
# Development build with race detection
make build

# Release build (optimized)
make build-release

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Show all available targets
make help
```

### Project Structure

```
branch-wrangler/
├── cmd/
│   └── branch-wrangler/     # Main application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── git/                 # Git operations and branch classification
│   ├── github/              # GitHub API client and authentication
│   ├── ui/                  # TUI components and interactions
│   └── version/             # Version information
├── docs/                    # Documentation
├── .github/                 # GitHub Actions workflows
└── Makefile                 # Build automation
```

### Coding Standards

**Go Style:**
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting (run `go fmt ./...`)
- Use `go vet` for static analysis
- Follow Go naming conventions

**Code Organization:**
- Keep functions focused and single-purpose
- Use descriptive variable and function names
- Add comments for exported functions and complex logic
- Group related functionality in appropriate packages

**Error Handling:**
- Always handle errors explicitly
- Provide meaningful error messages with context
- Use error wrapping (`fmt.Errorf("context: %w", err)`)

## Making Changes

### Branch Strategy

1. **Create a feature branch**
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/description
   ```

2. **Keep branches focused**
   - One feature or fix per branch
   - Small, atomic commits with clear messages
   - Rebase on main before submitting PR

3. **Commit message format**
   ```
   type: brief description
   
   Optional longer explanation of the change, including
   why it was made and what problem it solves.
   
   Fixes #issue-number (if applicable)
   ```

   **Types:** `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

### Development Workflow

1. **Make your changes**
   - Write code following project conventions
   - Add tests for new functionality
   - Update documentation as needed

2. **Test thoroughly**
   ```bash
   # Run all tests
   make test
   
   # Test build across platforms
   make build-all
   
   # Manual testing
   ./build/branch-wrangler --help
   cd /path/to/test/git/repo
   ./path/to/branch-wrangler
   ```

3. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: add fuzzy search for branch names"
   git push origin feature/fuzzy-search
   ```

## Testing

### Testing Requirements

- **Unit tests** for all new functionality
- **Integration tests** for GitHub API interactions
- **Manual testing** for TUI components
- **Cross-platform testing** when possible

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/git/

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

### Writing Tests

**Unit Test Example:**
```go
func TestBranchClassification(t *testing.T) {
    // Use testify/assert for assertions
    assert := assert.New(t)
    
    // Test data setup
    branch := &Branch{
        Name: "feature/test",
        Ahead: 1,
        Behind: 0,
    }
    
    // Test classification
    classifier := NewClassifier(mockGitClient, mockGitHubClient, []string{"main"})
    err := classifier.ClassifyBranch(context.Background(), branch)
    
    assert.NoError(err)
    assert.Equal(UnpushedAhead, branch.State)
}
```

**Mock Example:**
```go
type mockGitHubClient struct{}

func (m *mockGitHubClient) GetPullRequestsForBranch(ctx context.Context, branch string) ([]PullRequest, error) {
    // Return test data
    return []PullRequest{}, nil
}
```

## Pull Request Process

### Before Submitting

1. **Ensure your branch is up to date**
   ```bash
   git checkout main
   git pull upstream main
   git checkout feature/your-branch
   git rebase main
   ```

2. **Run the full test suite**
   ```bash
   make test
   make build-all
   ```

3. **Check code quality**
   ```bash
   go fmt ./...
   go vet ./...
   ```

### Pull Request Guidelines

**Title Format:**
- Use descriptive, concise titles
- Start with type: `feat:`, `fix:`, `docs:`, etc.
- Example: `feat: implement OAuth device flow authentication`

**Description Requirements:**
- **What**: Clear description of changes made
- **Why**: Explanation of motivation and context
- **How**: Brief overview of implementation approach
- **Testing**: Description of testing performed
- **Screenshots**: For UI changes (if applicable)

**PR Checklist:**
- [ ] Tests pass locally
- [ ] Code follows project style guidelines
- [ ] Documentation updated (if applicable)
- [ ] No breaking changes (or clearly marked)
- [ ] Linked to relevant issue (if applicable)

### Review Process

1. **Automated checks** must pass:
   - All tests passing
   - Build successful on all platforms
   - Code quality checks passed

2. **Maintainer review**:
   - Code quality and style
   - Functionality and correctness
   - Test coverage and quality
   - Documentation completeness

3. **Approval and merge**:
   - At least one maintainer approval required
   - All feedback addressed
   - Squash merge for clean history

## Issue Guidelines

### Reporting Bugs

Use the bug report template and include:

- **Environment**: OS, Go version, Branch Wrangler version
- **Steps to reproduce**: Clear, minimal reproduction steps
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Logs/screenshots**: Error messages or visual evidence

### Feature Requests

Use the feature request template and include:

- **Problem statement**: What problem does this solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other approaches evaluated
- **Additional context**: Use cases, examples, references

### Issue Labels

- `bug`: Something isn't working correctly
- `enhancement`: New feature or improvement
- `documentation`: Documentation improvements
- `good first issue`: Suitable for new contributors
- `help wanted`: Extra attention needed
- `priority/high`: High priority issues
- `type/question`: Questions about usage

## Release Process

### For Contributors

- **Feature development**: Target next minor version
- **Bug fixes**: May be included in patch releases
- **Breaking changes**: Require major version bump and RFC process

### Release Contributions

Contributors can help with releases by:
- Testing release candidates
- Updating documentation
- Reporting issues in pre-release versions
- Contributing to changelog and release notes

See [Release Process Documentation](releasing.md) for detailed release procedures.

## Getting Help

### Documentation

- **Installation**: [Installation Guide](installation.md)
- **Building**: [Building from Source](building-from-source.md)
- **Troubleshooting**: [Troubleshooting Guide](troubleshooting.md)
- **Releases**: [Release Process](releasing.md)

### Community Support

- **GitHub Discussions**: General questions and community chat
- **GitHub Issues**: Bug reports and feature requests  
- **Code review**: Ask for feedback on your implementation approach

### Development Questions

For development-specific questions:

1. **Check existing issues** and discussions first
2. **Search documentation** for relevant information
3. **Ask in GitHub Discussions** with `development` tag
4. **Join the conversation** on existing relevant issues

## Recognition

Contributors are recognized in several ways:

- **Contributor list**: Listed in README.md
- **Release notes**: Mentioned in changelog for significant contributions
- **GitHub insights**: Contribution activity tracked automatically

## Development Tips

### Debugging TUI Applications

```bash
# Use log files for debugging TUI
branch-wrangler 2> debug.log

# Test TUI components in isolation
go test ./internal/ui/ -v
```

### Working with Git Operations

```bash
# Create test repositories for development
mkdir test-repo && cd test-repo
git init
git remote add origin https://github.com/user/repo.git
# Create test branches and scenarios
```

### GitHub API Development

```bash
# Set up test token for development
export GITHUB_TOKEN=your_test_token

# Use GitHub CLI for testing API interactions
gh api /rate_limit
```

## Common Development Tasks

### Adding a New CLI Command

1. Add flag definition in `cmd/branch-wrangler/main.go`
2. Implement handler function
3. Add tests for new functionality
4. Update help documentation
5. Add integration test

### Adding a New Branch State

1. Define state constant in `internal/git/types.go`
2. Add display name in `DisplayName()` method
3. Implement detection logic in classifier
4. Add tests for new state
5. Update documentation

### Adding a New TUI Feature

1. Implement feature in appropriate UI component
2. Add keyboard shortcuts and help text
3. Add tests for testable components
4. Update help screen
5. Test manually in TUI

Thank you for contributing to Branch Wrangler! Your efforts help make Git branch management better for everyone. 🚀