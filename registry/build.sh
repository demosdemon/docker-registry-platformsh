#!/usr/bin/env bash

set -eux

go build -o bin/config ./cmd/config

# unset GO111MODULE as `distribution` doesn't support it yet
unset GO111MODULE

GOPATH=$(mktemp -d --tmpdir go.XXXXXX)
export GOPATH

# cleanup, not really necessary on platform but nice when testing locally
trap 'rm -rf "$GOPATH"' EXIT

DISTRIBUTION_DIR=${GOPATH}/src/github.com/docker/distribution
mkdir -p "$(dirname "${DISTRIBUTION_DIR}")"

git clone -b v2.7.1 https://github.com/docker/distribution "$DISTRIBUTION_DIR"

(
	cd "$DISTRIBUTION_DIR" &&
		CGO_ENABLED=0 make PREFIX="$GOPATH" clean binaries &&
		file ./bin/registry
)

mkdir -p bin etc/docker/registry var/lib/registry
cp -a "$DISTRIBUTION_DIR/bin/registry" bin/
