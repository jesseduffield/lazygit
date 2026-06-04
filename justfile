default:
    just --list

# Build lazygit with optimizations disabled (to make debugging easier).
build:
    go build -gcflags='all=-N -l'

install:
    go install

run: build
    ./lazygit

# Run `just debug` in one terminal tab and `just print-log` in another to view the program and its log output side by side
debug: build
    ./lazygit -debug

print-log: build
    ./lazygit --logs

unit-test:
    go test ./... -short

# Run both unit tests and integration tests.
test: unit-test e2e-all

# Generate all our auto-generated files (test list, cheatsheets, json schema, maybe other things in the future)
generate:
    go generate ./...

format:
    gofumpt -l -w .

lint:
    ./scripts/golangci-lint-shim.sh run

# Run integration tests with a visible UI. Most useful for running a single test; for running all tests, use `e2e-all` instead.
e2e *args:
    go run cmd/integration_test/main.go cli {{ args }}

# Open the TUI for running integration tests.
e2e-tui *args:
    go run cmd/integration_test/main.go tui {{ args }}

# Run all integration tests headlessly (without a visible UI).
e2e-all:
    go test pkg/integration/clients/*.go

bump-gocui:
    scripts/bump_gocui.sh

# Record a demo
demo *args:
    demo/record_demo.sh {{ args }}

vendor:
    go mod vendor && go mod tidy
