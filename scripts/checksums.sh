#!/bin/bash
set -euo pipefail

# Branch Wrangler Checksums Script
# Generates SHA256 checksums for built binaries

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Colors for output
COLOR_RESET='\033[0m'
COLOR_GREEN='\033[32m'
COLOR_YELLOW='\033[33m'
COLOR_BLUE='\033[34m'
COLOR_RED='\033[31m'

# Default values
INPUT_DIR="dist"
OUTPUT_FILE="checksums.txt"
VERBOSE=false
VERIFY=false

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

Generate SHA256 checksums for Branch Wrangler binaries

OPTIONS:
    -h, --help              Show this help message
    -d, --dir DIR           Input directory containing binaries (default: dist)
    -o, --output FILE       Output checksums file (default: checksums.txt)
    -v, --verbose           Enable verbose output
    --verify                Verify existing checksums instead of generating new ones

EXAMPLES:
    $0                          # Generate checksums for dist/ directory
    $0 --dir build              # Generate checksums for build/ directory
    $0 --verify                 # Verify existing checksums
    $0 --dir release --output release-checksums.txt  # Custom input/output

EOF
}

# Function to find binary files
find_binaries() {
    local dir="$1"
    
    if [ ! -d "$dir" ]; then
        print_error "Directory not found: $dir"
        return 1
    fi
    
    # Find files that look like binaries (no extension, executable, contain "branch-wrangler")
    find "$dir" -type f -name "*branch-wrangler*" ! -name "*.txt" ! -name "*.md" ! -name "*.json" | sort
}

# Function to generate checksums
generate_checksums() {
    local input_dir="$1"
    local output_file="$2"
    
    print_info "Generating checksums for binaries in $input_dir..."
    
    local binaries=()
    while IFS= read -r binary; do
        [ -n "$binary" ] && binaries+=("$binary")
    done < <(find_binaries "$input_dir")
    
    if [ ${#binaries[@]} -eq 0 ]; then
        print_error "No binary files found in $input_dir"
        print_info "Looking for files matching pattern: *branch-wrangler* (excluding .txt, .md, .json)"
        return 1
    fi
    
    print_info "Found ${#binaries[@]} binary file(s):"
    for binary in "${binaries[@]}"; do
        print_info "  $(basename "$binary")"
    done
    
    # Generate checksums
    local checksums_path="$input_dir/$output_file"
    > "$checksums_path" # Clear the file
    
    local success_count=0
    for binary in "${binaries[@]}"; do
        local relative_path
        relative_path=$(basename "$binary")
        
        if [ "$VERBOSE" = true ]; then
            print_info "Generating checksum for: $relative_path"
        fi
        
        if cd "$input_dir" && shasum -a 256 "$relative_path" >> "$output_file"; then
            ((success_count++))
        else
            print_error "Failed to generate checksum for: $relative_path"
        fi
    done
    
    if [ $success_count -eq ${#binaries[@]} ]; then
        print_success "Generated checksums for $success_count binary file(s)"
        print_success "Checksums saved to: $checksums_path"
        
        # Display checksums
        print_info "Generated checksums:"
        while IFS= read -r line; do
            echo "  $line"
        done < "$checksums_path"
        
        return 0
    else
        print_error "Failed to generate checksums for some files"
        return 1
    fi
}

# Function to verify checksums
verify_checksums() {
    local input_dir="$1"
    local checksums_file="$2"
    
    local checksums_path="$input_dir/$checksums_file"
    
    if [ ! -f "$checksums_path" ]; then
        print_error "Checksums file not found: $checksums_path"
        return 1
    fi
    
    print_info "Verifying checksums from $checksums_path..."
    
    # Change to input directory and verify
    if cd "$input_dir" && shasum -a 256 -c "$checksums_file"; then
        print_success "All checksums verified successfully"
        return 0
    else
        print_error "Checksum verification failed"
        return 1
    fi
}

# Function to show file information
show_file_info() {
    local input_dir="$1"
    
    print_info "File information for $input_dir:"
    
    local binaries=()
    while IFS= read -r binary; do
        [ -n "$binary" ] && binaries+=("$binary")
    done < <(find_binaries "$input_dir")
    
    for binary in "${binaries[@]}"; do
        local size
        size=$(du -h "$binary" | cut -f1)
        local modified
        modified=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "$binary" 2>/dev/null || stat -c "%y" "$binary" 2>/dev/null || echo "unknown")
        
        print_info "  $(basename "$binary"):"
        print_info "    Size: $size"
        print_info "    Modified: $modified"
        
        # Test if binary has version info
        if [ "$VERBOSE" = true ] && [ -x "$binary" ]; then
            print_info "    Version info:"
            if "$binary" --version 2>/dev/null | sed 's/^/      /'; then
                :
            else
                print_warning "      Version command failed or not available"
            fi
        fi
    done
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
            -d|--dir)
                INPUT_DIR="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --verify)
                VERIFY=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    # Show file info if verbose
    if [ "$VERBOSE" = true ]; then
        show_file_info "$INPUT_DIR"
        echo
    fi
    
    # Perform operation
    if [ "$VERIFY" = true ]; then
        verify_checksums "$INPUT_DIR" "$OUTPUT_FILE"
    else
        generate_checksums "$INPUT_DIR" "$OUTPUT_FILE"
    fi
}

# Run main function
main "$@"