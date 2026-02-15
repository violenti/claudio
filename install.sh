#!/bin/bash

set -e

REPO="violenti/claudio"
VERSION="v1.0.0"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="claudio"

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

main() {
    detect_platform
    download_binary
    install_binary
}

main
