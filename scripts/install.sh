#!/bin/bash

# DBZ Installation Script
# This script downloads and installs the dbz CLI tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="georgeglessner/dbz"
BINARY_NAME="dbz"
INSTALL_DIR="/usr/local/bin"

# Print colored output
print_error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    case "$OS" in
        linux|darwin)
            PLATFORM="${OS}_${ARCH}"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    print_info "Detected platform: $PLATFORM"
}

# Check dependencies
check_dependencies() {
    print_info "Checking dependencies..."
    
    if ! command_exists curl; then
        print_error "curl is required but not installed"
        exit 1
    fi

    if ! command_exists docker; then
        print_info "Docker not found. Please install Docker to use dbz."
        print_info "Visit https://docs.docker.com/get-docker/ for installation instructions."
    fi
}

# Get latest release version
get_latest_version() {
    print_info "Fetching latest version..."
    
    LATEST_VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    print_info "Latest version: $LATEST_VERSION"
}

# Download binary
download_binary() {
    local version=$1
    local platform=$2
    local temp_dir=$(mktemp -d)
    
    print_info "Downloading dbz $version for $platform..."
    
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/dbz_${platform}"
    local binary_path="${temp_dir}/${BINARY_NAME}"
    
    if ! curl -fsSL "$download_url" -o "$binary_path"; then
        print_error "Failed to download binary from $download_url"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    echo "$binary_path"
}

# Install binary
install_binary() {
    local binary_path=$1
    
    print_info "Installing dbz to $INSTALL_DIR..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    else
        print_info "Root privileges required to install to $INSTALL_DIR"
        sudo mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    print_success "Successfully installed dbz to $INSTALL_DIR"
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local version=$("$BINARY_NAME" --version)
        print_success "dbz installed successfully: $version"
        
        # Show usage
        echo ""
        print_info "Getting started:"
        echo "  dbz create postgres    # Create a PostgreSQL database"
        echo "  dbz list              # List all databases"
        echo "  dbz --help            # Show help"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Clean up
cleanup() {
    if [ -n "$temp_dir" ] && [ -d "$temp_dir" ]; then
        rm -rf "$temp_dir"
    fi
}

# Main installation function
main() {
    print_info "Installing dbz - Database CLI Tool"
    echo ""
    
    # Trap to clean up on exit
    trap cleanup EXIT
    
    # Installation steps
    detect_platform
    check_dependencies
    get_latest_version
    
    binary_path=$(download_binary "$LATEST_VERSION" "$PLATFORM")
    install_binary "$binary_path"
    verify_installation
    
    echo ""
    print_success "Installation complete! 🎉"
    print_info "Run 'dbz --help' to get started."
}

# Run main function
main "$@"