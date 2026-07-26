#!/bin/bash

set -e

REPO="violenti/claudio"
VERSION="${CLAUDIO_VERSION:-latest}"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="claudio"

resolve_version() {
    if [ "$VERSION" != "latest" ]; then
        return
    fi

    echo "Looking up latest release..."

    if command -v curl &> /dev/null; then
        VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
    else
        echo "Error: curl or wget is required"
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest ${BINARY_NAME} version"
        exit 1
    fi
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        *)
            echo "Error: Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            echo "Error: Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    if [ "$OS" = "linux" ] && [ "$ARCH" = "arm64" ]; then
        echo "Error: linux-arm64 is not available yet"
        exit 1
    fi

    PLATFORM="${OS}-${ARCH}"
}

download_binary() {
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}"

    echo "Downloading ${BINARY_NAME} ${VERSION} for ${PLATFORM}..."

    if command -v curl &> /dev/null; then
        curl -fsSL "$URL" -o "/tmp/${BINARY_NAME}"
    elif command -v wget &> /dev/null; then
        wget -q "$URL" -O "/tmp/${BINARY_NAME}"
    else
        echo "Error: curl or wget is required"
        exit 1
    fi
}

install_binary() {
    echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."

    chmod +x "/tmp/${BINARY_NAME}"

    if [ -w "$INSTALL_DIR" ]; then
        mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    echo "Successfully installed ${BINARY_NAME} ${VERSION}"
    echo "Run 'claudio --help' to get started"
}

install_quotes() {
    QUOTES_DIR="${HOME}/.claudio"
    QUOTES_FILE="${QUOTES_DIR}/foghorn_quotes.json"
    QUOTES_URL="https://raw.githubusercontent.com/${REPO}/main/assets/foghorn_quotes.json"

    if [ -f "$QUOTES_FILE" ]; then
        echo "Foghorn quotes already installed at ${QUOTES_FILE}, keeping your version"
        return 0
    fi

    echo "Installing Foghorn Leghorn quotes to ${QUOTES_FILE}..."
    mkdir -p "$QUOTES_DIR"

    if command -v curl &> /dev/null; then
        curl -fsSL "$QUOTES_URL" -o "$QUOTES_FILE" || echo "Warning: could not download quotes file, skipping (Claudio works fine without it)"
    elif command -v wget &> /dev/null; then
        wget -q "$QUOTES_URL" -O "$QUOTES_FILE" || { rm -f "$QUOTES_FILE"; echo "Warning: could not download quotes file, skipping (Claudio works fine without it)"; }
    fi
}

install_config() {
    CONFIG_DIR="${HOME}/.claudio"
    CONFIG_FILE="${CONFIG_DIR}/config.json"
    CONFIG_URL="https://raw.githubusercontent.com/${REPO}/main/config.example.json"

    if [ -f "$CONFIG_FILE" ]; then
        echo "Config already installed at ${CONFIG_FILE}, keeping your version"
        return 0
    fi

    echo "Installing default config to ${CONFIG_FILE}..."
    mkdir -p "$CONFIG_DIR"

    if command -v curl &> /dev/null; then
        curl -fsSL "$CONFIG_URL" -o "$CONFIG_FILE" || { rm -f "$CONFIG_FILE"; echo "Warning: could not download config file, copy config.example.json to ${CONFIG_FILE} manually"; }
    elif command -v wget &> /dev/null; then
        wget -q "$CONFIG_URL" -O "$CONFIG_FILE" || { rm -f "$CONFIG_FILE"; echo "Warning: could not download config file, copy config.example.json to ${CONFIG_FILE} manually"; }
    fi
}

main() {
    detect_platform
    resolve_version
    download_binary
    install_binary
    install_config
    install_quotes
}

main
