#!/bin/bash

OUTPUT_NAME="shamir"
GOOS="${1:-linux}"
GOARCH="${2:-amd64}"
VERSION="${3:v0_0_0}"

# Build the executable
filename="${OUTPUT_NAME}-${GOOS}-${GOARCH}-${VERSION}"
if [ "$GOOS" = "windows" ]; then
    filename="${filename}.exe"
fi

GOARCH=$GOARCH GOOS=$GOOS go build -o "$filename"

echo "Build complete: $filename"
