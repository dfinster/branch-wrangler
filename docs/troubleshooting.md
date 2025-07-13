# Troubleshooting Guide

This guide covers common issues you might encounter when installing, building, or using Branch Wrangler.

## Installation Issues

### Homebrew Installation Problems

#### Error: "Formula not found"
```
Error: No available formula with name "branch-wrangler"
```

**Solution:**
```bash
# Make sure the tap is added first
brew tap dfinster/tap
brew install branch-wrangler
```

#### Error: "Permission denied"
```
Error: Permission denied @ apply2files
```

**Solution:**
```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*
```

#### Error: Outdated formula
```
Warning: branch-wrangler X.X.X is available and more recent than version Y.Y.Y
```

**Solution:**
```bash
brew update
brew upgrade branch-wrangler
```

### Binary Download Issues

#### Error: "Permission denied" when running binary
```bash
chmod +x branch-wrangler
sudo mv branch-wrangler /usr/local/bin/
```

#### Error: "quarantine" on macOS
```
"branch-wrangler" cannot be opened because the developer cannot be verified
```

**Solution:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine branch-wrangler

# Or allow in System Preferences > Security & Privacy
```

#### Error: Binary not found in PATH
```bash
# Add /usr/local/bin to PATH if not present
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Or use full path
/usr/local/bin/branch-wrangler
```

## Build Issues

### Go Version Problems

#### Error: "go.mod requires Go 1.24.2"
```
go.mod requires Go 1.24.2, but you are using Go 1.21.0
```

**Solution:**
```bash
# Update Go to latest version
# On macOS with Homebrew:
brew install go

# On Linux, download from https://golang.org/dl/
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
```

#### Error: "go: command not found"
```bash
# Add Go to PATH
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Dependency Issues

#### Error: "module not found"
```
go: module github.com/charmbracelet/bubbletea: not found
```

**Solution:**
```bash
# Download dependencies
go mod download
go mod tidy

# If still failing, clean module cache
go clean -modcache
go mod download
```

#### Error: "checksum mismatch"
```
verifying module: checksum mismatch
```

**Solution:**
```bash
# Clear module cache and retry
go clean -modcache
rm go.sum
go mod tidy
```

### Build Failures

#### Error: "race detector not supported"
```
-race is only supported on linux/amd64, linux/ppc64le, linux/arm64, freebsd/amd64, netbsd/amd64, darwin/amd64, darwin/arm64, and windows/amd64
```

**Solution:**
```bash
# Use release build instead
make build-release

# Or build without race detector
go build ./cmd/branch-wrangler
```

#### Error: "permission denied" during install
```
cp: cannot create regular file '/usr/local/bin/branch-wrangler': Permission denied
```

**Solution:**
```bash
# Use sudo for system installation
sudo make install

# Or install to user directory
mkdir -p ~/bin
cp dist/branch-wrangler-* ~/bin/branch-wrangler
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
```

## Runtime Issues

### Git Repository Problems

#### Error: "not a git repository"
```
fatal: not a git repository (or any of the parent directories): .git
```

**Solution:**
- Run Branch Wrangler from within a Git repository
- Initialize a Git repository: `git init`
- Clone an existing repository: `git clone <url>`

#### Error: "no branches found"
```
No local branches found in repository
```

**Solution:**
```bash
# Create a branch if repository is empty
git checkout -b main
git commit --allow-empty -m "Initial commit"

# Or check if branches exist
git branch -a
```

#### Error: "detached HEAD state"
```
You are in 'detached HEAD' state
```

**This is normal** - Branch Wrangler classifies this as `DETACHED_HEAD` state. You can:
- Switch to a branch: `git checkout main`
- Create a new branch: `git checkout -b new-branch`

### GitHub API Issues

#### Error: "API rate limit exceeded"
```
GitHub API rate limit exceeded. Reset at: 2024-01-01T00:00:00Z
```

**Solution:**
```bash
# Set GitHub token to increase rate limit
export GITHUB_TOKEN=your_token_here

# Or configure in config file
~/.config/branch-wrangler/config.yml
```

#### Error: "authentication failed"
```
GitHub authentication failed: 401 Unauthorized
```

**Solutions:**
1. **Check token permissions:**
   - Token needs `repo` scope for private repositories
   - Token needs `workflow` scope for Actions integration

2. **Re-authenticate:**
   ```bash
   # Remove saved credentials
   rm ~/.config/branch-wrangler/config.yml
   
   # Restart Branch Wrangler to trigger OAuth flow
   branch-wrangler
   ```

3. **Use environment variable:**
   ```bash
   export GITHUB_TOKEN=ghp_your_token_here
   branch-wrangler
   ```

#### Error: "repository not found"
```
GitHub repository not found: 404 Not Found
```

**Solution:**
- Verify the repository exists and is accessible
- Check if it's a private repository and you have access
- Ensure the remote URL is correct: `git remote -v`

### TUI Interface Issues

#### Error: Terminal too small
```
Terminal size too small. Minimum: 80x24
```

**Solution:**
- Resize your terminal window
- Use full-screen terminal application
- Check terminal size: `echo $COLUMNS x $LINES`

#### Error: Colors not displaying correctly
```
TUI appears in black and white or with wrong colors
```

**Solutions:**
```bash
# Check terminal color support
echo $TERM

# Set color terminal
export TERM=xterm-256color

# For tmux users
export TERM=screen-256color
```

#### Error: Characters not displaying properly
```
Unicode characters appear as boxes or question marks
```

**Solution:**
```bash
# Ensure UTF-8 locale
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

# Install required fonts (for powerline symbols)
# On Ubuntu/Debian:
sudo apt install fonts-powerline

# On macOS:
brew install font-powerline-symbols
```

### Configuration Issues

#### Error: "config file not found"
```
Config file not found at ~/.config/branch-wrangler/config.yml
```

**This is normal** - Branch Wrangler creates the config file automatically on first run.

#### Error: "invalid YAML syntax"
```
Error parsing config file: yaml: line X: mapping values are not allowed in this context
```

**Solution:**
```bash
# Check YAML syntax
yamllint ~/.config/branch-wrangler/config.yml

# Or recreate config file
rm ~/.config/branch-wrangler/config.yml
# Restart Branch Wrangler to regenerate
```

#### Error: "permission denied" on config file
```
permission denied: ~/.config/branch-wrangler/config.yml
```

**Solution:**
```bash
# Fix file permissions
chmod 600 ~/.config/branch-wrangler/config.yml
chmod 700 ~/.config/branch-wrangler/
```

## Performance Issues

### Slow Branch Scanning

#### Large repositories with many branches
**Symptoms:** Takes >10 seconds to scan branches

**Solutions:**
1. **Limit branch scope:**
   - Focus on recent branches only
   - Archive old branches to reduce scan time

2. **Check GitHub API rate limits:**
   ```bash
   # Use authenticated requests for higher limits
   export GITHUB_TOKEN=your_token_here
   ```

3. **Network connectivity:**
   - Check internet connection
   - GitHub API may be slow or experiencing issues

### High Memory Usage

#### Memory consumption with large repositories
**Solutions:**
1. **Monitor memory usage:**
   ```bash
   # Check memory usage
   ps aux | grep branch-wrangler
   ```

2. **Reduce caching:**
   - Restart the application periodically
   - Clear GitHub API cache (restart application)

## Platform-Specific Issues

### Linux Issues

#### Error: "GLIBC version mismatch"
```
version `GLIBC_2.XX' not found
```

**Solution:**
Build from source on your specific Linux distribution:
```bash
git clone https://github.com/dfinster/branch-wrangler.git
cd branch-wrangler
make build-release
```

#### Error: Missing libraries on Alpine Linux
```
/lib64/ld-linux-x86-64.so.2: No such file or directory
```

**Solution:**
```bash
# Install glibc compatibility on Alpine
apk add gcompat

# Or build with musl target
GOOS=linux GOARCH=amd64 go build -tags musl ./cmd/branch-wrangler
```

### macOS Issues

#### Error: "branch-wrangler" is damaged and can't be opened
**Solution:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine branch-wrangler

# Re-download if file is actually corrupted
```

#### Error: Rosetta 2 required (Apple Silicon)
```
Bad CPU type in executable
```

**Solution:**
Download the correct architecture:
- Apple Silicon: `branch-wrangler-darwin-arm64.tar.gz`
- Intel: `branch-wrangler-darwin-amd64.tar.gz`

## Getting Help

### Diagnostic Information

When reporting issues, include:

```bash
# System information
uname -a
go version 2>/dev/null || echo "Go not installed"

# Branch Wrangler version
branch-wrangler --version

# Git repository status
git status
git remote -v

# Terminal information
echo "TERM: $TERM"
echo "Terminal size: $(tput cols)x$(tput lines)"
```

### Log Files

Branch Wrangler doesn't create log files by default. For debugging:

```bash
# Run with verbose output (if available)
branch-wrangler --verbose

# Or capture output
branch-wrangler 2> debug.log
```

### Support Channels

1. **GitHub Issues**: [Report bugs and feature requests](https://github.com/dfinster/branch-wrangler/issues)
2. **GitHub Discussions**: [Ask questions and get help](https://github.com/dfinster/branch-wrangler/discussions)
3. **Documentation**: Check other docs in the `docs/` directory

### Before Reporting Issues

1. Update to the latest version
2. Check existing GitHub issues
3. Try building from source if using pre-built binaries
4. Test in a minimal Git repository
5. Include diagnostic information in your report