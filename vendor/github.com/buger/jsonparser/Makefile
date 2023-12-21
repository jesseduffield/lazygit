SOURCE = parser.go
CONTAINER = jsonparser
SOURCE_PATH = /go/src/github.com/buger/jsonparser
BENCHMARK = JsonParser
BENCHTIME = 5s
TEST = .
DRUN = docker run -v `pwd`:$(SOURCE_PATH) -i -t $(CONTAINER)

build:
	docker build -t $(CONTAINER) .

race:
	$(DRUN) --env GORACE="halt_on_error=1" go test ./. $(ARGS) -v -race -timeout 15s

bench:
	$(DRUN) go test $(LDFLAGS) -test.benchmem -bench $(BENCHMARK) ./benchmark/ $(ARGS) -benchtime $(BENCHTIME) -v

bench_local:
	$(DRUN) go test $(LDFLAGS) -test.benchmem -bench . $(ARGS) -benchtime $(BENCHTIME) -v

profile:
	$(DRUN) go test $(LDFLAGS) -test.benchmem -bench $(BENCHMARK) ./benchmark/ $(ARGS) -memprofile mem.mprof -v
	$(DRUN) go test $(LDFLAGS) -test.benchmem -bench $(BENCHMARK) ./benchmark/ $(ARGS) -cpuprofile cpu.out -v
	$(DRUN) go test $(LDFLAGS) -test.benchmem -bench $(BENCHMARK) ./benchmark/ $(ARGS) -c

test:
	$(DRUN) go test $(LDFLAGS) ./ -run $(TEST) -timeout 10s $(ARGS) -v

fmt:
	$(DRUN) go fmt ./...

vet:
	$(DRUN) go vet ./.

bash:
	$(DRUN) /bin/bash