#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $( find ./* -maxdepth 10 ! -path "./vendor*" ! -path "./.git*" -type d); do
    if ls $d/*.go &> /dev/null; then
        go test -v -race -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done
