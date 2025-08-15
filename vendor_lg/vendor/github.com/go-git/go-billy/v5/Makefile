# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test 
WASIRUN_WRAPPER := $(CURDIR)/scripts/wasirun-wrapper

.PHONY: test
test:
	$(GOTEST) -race ./...

test-coverage:
	echo "" > $(COVERAGE_REPORT); \
	$(GOTEST) -coverprofile=$(COVERAGE_REPORT) -coverpkg=./... -covermode=$(COVERAGE_MODE) ./...

.PHONY: wasitest
wasitest: export GOARCH=wasm
wasitest: export GOOS=wasip1
wasitest:
	$(GOTEST) -exec $(WASIRUN_WRAPPER) ./...
