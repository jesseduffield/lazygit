#!/bin/sh

set -e

# Must be kept in sync with the version in .github/workflows/ci.yml
version="v2.4.0"

go run "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$version" "$@"
