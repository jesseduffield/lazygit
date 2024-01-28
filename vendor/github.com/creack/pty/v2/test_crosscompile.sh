#!/usr/bin/env sh

# Test script checking that all expected os/arch compile properly.
# Does not actually test the logic, just the compilation so we make sure we don't break code depending on the lib.

echo2() {
  echo $@ >&2
}

trap end 0
end() {
  [ "$?" = 0 ] && echo2 "Pass." || (echo2 "Fail."; exit 1)
}

cross() {
  os=$1
  shift
  echo2 "Build for $os."
  for arch in $@; do
    echo2 "  - $os/$arch"
    GOOS=$os GOARCH=$arch go build
  done
  echo2
}

set -e

cross linux     amd64 386 arm arm64 ppc64 ppc64le s390x mips mipsle mips64 mips64le riscv64
cross darwin    amd64 arm64
cross freebsd   amd64 386 arm arm64 riscv64
cross netbsd    amd64 386 arm arm64
cross openbsd   amd64 386 arm arm64
cross dragonfly amd64
cross solaris   amd64
cross windows   amd64 386 arm

# TODO: Fix compilation error on openbsd/arm.
# TODO: Merge the solaris PR.

# Some os/arch require a different compiler. Run in docker.
if ! hash docker; then
  # If docker is not present, stop here.
  return
fi

# Golang dropped support for darwin 32bits since go1.15. Make sure the lib still compile with go1.14 on those archs.
echo2 "Build for darwin (32bits)."
echo2 "  - darwin/386"
docker build -t creack-pty-test -f Dockerfile.golang --build-arg=GOVERSION=1.14 --build-arg=GOOS=darwin --build-arg=GOARCH=386 .
echo2 "  - darwin/arm"
docker build -t creack-pty-test -f Dockerfile.golang --build-arg=GOVERSION=1.14 --build-arg=GOOS=darwin --build-arg=GOARCH=arm .

# Run a single test for an old go version. Would be best with go1.0, but not available on Dockerhub.
# Using 1.6 as it is the base version for the RISCV compiler.
# Would also be better to run all the tests, not just one, need to refactor this file to allow for specifc archs per version.
echo2 "Build for linux - go1.6."
echo2 "  - linux/amd64"
docker build -t creack-pty-test -f Dockerfile.golang --build-arg=GOVERSION=1.6 --build-arg=GOOS=linux --build-arg=GOARCH=amd64 .
