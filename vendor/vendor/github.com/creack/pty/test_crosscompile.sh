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

cross linux     amd64 386 arm arm64 ppc64 ppc64le s390x mips mipsle mips64 mips64le
cross darwin    amd64 386 arm arm64
cross freebsd   amd64 386 arm
cross netbsd    amd64 386 arm
cross openbsd   amd64 386 arm arm64
cross dragonfly amd64
cross solaris   amd64

# Not expected to work but should still compile.
cross windows amd64 386 arm

# TODO: Fix compilation error on openbsd/arm.
# TODO: Merge the solaris PR.

# Some os/arch require a different compiler. Run in docker.
if ! hash docker; then
    # If docker is not present, stop here.
    return
fi

echo2 "Build for linux."
echo2 "  - linux/riscv"
docker build -t test -f Dockerfile.riscv .
