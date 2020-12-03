#!/bin/bash

function die() {
  echo "$@" 1>&2
  exit 1
}

GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}

case "${GOOS}_${GOARCH}" in
  linux_amd64) CC=musl-gcc ;;
  linux_arm64) CC=aarch64-linux-musl-gcc ;;
  linux_arm)   CC=arm-linux-musleabihf-gcc ;;
  darwin_amd64) CC=/opt/osxcross/target/bin/o64-clang ;;
  *) die "No cross-compiler set for ${GOOS}_${GOARCH}" ;;
esac

exec ${CC} "$@"
