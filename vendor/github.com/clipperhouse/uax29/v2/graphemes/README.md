An implementation of grapheme cluster boundaries from [Unicode text segmentation](https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries) (UAX 29), for Unicode 17.

[![Documentation](https://pkg.go.dev/badge/github.com/clipperhouse/uax29/v2/graphemes.svg)](https://pkg.go.dev/github.com/clipperhouse/uax29/v2/graphemes)
![Tests](https://github.com/clipperhouse/uax29/actions/workflows/gotest.yml/badge.svg)
![Fuzz](https://github.com/clipperhouse/uax29/actions/workflows/gofuzz.yml/badge.svg)

## Quick start

```
go get github.com/clipperhouse/uax29/v2/graphemes
```

```go
import "github.com/clipperhouse/uax29/v2/graphemes"

text := "Hello, 世界. Nice dog! 👍🐶"
g := graphemes.FromString(text)

for g.Next() {                     // Next() returns true until end of data
	fmt.Println(g.Value())         // Do something with the current grapheme
}
```

_A grapheme is a “single visible character”, which might be a simple as a single letter, or a complex emoji that consists of several Unicode code points._

## Conformance

We use the Unicode [test suite](https://unicode.org/reports/tr41/tr41-36.html#Tests29).

![Tests](https://github.com/clipperhouse/uax29/actions/workflows/gotest.yml/badge.svg)
![Fuzz](https://github.com/clipperhouse/uax29/actions/workflows/gofuzz.yml/badge.svg)

## APIs

### If you have a `string`

```go
text := "Hello, 世界. Nice dog! 👍🐶"
g := graphemes.FromString(text)

for g.Next() {                     // Next() returns true until end of data
	fmt.Println(g.Value())         // Do something with the current grapheme
}
```

### If you have an `io.Reader`

`FromReader` embeds a [`bufio.Scanner`](https://pkg.go.dev/bufio#Scanner), so just use those methods.

```go
r := getYourReader()                    // from a file or network maybe
g := graphemes.FromReader(r)

for g.Scan() {                         // Scan() returns true until error or EOF
	fmt.Println(g.Text())              // Do something with the current grapheme
}

if g.Err() != nil {                    // Check the error
	log.Fatal(g.Err())
}
```

### If you have a `[]byte`

```go
b := []byte("Hello, 世界. Nice dog! 👍🐶")

g := graphemes.FromBytes(b)

for g.Next() {                     // Next() returns true until end of data
	fmt.Println(g.Value())         // Do something with the current grapheme
}
```

### ANSI escape sequences

By the UAX 29 specification, ANSI escape sequences are not grapheme clusters. To treat 7-bit ANSI escape sequences as a single cluster, set `AnsiEscapeSequences` to true.

```go
text := "Hello, \x1b[31mworld\x1b[0m!"
g := graphemes.FromString(text)
g.AnsiEscapeSequences = true

for g.Next() {
	fmt.Println(g.Value())
}
```

To also parse 8-bit C1 controls (non-UTF-8 bytes), set `AnsiEscapeSequences8Bit` to true.

```go
g.AnsiEscapeSequences = true     // 7-bit forms (ESC ...)
g.AnsiEscapeSequences8Bit = true // 8-bit C1 forms (0x80-0x9F), not valid UTF-8
```

For ESC-initiated (7-bit) control strings, only 7-bit terminators are recognized.
For C1-initiated (8-bit) control strings, only C1 ST (`0x9C`) is recognized as ST.

We implement [ECMA-48](https://ecma-international.org/publications-and-standards/standards/ecma-48/) control codes in both 7-bit and 8-bit representations. 8-bit control codes are not UTF-8 encoded and are not valid UTF-8, caveat emptor.

### Benchmarks

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/uax29/graphemes/comparative
cpu: Apple M2

BenchmarkGraphemesMixed/clipperhouse/uax29-8  	    142635 ns/op	 245.12 MB/s    0 B/op	   0 allocs/op
BenchmarkGraphemesMixed/rivo/uniseg-8         	   2018284 ns/op	  17.32 MB/s    0 B/op	   0 allocs/op

BenchmarkGraphemesASCII/clipperhouse/uax29-8  	      8846 ns/op	 508.73 MB/s    0 B/op	   0 allocs/op
BenchmarkGraphemesASCII/rivo/uniseg-8         	    366760 ns/op	  12.27 MB/s    0 B/op	   0 allocs/op
```

### Invalid inputs

Invalid UTF-8 input is considered undefined behavior. We test to ensure that bad inputs will not cause pathological outcomes, such as a panic or infinite loop. Callers should expect “garbage-in, garbage-out”.

Your pipeline should probably include a call to [`utf8.Valid()`](https://pkg.go.dev/unicode/utf8#Valid).
