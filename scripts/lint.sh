#!/bin/sh

set -e

if [ ! -x ./.bin/golangci-lint ]; then
  echo 'You need to install golangci-lint into .bin'
  echo 'One way to do this is to run'
  echo '  GOBIN=$(pwd)/.bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1'
  exit 1
fi

./.bin/golangci-lint run
