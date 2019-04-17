#!/bin/sh
go build && ./codegen -cout ../pluralspec_gen.go -tout ../pluralspec_gen_test.go && \
    gofmt -w=true ../pluralspec_gen.go && \
    gofmt -w=true ../pluralspec_gen_test.go && \
    rm codegen
