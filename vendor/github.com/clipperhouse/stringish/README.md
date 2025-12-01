# stringish

A small Go module that provides a generic type constraint for “string-like”
data, and a utf8 package that works with both strings and byte slices
without conversions.

```go
type Interface interface {
	~[]byte | ~string
}
```

[![Go Reference](https://pkg.go.dev/badge/github.com/clipperhouse/stringish/utf8.svg)](https://pkg.go.dev/github.com/clipperhouse/stringish/utf8)
[![Test Status](https://github.com/clipperhouse/stringish/actions/workflows/gotest.yml/badge.svg)](https://github.com/clipperhouse/stringish/actions/workflows/gotest.yml)

## Install

```
go get github.com/clipperhouse/stringish
```

## Examples

```go
import (
    "github.com/clipperhouse/stringish"
    "github.com/clipperhouse/stringish/utf8"
)

s := "Hello, 世界"
r, size := utf8.DecodeRune(s)   // not DecodeRuneInString 🎉

b := []byte("Hello, 世界")
r, size = utf8.DecodeRune(b)    // same API!

func MyFoo[T stringish.Interface](s T) T {
    // pass a string or a []byte
    // iterate, slice, transform, whatever
}
```

## Motivation

Sometimes we want APIs to accept `string` or `[]byte` without having to convert
between those types. That conversion usually allocates!

By implementing with `stringish.Interface`, we can have a single API, and
single implementation for both types: one `Foo` instead of `Foo` and
`FooString`.

We have converted the
[`unicode/utf8` package](https://github.com/clipperhouse/stringish/blob/main/utf8/utf8.go)
as an example -- note the absence of`*InString` funcs. We might look at `x/text`
next.

## Used by

- clipperhouse/uax29: [stringish trie](https://github.com/clipperhouse/uax29/blob/master/graphemes/trie.go#L27), [stringish iterator](https://github.com/clipperhouse/uax29/blob/master/internal/iterators/iterator.go#L9), [stringish SplitFunc](https://github.com/clipperhouse/uax29/blob/master/graphemes/splitfunc.go#L21)

- [clipperhouse/displaywidth](https://github.com/clipperhouse/displaywidth)

## Prior discussion

- [Consideration of similar by the Go team](https://github.com/golang/go/issues/48643)
