<img src="assets/search-gopher-1.png" alt="gopher looking for stuff">  <img src="assets/search-gopher-2.png" alt="gopher found stuff">

# fuzzy
[![Build Status](https://travis-ci.org/sahilm/fuzzy.svg?branch=master)](https://travis-ci.org/sahilm/fuzzy)
[![Documentation](https://godoc.org/github.com/sahilm/fuzzy?status.svg)](https://godoc.org/github.com/sahilm/fuzzy)

Go library that provides fuzzy string matching optimized for filenames and code symbols in the style of Sublime Text, 
VSCode, IntelliJ IDEA et al. This library is external dependency-free. It only depends on the Go standard library.

## Features

- Intuitive matching. Results are returned in descending order of match quality. Quality is determined by:
  - The first character in the pattern matches the first character in the match string.
  - The matched character is camel cased.
  - The matched character follows a separator such as an underscore character.
  - The matched character is adjacent to a previous match.

- Speed. Matches are returned in milliseconds. It's perfect for interactive search boxes.

- The positions of matches are returned. Allows you to highlight matching characters.

- Unicode aware.

## Demo

Here is a [demo](_example/main.go) of matching various patterns against ~16K files from the Unreal Engine 4 codebase.

![demo](assets/demo.gif)

You can run the demo yourself like so:

```
cd _example/
go get github.com/jroimartin/gocui
go run main.go
```

## Usage

The following example prints out matches with the matched chars in bold.

```go
package main

import (
	"fmt"

	"github.com/sahilm/fuzzy"
)

func main() {
	const bold = "\033[1m%s\033[0m"
	pattern := "mnr"
	data := []string{"game.cpp", "moduleNameResolver.ts", "my name is_Ramsey"}

	matches := fuzzy.Find(pattern, data)

	for _, match := range matches {
		for i := 0; i < len(match.Str); i++ {
			if contains(i, match.MatchedIndexes) {
				fmt.Print(fmt.Sprintf(bold, string(match.Str[i])))
			} else {
				fmt.Print(string(match.Str[i]))
			}
		}
		fmt.Println()
	}
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}
``` 
If the data you want to match isn't a slice of strings, you can use `FindFrom` by implementing
the provided `Source` interface. Here's an example:

```go
package main

import (
	"fmt"

	"github.com/sahilm/fuzzy"
)

type employee struct {
	name string
	age  int
}

type employees []employee

func (e employees) String(i int) string {
	return e[i].name
}

func (e employees) Len() int {
	return len(e)
}

func main() {
	emps := employees{
		{
			name: "Alice",
			age:  45,
		},
		{
			name: "Bob",
			age:  35,
		},
		{
			name: "Allie",
			age:  35,
		},
	}
	results := fuzzy.FindFrom("al", emps)
	for _, r := range results {
		fmt.Println(emps[r.Index])
	}
}
```

Check out the [godoc](https://godoc.org/github.com/sahilm/fuzzy) for detailed documentation.

## Installation

`go get github.com/sahilm/fuzzy` or use your favorite dependency management tool.

## Speed

Here are a few benchmark results on a normal laptop.

```
BenchmarkFind/with_unreal_4_(~16K_files)-4         	     100	  12915315 ns/op
BenchmarkFind/with_linux_kernel_(~60K_files)-4     	      50	  30885038 ns/op
```

Matching a pattern against ~60K files from the Linux kernel takes about 30ms.

## Contributing

Everyone is welcome to contribute. Please send me a pull request or file an issue. I promise
to respond promptly.

## Credits

* [@ericpauley](https://github.com/ericpauley) & [@lunixbochs](https://github.com/lunixbochs) contributed Unicode awareness and various performance optimisations.

* The algorithm is based of the awesome work of [forrestthewoods](https://github.com/forrestthewoods/lib_fts/blob/master/code/fts_fuzzy_match.js). 
See [this](https://blog.forrestthewoods.com/reverse-engineering-sublime-text-s-fuzzy-match-4cffeed33fdb#.d05n81yjy)
blog post for details of the algorithm.

* The artwork is by my lovely wife Sanah. It's based on the Go Gopher.

* The Go gopher was designed by Renee French (http://reneefrench.blogspot.com/). 
The design is licensed under the Creative Commons 3.0 Attributions license.

## License

The MIT License (MIT)

Copyright (c) 2017 Sahil Muthoo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

