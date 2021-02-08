# Unicode Text Segmentation for Go

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/rivo/uniseg)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/rivo/uniseg)

This Go package implements Unicode Text Segmentation according to [Unicode Standard Annex #29](http://unicode.org/reports/tr29/) (Unicode version 12.0.0).

At this point, only the determination of grapheme cluster boundaries is implemented.

## Background

In Go, [strings are read-only slices of bytes](https://blog.golang.org/strings). They can be turned into Unicode code points using the `for` loop or by casting: `[]rune(str)`. However, multiple code points may be combined into one user-perceived character or what the Unicode specification calls "grapheme cluster". Here are some examples:

|String|Bytes (UTF-8)|Code points (runes)|Grapheme clusters|
|-|-|-|-|
|Käse|6 bytes: `4b 61 cc 88 73 65`|5 code points: `4b 61 308 73 65`|4 clusters: `[4b],[61 308],[73],[65]`|
|🏳️‍🌈|14 bytes: `f0 9f 8f b3 ef b8 8f e2 80 8d f0 9f 8c 88`|4 code points: `1f3f3 fe0f 200d 1f308`|1 cluster: `[1f3f3 fe0f 200d 1f308]`|
|🇩🇪|8 bytes: `f0 9f 87 a9 f0 9f 87 aa`|2 code points: `1f1e9 1f1ea`|1 cluster: `[1f1e9 1f1ea]`|

This package provides a tool to iterate over these grapheme clusters. This may be used to determine the number of user-perceived characters, to split strings in their intended places, or to extract individual characters which form a unit.

## Installation

```bash
go get github.com/rivo/uniseg
```

## Basic Example

```go
package uniseg

import (
	"fmt"

	"github.com/rivo/uniseg"
)

func main() {
	gr := uniseg.NewGraphemes("👍🏼!")
	for gr.Next() {
		fmt.Printf("%x ", gr.Runes())
	}
	// Output: [1f44d 1f3fc] [21]
}
```

## Documentation

Refer to https://godoc.org/github.com/rivo/uniseg for the package's documentation.

## Dependencies

This package does not depend on any packages outside the standard library.

## Your Feedback

Add your issue here on GitHub. Feel free to get in touch if you have any questions.

## Version

Version tags will be introduced once Golang modules are official. Consider this version 0.1.
