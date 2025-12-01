An implementation of grapheme cluster boundaries from [Unicode text segmentation](https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries) (UAX 29), for Unicode version 15.0.0.

[![Documentation](https://pkg.go.dev/badge/github.com/clipperhouse/uax29/v2/graphemes.svg)](https://pkg.go.dev/github.com/clipperhouse/uax29/v2/graphemes)
![Tests](https://github.com/clipperhouse/uax29/actions/workflows/gotest.yml/badge.svg)
![Fuzz](https://github.com/clipperhouse/uax29/actions/workflows/gofuzz.yml/badge.svg)

## Quick start

```
go get "github.com/clipperhouse/uax29/v2/graphemes"
```

```go
import "github.com/clipperhouse/uax29/v2/graphemes"

text := "Hello, 世界. Nice dog! 👍🐶"

tokens := graphemes.FromString(text)

for tokens.Next() {                     // Next() returns true until end of data
	fmt.Println(tokens.Value())         // Do something with the current grapheme
}
```

_A grapheme is a “single visible character”, which might be a simple as a single letter, or a complex emoji that consists of several Unicode code points._

## Conformance

We use the Unicode [test suite](https://unicode.org/reports/tr41/tr41-26.html#Tests29).

![Tests](https://github.com/clipperhouse/uax29/actions/workflows/gotest.yml/badge.svg)
![Fuzz](https://github.com/clipperhouse/uax29/actions/workflows/gofuzz.yml/badge.svg)

## APIs

### If you have a `string`

```go
text := "Hello, 世界. Nice dog! 👍🐶"

tokens := graphemes.FromString(text)

for tokens.Next() {                     // Next() returns true until end of data
	fmt.Println(tokens.Value())         // Do something with the current grapheme
}
```

### If you have an `io.Reader`

`FromReader` embeds a [`bufio.Scanner`](https://pkg.go.dev/bufio#Scanner), so just use those methods.

```go
r := getYourReader()                        // from a file or network maybe
tokens := graphemes.FromReader(r)

for tokens.Scan() {                         // Scan() returns true until error or EOF
	fmt.Println(tokens.Text())              // Do something with the current grapheme
}

if tokens.Err() != nil {                    // Check the error
	log.Fatal(tokens.Err())
}
```

### If you have a `[]byte`

```go
b := []byte("Hello, 世界. Nice dog! 👍🐶")

tokens := graphemes.FromBytes(b)

for tokens.Next() {                     // Next() returns true until end of data
	fmt.Println(tokens.Value())         // Do something with the current grapheme
}
```

### Benchmarks

On a Mac M2 laptop, we see around 200MB/s, or around 100 million graphemes per second, and no allocations.

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/uax29/graphemes/comparative
cpu: Apple M2
BenchmarkGraphemes/clipperhouse/uax29-8    	    173805 ns/op	 201.16 MB/s      0 B/op	   0 allocs/op
BenchmarkGraphemes/rivo/uniseg-8           	   2045128 ns/op	  17.10 MB/s      0 B/op	   0 allocs/op
```

### Invalid inputs

Invalid UTF-8 input is considered undefined behavior. We test to ensure that bad inputs will not cause pathological outcomes, such as a panic or infinite loop. Callers should expect “garbage-in, garbage-out”.

Your pipeline should probably include a call to [`utf8.Valid()`](https://pkg.go.dev/unicode/utf8#Valid).
