
#!/usr/bin/make -f

VERSION := $(shell git describe)

all: build

build:
	go build -a -v

install:
	go install ./cmd/escargs

escargs: build
	go build -v \
          -ldflags="-X 'main.version=$(VERSION)'" \
          ./cmd/escargs

clean:
	rm -rfv escargs

uninstall:
	rm -v $(shell go env GOPATH)/bin/escargs

.PHONY: build clean install uninstall
