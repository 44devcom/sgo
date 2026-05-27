#!/usr/bin/env bash

set -e

APP_NAME="sgo"
REPO="44devcom/sgo"

RAW_OS="$(uname -s)"
RAW_ARCH="$(uname -m)"

echo "Installing $APP_NAME..."

# Normalize architecture (Linux reports aarch64; Go dist uses arm64)
case "$RAW_ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $RAW_ARCH"
    echo "Supported: amd64 (x86_64), arm64 (aarch64). See https://github.com/$REPO#download"
    exit 1
    ;;
esac

# Normalize OS
case "$RAW_OS" in
  Darwin) OS="darwin" ;;
  Linux) OS="linux" ;;
  *)
    echo "Unsupported OS: $RAW_OS"
    exit 1
    ;;
esac

DIST_ID="$OS-$ARCH"
# Termux is Linux+aarch64 but needs the Android/Bionic binary (TLS alignment).
if [ "$ARCH" = "arm64" ] && { [ -n "${TERMUX_VERSION:-}" ] || [ "${PREFIX:-}" = "/data/data/com.termux/files/usr" ]; }; then
  DIST_ID="android-arm64"
fi
case "$DIST_ID" in
  linux-amd64|linux-arm64|android-arm64|darwin-amd64|darwin-arm64) ;;
  *)
    echo "No binary for $DIST_ID"
    exit 1
    ;;
esac

URL="https://github.com/$REPO/raw/refs/heads/master/dist/$DIST_ID/$APP_NAME"

TMP_FILE="$APP_NAME"

echo "Downloading from: $URL"

curl -fsSL "$URL" -o "$TMP_FILE"

chmod +x "$TMP_FILE"

# macOS quarantine fix (safe even if not needed)
if [ "$OS" = "darwin" ]; then
  xattr -dr com.apple.quarantine "$TMP_FILE" 2>/dev/null || true
fi

DOWNLOADS_DIR="$HOME/Downloads"

if [ ! -d "$DOWNLOADS_DIR" ]; then
  if ! mkdir -p "$DOWNLOADS_DIR" 2>/dev/null; then
    echo "Cannot create Downloads directory: $DOWNLOADS_DIR"
    exit 1
  fi
fi

if [ ! -w "$DOWNLOADS_DIR" ]; then
  echo "Downloads directory is not writable: $DOWNLOADS_DIR"
  exit 1
fi

mv "$TMP_FILE" "$DOWNLOADS_DIR"

echo ""
echo "Done ✔ Saved to: $DOWNLOADS_DIR/$APP_NAME"