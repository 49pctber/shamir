# Shamir Secret Sharing

This is a practical example of [Shamir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing).

This tool allows you to divide up a secret string `S` into `n` shares where any `k` of the shares will allow you to reconstruct the original secret `S`.
If you have fewer than `k` of the shares, you will have no information (other than the length of the secret) about what that secret is.

## Usage

To distribute a string, use the command

```
shamir distribute string "<secret string>" -n <number of shares to produce> -k <minimum number of shares to reconstruct the secret>
```

`n` shares will be printed to the terminal.

We can reconstruct the secret by using

```
shamir reconstruct string "<secret 1>" "<secret 2>" ...
```

As long as we provide `k` shares, the secret will then be printed to the terminal. Otherwise, garbage will be printed to the terminal.

### Example

Say we run `shamir distribute string "This is a secret." -n=5 -k=3`.
Something like the following could be printed to the screen:

```
Generating 3-of-5 secret sharing scheme...
shamir-ZDUIQPAX-11d-1-5dxDU36YEiGQmcX3tCfJlHg
shamir-ZDUIQPAX-11d-2-wzJianU+jsC3xr1O5qPj6W8
shamir-ZDUIQPAX-11d-3-coZISivP78FGfwvcMfZPCTk
shamir-ZDUIQPAX-11d-4-4GAeOAldaFJ07zR7rDxfTj8
shamir-ZDUIQPAX-11d-5-UdQ0GFesCVOFVoLpe2nzrmk
```

The prefix `shamir-` indicates that this is a share in Shamir's secret sharing scheme.
`ZDUIQPAX` is a randomly generated ID that allows you to correlate shares to the same secret.
`11d` is the primitive polynomial used to construct the underlying Galois field.
`1` is the x coordinate of the share.
The remaining text is base64-encoded data.
Each byte corresponds to the value of a polynomial evaluated at the corresponding x-coordinate of the share.
(Note that each byte is encoded separately, each with a randomly-generated polynomial.)

To reconstruct the message, simply run `shamir reconstruct` with any three of the five shares.

For example, using shares 1, 4, and 5, we would run the command

```
shamir reconstruct string "shamir-ZDUIQPAX-11d-4-4GAeOAldaFJ07zR7rDxfTj8" "shamir-ZDUIQPAX-11d-5-UdQ0GFesCVOFVoLpe2nzrmk" "shamir-ZDUIQPAX-11d-1-5dxDU36YEiGQmcX3tCfJlHg"
```

We get the following output:

```
...
ZDUIQPAX:
This is a secret.
```

Note that the order in which we supply the shares is irrelevant.

If we try to only use shares 4 and 5, we cannot reconstruct the message, and we get gibberish:

``` text
...
ZDUIQPAX:
���l��V�1�      �u��z
```

### QR Code Support

To make distribution of shares easier in the real world, you can use the `--qr` flag to save each share as a unique QR code.

### Sharing Files

A similar process to the one described above can be performed to use this scheme on a file.

To distribute a file, use `shamir distribute file <filename>` with an appropriate `-n` and `-k`.
This will produce `n` `.txt` files in your current directory.

To reconstruct the secret, run `shamir reconstruct file` in a directory with at least `k` of the shares to reconstruct the secret.
The secret will be saved to a file.
The filename is `secret-<secret id>` with no extension.
**Note that the original filename will be lost!**

## Build Notes

The following scripts are what I use to cross-compile this software.

### From Windows 11 PowerShell

```
$version = "v<major>_<minor>_<patch>"

$env:GOOS = "windows"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION.exe"; go build -o $filename

$env:GOOS = "linux"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
$env:GOOS = "linux"; $env:GOARCH = "arm64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename

$env:GOOS = "darwin"; $env:GOARCH = "arm64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
$env:GOOS = "darwin"; $env:GOARCH = "amd64"; $filename = "shamir-$env:GOOS-$env:GOARCH-$VERSION"; go build -o $filename
```

### From Linux Bash

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
