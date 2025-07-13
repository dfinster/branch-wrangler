# Installation Guide

Branch Wrangler is a cross-platform terminal application for managing Git branches. This guide covers installation methods for all supported platforms.

## Supported Platforms

- **macOS**: Intel and Apple Silicon (via Homebrew or binary download)
- **Linux**: Intel x86_64 and ARM64 (via Homebrew, binary download, or source compilation)
- **Windows**: Not officially supported (use WSL or build from source)

## Quick Installation

### macOS and Linux with Homebrew (Recommended)

```bash
brew install dfinster/tap/branch-wrangler
```

### Download Pre-built Binary

1. Visit the [releases page](https://github.com/dfinster/branch-wrangler/releases)
2. Download the appropriate binary for your platform
3. Extract and move to your PATH

## Detailed Installation Methods

### Method 1: Homebrew Installation

Homebrew is the recommended installation method for macOS and Linux users.

#### Prerequisites
- [Homebrew](https://brew.sh/) installed on your system

#### Installation Steps

1. **Add the tap** (first time only):
   ```bash
   brew tap dfinster/tap
   ```

2. **Install Branch Wrangler**:
   ```bash
   brew install branch-wrangler
   ```

3. **Verify installation**:
   ```bash
   branch-wrangler --version
   branch-wrangler --help
   ```

#### Updating
```bash
brew update
brew upgrade branch-wrangler
```

#### Uninstalling
```bash
brew uninstall branch-wrangler
brew untap dfinster/tap  # Optional: remove the tap
```

### Method 2: Binary Download

Download pre-compiled binaries from GitHub releases.

#### macOS Binary Installation

1. **Download the binary**:
   ```bash
   # Apple Silicon (M1/M2/M3/M4)
   curl -L -o branch-wrangler.tar.gz https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-darwin-arm64.tar.gz

   # Intel
   curl -L -o branch-wrangler.tar.gz https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-darwin-amd64.tar.gz
   ```

2. **Extract and install**:
   ```bash
   tar -xzf branch-wrangler.tar.gz
   sudo mv branch-wrangler /usr/local/bin/
   sudo chmod +x /usr/local/bin/branch-wrangler
   ```

3. **Verify installation**:
   ```bash
   branch-wrangler --version
   ```

#### Linux Binary Installation

1. **Download the binary**:
   ```bash
   # x86_64
   curl -L -o branch-wrangler.tar.gz https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-linux-amd64.tar.gz

   # ARM64
   curl -L -o branch-wrangler.tar.gz https://github.com/dfinster/branch-wrangler/releases/latest/download/branch-wrangler-linux-arm64.tar.gz
   ```

2. **Extract and install**:
   ```bash
   tar -xzf branch-wrangler.tar.gz
   sudo mv branch-wrangler /usr/local/bin/
   sudo chmod +x /usr/local/bin/branch-wrangler
   ```

3. **Verify installation**:
   ```bash
   branch-wrangler --version
   ```

### Method 3: Build from Source

For Linux users or those who prefer to compile from source, see the [Building from Source Guide](building-from-source.md).

## Binary Verification

All releases include SHA256 checksums for verification.

1. **Download the checksum file**:
   ```bash
   curl -L -o checksums.txt https://github.com/dfinster/branch-wrangler/releases/latest/download/checksums.txt
   ```

2. **Verify the binary**:
   ```bash
   # macOS
   shasum -a 256 -c checksums.txt --ignore-missing

   # Linux
   sha256sum -c checksums.txt --ignore-missing
   ```

## Platform-Specific Requirements

### macOS
- **OS Version**: macOS 10.15 (Catalina) or later
- **Architecture**: Intel x86_64 or Apple Silicon (ARM64)
- **Dependencies**: None (statically linked)

### Linux
- **Kernel**: Linux 3.2 or later
- **Architecture**: x86_64 or ARM64
- **Dependencies**: None (statically linked)
- **libc**: Compatible with both glibc and musl

### Supported Linux Distributions
- Ubuntu 18.04 LTS and later
- Debian 10 and later
- CentOS/RHEL 8 and later
- Fedora 32 and later
- Alpine Linux 3.12 and later
- Arch Linux (rolling)

## Quick Start

After installation:

1. **Navigate to a Git repository**:
   ```bash
   cd /path/to/your/git/repo
   ```

2. **Run Branch Wrangler**:
   ```bash
   branch-wrangler
   ```

3. **First-time setup** (if using GitHub features):
   - The app will guide you through GitHub authentication
   - Follow the OAuth device flow prompts

## Configuration

Branch Wrangler creates a configuration file at:
- **macOS**: `~/Library/Application Support/branch-wrangler/config.yml`
- **Linux**: `~/.config/branch-wrangler/config.yml`

The configuration file is created automatically on first run.

## Environment Variables

- `GITHUB_TOKEN`: Personal Access Token for GitHub API (optional)
- `XDG_CONFIG_HOME`: Override default config directory on Linux

## Next Steps

- Read the [User Guide](user/README.md) for detailed usage instructions
- Check the [Troubleshooting Guide](troubleshooting.md) if you encounter issues
- See [Building from Source](building-from-source.md) for compilation instructions

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/dfinster/branch-wrangler/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dfinster/branch-wrangler/discussions)
