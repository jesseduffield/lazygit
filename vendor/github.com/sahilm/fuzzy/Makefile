.PHONY: all
all: setup lint test

PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
test: setup
	go test $(PKGS)

sources = $(shell find . -name '*.go' -not -path './vendor/*')
.PHONY: goimports
goimports: setup
	goimports -w $(sources)

.PHONY: lint
lint: setup
	$(BIN_DIR)/golangci-lint run

COVERAGE := $(CURDIR)/coverage
COVER_PROFILE :=$(COVERAGE)/cover.out
TMP_COVER_PROFILE :=$(COVERAGE)/cover.tmp
.PHONY: cover
cover: setup
	rm -rf $(COVERAGE)
	mkdir -p $(COVERAGE)
	echo "mode: set" > $(COVER_PROFILE)
	for pkg in $(PKGS); do \
		go test -v -coverprofile=$(TMP_COVER_PROFILE) $$pkg; \
		if [ -f $(TMP_COVER_PROFILE) ]; then \
			grep -v 'mode: set' $(TMP_COVER_PROFILE) >> $(COVER_PROFILE); \
			rm $(TMP_COVER_PROFILE); \
		fi; \
	done
	go tool cover -html=$(COVER_PROFILE) -o $(COVERAGE)/index.html

.PHONY: ci
ci: setup lint test

.PHONY: install
install: setup
	go install $(PKGS)

.PHONY: build
build: setup
	go build $(PKGS)

GOPATH ?= $(HOME)/go
BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOLANG_CI_LINT := $(BIN_DIR)/golangci-lint

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOLANG_CI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(BIN_DIR) v2.12.2

tools: $(GOIMPORTS) $(GOLANG_CI_LINT)

setup: tools
