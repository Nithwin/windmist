#!/usr/bin/env bash
set -e

# WindMist Installer
echo "Installing WindMist..."

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

if [ "${OS}" != "linux" ] && [ "${OS}" != "darwin" ]; then
  echo "Unsupported OS: ${OS}"
  exit 1
fi

# Fetch latest release version from GitHub API
echo "Fetching latest release..."
TAG=$(curl -sL https://api.github.com/repos/Nithwin/WindMist/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$TAG" ]; then
  echo "Failed to fetch latest release. Please install manually."
  exit 1
fi

FILENAME="windmist_${TAG}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/Nithwin/WindMist/releases/download/v${TAG}/${FILENAME}"

TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo "Downloading ${FILENAME}..."
curl -sL -o windmist.tar.gz "${DOWNLOAD_URL}"

echo "Extracting..."
tar -xzf windmist.tar.gz

echo "Installing to /usr/local/bin (requires sudo)..."
sudo mv windmist /usr/local/bin/windmist

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "✅ WindMist installed successfully!"
echo "Run 'windmist' to get started."
