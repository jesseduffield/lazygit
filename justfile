#!/usr/bin/env just --justfile

set shell := ["pwsh", "-Command"]

alias d := default
alias b := build
alias i := install
alias r := run
alias rd := run-debug
alias pl := print-log
alias ut := unit-test
alias t := test
alias g := generate
alias f := format

default:
  @just --list

build target='main':
  go build -gcflags='all=-N -l' {{target}}

install target='main':
  go install

run:
  build
  ./lazytask

run-debug:
  go run main.go -debug

print-log:
  go run main.go --logs

unit-test:
  go test ./... -short

# Run all tests, including unit and integration tests
test:
  unit-test
  integration-test-all



# Generate auto-generated files
generate:
  go generate ./...



# Format the project using gofumpt
format:
  gofumpt -l -w .



# Run TUI integration tests
integration-test-tui ARGS:
  go run cmd/integration_test/main.go tui {{ARGS}}

# Run CLI integration tests
integration-test-cli ARGS:
  go run cmd/integration_test/main.go cli {{ARGS}}

# Run all integration tests
integration-test-all:
  go test pkg/integration/clients/*.go

# Bump gocui version
bump-gocui:
  scripts/bump_gocui.sh

# Bump lazycore version
bump-lazycore:
  scripts/bump_lazycore.sh

# Record a demo
record-demo ARGS:
  demo/record_demo.sh {{ARGS}}

# Vendor dependencies
vendor:
  go mod vendor && go mod tidy
