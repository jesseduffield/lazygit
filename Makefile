.PHONY: all
all: build

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install

.PHONY: run
run:
	go run main.go

# Run `make run-debug` in one terminal tab and `make print-log` in another to view the program and its log output side by side
.PHONY: run-debug
run-debug:
	go run main.go -debug

.PHONY: print-log
print-log:
	go run main.go --logs

.PHONY: unit-test
unit-test:
	go test ./... -short

.PHONY: test
test: unit-test integration-test-all

.PHONY: generate
generate:
	go generate ./...

.PHONY: format
format:
	gofumpt -l -w .

.PHONY: update-cheatsheet
update-cheatsheet:
	go run scripts/cheatsheet/main.go generate

# For more details about integration test, see https://github.com/jesseduffield/lazygit/blob/master/pkg/integration/README.md.
.PHONY: integration-test-tui
integration-test-tui:
	go run cmd/integration_test/main.go tui

.PHONY: integration-test-cli
integration-test-cli:
	go run cmd/integration_test/main.go cli $(filter-out $@,$(MAKECMDGOALS))

.PHONY: integration-test-all
integration-test-all:
	go test pkg/integration/clients/*.go

.PHONY: bump-gocui
bump-gocui:
	scripts/bump_gocui.sh

.PHONY: bump-lazycore
bump-lazycore:
	scripts/bump_lazycore.sh

.PHONY: record-demo
record-demo:
	demo/record_demo.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: vendor
vendor:
	go mod vendor && go mod tidy
