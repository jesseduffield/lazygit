
FROM golang:1.18

WORKDIR /go/src/github.com/samber/lo

COPY Makefile go.* ./

RUN make tools
