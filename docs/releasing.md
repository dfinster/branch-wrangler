# Release Process Documentation

This document provides comprehensive guidelines for releasing Branch Wrangler, including semantic versioning, automation, security considerations, and rollback procedures.

## Release Overview

Branch Wrangler follows a GitOps-based release process with full automation via GitHub Actions. Releases are triggered by semantic version tags and produce cross-platform binaries with automatic Homebrew formula updates.

## Semantic Versioning Guidelines

Branch Wrangler follows [Semantic Versioning (semver)](https://semver.org/) with the format `MAJOR.MINOR.PATCH`.

### Version Increment Rules

**MAJOR version** (e.g., `v1.0.0` → `v2.0.0`):
- Breaking API changes in CLI interface
- Removal of CLI flags or commands
- Changes that require user action to upgrade
- Major architectural changes affecting compatibility

**MINOR version** (e.g., `v1.0.0` → `v1.1.0`):
- New features and capabilities
- New CLI commands or flags (backward compatible)
- New TUI functionality
- Performance improvements
- New supported platforms

**PATCH version** (e.g., `v1.0.0` → `v1.0.1`):
- Bug fixes and security patches
- Documentation improvements
- Internal refactoring without user-visible changes
- Dependency updates

### Version Decision Matrix

| Change Type | Examples | Version Bump |
|-------------|----------|--------------|
| New CLI command | `--export`, `--import` | MINOR |
| Bug fix | Fix branch classification | PATCH |
| CLI flag removal | Remove deprecated `--old-flag` | MAJOR |
| New TUI feature | Fuzzy search, bulk operations | MINOR |
| Security patch | Fix token exposure | PATCH |
| Breaking config | Change config file format | MAJOR |
| Performance improvement | Faster branch scanning | MINOR |
| Documentation | Update README, add guides | PATCH |

## Release Process

### Prerequisites

1. **Testing Requirements:**
   - All tests passing on main branch
   - Manual testing on macOS and Linux
   - Performance testing with large repositories (200+ branches)
   - Security review completed

2. **Documentation Requirements:**
   - CHANGELOG.md updated with release notes
   - README.md reflects new features
   - Installation documentation current
   - API documentation updated (if applicable)

3. **Branch Requirements:**
   - All changes merged to `main` branch
   - No outstanding critical issues
   - Dependencies up to date

### Step-by-Step Release Process

#### 1. Pre-Release Preparation

```bash
# 1. Ensure you're on main branch and up to date
git checkout main
git pull origin main

# 2. Verify all tests pass
make test
make build-all

# 3. Manual testing checklist
branch-wrangler --version
branch-wrangler --help
branch-wrangler --list
# Test in actual git repository
cd /path/to/test/repo
branch-wrangler
```

#### 2. Version Decision and Tagging

```bash
# 1. Determine next version using guidelines above
CURRENT_VERSION=$(git describe --tags --abbrev=0)
echo "Current version: $CURRENT_VERSION"

# 2. Choose next version (replace with actual version)
NEXT_VERSION="v1.2.3"

# 3. Update version references if needed
# (Most version info is injected at build time)

# 4. Create and push the release tag
git tag -a $NEXT_VERSION -m "Release $NEXT_VERSION"
git push origin $NEXT_VERSION
```

#### 3. Automated Release Process

Once the tag is pushed, GitHub Actions automatically:

1. **Validates the tag format** (`v*.*.*`)
2. **Builds cross-platform binaries:**
   - macOS ARM64 (`darwin-arm64`)
   - macOS Intel (`darwin-amd64`) 
   - Linux x86_64 (`linux-amd64`)
   - Linux ARM64 (`linux-arm64`)
3. **Generates SHA256 checksums**
4. **Verifies binary functionality**
5. **Creates GitHub release with:**
   - Generated release notes from commits
   - All platform binaries
   - Checksums file
   - Installation instructions
6. **Updates Homebrew formula** automatically

#### 4. Post-Release Verification

```bash
# 1. Verify GitHub release created successfully
open "https://github.com/dfinster/branch-wrangler/releases/latest"

# 2. Test Homebrew installation (after formula update)
brew uninstall branch-wrangler
brew install dfinster/tap/branch-wrangler
branch-wrangler --version

# 3. Test binary downloads
curl -L -o test-binary https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-$NEXT_VERSION-darwin-arm64
chmod +x test-binary
./test-binary --version

# 4. Verify checksums
curl -L -o checksums.txt https://github.com/dfinster/branch-wrangler/releases/latest/download/checksums.txt
shasum -a 256 -c checksums.txt --ignore-missing
```

#### 5. Announcement and Communication

1. **Update project documentation:**
   - Update README.md if needed
   - Verify installation instructions work
   - Check all documentation links

2. **Community communication:**
   - GitHub Discussions announcement
   - Update project status documentation

## Release Checklist Template

Use this checklist for each release:

### Pre-Release Checklist
- [ ] All tests passing on main branch
- [ ] Manual testing completed on macOS and Linux
- [ ] Performance testing with large repositories completed
- [ ] Security review completed (if applicable)
- [ ] CHANGELOG.md updated
- [ ] Documentation updated
- [ ] Version number decided using semantic versioning guidelines
- [ ] No critical issues outstanding

### Release Execution Checklist
- [ ] On main branch and up to date
- [ ] Version tag created and pushed
- [ ] GitHub Actions release workflow completed successfully
- [ ] All platform binaries present in release
- [ ] Checksums file generated and valid
- [ ] Homebrew formula updated automatically

### Post-Release Checklist
- [ ] GitHub release verified and accessible
- [ ] Homebrew installation tested
- [ ] Binary downloads tested
- [ ] Checksum verification working
- [ ] Documentation links functional
- [ ] Community announcement posted
- [ ] Project status documentation updated

## Rollback Procedures

### Failed Release Recovery

If a release fails or has critical issues:

#### 1. Immediate Actions
```bash
# 1. Delete the problematic tag
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3

# 2. Delete GitHub release if created
gh release delete v1.2.3 --yes

# 3. Revert Homebrew formula if updated
# (Contact maintainer or create manual PR to revert)
```

#### 2. Critical Issue Hotfix
```bash
# 1. Create hotfix branch from problematic tag
git checkout -b hotfix/v1.2.4 v1.2.3

# 2. Apply minimal fix
# ... make necessary changes ...

# 3. Test thoroughly
make test
make build-all

# 4. Create new patch version
git tag -a v1.2.4 -m "Hotfix release v1.2.4"
git push origin v1.2.4
```

#### 3. Major Issue Recovery
For major issues requiring version rollback:

```bash
# 1. Communicate issue to users immediately
# 2. Identify last known good version
LAST_GOOD_VERSION="v1.2.2"

# 3. Create emergency documentation
echo "⚠️ **Critical Issue in v1.2.3**" > emergency-notice.md
echo "Please use $LAST_GOOD_VERSION until fixed" >> emergency-notice.md

# 4. Plan and execute proper fix
# 5. Increment PATCH version for fix release
```

## Security Considerations

### Release Security Checklist
- [ ] No secrets or tokens in release artifacts
- [ ] All binaries signed with consistent checksums
- [ ] Dependencies scanned for vulnerabilities
- [ ] GitHub Actions secrets properly configured
- [ ] Homebrew tap repository secured

### Security Best Practices
1. **Binary Integrity:**
   - SHA256 checksums for all binaries
   - Consistent build environment (GitHub Actions)
   - Reproducible builds when possible

2. **Access Control:**
   - Release process requires maintainer privileges
   - Homebrew tap updates use dedicated token
   - No manual intervention in automated process

3. **Vulnerability Management:**
   - Regular dependency updates
   - Security scanning in CI/CD
   - Responsible disclosure process

## Homebrew Formula Maintenance

The Homebrew formula is automatically updated by the release workflow.

### Manual Formula Updates
If manual updates are needed:

```ruby
# Formula location: homebrew-tap/Formula/branch-wrangler.rb
class BranchWrangler < Formula
  desc "Cross-platform TUI for managing local Git branches with GitHub integration"
  homepage "https://github.com/dfinster/branch-wrangler"
  version "1.2.3"
  license "MIT"
  
  # Platform-specific URLs and checksums...
end
```

### Formula Testing
```bash
# Test formula before release
brew install --build-from-source dfinster/tap/branch-wrangler
brew test branch-wrangler
brew audit --strict branch-wrangler
```

## Monitoring and Verification

### Release Health Monitoring
- GitHub release download metrics
- Homebrew installation analytics
- Issue reports post-release
- Performance regression monitoring

### Verification Procedures
1. **Automated verification:**
   - Binary functionality tests in CI
   - Checksum validation
   - Version string verification

2. **Manual verification:**
   - Cross-platform installation testing
   - User workflow testing
   - Documentation accuracy

## Troubleshooting Common Issues

### Release Workflow Failures

**Issue**: Build fails on specific platform
```bash
# Solution: Check build logs and platform-specific requirements
# Often caused by dependency issues or platform-specific code
```

**Issue**: Homebrew formula update fails
```bash
# Solution: Check HOMEBREW_TAP_TOKEN permissions
# Verify tap repository access and formula syntax
```

**Issue**: Checksums don't match
```bash
# Solution: Rebuild binaries or regenerate checksums
make clean
make build-all
make checksums
```

### Version Tag Issues

**Issue**: Tag already exists
```bash
# Delete local and remote tag, then recreate
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Release Announcement Templates

### GitHub Release Notes Template
```markdown
## Changes in v1.2.3

### New Features
- Feature 1 description
- Feature 2 description

### Bug Fixes  
- Fix 1 description
- Fix 2 description

### Performance Improvements
- Improvement 1 description

## Installation

See the [Installation Guide](docs/installation.md) for detailed instructions.

### Quick Install
```bash
# Homebrew
brew install dfinster/tap/branch-wrangler

# Binary download
curl -L -o branch-wrangler https://github.com/dfinster/branch-wrangler/releases/download/v1.2.3/branch-wrangler-v1.2.3-darwin-arm64
chmod +x branch-wrangler
```

## Verification
All binaries include SHA256 checksums for verification.
```

### Community Announcement Template
```markdown
🎉 **Branch Wrangler v1.2.3 Released!**

We're excited to announce the release of Branch Wrangler v1.2.3 with [key features/fixes].

**Installation:**
- Homebrew: `brew install dfinster/tap/branch-wrangler`
- Binary: Download from [releases page]
- Source: See [building guide]

**What's New:**
- [Major feature 1]
- [Major feature 2]
- [Important fixes]

**Getting Help:**
- 📖 [Documentation](docs/README.md)
- 🐛 [Issues](https://github.com/dfinster/branch-wrangler/issues)
- 💬 [Discussions](https://github.com/dfinster/branch-wrangler/discussions)

Thank you to all contributors and users! 🙏
```

## Contact and Escalation

For release-related issues:
1. **Technical issues**: Open GitHub issue with `release` label
2. **Security concerns**: Follow responsible disclosure process
3. **Critical failures**: Contact maintainers directly

## Related Documentation

- [Contributing Guidelines](CONTRIBUTING.md)
- [Installation Guide](installation.md)
- [Building from Source](building-from-source.md)
- [Troubleshooting](troubleshooting.md)