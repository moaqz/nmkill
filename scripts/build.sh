#!/bin/bash

VERSION=${VERSION:-"0.0.0"}
BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2> /dev/null)

PLATFORMS=${PLATFORMS:-"linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"}
LDFLAGS="-X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT-:"unknown-commit"}"

OUTPUT_DIR=${OUTPUT_DIR:-"./dist"}
mkdir -p "${OUTPUT_DIR}"

for platform in $PLATFORMS; do
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"

    out="${OUTPUT_DIR}/nmkill-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        out+='.exe'
    fi

    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "${LDFLAGS}" -o "$out"
    echo "✓ Built ${out}"
done

echo "✓ Build complete for version ${VERSION} (${GIT_COMMIT})"
