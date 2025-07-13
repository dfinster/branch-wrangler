---
name: Release Checklist
about: Template for planning and tracking releases
title: 'Release v[VERSION] Checklist'
labels: ['release', 'process']
assignees: ['']

---

## Release Information

**Version**: v[MAJOR.MINOR.PATCH]  
**Release Type**: [ ] Major [ ] Minor [ ] Patch  
**Target Date**: [DATE]  
**Release Manager**: @[USERNAME]

## Version Decision

**Changes since last release:**
- 
- 
- 

**Semantic versioning decision:**
- [ ] **MAJOR**: Breaking changes requiring user action
- [ ] **MINOR**: New features, backward compatible
- [ ] **PATCH**: Bug fixes, security patches

## Pre-Release Checklist

### Development Readiness
- [ ] All planned features/fixes merged to main
- [ ] No critical issues outstanding
- [ ] Dependencies up to date
- [ ] Security review completed (if applicable)

### Testing Requirements
- [ ] All automated tests passing
- [ ] Manual testing completed on macOS
- [ ] Manual testing completed on Linux
- [ ] Performance testing with large repositories (200+ branches)
- [ ] Cross-platform binary verification

### Documentation Updates
- [ ] CHANGELOG.md updated with release notes
- [ ] README.md reflects new features/changes
- [ ] Installation documentation current
- [ ] API documentation updated (if applicable)
- [ ] Breaking changes documented (if applicable)

### Version Preparation
- [ ] Version number decided using semantic versioning
- [ ] Release notes drafted
- [ ] Announcement content prepared

## Release Execution

### Tag Creation and Push
- [ ] On main branch and up to date
- [ ] Release tag created: `git tag -a v[VERSION] -m "Release v[VERSION]"`
- [ ] Tag pushed: `git push origin v[VERSION]`
- [ ] GitHub Actions release workflow triggered

### Automated Release Verification
- [ ] GitHub Actions workflow completed successfully
- [ ] All platform binaries built:
  - [ ] macOS ARM64 (`darwin-arm64`)
  - [ ] macOS Intel (`darwin-amd64`)
  - [ ] Linux x86_64 (`linux-amd64`)
  - [ ] Linux ARM64 (`linux-arm64`)
- [ ] SHA256 checksums generated
- [ ] GitHub release created with artifacts
- [ ] Homebrew formula updated automatically

## Post-Release Verification

### Installation Testing
- [ ] GitHub release accessible and complete
- [ ] Homebrew installation working:
  ```bash
  brew uninstall branch-wrangler
  brew install dfinster/tap/branch-wrangler
  branch-wrangler --version
  ```
- [ ] Binary downloads functional:
  ```bash
  curl -L -o test-binary https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-v[VERSION]-darwin-arm64
  chmod +x test-binary
  ./test-binary --version
  ```
- [ ] Checksum verification working:
  ```bash
  curl -L -o checksums.txt https://github.com/dfinster/branch-wrangler/releases/latest/download/checksums.txt
  shasum -a 256 -c checksums.txt --ignore-missing
  ```

### Functionality Testing
- [ ] Version command shows correct version
- [ ] Help command functional
- [ ] Basic TUI functionality working
- [ ] CLI headless commands working (if implemented)
- [ ] GitHub integration functional

### Documentation Verification
- [ ] All documentation links functional
- [ ] Installation instructions work as documented
- [ ] New features documented properly
- [ ] Examples and guides current

## Communication and Announcement

### Community Communication
- [ ] GitHub Discussions announcement posted
- [ ] Release highlights communicated
- [ ] Breaking changes prominently documented (if applicable)

### Documentation Updates
- [ ] Project status documentation updated
- [ ] README badges updated (if applicable)
- [ ] Installation guides verified

## Release Health Monitoring

### Metrics to Track
- [ ] Download statistics from GitHub releases
- [ ] Homebrew installation metrics (if available)
- [ ] Issue reports post-release
- [ ] User feedback and discussions

### Follow-up Actions
- [ ] Monitor for critical issues in first 24-48 hours
- [ ] Respond to user questions and issues
- [ ] Plan hotfix if critical issues discovered

## Rollback Plan

**In case of critical issues:**

### Immediate Actions
- [ ] Delete problematic tag: `git tag -d v[VERSION] && git push origin :refs/tags/v[VERSION]`
- [ ] Delete GitHub release: `gh release delete v[VERSION] --yes`
- [ ] Communicate issue to users
- [ ] Identify last known good version

### Recovery Process
- [ ] Create hotfix branch if possible
- [ ] Plan proper fix and new patch version
- [ ] Document lessons learned

## Release Notes

**Key Changes:**
- 

**Breaking Changes:**
- None / List breaking changes

**Security Updates:**
- None / List security fixes

**Performance Improvements:**
- 

**Bug Fixes:**
- 

**New Features:**
- 

## Sign-off

### Release Manager
- [ ] All checklist items completed
- [ ] Release quality verified
- [ ] Ready for production use

**Release Manager Signature**: @[USERNAME]  
**Release Date**: [ACTUAL DATE]  
**Release URL**: https://github.com/dfinster/branch-wrangler/releases/tag/v[VERSION]

---

## Additional Notes

<!-- Any additional context, concerns, or notes for this release -->

## Related Issues

<!-- Link to related issues, features, or fixes included in this release -->
- Closes #
- Fixes #
- Implements #