# terminfo

`terminfo` reads the terminfo database in pure Go. The terminfo database
describes what each kind of terminal can do, and it holds one compiled file
for each terminal name. The package finds the file for a terminal name,
decodes it, and evaluates its parameterized strings. It is meant as a
replacement for `ncurses` in simple Go programs.

The module is one package, `github.com/xo/terminfo`, and one command,
`cmd/infocmp`. It has no cgo and it does not link `ncurses`.

## Standing rules

These rules hold in every `xo` repository, for every coding agent (D22).

1. Stage changes for review. Commit and push only when Ken says so.
2. Load the `simple-english` skill before you write any text that a person
   reads: project documentation, a code comment, an error message or a
   commit message.
3. In a Go project, load the `go-pedantry` skill before you write or review
   Go code. A rule in this file wins where the two conflict.

`CLAUDE.md` holds one line that imports this file, so that Claude Code and
every other agent read the same rules. Edit this file, not that one.

## Which document to read

| If you are | Read |
| --- | --- |
| asking why something is the way it is | the table at the top of [docs/PLAN.md](docs/PLAN.md) |
| about to decide something that is not written down | "Open questions for Ken" at the end of [docs/PLAN.md](docs/PLAN.md). Ask Ken |
| looking for work that is known and not done | [docs/BACKLOG.md](docs/BACKLOG.md) |
| changing how a file is decoded | `dec.go` and `terminfo.go`, then D3, D13 and D15 |
| changing how a parameterized string is evaluated | `param.go` and `stack.go`, then D9 and D16 |
| changing a capability name or constant | `gen.go`, then D4 and D26. Never edit `capvals.go` by hand |
| changing color detection | `color.go`, then D10 and D24 |
| changing how padding is written | `Puts` in `terminfo.go`, then D23 |
| writing a document, a code comment, an error message or a commit message | the `simple-english` skill. Load it first |
| writing or reviewing Go code | the `go-pedantry` skill. Load it first |
| installing or updating an agent skill | [CONTRIBUTING.md](CONTRIBUTING.md), under Agent skills |
| looking for what the library does | [README.md](README.md) |

A document that is not in this table does not exist. If you cannot find where
something is written down, it is not written down. Ask Ken. Do not decide it
yourself, and do not write it as though it were settled.

## Decision numbers

A bare number, such as D3, names a decision of this repository. The decisions
are in [docs/PLAN.md](docs/PLAN.md), and the table at the top of that file
lists each one with its status. A decision of another repository names that
repository first, such as "dbmeta D110" or "tblfmt D37".

## Layout

- `terminfo.go` holds the `Terminfo` type, `Decode`, `Open`, the error
  constants, and the methods that read and print capabilities. `Puts`
  writes a capability with its padding applied (D23).
- `dec.go` holds the decoder for the compiled file format.
- `load.go` holds `Load` and `LoadFromEnv`. `Load` searches the directories
  that terminfo(5) names, and caches each entry that it opens (D5). A caller
  must not change a cached entry (D27).
- `param.go` and `stack.go` evaluate a parameterized string, such as the
  value of `cup`. `Printf` and `Fprintf` are the exported entry points.
- `color.go` holds `ColorLevel` and `ColorLevelFromEnv` (D10).
- `caps.go` returns the long and short name of each capability.
- `capvals.go` holds the capability constants and names. `gen.go` writes it
  from the `Caps` file of `ncurses` 6.6 (D4, D26). The file says
  "DO NOT EDIT", and names the release.
- `gen.go` is the generator. It has the `ignore` build tag, so the package
  build skips it.
- `cmd/infocmp` prints an entry in the form of `infocmp -1 -L`.
  `cmd/infocmp/test.sh` compares its output with the system `infocmp` for
  every entry in `/usr/share/terminfo`.
- `_examples/main.go` is the example that `README.md` links. The go command
  skips a folder whose name starts with `_` (D12), so `go vet ./...` does not
  check it. Run `go vet ./_examples/` after you change it.
- `skills_test.go` and `docs_test.go` hold the checks on the agent files and
  on `docs/PLAN.md`.

## Tests and lint

Run these in the repository root. Each must pass, `gofmt -l .` must print
nothing, and `golangci-lint` must report 0 issues:

```sh
gofmt -l .
go vet ./...
go test -count=1 ./...
golangci-lint run ./...
```

CI runs the same commands in GitHub Actions. It runs the tests on Go 1.25
and on the latest stable Go (D17, D18). It runs golangci-lint v2.14.0 with
`.golangci.yml`, which is the configuration of dburl (D25). A `//nolint`
comment names the linter and gives the reason.

The tests read the real terminfo database. `TestLoad`, `TestOpen` and
`TestValues` walk every entry in `/lib/terminfo` and `/usr/share/terminfo`.
`TestValues` runs `/usr/bin/infocmp -E -x` on each entry and compares each
value with the decoded one (D3). A machine that runs the tests needs
`infocmp` at that path and a terminfo database. On Debian and Ubuntu, the
packages are `ncurses-bin` and `ncurses-term`.

## Go conventions

These are the conventions that the code follows now. Match them.

- A sentinel error is a constant of `type Error string`, declared in
  `terminfo.go`, with an `Err` prefix (D11). The code does not use
  `errors.New` in the package.
- An error string starts with a lowercase letter.
- A receiver is one or two letters: `ti` for `*Terminfo`, `d` for
  `*decoder`, `p` for `*parametizer`, `s` for `*stack`, `c` for
  `ColorLevel`.
- A comment inside a function is short and starts with a lowercase letter,
  such as `// read header`. A doc comment is a full sentence that starts with
  the name.
- The code uses `any`, `range` over an integer, and the `slices` and `maps`
  packages. The module declares `go 1.25` (D18).
- A test is table driven where it has cases, and it names each case with
  `t.Run`.
- The package imports only the standard library (D20).

## Writing documentation

Standing rule 2 applies: load `simple-english` first. Write short sentences
in the active voice. Do not use contractions, semicolons or em dashes. Use
"must" and "can", not "should" and "may". State the fact and not its
importance.
