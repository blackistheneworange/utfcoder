#!/bin/bash
set -e

echo "Building..."

mkdir -p build

export CGO_ENABLED=0

targets=(
    "windows amd64"
    "windows arm64"
    "linux amd64"
    "linux arm64"
    "darwin amd64"
    "darwin arm64"
)

for t in "${targets[@]}"; do
    set -- $t

    ext=""
    if [[ $1 == "windows" ]]; then
        ext=".exe"
    fi

    echo "Building for OS=${1}, Arch=${2}..."
    GOOS=$1 GOARCH=$2 go build -o "build/utfcoder_${1}-${2}${ext}"
done

echo "Build finished"