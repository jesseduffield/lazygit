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
[unix]
test: unit-test e2e

# On Windows, integration tests are not supported right now
[windows]
test: unit-test

# Generate all our auto-generated files (test list, cheatsheets, json schema, maybe other things in the future)
generate:
    go generate ./...

format:
    go tool gofumpt -l -w .

lint:
    ./scripts/gofumpt-check.sh
    ./scripts/golangci-lint-shim.sh run

e2e-test-command := "go test pkg/integration/clients/*.go"

# Run integration tests headlessly: no args runs all tests, a test name (or path) runs just that one. Use e2e-cli for a visible UI.
e2e *args:
    {{ if args == "" { e2e-test-command } else { \
        e2e-test-command + " -run 'TestIntegration/" + \
        replace( \
            replace_regex( \
                replace_regex(args, '\S*pkg/integration/tests/', ''), \
                '\.go( |$)', '${1}' \
            ), \
            " ", "$' && " + e2e-test-command + " -run 'TestIntegration/" \
        ) + "$'" \
    } }}

# Run a single integration test with a visible UI; most useful with --sandbox or --slow.
e2e-cli *args:
    go run cmd/integration_test/main.go cli {{ args }}

# Open the TUI for running integration tests.
e2e-tui *args:
    go run cmd/integration_test/main.go tui {{ args }}

# Run some tests on the current commit, similar to what CI does.
check:
    ./scripts/check_commit.sh

bump-gocui:
    scripts/bump_gocui.sh

# Record a demo
demo *args:
    demo/record_demo.sh {{ args }}

vendor:
    go mod tidy && go mod vendor
