# Only build/test/lint exp/simd when Go version is >= 1.26 (requires goexperiment.simd)
GO_VERSION := $(shell go version 2>/dev/null | sed -n 's/.*go\([0-9]*\)\.\([0-9]*\).*/\1.\2/p')
GO_SIMD_SUPPORT := $(shell ver="$(GO_VERSION)"; [ -n "$$ver" ] && [ "$$(printf '%s\n1.26\n' "$$ver" | sort -V | tail -1)" = "$$ver" ] && echo yes)

build:
	go build -v ./...
	@if [ -n "$(GO_SIMD_SUPPORT)" ]; then cd ./exp/simd && GOEXPERIMENT=simd go build -v ./; fi

test:
	go test -race ./...
	@if [ -n "$(GO_SIMD_SUPPORT)" ]; then cd ./exp/simd && GOEXPERIMENT=simd go test -race ./; fi
watch-test:
	reflex -t 50ms -s -- sh -c 'gotest -race ./...'

bench:
	go test -v -run=^Benchmark -benchmem -count 3 -bench ./...
watch-bench:
	reflex -t 50ms -s -- sh -c 'go test -v -run=^Benchmark -benchmem -count 3 -bench ./...'

coverage:
	go test -v -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html

tools:
	go install github.com/cespare/reflex@latest
	go install github.com/rakyll/gotest@latest
	go install github.com/psampaz/go-mod-outdated@latest
	go install github.com/jondot/goweight@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get -t -u golang.org/x/tools/cmd/cover
	go install github.com/sonatype-nexus-community/nancy@latest
	go install golang.org/x/perf/cmd/benchstat@latest
	go install github.com/cespare/prettybench@latest
	go mod tidy

	# brew install hougesen/tap/mdsf

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...
	@if [ -n "$(GO_SIMD_SUPPORT)" ]; then cd ./exp/simd && golangci-lint run --timeout 60s --max-same-issues 50 ./; fi
	# mdsf verify --debug --log-level warn docs/
lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...
	@if [ -n "$(GO_SIMD_SUPPORT)" ]; then cd ./exp/simd && golangci-lint run --timeout 60s --max-same-issues 50 --fix ./; fi
	# mdsf format --debug --log-level warn docs/

audit:
	go list -json -m all | nancy sleuth

outdated:
	go list -u -m -json all | go-mod-outdated -update -direct

weight:
	goweight

doc:
	cd docs && npm install && npm start
