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
