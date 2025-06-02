#!/bin/bash

set -e

echo "🔨 Building KubeStellar Cluster Operations Plugin..."

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_NAME="cluster-ops-plugin"

# Create local build directory
BUILD_DIR="${SCRIPT_DIR}/build"
mkdir -p "${BUILD_DIR}"

echo "📦 Compiling Go plugin..."
cd "${SCRIPT_DIR}"

# Build with optimizations for production
go build -buildmode=plugin \
    -ldflags='-w -s' \
    -o "${BUILD_DIR}/${PLUGIN_NAME}.so" \
    main.go

echo "✅ Plugin built successfully: ${BUILD_DIR}/${PLUGIN_NAME}.so"

# Copy manifest file
echo "📋 Copying plugin manifest..."
cp plugin.yaml "${BUILD_DIR}/${PLUGIN_NAME}.yaml"

echo "✅ Manifest copied: ${BUILD_DIR}/${PLUGIN_NAME}.yaml"

# Verify the plugin
echo "🔍 Verifying plugin..."
if [ -f "${BUILD_DIR}/${PLUGIN_NAME}.so" ]; then
    echo "✅ Plugin file exists and is readable"
    echo "📊 Plugin size: $(du -h "${BUILD_DIR}/${PLUGIN_NAME}.so" | cut -f1)"
else
    echo "❌ Plugin file not found!"
    exit 1
fi

echo "🎉 Build completed successfully!"
echo ""
echo "📁 Plugin files are in: ${BUILD_DIR}/"
echo "🔌 Plugin binary: ${BUILD_DIR}/${PLUGIN_NAME}.so"
echo "📋 Plugin manifest: ${BUILD_DIR}/${PLUGIN_NAME}.yaml" 