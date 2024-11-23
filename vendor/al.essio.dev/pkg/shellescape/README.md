![Build](https://github.com/alessio/shellescape/workflows/Build/badge.svg)
[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/alessio/shellescape?tab=overview)
[![sourcegraph](https://sourcegraph.com/github.com/alessio/shellescape/-/badge.svg)](https://sourcegraph.com/github.com/alessio/shellescape)
[![codecov](https://codecov.io/gh/alessio/shellescape/branch/master/graph/badge.svg)](https://codecov.io/gh/alessio/shellescape)
[![Coverage](https://gocover.io/_badge/github.com/alessio/shellescape)](https://gocover.io/github.com/alessio/shellescape)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessio/shellescape)](https://goreportcard.com/report/github.com/alessio/shellescape)

# shellescape
Escape arbitrary strings for safe use as command line arguments.
## Contents of the package

This package provides the `shellescape.Quote()` function that returns a
shell-escaped copy of a string. This functionality could be helpful
in those cases where it is known that the output of a Go program will
be appended to/used in the context of shell programs' command line arguments.

This work was inspired by the Python original package
[shellescape](https://pypi.python.org/pypi/shellescape).

## Usage

The following snippet shows a typical unsafe idiom:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("ls -l %s\n", os.Args[1])
}
```
_[See in Go Playground](https://play.golang.org/p/Wj2WoUfH_d)_

Especially when creating pipeline of commands which might end up being
executed by a shell interpreter, it is particularly unsafe to not
escape arguments.

`shellescape.Quote()` comes in handy and to safely escape strings:

```go
package main

import (
        "fmt"
        "os"

        "al.essio.dev/pkg/shellescape"
)

func main() {
        fmt.Printf("ls -l %s\n", shellescape.Quote(os.Args[1]))
}
```
_[See in Go Playground](https://go.dev/play/p/GeguukpSUTk)_

## The escargs utility
__escargs__ reads lines from the standard input and prints shell-escaped versions. Unlinke __xargs__, blank lines on the standard input are not discarded.
