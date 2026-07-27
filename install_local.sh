#!/usr/bin/env bash
set -e

# Determine executable name (Windows needs .exe)
EXE_NAME="windmist"
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    EXE_NAME="windmist.exe"
fi

echo "🔨 Building WindMist from source..."
go build -o "$EXE_NAME" .

# Find if windmist is already installed in PATH
EXISTING_PATH=$(command -v "$EXE_NAME" || true)

if [ -n "$EXISTING_PATH" ]; then
    echo "📦 Existing installation found at $EXISTING_PATH. Overwriting..."
    
    # Check if we need sudo to overwrite
    if [ -w "$(dirname "$EXISTING_PATH")" ]; then
        mv "$EXE_NAME" "$EXISTING_PATH"
    else
        echo "🔒 Requires sudo to overwrite $EXISTING_PATH"
        sudo mv "$EXE_NAME" "$EXISTING_PATH"
    fi
else
    # Default local installation folder if not in PATH
    INSTALL_DIR="$HOME/.local/bin"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
        INSTALL_DIR="$HOME/bin"
    fi
    
    echo "📦 Installing WindMist to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
    mv "$EXE_NAME" "$INSTALL_DIR/$EXE_NAME"
    
    echo "⚠️  NOTE: Please ensure $INSTALL_DIR is in your system's PATH!"
fi

echo "✅ WindMist successfully built and installed!"
echo "Run 'windmist --version' to verify."
