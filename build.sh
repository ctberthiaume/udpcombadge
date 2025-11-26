#!/bin/bash
# Build command-line tool for 64-bit MacOS and Linux

VERSION=$(git describe --tags)

[[ -d build ]] || mkdir build
GOOS=darwin GOARCH=amd64 go build -o build/udpcombadge.${VERSION}.darwin-amd64 main.go || exit 1
GOOS=linux GOARCH=amd64 go build -o build/udpcombadge.${VERSION}.linux-amd64 main.go || exit 1
GOOS=darwin GOARCH=arm64 go build -o build/udpcombadge.${VERSION}.darwin-arm64 main.go || exit 1
openssl dgst -sha256 build/*.${VERSION}.* | sed -e 's|build/||g'