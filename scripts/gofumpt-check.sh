#!/bin/sh

# Checks that all Go files are gofumpt-formatted, and fails if any aren't.
# Used by `just lint`, `make lint`, and CI. We run gofumpt with the version
# pinned in go.mod (via `go tool`) rather than the one bundled with
# golangci-lint, so that formatting is identical across all of them and the
# editor.

set -e

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(dirname "$script_dir")

cd "$repo_root"

unformatted=$(go tool gofumpt -l .)
if [ -n "$unformatted" ]; then
  echo "The following files are not formatted correctly:"
  echo "$unformatted"
  echo "Run 'just format' (or 'make format') and commit the result."
  exit 1
fi
