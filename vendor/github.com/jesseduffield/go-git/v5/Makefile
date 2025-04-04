# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test 

# Git config
GIT_VERSION ?=
GIT_DIST_PATH ?= $(PWD)/.git-dist
GIT_REPOSITORY = http://github.com/git/git.git

# Coverage
COVERAGE_REPORT = coverage.out
COVERAGE_MODE = count

build-git:
	@if [ -f $(GIT_DIST_PATH)/git ]; then \
		echo "nothing to do, using cache $(GIT_DIST_PATH)"; \
	else \
		git clone $(GIT_REPOSITORY) -b $(GIT_VERSION) --depth 1 --single-branch $(GIT_DIST_PATH); \
		cd $(GIT_DIST_PATH); \
		make configure; \
		./configure; \
		make all; \
	fi

test:
	@echo "running against `git version`"; \
	$(GOTEST) -race ./...
	$(GOTEST) -v _examples/common_test.go _examples/common.go --examples

TEMP_REPO := $(shell mktemp)
test-sha256:
	$(GOCMD) run -tags sha256 _examples/sha256/main.go $(TEMP_REPO)
	cd $(TEMP_REPO) && git fsck
	rm -rf $(TEMP_REPO)

test-coverage:
	@echo "running against `git version`"; \
	echo "" > $(COVERAGE_REPORT); \
	$(GOTEST) -coverprofile=$(COVERAGE_REPORT) -coverpkg=./... -covermode=$(COVERAGE_MODE) ./...

clean:
	rm -rf $(GIT_DIST_PATH)

fuzz:
	@go test -fuzz=FuzzParser				$(PWD)/internal/revision
	@go test -fuzz=FuzzDecoder				$(PWD)/plumbing/format/config
	@go test -fuzz=FuzzPatchDelta			$(PWD)/plumbing/format/packfile
	@go test -fuzz=FuzzParseSignedBytes		$(PWD)/plumbing/object
	@go test -fuzz=FuzzDecode				$(PWD)/plumbing/object
	@go test -fuzz=FuzzDecoder				$(PWD)/plumbing/protocol/packp
	@go test -fuzz=FuzzNewEndpoint			$(PWD)/plumbing/transport
