#!/bin/sh

# This is used by VSCode; it is not very useful otherwise, since it's easy
# enough to just run `go tool gofumpt` directly, or use `just format`.

set -e

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(dirname "$script_dir")

cd "$repo_root"
exec go tool gofumpt "$@"
