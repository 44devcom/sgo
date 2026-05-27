#!/usr/bin/env bash

set -e

APP_NAME="sgo"
REPO="44devcom/sgo"

OS="$(uname -s)"
ARCH="$(uname -m)"

echo "Installing $APP_NAME..."

# Normalize architecture
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
fi

# Normalize OS
if [ "$OS" = "Darwin" ]; then
  OS="darwin"
elif [ "$OS" = "Linux" ]; then
  OS="linux"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

URL="https://github.com/$REPO/raw/refs/heads/master/dist/$OS-$ARCH/$APP_NAME"

TMP_FILE="$APP_NAME"

echo "Downloading from: $URL"

curl -fsSL "$URL" -o "$TMP_FILE"

chmod +x "$TMP_FILE"

# macOS quarantine fix (safe even if not needed)
if [ "$OS" = "darwin" ]; then
  xattr -dr com.apple.quarantine "$TMP_FILE" 2>/dev/null || true
fi

mv "$TMP_FILE" "~/Downloads/$APP_NAME"

echo ""
echo "Done ✔"