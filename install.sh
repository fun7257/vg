#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="vg"
GITHUB_REPO="fun7257/vg"
TEMP_DIR=$(mktemp -d)

# Cleanup function
cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Function to print colored messages
error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

success() {
    echo -e "${GREEN}$1${NC}"
}

warning() {
    echo -e "${YELLOW}$1${NC}"
}

info() {
    echo -e "$1"
}

# Detect OS and Architecture
detect_platform() {
    OS=""
    ARCH=""
    
    case "$(uname -s)" in
        Linux*)
            OS="linux"
            ;;
        Darwin*)
            OS="darwin"
            ;;
        *)
            error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    
    info "Detected platform: ${OS}-${ARCH}"
}

# Get latest release version from GitHub
get_latest_version() {
    local api_url="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    local version=$(curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')
    
    if [ -z "$version" ]; then
        error "Failed to fetch latest version from GitHub"
        exit 1
    fi
    
    echo "$version"
}

# Get latest version (no user input needed)
get_version() {
    local latest_version=$(get_latest_version)
    VERSION="$latest_version"
    info "Installing vg version: ${VERSION}"
}

# Download binary from GitHub Releases
download_binary() {
    local version="$1"
    local os="$2"
    local arch="$3"
    
    local filename="${BINARY_NAME}-${os}-${arch}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/v${version}/${filename}"
    
    info "Downloading from: ${download_url}"
    
    # Check if curl or wget is available
    if command -v curl &> /dev/null; then
        if ! curl -L -f -o "${TEMP_DIR}/${BINARY_NAME}" "$download_url"; then
            error "Failed to download binary. Please check if version ${version} exists for ${os}-${arch}"
            exit 1
        fi
    elif command -v wget &> /dev/null; then
        if ! wget -q -O "${TEMP_DIR}/${BINARY_NAME}" "$download_url"; then
            error "Failed to download binary. Please check if version ${version} exists for ${os}-${arch}"
            exit 1
        fi
    else
        error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Make binary executable
    chmod +x "${TEMP_DIR}/${BINARY_NAME}"
    
    success "Download successful!"
}

# Install binary
install_binary() {
    # Check if install directory exists and is writable
    if [ ! -d "$(dirname "${INSTALL_DIR}")" ]; then
        error "Install directory parent does not exist: $(dirname "${INSTALL_DIR}")"
        exit 1
    fi
    
    # Check if we need sudo
    NEED_SUDO=false
    if [ ! -w "${INSTALL_DIR}" ]; then
        NEED_SUDO=true
        info "Install directory requires sudo permissions"
    fi
    
    # Install the binary
    info "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
    
    if [ "$NEED_SUDO" = true ]; then
        if sudo cp "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"; then
            sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
            success "vg has been installed to ${INSTALL_DIR}/${BINARY_NAME}"
        else
            error "Failed to install vg"
            exit 1
        fi
    else
        if cp "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"; then
            chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
            success "vg has been installed to ${INSTALL_DIR}/${BINARY_NAME}"
        else
            error "Failed to install vg"
            exit 1
        fi
    fi
}

# Check if this is first time using vg
is_first_time() {
    local vg_home="${HOME}/.vg"
    local sdks_dir="${vg_home}/sdks"
    
    # Check if .vg directory doesn't exist
    if [ ! -d "$vg_home" ]; then
        return 0  # First time
    fi
    
    # Check if sdks directory doesn't exist or is empty
    if [ ! -d "$sdks_dir" ]; then
        return 0  # First time
    fi
    
    # Check if there are any installed versions
    local count=$(find "$sdks_dir" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | wc -l)
    if [ "$count" -eq 0 ]; then
        return 0  # First time
    fi
    
    return 1  # Not first time
}

# Helper function to read from terminal even when stdin is piped
read_from_tty() {
    local prompt="$1"
    local var_name="$2"
    
    # Try to read from /dev/tty if available, otherwise use stdin
    if [ -r /dev/tty ]; then
        # stdin might be piped, read from /dev/tty
        echo -n "$prompt" > /dev/tty
        read "$var_name" < /dev/tty
    elif [ -t 0 ]; then
        # stdin is a terminal, use it
        read -p "$prompt" "$var_name"
    else
        # Fallback: try to read from stdin anyway
        read -p "$prompt" "$var_name" || true
    fi
}

# Interactive Go version installation
install_initial_go_version() {
    info "${BLUE}=== Welcome to vg! ===${NC}"
    echo ""
    info "This appears to be your first time using vg."
    info "Let's install your first Go version to get started."
    echo ""
    
    while true; do
        GO_VERSION_INPUT=""
        read_from_tty "Enter Go version to install (e.g., 1.24.0 or go1.24.0, or 'skip' to skip): " GO_VERSION_INPUT
        
        if [ -z "$GO_VERSION_INPUT" ]; then
            warning "Go version cannot be empty. Please try again."
            continue
        fi
        
        # Allow user to skip
        if [ "$GO_VERSION_INPUT" = "skip" ] || [ "$GO_VERSION_INPUT" = "Skip" ] || [ "$GO_VERSION_INPUT" = "SKIP" ]; then
            info "Skipping initial Go version installation."
            info "You can install a Go version later using: vg install <version>"
            return 0
        fi
        
        # Normalize version (remove 'go' prefix if present)
        GO_VERSION="${GO_VERSION_INPUT#go}"
        
        # Basic version format validation
        if [[ ! "$GO_VERSION" =~ ^[0-9]+\.[0-9]+(\.[0-9]+)?$ ]]; then
            warning "Invalid version format. Please use format like '1.24.0' or '1.24'"
            continue
        fi
        
        break
    done
    
    echo ""
    info "Installing Go ${GO_VERSION}..."
    echo ""
    
    # Check if vg command is available
    if ! command -v "${BINARY_NAME}" &> /dev/null; then
        warning "vg command not found in PATH. Trying to use installed binary directly..."
        VG_CMD="${INSTALL_DIR}/${BINARY_NAME}"
    else
        VG_CMD="${BINARY_NAME}"
    fi
    
    # Install Go version
    if $VG_CMD install "$GO_VERSION"; then
        echo ""
        success "Go ${GO_VERSION} installed successfully!"
        
        # Automatically switch to the installed version
        if $VG_CMD use "$GO_VERSION"; then
            echo ""
            success "Switched to Go ${GO_VERSION}!"
        else
            warning "Failed to switch to Go ${GO_VERSION}"
            info "You can switch manually using: vg use ${GO_VERSION}"
        fi
    else
        error "Failed to install Go ${GO_VERSION}"
        warning "You can try installing it later using: vg install ${GO_VERSION}"
    fi
}

# Detect shell and find configuration file
detect_shell_config() {
    local shell_name=$(basename "$SHELL" 2>/dev/null || echo "bash")
    local config_file=""
    
    case "$shell_name" in
        zsh)
            # zsh: prefer .zshrc, fallback to .zprofile
            if [ -f "${HOME}/.zshrc" ]; then
                config_file="${HOME}/.zshrc"
            elif [ -f "${HOME}/.zprofile" ]; then
                config_file="${HOME}/.zprofile"
            else
                config_file="${HOME}/.zshrc"
            fi
            ;;
        bash)
            # bash: prefer .bashrc (Linux), fallback to .bash_profile (macOS)
            if [ -f "${HOME}/.bashrc" ]; then
                config_file="${HOME}/.bashrc"
            elif [ -f "${HOME}/.bash_profile" ]; then
                config_file="${HOME}/.bash_profile"
            else
                config_file="${HOME}/.bashrc"
            fi
            ;;
        sh|dash|ash)
            # POSIX shells: use .profile
            config_file="${HOME}/.profile"
            ;;
        *)
            # Default: try .profile (most universal)
            config_file="${HOME}/.profile"
            ;;
    esac
    
    echo "$config_file"
}

# Check if vg init is already in config file
has_vg_init() {
    local config_file="$1"
    
    if [ ! -f "$config_file" ]; then
        return 1  # File doesn't exist, so not present
    fi
    
    # Check if eval "$(vg init)" or similar exists
    if grep -q 'eval.*vg init' "$config_file" 2>/dev/null; then
        return 0  # Found
    fi
    
    return 1  # Not found
}

# Add vg init to shell configuration
add_vg_init_to_config() {
    local config_file="$1"
    local vg_cmd="${BINARY_NAME}"
    
    # If vg is not in PATH, use full path
    if ! command -v "${BINARY_NAME}" &> /dev/null; then
        vg_cmd="${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    local init_line="eval \"\$(${vg_cmd} init)\""
    
    # Create file if it doesn't exist
    if [ ! -f "$config_file" ]; then
        touch "$config_file"
    fi
    
    # Add vg init to config file
    echo "" >> "$config_file"
    echo "# vg - Virtual Go environment manager" >> "$config_file"
    echo "$init_line" >> "$config_file"
    
    success "Added vg initialization to ${config_file}"
}

# Setup shell configuration
setup_shell_config() {
    local config_file=$(detect_shell_config)
    local shell_name=$(basename "$SHELL" 2>/dev/null || echo "unknown")
    
    info "Detected shell: ${shell_name}"
    info "Configuration file: ${config_file}"
    echo ""
    
    # Check if already configured
    if has_vg_init "$config_file"; then
        info "vg is already configured in ${config_file}"
        return 0
    fi
    
    # Always prompt user (even in non-interactive mode, we can read from /dev/tty)
    echo ""
    info "📝 Shell Configuration"
    echo ""
    
    # Interactive: ask user if they want to add it
    REPLY=""
    read_from_tty "Add 'eval \"\$(vg init)\"' to ${config_file}? (Y/n) " REPLY
    echo ""
    if [[ "$REPLY" =~ ^[Nn]$ ]]; then
        info "Skipped. You can manually add 'eval \"\$(vg init)\"' to your shell configuration."
        info "Or run it manually in each new shell session."
        return 0
    fi
    
    # Add to config file
    if add_vg_init_to_config "$config_file"; then
        echo ""
        success "Configuration updated!"
        info "The changes will take effect in new shell sessions."
        info "To use it in the current session, run: source ${config_file}"
        echo ""
        REPLY=""
        read_from_tty "Apply configuration to current session now? (Y/n) " REPLY
        echo ""
        if [ -z "$REPLY" ] || [[ ! "$REPLY" =~ ^[Nn]$ ]]; then
            if [ -f "$config_file" ]; then
                # Source the file to apply changes
                if source "$config_file" 2>/dev/null || . "$config_file" 2>/dev/null; then
                    success "Configuration applied to current session!"
                else
                    warning "Could not apply configuration to current session"
                    info "Please run: source ${config_file}"
                fi
            fi
        fi
    else
        error "Failed to add configuration to ${config_file}"
    fi
}

# Verify installation
verify_installation() {
    if command -v "${BINARY_NAME}" &> /dev/null; then
        success "Installation verified! vg is available in PATH"
    else
        warning "vg was installed but may not be in your PATH"
        info "Make sure ${INSTALL_DIR} is in your PATH"
        info "You can add it by running: export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Main installation flow
main() {
    info "${BLUE}=== vg Installation Script ===${NC}"
    echo ""
    
    # Detect platform
    detect_platform
    echo ""
    
    # Get version
    get_version
    echo ""
    
    # Download binary
    download_binary "$VERSION" "$OS" "$ARCH"
    echo ""
    
    # Install binary
    install_binary
    echo ""
    
    # Verify installation
    verify_installation
    echo ""
    
    # Check if first time and offer to install Go version
    # This must be done before shell configuration
    if is_first_time; then
        echo ""
        install_initial_go_version
        echo ""
    fi
    
    # Setup shell configuration (after Go version is installed)
    setup_shell_config
    echo ""
    
    if ! is_first_time; then
        info "Run 'vg --help' to see available commands"
    fi
}

# Run main function
main
