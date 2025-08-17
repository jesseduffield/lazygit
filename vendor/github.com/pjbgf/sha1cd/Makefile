FUZZ_TIME ?= 1m

export CGO_ENABLED := 1

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -benchmem -run=^$$ -bench ^Benchmark ./...

.PHONY: fuzz
fuzz:
	go test -tags gofuzz -fuzz=. -fuzztime=$(FUZZ_TIME) ./test/

# Cross build project in arm/v7.
build-arm:
	docker build -t sha1cd-arm -f Dockerfile.arm .
	docker run --rm sha1cd-arm

# Cross build project in arm64.
build-arm64:
	docker build -t sha1cd-arm64 -f Dockerfile.arm64 .
	docker run --rm sha1cd-arm64

# Build with cgo disabled.
build-nocgo:
	CGO_ENABLED=0 go build ./cgo

# Run cross-compilation to assure supported architectures.
cross-build: build-arm build-arm64 build-nocgo

generate:
	go generate -x ./...

verify: generate
	git diff --exit-code
	go vet ./...
