<div align="center">
  <a href="#installing" title="Installing">Installing</a> |
  <a href="#using" title="Using">Using</a> |
  <a href="#the-example" title="The example">Example</a> |
  <a href="#infocmp" title="infocmp">infocmp</a> |
  <a href="https://pkg.go.dev/github.com/xo/terminfo" title="Go Reference">Reference</a> |
  <a href="#documents" title="Documents">Documents</a> |
  <a href="#about" title="About">About</a>
</div>

<br/>

[![Unit Tests][terminfo-ci-status]][terminfo-ci]
[![Go Reference][goref-terminfo-status]][goref-terminfo]
[![Releases][release-status]][releases]
[![Discord Discussion][discord-status]][discord]

[terminfo-ci]: https://github.com/xo/terminfo/actions/workflows/test.yml "Test CI"
[terminfo-ci-status]: https://github.com/xo/terminfo/actions/workflows/test.yml/badge.svg "Test CI"
[goref-terminfo]: https://pkg.go.dev/github.com/xo/terminfo "Go Reference"
[goref-terminfo-status]: https://pkg.go.dev/badge/github.com/xo/terminfo.svg "Go Reference"
[release-status]: https://img.shields.io/github/v/release/xo/terminfo?display_name=tag "Latest Release"
[releases]: https://github.com/xo/terminfo/releases "Releases"
[discord]: https://discord.gg/WDWAgXwJqN "Discord Discussion"
[discord-status]: https://img.shields.io/discord/829150509658013727.svg?label=Discord&logo=Discord&colorB=7289da&style=flat-square "Discord Discussion"

# terminfo

`terminfo` reads the terminfo database in Go. The terminfo database describes
what each kind of terminal can do, and the byte sequences that make it do
each thing. A program uses those sequences to move the cursor, set colors,
clear the screen or flash it.

This package is written in pure Go. It does not use cgo, and it does not link
the `ncurses` library. It replaces `ncurses` in a Go program that needs the
terminal capabilities and not a full screen library.

The package can:

- Find the entry for a terminal name in the directories that terminfo(5)
  names, and cache it.
- Decode the legacy format and the 32-bit number format of `ncurses` 6.1, with
  the extended capabilities.
- Evaluate a parameterized string, such as the value of `cup`.
- Write a capability with its padding applied, as `tputs` does.
- Find the color level of the terminal from the environment.

The system `infocmp` is the reference. The tests compare each decoded value
with the output of `infocmp` for every entry on the machine.

## Installing

Install `terminfo` with the [Go][go-project] tool:

```sh
go get github.com/xo/terminfo@latest
```

## Using

`LoadFromEnv` loads the entry that `$TERM` names. `Printf` evaluates a
capability with its parameters, and `Puts` writes a capability with its
padding applied:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xo/terminfo"
)

func main() {
	// load the entry that $TERM names
	ti, err := terminfo.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// move the cursor to row 5, column 10, and write in bold
	fmt.Print(ti.Printf(terminfo.CursorAddress, 5, 10))
	fmt.Print(ti.Printf(terminfo.EnterBoldMode) + "hello" + ti.Printf(terminfo.ExitAttributeMode))
	// flash the screen, with the delay that the entry asks for
	if _, err := ti.Puts(os.Stdout, ti.Printf(terminfo.FlashScreen), 1, 38400); err != nil {
		log.Fatal(err)
	}
}
```

A capability that holds a padding marker, such as the `$<100/>` in the
`flash` of `xterm`, must go through `Puts`. If you write the value of
`Strings` or of `Printf` directly, the terminal shows the marker.

`ColorLevelFromEnv` returns the color level of the terminal. It reads
`COLORTERM`, `FORCE_COLOR` and `TERM_PROGRAM` first, then the `max_colors`
capability of the entry.

`Load` and `Open` return the same `*Terminfo` to every caller. Do not change
the returned value or its maps.

## The example

[`_examples/main.go`](_examples/main.go) draws the color table of the
terminal and waits for Ctrl-C. Run it in the repository root:

```sh
go run ./_examples
```

## infocmp

[`cmd/infocmp`](cmd/infocmp) prints an entry in the same form as
`infocmp -1 -L`. It shows what the package decodes. Install it and run it:

```sh
go install github.com/xo/terminfo/cmd/infocmp@latest
TERM=xterm-256color infocmp -x
```

## Documents

- [CONTRIBUTING.md](CONTRIBUTING.md): how to change the project, the checks
  to run, and how to install the agent skills.
- [AGENTS.md](AGENTS.md): the rules for a coding agent. `CLAUDE.md` imports
  it.
- [docs/PLAN.md](docs/PLAN.md): the plan, every decision, and the open
  questions.
- [docs/BACKLOG.md](docs/BACKLOG.md): the work that is known and not done.

## About

`terminfo` supports [rline][rline], the line editor of [usql][usql].

<br/>

<div align="center">
  <a href="https://github.com/xo/usql" title="A command line client for many databases">usql</a> |
  <a href="https://github.com/xo/dburl" title="Database connection URLs">dburl</a> |
  <a href="https://github.com/xo/dbmeta" title="Database metadata">dbmeta</a> |
  <a href="https://github.com/xo/dbimp" title="Database drivers in pure Go">dbimp</a> |
  <a href="https://github.com/xo/cassandra" title="A database/sql driver for Cassandra">cassandra</a> |
  <a href="https://github.com/xo/dbtpl" title="Go code generated from a database">dbtpl</a> |
  <a href="https://github.com/xo/tblfmt" title="Tables of database results">tblfmt</a> |
  <a href="https://github.com/xo/rline" title="The line editor of usql">rline</a> |
  <a href="https://github.com/xo/terminfo" title="The terminfo database in pure Go, this project">terminfo</a> |
  <a href="https://github.com/xo/transit" title="tree-sitter in pure Go">transit</a>
</div>

[go-project]: https://go.dev/project
[rline]: https://github.com/xo/rline
[usql]: https://github.com/xo/usql
