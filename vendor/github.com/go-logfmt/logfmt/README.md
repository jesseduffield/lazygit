[![GoDoc](https://godoc.org/github.com/go-logfmt/logfmt?status.svg)](https://godoc.org/github.com/go-logfmt/logfmt)
[![Go Report Card](https://goreportcard.com/badge/go-logfmt/logfmt)](https://goreportcard.com/report/go-logfmt/logfmt)
[![TravisCI](https://travis-ci.org/go-logfmt/logfmt.svg?branch=master)](https://travis-ci.org/go-logfmt/logfmt)
[![Coverage Status](https://coveralls.io/repos/github/go-logfmt/logfmt/badge.svg?branch=master)](https://coveralls.io/github/go-logfmt/logfmt?branch=master)

# logfmt

Package logfmt implements utilities to marshal and unmarshal data in the [logfmt
format](https://brandur.org/logfmt). It provides an API similar to
[encoding/json](http://golang.org/pkg/encoding/json/) and
[encoding/xml](http://golang.org/pkg/encoding/xml/).

The logfmt format was first documented by Brandur Leach in [this
article](https://brandur.org/logfmt). The format has not been formally
standardized. The most authoritative public specification to date has been the
documentation of a Go Language [package](http://godoc.org/github.com/kr/logfmt)
written by Blake Mizerany and Keith Rarick.

## Goals

This project attempts to conform as closely as possible to the prior art, while
also removing ambiguity where necessary to provide well behaved encoder and
decoder implementations.

## Non-goals

This project does not attempt to formally standardize the logfmt format. In the
event that logfmt is standardized this project would take conforming to the
standard as a goal.

## Versioning

Package logfmt publishes releases via [semver](http://semver.org/) compatible Git tags prefixed with a single 'v'.
