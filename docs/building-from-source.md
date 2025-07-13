# Building from Source

This guide provides comprehensive instructions for building Branch Wrangler from source code on Linux, macOS, and other Unix-like systems.

## Prerequisites

### Required Dependencies

- **Go 1.24.2 or later**: The project requires Go 1.24.2+ (as specified in go.mod)
- **Git**: For cloning the repository and version information
- **Make**: For using the build system (optional but recommended)

### Installing Go

#### Linux

**Ubuntu/Debian:**
```bash
# Install from official packages (may be older version)
sudo apt update
sudo apt install golang-go

# Or install latest version manually
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**CentOS/RHEL/Fedora:**
```bash
# RHEL/CentOS 8+
sudo dnf install golang

# Or install latest version manually
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Arch Linux:**
```bash
sudo pacman -S go
```

**Alpine Linux:**
```bash
sudo apk add go
```

#### macOS

```bash
# Using Homebrew (recommended)
brew install go

# Or download from official site
# Visit https://golang.org/dl/ and download the macOS installer
```

### Verify Go Installation

```bash
go version
# Should output: go version go1.24.2 linux/amd64 (or similar)
```

## Quick Build

For users who want to build quickly:

```bash
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release
./dist/branch-wrangler-*
```

## Detailed Build Instructions

### Step 1: Clone the Repository

```bash
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
```

### Step 2: Verify Dependencies

```bash
make check-deps
```

This will verify that Go and Git are properly installed.

### Step 3: Build Options

#### Development Build (with race detection)

```bash
make build
# or
make build-dev
```

The development build includes:
- Race condition detection (`-race` flag)
- Debug symbols
- Output: `build/branch-wrangler`

#### Release Build (optimized)

```bash
make build-release
```

The release build includes:
- Optimized for size and performance
- Trimmed paths for reproducible builds
- Version information embedded
- Output: `dist/branch-wrangler-<version>-<os>-<arch>`

#### Cross-platform Builds

Build for all supported platforms:

```bash
make build-all
```

Build for specific platforms:

```bash
make build-linux-amd64    # Linux x86_64
make build-linux-arm64    # Linux ARM64
make build-darwin-amd64   # macOS Intel
make build-darwin-arm64   # macOS Apple Silicon
```

### Step 4: Generate Distribution Package

```bash
make dist
```

This creates:
- Optimized binary
- SHA256 checksums
- All files in the `dist/` directory

For all platforms:

```bash
make dist-all
```

## Build System Details

### Makefile Targets

| Target               | Description                                  |
|----------------------|----------------------------------------------|
| `make help`          | Show all available targets                   |
| `make build`         | Development build with race detection        |
| `make build-release` | Optimized release build                      |
| `make build-all`     | Build for all platforms                      |
| `make dist`          | Build with checksums (current platform)      |
| `make dist-all`      | Build for all platforms with checksums       |
| `make test`          | Run all tests                                |
| `make clean`         | Clean build artifacts                        |
| `make install`       | Install to `/usr/local/bin` (requires sudo)  |
| `make uninstall`     | Remove from `/usr/local/bin` (requires sudo) |
| `make version`       | Show version information                     |
| `make checksums`     | Generate SHA256 checksums                    |

### Manual Build (without Make)

If you prefer not to use Make:

```bash
# Development build
go build -race -o branch-wrangler ./cmd/branch-wrangler

# Release build with version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT_HASH=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

go build -trimpath \
  -ldflags "-X github.com/dfinster/branch-wrangler/internal/version.Version=$VERSION \
           -X github.com/dfinster/branch-wrangler/internal/version.BuildDate=$BUILD_DATE \
           -X github.com/dfinster/branch-wrangler/internal/version.CommitHash=$COMMIT_HASH" \
  -o branch-wrangler ./cmd/branch-wrangler
```

### Cross-compilation

Go's cross-compilation makes it easy to build for different platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o branch-wrangler-linux-amd64 ./cmd/branch-wrangler

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o branch-wrangler-linux-arm64 ./cmd/branch-wrangler

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o branch-wrangler-darwin-amd64 ./cmd/branch-wrangler

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o branch-wrangler-darwin-arm64 ./cmd/branch-wrangler
```

## Linux Distribution-Specific Instructions

### Ubuntu/Debian

```bash
# Install dependencies
sudo apt update
sudo apt install golang-go git make

# Clone and build
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release

# Install system-wide (optional)
sudo make install
```

### CentOS/RHEL 8+

```bash
# Install dependencies
sudo dnf install golang git make

# Clone and build
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release

# Install system-wide (optional)
sudo make install
```

### Arch Linux

```bash
# Install dependencies
sudo pacman -S go git make

# Clone and build
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release

# Install system-wide (optional)
sudo make install
```

### Alpine Linux

```bash
# Install dependencies
sudo apk add go git make

# Clone and build
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release

# Install system-wide (optional)
sudo make install
```

## Build Configuration

### Environment Variables

- `VERSION`: Override version (defaults to git describe)
- `GOOS`: Target operating system
- `GOARCH`: Target architecture
- `CGO_ENABLED`: Disable CGO (defaults to 0 for static builds)

### Go Build Flags

The build system uses these Go build flags:

**Development builds:**
- `-race`: Enable race detector
- `-ldflags`: Embed version information

**Release builds:**
- `-trimpath`: Remove absolute paths for reproducible builds
- `-ldflags`: Embed version information with optimizations

## Testing the Build

After building, verify the binary works:

```bash
# Check version
./branch-wrangler --version

# Run help
./branch-wrangler --help

# Test in a git repository
cd /path/to/git/repo
/path/to/branch-wrangler
```

## Installing System-wide

### Using Make (recommended)

```bash
sudo make install
```

This installs to `/usr/local/bin/branch-wrangler`.

### Manual Installation

```bash
# Copy binary
sudo cp dist/branch-wrangler-* /usr/local/bin/branch-wrangler
sudo chmod +x /usr/local/bin/branch-wrangler

# Verify installation
branch-wrangler --version
```

### Creating a Package

For distribution-specific packages, you can:

1. **DEB package** (Ubuntu/Debian):
   ```bash
   # Install fpm
   gem install fpm

   # Create package
   fpm -s dir -t deb -n branch-wrangler -v $(git describe --tags) \
       --description "Git branch management TUI" \
       dist/branch-wrangler-linux-amd64=/usr/local/bin/branch-wrangler
   ```

2. **RPM package** (CentOS/RHEL/Fedora):
   ```bash
   # Install fpm
   gem install fpm

   # Create package
   fpm -s dir -t rpm -n branch-wrangler -v $(git describe --tags) \
       --description "Git branch management TUI" \
       dist/branch-wrangler-linux-amd64=/usr/local/bin/branch-wrangler
   ```

## Troubleshooting Build Issues

### Common Issues

1. **Go version too old**:
   ```
   go.mod requires Go 1.24.2, but you're using Go 1.21
   ```
   Solution: Update Go to version 1.24.2 or later.

2. **Missing dependencies**:
   ```bash
   go mod download
   go mod tidy
   ```

3. **Permission denied during install**:
   ```bash
   sudo make install
   ```

4. **Build fails with "command not found"**:
   Ensure Go is in your PATH:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   ```

### Getting Help

- Check the [Troubleshooting Guide](troubleshooting.md)
- Open an issue on [GitHub](https://github.com/dfinster/branch-wrangler/issues)
- Verify your environment with `make check-deps`

## Contributing

When building for development:

1. Use `make build` for development builds
2. Run `make test` to execute tests
3. Use `make fmt` to format code
4. Run `make lint` if you have golangci-lint installed

See the project's contributing guidelines for more information.
