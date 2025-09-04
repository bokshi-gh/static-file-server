#!/bin/bash
set -e

PROJECT_NAME=goserve
VERSION=1.0.0
OUTPUT_DIR=release

rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

declare -A PLATFORMS
PLATFORMS=( ["linux"]="amd64" ["darwin"]="amd64" ["windows"]="amd64" )

for OS in "${!PLATFORMS[@]}"; do
    ARCH=${PLATFORMS[$OS]}
    echo "Building for $OS/$ARCH..."

    BIN_NAME=$PROJECT_NAME
    [ "$OS" == "windows" ] && BIN_NAME=$PROJECT_NAME.exe

    BUILD_DIR=$OUTPUT_DIR/${OS}
    mkdir -p $BUILD_DIR

    GOOS=$OS GOARCH=$ARCH go build -o $BUILD_DIR/$BIN_NAME main.go

    pushd $OUTPUT_DIR > /dev/null
    zip -r ${PROJECT_NAME}-${VERSION}-${OS}.zip ${OS}
    popd > /dev/null

    echo "Packaged ${PROJECT_NAME}-${VERSION}-${OS}.zip"
done

echo "All builds completed!"
