# Branch Wrangler ${VERSION}

## Changes

<!-- Automatically generated changelog will be inserted here -->

## Installation

### macOS Apple Silicon

1. Download the `branch-wrangler-${VERSION}-darwin-arm64` binary
2. Verify the checksum against `checksums.txt`
3. Make it executable: `chmod +x branch-wrangler-${VERSION}-darwin-arm64`
4. Move to your PATH: `mv branch-wrangler-${VERSION}-darwin-arm64 /usr/local/bin/branch-wrangler`

### Build from Source

```bash
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
git checkout ${VERSION}
make build-release
```

## What's New

Branch Wrangler is a cross-platform, full-screen terminal application (TUI) that helps manage local Git branches by reconciling them with GitHub. The tool provides a comprehensive branch taxonomy system and interactive UI for cleaning up stale branches.

### Key Features

- **16-State Branch Classification**: Comprehensive taxonomy covering all possible branch states
- **GitHub Integration**: OAuth authentication with PR status and remote branch tracking
- **Interactive TUI**: Split-pane interface built with Bubble Tea framework
- **Safe Operations**: Confirmation dialogs and dry-run modes for destructive actions
- **Cross-Platform**: Native binaries for macOS, Linux, and Windows

## Verification

All binaries are checksummed for security verification.

```bash
# Verify checksum (macOS)
shasum -a 256 -c checksums.txt

# Verify signature (if available)
# gpg --verify branch-wrangler-${VERSION}-darwin-arm64.sig
```

## Support

- **Documentation**: [README.md](https://github.com/dfinster/branch-wrangler/blob/main/README.md)
- **Issues**: [GitHub Issues](https://github.com/dfinster/branch-wrangler/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dfinster/branch-wrangler/discussions)

## Compatibility

- **Go Version**: Built with Go 1.23
- **macOS**: macOS 11+ (Apple Silicon native)
- **Git**: Git 2.20+ required
- **GitHub**: GitHub.com and GitHub Enterprise Server supported

---

**Full Changelog**: https://github.com/dfinster/branch-wrangler/compare/${PREVIOUS_TAG}...${VERSION}
