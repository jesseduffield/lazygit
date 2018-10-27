SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

test: dep
	dep ensure
	go test .

race: dep
	dep ensure
	go test -race .


dep:
ifeq (, $(shell which dep))
	go get -u github.com/golang/dep/cmd/dep
endif

all: prelog cmdtmpl procmon

prelog: $(SOURCES)
	go build ./cmd/prelog

cmdtmpl: $(SOURCES)
	go build ./cmd/cmdtmpl

procmon: $(SOURCES)
	go build ./cmd/procmon

clean:
	rm -f prelog
	rm -f cmdtmpl
	rm -f procmon
