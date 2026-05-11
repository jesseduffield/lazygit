.DEFAULT_GOAL := all

.PHONY: all
all: test_with_fuzz lint

# the TEST_FLAGS env var can be set to eg run only specific tests
TEST_COMMAND = go test -v -count=1 -race -cover $(TEST_FLAGS)

.PHONY: test
test:
	$(TEST_COMMAND)

.PHONY: bench
bench:
	go test -bench=.

FUZZ_TIME ?= 10s

# see https://github.com/golang/go/issues/46312
# and https://stackoverflow.com/a/72673487/4867444
# if we end up having more fuzz tests
.PHONY: test_with_fuzz
test_with_fuzz:
	$(TEST_COMMAND) -fuzz=FuzzRoundTripJSON -fuzztime=$(FUZZ_TIME)
	$(TEST_COMMAND) -fuzz=FuzzRoundTripYAML -fuzztime=$(FUZZ_TIME)

.PHONY: fuzz
fuzz: test_with_fuzz

.PHONY: lint
lint:
	golangci-lint run
