#!/bin/bash
set -euo pipefail

# Branch Wrangler Build Script
# Provides flexible building for different scenarios

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
APP_NAME="branch-wrangler"
MAIN_PATH="./cmd/branch-wrangler"

# Colors for output
COLOR_RESET='\033[0m'
COLOR_GREEN='\033[32m'
COLOR_YELLOW='\033[33m'
COLOR_BLUE='\033[34m'
COLOR_RED='\033[31m'

# Default values
OUTPUT_DIR="dist"
BUILD_TYPE="release"
GOOS="darwin"
GOARCH="arm64"
VERBOSE=false

# Function to print colored output
print_info() {
    echo -e "${COLOR_BLUE}$1${COLOR_RESET}"
}

print_success() {
    echo -e "${COLOR_GREEN}✓ $1${COLOR_RESET}"
}

print_warning() {
    echo -e "${COLOR_YELLOW}⚠ $1${COLOR_RESET}"
}

print_error() {
    echo -e "${COLOR_RED}✗ $1${COLOR_RESET}" >&2
}

# Function to show usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Build script for Branch Wrangler

OPTIONS:
    -h, --help          Show this help message
    -o, --output DIR    Output directory (default: dist)
    -t, --type TYPE     Build type: release|dev (default: release)
    -v, --verbose       Enable verbose output
    --os OS             Target OS (default: darwin)
    --arch ARCH         Target architecture (default: arm64)
    --version VERSION   Override version (default: auto-detect)

EXAMPLES:
    $0                                    # Default release build
    $0 --type dev --verbose              # Development build with verbose output
    $0 --output build --type dev         # Development build to custom directory
    $0 --version v1.0.0                  # Release build with specific version

EOF
}

# Function to get version information
get_version() {
    if [ -n "${VERSION:-}" ]; then
        echo "$VERSION"
    elif git describe --tags --exact-match HEAD 2>/dev/null; then
        # We're on a tagged commit
        git describe --tags --exact-match HEAD
    elif git describe --tags --always --dirty 2>/dev/null; then
        # Use git describe for development versions
        git describe --tags --always --dirty
    else
        echo "dev"
    fi
}

get_build_date() {
    date -u +%Y-%m-%dT%H:%M:%SZ
}

get_commit_hash() {
    git rev-parse HEAD 2>/dev/null || echo "unknown"
}

# Function to validate environment
check_requirements() {
    print_info "Checking build requirements..."
    
    if ! go version >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v git >/dev/null 2>&1; then
        print_error "Git is not installed or not in PATH"
        exit 1
    fi
    
    # Check if we're in a Go module
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        print_error "go.mod not found. Are you in the project root?"
        exit 1
    fi
    
    print_success "Build requirements satisfied"
}

# Function to build the binary
build_binary() {
    local version="$1"
    local build_date="$2"
    local commit_hash="$3"
    local output_path="$4"
    
    print_info "Building binary..."
    print_info "  Version: $version"
    print_info "  Build Date: $build_date"
    print_info "  Commit: $commit_hash"
    print_info "  Target: ${GOOS}/${GOARCH}"
    print_info "  Output: $output_path"
    
    # Create output directory
    mkdir -p "$(dirname "$output_path")"
    
    # Prepare ldflags
    local ldflags="-X github.com/dfinster/branch-wrangler/internal/version.Version=$version"
    ldflags="$ldflags -X github.com/dfinster/branch-wrangler/internal/version.BuildDate=$build_date"
    ldflags="$ldflags -X github.com/dfinster/branch-wrangler/internal/version.CommitHash=$commit_hash"
    
    # Prepare build flags
    local build_flags=()
    if [ "$BUILD_TYPE" = "release" ]; then
        build_flags+=("-trimpath")
    else
        build_flags+=("-race")
    fi
    
    # Execute build
    cd "$PROJECT_ROOT"
    
    local cmd_args=()
    if [ "$VERBOSE" = true ]; then
        cmd_args+=("-v")
    fi
    cmd_args+=("${build_flags[@]}" -ldflags "$ldflags" -o "$output_path" "$MAIN_PATH")
    
    if [ "$VERBOSE" = true ]; then
        print_info "Build command: go build ${cmd_args[*]}"
    fi
    
    GOOS="$GOOS" GOARCH="$GOARCH" go build "${cmd_args[@]}"
    
    # Verify binary was created
    if [ ! -f "$output_path" ]; then
        print_error "Build failed: binary not found at $output_path"
        exit 1
    fi
    
    print_success "Binary built successfully: $output_path"
    
    # Show binary info
    local binary_size
    binary_size=$(du -h "$output_path" | cut -f1)
    print_info "Binary size: $binary_size"
}

# Function to test the built binary
test_binary() {
    local binary_path="$1"
    
    print_info "Testing binary..."
    
    if ! "$binary_path" --version >/dev/null 2>&1; then
        print_error "Binary test failed: --version command failed"
        return 1
    fi
    
    print_success "Binary test passed"
    
    if [ "$VERBOSE" = true ]; then
        print_info "Version output:"
        "$binary_path" --version
    fi
}

# Main function
main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -t|--type)
                BUILD_TYPE="$2"
                if [[ "$BUILD_TYPE" != "release" && "$BUILD_TYPE" != "dev" ]]; then
                    print_error "Invalid build type: $BUILD_TYPE. Must be 'release' or 'dev'"
                    exit 1
                fi
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --os)
                GOOS="$2"
                shift 2
                ;;
            --arch)
                GOARCH="$2"
                shift 2
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            *)
                print_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    # Get build information
    local version
    version=$(get_version)
    local build_date
    build_date=$(get_build_date)
    local commit_hash
    commit_hash=$(get_commit_hash)
    
    # Determine output filename
    local binary_name
    if [ "$BUILD_TYPE" = "release" ]; then
        binary_name="${APP_NAME}-${version}-${GOOS}-${GOARCH}"
    else
        binary_name="${APP_NAME}"
    fi
    local output_path="${OUTPUT_DIR}/${binary_name}"
    
    # Run build process
    check_requirements
    build_binary "$version" "$build_date" "$commit_hash" "$output_path"
    test_binary "$output_path"
    
    print_success "Build completed successfully!"
    print_info "Binary location: $output_path"
}

# Run main function
main "$@"