#!/usr/bin/env bash
# Universal 1-Line Installation Script for WindMist
# Usage: curl -sSL https://raw.githubusercontent.com/Nithwin/windmist/main/scripts/install.sh | bash

set -euo pipefail

REPO="Nithwin/windmist"
INSTALL_DIR="/usr/local/bin"

echo "🌪️  Installing WindMist CLI..."

# 1. Detect Operating System
OS="$(uname -s)"
case "$OS" in
  Linux*)     OS_NAME="Linux" ;;
  Darwin*)    OS_NAME="macOS" ;;
  *)          echo "❌ Unsupported operating system: $OS. Please install manually or via 'go install'." && exit 1 ;;
esac

# 2. Detect CPU Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)   ARCH_NAME="x86_64" ;;
  arm64|aarch64)  ARCH_NAME="arm64" ;;
  *)              echo "❌ Unsupported architecture: $ARCH. Please install manually or via 'go install'." && exit 1 ;;
esac

# 3. Get Latest Release Tag
echo "🔍 Fetching latest release tag from GitHub..."
LATEST_TAG="$(curl -sSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/' || true)"

if [ -z "$LATEST_TAG" ]; then
  # Fallback if API rate limited or offline, default to v1.0.1
  LATEST_TAG="v1.0.1"
fi

VERSION="${LATEST_TAG#v}" # strip leading 'v'

# Handle transition between earlier 'Darwin' template in v1.0.1 and new 'macOS' template in v1.0.2+
DOWNLOAD_OS="$OS_NAME"
if [ "$OS_NAME" = "macOS" ] && [ "$LATEST_TAG" = "v1.0.1" ]; then
  DOWNLOAD_OS="Darwin"
fi

TARBALL="windmist_${VERSION}_${DOWNLOAD_OS}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TARBALL"

TMP_DIR="$(mktemp -d -t windmist-install.XXXXXX)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

echo "📦 Downloading WindMist $LATEST_TAG ($DOWNLOAD_OS $ARCH_NAME)..."
if ! curl -sSL -f -o "$TMP_DIR/$TARBALL" "$DOWNLOAD_URL"; then
  echo "❌ Error downloading from $DOWNLOAD_URL"
  echo "Please verify that release $LATEST_TAG has an archive named '$TARBALL'."
  exit 1
fi

echo "📂 Extracting archive..."
tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/windmist" ]; then
  echo "❌ Error: 'windmist' executable not found inside archive."
  exit 1
fi

echo "🚀 Installing WindMist binary..."
if [ -w "$INSTALL_DIR" ]; then
  install -m 755 "$TMP_DIR/windmist" "$INSTALL_DIR/windmist"
elif command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
  sudo install -m 755 "$TMP_DIR/windmist" "$INSTALL_DIR/windmist"
elif [ -t 0 ] && command -v sudo >/dev/null 2>&1; then
  echo "🔑 Prompting for sudo password/fingerprint to install to $INSTALL_DIR..."
  sudo install -m 755 "$TMP_DIR/windmist" "$INSTALL_DIR/windmist" || {
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    install -m 755 "$TMP_DIR/windmist" "$INSTALL_DIR/windmist"
    echo "ℹ️ Installed to $INSTALL_DIR as fallback."
  }
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  install -m 755 "$TMP_DIR/windmist" "$INSTALL_DIR/windmist"
  echo "ℹ️ Installed to $INSTALL_DIR without sudo (user-local installation)."
fi

echo ""
echo "✅ WindMist $LATEST_TAG successfully installed!"
echo "Try running:"
echo "  windmist version"
