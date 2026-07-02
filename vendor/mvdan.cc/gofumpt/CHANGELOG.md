# Changelog

## [v0.9.1] - 2025-09-07

This is a bugfix release to address a regression in detecting
comment directives with special characters such as `//golangcitest:config_path`.

## [v0.9.0] - 2025-09-02

This release is based on Go 1.25's gofmt, and requires Go 1.24 or later.

A new rule is introduced to "clothe" naked returns for the sake of clarity.
While there is nothing wrong with naming results in function signatures,
using lone `return` statements can be confusing to the reader.

Go 1.25's `ignore` directives in `go.mod` files are now obeyed;
any directories within the module matching any of the patterns
are now omitted when walking directories, such as with `gofumpt -w .`.

Module information is now loaded via Go's [`x/mod/modfile` package](https://pkg.go.dev/golang.org/x/mod/modfile)
rather than executing `go mod edit -json`, which is way faster.
This should result in moderate speed-ups when formatting many directories.

## [v0.8.0] - 2025-04-13

This release is based on Go 1.24's gofmt, and requires Go 1.23 or later.

The following changes are included:

* Fail with `-d` if formatting any file resulted in a diff - #114
* Do not panic when a `go.mod` file is missing a `go` directive - #317

## [v0.7.0] - 2024-08-16

This release is based on Go 1.23.0's gofmt, and requires Go 1.22 or later.

The following changes are included:

* Group `internal/...` imported packages as standard library - #307

## [v0.6.0] - 2024-01-28

This release is based on Go 1.21's gofmt, and requires Go 1.20 or later.

The following changes are included:

* Support `go` version strings from newer go.mod files - [#280]
* Consider simple error checks even if they use the `=` operator - [#271]
* Ignore `//line` directives to avoid panics - [#288]

## [v0.5.0] - 2023-04-09

This release is based on Go 1.20's gofmt, and requires Go 1.19 or later.

The biggest change in this release is that we now vendor copies of the packages
`go/format`, `go/printer`, and `go/doc/comment` on top of `cmd/gofmt` itself.
This allows for each gofumpt release to format code in exactly the same way
no matter what Go version is used to build it, as Go versions can change those
three packages in ways that alter formatting behavior.

This vendoring adds a small amount of duplication when using the
`mvdan.cc/gofumpt/format` library, but it's the only way to make gofumpt
versions consistent in their behavior and formatting, just like gofmt.

The jump to Go 1.20's `go/printer` should also bring a small performance
improvement, as we contributed patches to make printing about 25% faster:

* https://go.dev/cl/412555
* https://go.dev/cl/412557
* https://go.dev/cl/424924

The following changes are included as well:

* Skip `testdata` dirs by default like we already do for `vendor` - [#260]
* Avoid inserting newlines incorrectly in some func signatures - [#235]
* Avoid joining some comments with the previous line - [#256]
* Fix `gofumpt -version` for release archives - [#253]

## [v0.4.0] - 2022-09-27

This release is based on Go 1.19's gofmt, and requires Go 1.18 or later.
We recommend building gofumpt with Go 1.19 for the best formatting results.

The jump from Go 1.18 brings diffing in pure Go, removing the need to exec `diff`,
and a small parsing speed-up thanks to `go/parser.SkipObjectResolution`.

The following formatting fixes are included as well:

* Allow grouping declarations with comments - [#212]
* Properly measure the length of case clauses - [#217]
* Fix a few crashes found by Go's native fuzzing

## [v0.3.1] - 2022-03-21

This bugfix release resolves a number of issues:

* Avoid "too many open files" error regression introduced by [v0.3.0] - [#208]
* Use the `go.mod` relative to each Go file when deriving flag defaults - [#211]
* Remove unintentional debug prints when directly formatting files

## [v0.3.0] - 2022-02-22

This is gofumpt's third major release, based on Go 1.18's gofmt.
The jump from Go 1.17's gofmt should bring a noticeable speed-up,
as the tool can now format many files concurrently.
On an 8-core laptop, formatting a large codebase is 4x as fast.

The following [formatting rules](https://github.com/mvdan/gofumpt#Added-rules) are added:

* Functions should separate `) {` where the indentation helps readability
* Field lists should not have leading or trailing empty lines

The following changes are included as well:

* Generated files are now fully formatted when given as explicit arguments
* Prepare for Go 1.18's module workspaces, which could cause errors
* Import paths sharing a prefix with the current module path are no longer
  grouped with standard library imports
* `format.Options` gains a `ModulePath` field per the last bullet point

## [v0.2.1] - 2021-12-12

This bugfix release resolves a number of issues:

* Add deprecated flags `-s` and `-r` once again, now giving useful errors
* Avoid a panic with certain function declaration styles
* Don't group interface members of different kinds
* Account for leading comments in composite literals

## [v0.2.0] - 2021-11-10

This is gofumpt's second major release, based on Go 1.17's gofmt.
The jump from Go 1.15's gofmt should bring a mild speed-up,
as walking directories with `filepath.WalkDir` uses fewer syscalls.

gofumports is now removed, after being deprecated in [v0.1.0].
Its main purpose was IDE integration; it is now recommended to use gopls,
which in turn implements goimports and supports gofumpt natively.
IDEs which don't integrate with gopls (such as GoLand) implement goimports too,
so it is safe to use gofumpt as their "format on save" command.
See the [installation instructions](https://github.com/mvdan/gofumpt#Installation)
for more details.

The following [formatting rules](https://github.com/mvdan/gofumpt#Added-rules) are added:

* Composite literals should not have leading or trailing empty lines
* No empty lines following an assignment operator
* Functions using an empty line for readability should use a `) {` line instead
* Remove unnecessary empty lines from interfaces

Finally, the following changes are made to the gofumpt tool:

* Initial support for Go 1.18's type parameters is added
* The `-r` flag is removed in favor of `gofmt -r`
* The `-s` flag is removed as it is always enabled
* Vendor directories are skipped unless given as explicit arguments
* The added rules are not applied to generated Go files
* The `format` Go API now also applies the `gofmt -s` simplification
* Add support for `//gofumpt:diagnose` comments

## [v0.1.1] - 2021-03-11

This bugfix release backports fixes for a few issues:

* Keep leading empty lines in func bodies if they help readability
* Avoid breaking comment alignment on empty field lists
* Add support for `//go-sumtype:` directives

## [v0.1.0] - 2021-01-05

This is gofumpt's first release, based on Go 1.15.x. It solidifies the features
which have worked well for over a year.

This release will be the last to include `gofumports`, the fork of `goimports`
which applies `gofumpt`'s rules on top of updating the Go import lines. Users
who were relying on `goimports` in their editors or IDEs to apply both `gofumpt`
and `goimports` in a single step should switch to gopls, the official Go
language server. It is supported by many popular editors such as VS Code and
Vim, and already bundles gofumpt support. Instructions are available [in the
README](https://github.com/mvdan/gofumpt).

`gofumports` also added maintenance work and potential confusion to end users.
In the future, there will only be one way to use `gofumpt` from the command
line. We also have a [Go API](https://pkg.go.dev/mvdan.cc/gofumpt/format) for
those building programs with gofumpt.

Finally, this release adds the `-version` flag, to print the tool's own version.
The flag will work for "master" builds too.

[v0.9.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.9.0
[v0.8.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.8.0
[v0.7.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.7.0

[v0.6.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.6.0
[#271]: https://github.com/mvdan/gofumpt/issues/271
[#280]: https://github.com/mvdan/gofumpt/issues/280
[#288]: https://github.com/mvdan/gofumpt/issues/288

[v0.5.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.5.0
[#235]: https://github.com/mvdan/gofumpt/issues/235
[#253]: https://github.com/mvdan/gofumpt/issues/253
[#256]: https://github.com/mvdan/gofumpt/issues/256
[#260]: https://github.com/mvdan/gofumpt/issues/260

[v0.4.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.4.0
[#212]: https://github.com/mvdan/gofumpt/issues/212
[#217]: https://github.com/mvdan/gofumpt/issues/217

[v0.3.1]: https://github.com/mvdan/gofumpt/releases/tag/v0.3.1
[#208]: https://github.com/mvdan/gofumpt/issues/208
[#211]: https://github.com/mvdan/gofumpt/pull/211

[v0.3.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.3.0
[v0.2.1]: https://github.com/mvdan/gofumpt/releases/tag/v0.2.1
[v0.2.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.2.0
[v0.1.1]: https://github.com/mvdan/gofumpt/releases/tag/v0.1.1
[v0.1.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.1.0
