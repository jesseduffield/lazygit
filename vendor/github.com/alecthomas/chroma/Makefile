.PHONY: chromad upload all

VERSION ?= $(shell git describe --tags --dirty  --always)

all: README.md tokentype_string.go

README.md: lexers/*/*.go
	./table.py

tokentype_string.go: types.go
	go generate

chromad:
	rm -f chromad
	(export CGOENABLED=0 GOOS=linux GOARCH=amd64; cd ./cmd/chromad && go build -ldflags="-X 'main.version=$(VERSION)'" -o ../../chromad .)

upload: chromad
	scp chromad root@swapoff.org: && \
		ssh root@swapoff.org 'install -m755 ./chromad /srv/http/swapoff.org/bin && service chromad restart'
