#!/usr/bin/env bash
# Downloads the OpenObserve binary for the current platform.
# Run this on the HOST (not inside Docker) so it uses the system CA store
# (which includes the Zscaler cert).
#
# Usage:
#   bash observability/openobserve/download-openobserve.sh
#
# The binary is saved to observability/openobserve/openobserve and is
# .gitignored. The Dockerfile COPYs it into the container image.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT="$SCRIPT_DIR/openobserve"

# Detect platform
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
    linux)  PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    *)      echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64)   ARCH_SUFFIX="amd64" ;;
    aarch64|arm64)  ARCH_SUFFIX="arm64" ;;
    *)              echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

VERSION="${1:-latest}"
URL="https://github.com/openobserve/openobserve/releases/download/v${VERSION}/openobserve-v${VERSION}-${PLATFORM}-${ARCH_SUFFIX}.tar.gz"

# For "latest", use the download site which resolves to the latest version
if [ "$VERSION" = "latest" ]; then
    URL="https://downloads.openobserve.ai/releases/openobserve/latest/openobserve-latest-${PLATFORM}-${ARCH_SUFFIX}.tar.gz"
fi

echo "Downloading OpenObserve ($PLATFORM/$ARCH_SUFFIX) ..."
echo "  URL: $URL"

# Download and extract — the tarball contains a single 'openobserve' binary
curl -fSL "$URL" | tar -xz -C "$SCRIPT_DIR"

if [ -f "$OUTPUT" ]; then
    chmod +x "$OUTPUT"
    echo "✓ Downloaded to $OUTPUT ($(du -h "$OUTPUT" | cut -f1))"
else
    echo "✗ Binary not found after extraction" >&2
    ls -la "$SCRIPT_DIR/"
    exit 1
fi
