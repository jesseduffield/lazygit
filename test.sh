#!/usr/bin/env bash

set -e
echo "" > coverage.txt

export GOFLAGS=-mod=vendor

use_go_test=false
if command -v gotest; then
    use_go_test=true
fi

for d in $( find ./* -maxdepth 10 ! -path "./vendor*" ! -path "./.git*" ! -path "./scripts*" -type d); do
    if ls $d/*.go &> /dev/null; then
        args="-race -coverprofile=profile.out -covermode=atomic $d"
        if [ "$use_go_test" == true ]; then
            gotest $args
        else
            go test $args
        fi
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done
