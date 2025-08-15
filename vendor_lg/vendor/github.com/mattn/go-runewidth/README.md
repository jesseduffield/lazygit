go-runewidth
============

[![Build Status](https://github.com/mattn/go-runewidth/workflows/test/badge.svg?branch=master)](https://github.com/mattn/go-runewidth/actions?query=workflow%3Atest)
[![Codecov](https://codecov.io/gh/mattn/go-runewidth/branch/master/graph/badge.svg)](https://codecov.io/gh/mattn/go-runewidth)
[![GoDoc](https://godoc.org/github.com/mattn/go-runewidth?status.svg)](http://godoc.org/github.com/mattn/go-runewidth)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattn/go-runewidth)](https://goreportcard.com/report/github.com/mattn/go-runewidth)

Provides functions to get fixed width of the character or string.

Usage
-----

```go
runewidth.StringWidth("つのだ☆HIRO") == 12
```


Author
------

Yasuhiro Matsumoto

License
-------

under the MIT License: http://mattn.mit-license.org/2013
