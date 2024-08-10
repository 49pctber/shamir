# Instructions for Cross Compilation

## From Windows 11 PowerShell

```
$version = "v<major>_<minor>_<patch>"

$env:GOOS = "windows"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION.exe"; go build -o $filename

$env:GOOS = "linux"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
$env:GOOS = "linux"; $env:GOARCH = "arm64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename

$env:GOOS = "darwin"; $env:GOARCH = "arm64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
$env:GOOS = "darwin"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
```

## From Linux Bash

``` bash
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
```
