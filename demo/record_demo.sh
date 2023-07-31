#!/bin/sh

TEST=$1

set -e

if [ -z "$TEST" ]
then
    echo "Usage: $0 <test>"
    exit 1
fi

if ! command -v terminalizer &> /dev/null
then
    echo "terminalizer could not be found"
    echo "Install it with: npm install -g terminalizer"
    exit 1
fi

if ! command -v "gifsicle" &> /dev/null
then
    echo "gifsicle could not be found"
    echo "Install it with: npm install -g gifsicle"
    exit 1
fi

# get last part of the test path and set that as the output name
# example test path: pkg/integration/tests/01_basic_test.go
# For that we want: NAME=01_basic_test
NAME=$(echo "$TEST" | sed -e 's/.*\///' | sed -e 's/\..*//')

go generate pkg/integration/tests/tests.go

terminalizer -c demo/config.yml record --skip-sharing -d "go run cmd/integration_test/main.go cli --slow $TEST" "demo/output/$NAME"
terminalizer render "demo/output/$NAME" -o "demo/output/$NAME.gif"
gifsicle --colors 256 --use-col=web -O3 < "demo/output/$NAME.gif" > "demo/output/$NAME-compressed.gif"

echo "Demo recorded to demo/$NAME-compressed.gif"
