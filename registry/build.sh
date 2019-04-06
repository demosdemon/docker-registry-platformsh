#!/usr/bin/env bash

set -eux

mkdir -p bin var/lib/registry

GOPATH=$(mktemp -d --tmpdir go.XXXXXX)
export GOPATH
export PATH="$GOPATH/bin:$PATH"
# cleanup, not really necessary on platform but nice when testing locally
trap 'rm -rf "$GOPATH"' EXIT

go get github.com/demosdemon/docker-registry-platformsh/cmd/config
cp -a "$GOPATH/bin/config" bin/

# unset GO111MODULE as `docker_auth` doesn't support it yet
unset GO111MODULE

outdir=${GOPATH}/src/github.com/docker/distribution
mkdir -p "$(dirname "${outdir}")"

git clone -b v2.7.1 https://github.com/docker/distribution "$outdir"

(
	cd "$outdir" &&
		CGO_ENABLED=0 make PREFIX="$GOPATH" clean binaries &&
		file ./bin/registry
)

cp -a "$outdir/bin/registry" bin/
