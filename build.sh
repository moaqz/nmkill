#!/bin/bash

VERSION=${VERSION:-"0.0.0"}
BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2> /dev/null)

go build -ldflags "-X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT-:"unknown-commit"}" -o ./bin/nmkill

echo "✓ Built nmkill ${VERSION} (${GIT_COMMIT})"
