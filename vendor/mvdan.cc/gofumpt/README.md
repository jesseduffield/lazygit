# gofumpt

[![Go Reference](https://pkg.go.dev/badge/mvdan.cc/gofumpt/format.svg)](https://pkg.go.dev/mvdan.cc/gofumpt/format)

	go install mvdan.cc/gofumpt@latest

Enforce a stricter format than `gofmt`, while being backwards compatible.
That is, `gofumpt` is happy with a subset of the formats that `gofmt` is happy with.

The tool is a fork of `gofmt` as of Go 1.25.0, and requires Go 1.24 or later.
It can be used as a drop-in replacement to format your Go code,
and running `gofmt` after `gofumpt` should produce no changes.
For example:

	gofumpt -l -w .

Some of the Go source files in this repository belong to the Go project.
The project includes copies of `go/printer` and `go/doc/comment` as of Go 1.25.0
to ensure consistent formatting independent of what Go version is being used.
The [added formatting rules](#Added-rules) are implemented in the `format` package.

`vendor` and `testdata` directories are skipped unless given as explicit arguments.
Similarly, the added rules do not apply to generated Go files unless they are
given as explicit arguments.

[`ignore` directives](https://go.dev/ref/mod#go-mod-file-ignore) in `go.mod` files are obeyed as well,
unless directories or files within them are given as explicit arguments.

Finally, note that the `-r` rewrite flag is removed in favor of `gofmt -r`,
and the `-s` flag is hidden as it is always enabled.

### Added rules

**No empty lines following an assignment operator**

<details><summary><i>Example</i></summary>

```go
func foo() {
    foo :=
        "bar"
}
```

```go
func foo() {
	foo := "bar"
}
```

</details>

**No empty lines around function bodies**

<details><summary><i>Example</i></summary>

```go
func foo() {

	println("bar")

}
```

```go
func foo() {
	println("bar")
}
```

</details>

**Functions should separate `) {` where the indentation helps readability**

<details><summary><i>Example</i></summary>

```go
func foo(s string,
	i int) {
	println("bar")
}

// With an empty line it's slightly better, but still not great.
func bar(s string,
	i int) {

	println("bar")
}
```

```go
func foo(s string,
	i int,
) {
	println("bar")
}

// With an empty line it's slightly better, but still not great.
func bar(s string,
	i int,
) {
	println("bar")
}
```

</details>

**No empty lines around a lone statement (or comment) in a block**

<details><summary><i>Example</i></summary>

```go
if err != nil {

	return err
}
```

```go
if err != nil {
	return err
}
```

</details>

**No empty lines before a simple error check**

<details><summary><i>Example</i></summary>

```go
foo, err := processFoo()

if err != nil {
	return err
}
```

```go
foo, err := processFoo()
if err != nil {
	return err
}
```

</details>

**Composite literals should use newlines consistently**

<details><summary><i>Example</i></summary>

```go
// A newline before or after an element requires newlines for the opening and
// closing braces.
var ints = []int{1, 2,
	3, 4}

// A newline between consecutive elements requires a newline between all
// elements.
var matrix = [][]int{
	{1},
	{2}, {
		3,
	},
}
```

```go
var ints = []int{
	1, 2,
	3, 4,
}

var matrix = [][]int{
	{1},
	{2},
	{
		3,
	},
}
```

</details>

**Empty field lists should use a single line**

<details><summary><i>Example</i></summary>

```go
var V interface {
} = 3

type T struct {
}

func F(
)
```

```go
var V interface{} = 3

type T struct{}

func F()
```

</details>

**`std` imports must be in a separate group at the top**

<details><summary><i>Example</i></summary>

```go
import (
	"foo.com/bar"

	"io"

	"io/ioutil"
)
```

```go
import (
	"io"
	"io/ioutil"

	"foo.com/bar"
)
```

</details>

**Short case clauses should take a single line**

<details><summary><i>Example</i></summary>

```go
switch c {
case 'a', 'b',
	'c', 'd':
}
```

```go
switch c {
case 'a', 'b', 'c', 'd':
}
```

</details>

**Multiline top-level declarations must be separated by empty lines**

<details><summary><i>Example</i></summary>

```go
func foo() {
	println("multiline foo")
}
func bar() {
	println("multiline bar")
}
```

```go
func foo() {
	println("multiline foo")
}

func bar() {
	println("multiline bar")
}
```

</details>

**Single var declarations should not be grouped with parentheses**

<details><summary><i>Example</i></summary>

```go
var (
	foo = "bar"
)
```

```go
var foo = "bar"
```

</details>

**Contiguous top-level declarations should be grouped together**

<details><summary><i>Example</i></summary>

```go
var nicer = "x"
var with = "y"
var alignment = "z"
```

```go
var (
	nicer     = "x"
	with      = "y"
	alignment = "z"
)
```

</details>

**Simple var-declaration statements should use short assignments**

<details><summary><i>Example</i></summary>

```go
var s = "somestring"
```

```go
s := "somestring"
```

</details>

**The `-s` code simplification flag is enabled by default**

<details><summary><i>Example</i></summary>

```go
var _ = [][]int{[]int{1}}
```

```go
var _ = [][]int{{1}}
```

</details>

**Octal integer literals should use the `0o` prefix on modules using Go 1.13 and later**

<details><summary><i>Example</i></summary>

```go
const perm = 0755
```

```go
const perm = 0o755
```

</details>

**Comments which aren't Go directives should start with a whitespace**

<details><summary><i>Example</i></summary>

```go
//go:noinline

//Foo is awesome.
func Foo() {}
```

```go
//go:noinline

// Foo is awesome.
func Foo() {}
```

</details>

**Composite literals should not have leading or trailing empty lines**

<details><summary><i>Example</i></summary>

```go
var _ = []string{

	"foo",

}

var _ = map[string]string{

	"foo": "bar",

}
```

```go
var _ = []string{
	"foo",
}

var _ = map[string]string{
	"foo": "bar",
}
```

</details>

**Field lists should not have leading or trailing empty lines**

<details><summary><i>Example</i></summary>

```go
type Person interface {

	Name() string

	Age() int

}

type ZeroFields struct {

	// No fields are needed here.

}
```

```go
type Person interface {
	Name() string

	Age() int
}

type ZeroFields struct {
	// No fields are needed here.
}
```

</details>

### Extra rules behind `-extra`

**Adjacent parameters with the same type should be grouped together**

<details><summary><i>Example</i></summary>

```go
func Foo(bar string, baz string) {}
```

```go
func Foo(bar, baz string) {}
```

</details>

**Avoid naked returns for the sake of clarity**

<details><summary><i>Example</i></summary>

```go
func Foo() (err error) {
	return
}
```

```go
func Foo() (err error) {
	return err
}
```

</details>

### Installation

`gofumpt` is a replacement for `gofmt`, so you can simply `go install` it as
described at the top of this README and use it.

When using an IDE or editor with Go integration based on `gopls`,
it's best to configure the editor to use the `gofumpt` support built into `gopls`.

The instructions below show how to set up `gofumpt` for some of the
major editors out there.

#### Visual Studio Code

Enable the language server following [the official docs](https://github.com/golang/vscode-go#readme),
and then enable gopls's `gofumpt` option. Note that VS Code will complain about
the `gopls` settings, but they will still work.

```json
"go.useLanguageServer": true,
"gopls": {
	"formatting.gofumpt": true,
},
```

#### GoLand

GoLand doesn't use `gopls` so it should be configured to use `gofumpt` directly.
Once `gofumpt` is installed, follow the steps below:

- Open **Settings** (File > Settings)
- Open the **Tools** section
- Find the *File Watchers* sub-section
- Click on the `+` on the right side to add a new file watcher
- Choose *Custom Template*

When a window asks for settings, you can enter the following:

* File Types: Select all .go files
* Scope: Project Files
* Program: Select your `gofumpt` executable
* Arguments: `-w $FilePath$`
* Output path to refresh: `$FilePath$`
* Working directory: `$ProjectFileDir$`
* Environment variables: `GOROOT=$GOROOT$;GOPATH=$GOPATH$;PATH=$GoBinDirs$`

To avoid unnecessary runs, you should disable all checkboxes in the *Advanced* section.

#### Vim

The configuration depends on the plugin you are using: [vim-go](https://github.com/fatih/vim-go)
or [govim](https://github.com/govim/govim).

##### vim-go

To configure `gopls` to use `gofumpt`:

```vim
let g:go_fmt_command="gopls"
let g:go_gopls_gofumpt=1
```

##### govim

To configure `gopls` to use `gofumpt`:

```vim
call govim#config#Set("Gofumpt", 1)
```

#### Neovim

When using [`lspconfig`](https://github.com/neovim/nvim-lspconfig), pass the `gofumpt` setting to `gopls`:

```lua
require('lspconfig').gopls.setup({
    settings = {
        gopls = {
            gofumpt = true
        }
    }
})
```

#### Emacs

For [lsp-mode](https://emacs-lsp.github.io/lsp-mode/) users on version 8.0.0 or higher:

```elisp
(setq lsp-go-use-gofumpt t)
```

For users of `lsp-mode` before `8.0.0`:

```elisp
(lsp-register-custom-settings
 '(("gopls.gofumpt" t)))
```

For [eglot](https://github.com/joaotavora/eglot) users:

```elisp
(setq-default eglot-workspace-configuration
 '((:gopls . ((gofumpt . t)))))
```

#### Helix

When using the `gopls` language server, modify the Go settings in `~/.config/helix/languages.toml`:

```toml
[language-server.gopls.config]
"formatting.gofumpt" = true
```

#### Sublime Text

With ST4, install the Sublime Text LSP extension according to [the documentation](https://github.com/sublimelsp/LSP),
and enable `gopls`'s `gofumpt` option in the LSP package settings,
including setting `lsp_format_on_save` to `true`.

```json
"lsp_format_on_save": true,
"clients":
{
	"gopls":
	{
		"enabled": true,
		"initializationOptions": {
			"gofumpt": true,
		}
	}
}
```

### Zed
For `gofumpt` to be used in Zed, you need to set the `gofumpt` option in the LSP settings. This is done by providing the `"gofumpt": true` in `initialization_options`.

```json
"lsp": {
  "gopls": {
    "initialization_options": {
      "gofumpt": true
    }
  }
}
```

### Roadmap

This tool is a place to experiment. In the long term, the features that work
well might be proposed for `gofmt` itself.

The tool is also compatible with `gofmt` and is aimed to be stable, so you can
rely on it for your code as long as you pin a version of it.

### Frequently Asked Questions

> Why attempt to replace `gofmt` instead of building on top of it?

Our design is to build on top of `gofmt`, and we'll never add rules which
disagree with its formatting. So we extend `gofmt` rather than compete with it.

The tool is a modified copy of `gofmt`, for the purpose of allowing its use as a
drop-in replacement in editors and scripts.

> Why are my module imports being grouped with standard library imports?

Any import paths that don't start with a domain name like `foo.com` are
effectively [reserved by the Go toolchain](https://github.com/golang/go/issues/32819).
Third party modules should either start with a domain name,
even a local one like `foo.local`, or use [a reserved path prefix](https://github.com/golang/go/issues/37641).

For backwards compatibility with modules set up before these rules were clear,
`gofumpt` will treat any import path sharing a prefix with the current module
path as third party. For example, if the current module is `mycorp/mod1`, then
all import paths in `mycorp/...` will be considered third party.

> How can I use `gofumpt` if I already use `goimports` to replace `gofmt`?

Most editors have replaced the `goimports` program with the same functionality
provided by a language server like `gopls`. This mechanism is significantly
faster and more powerful, since the language server has more information that is
kept up to date, necessary to add missing imports.

As such, the general recommendation is to let your editor fix your imports -
either via `gopls`, such as VSCode or vim-go, or via their own custom
implementation, such as GoLand. Then follow the install instructions above to
enable the use of `gofumpt` instead of `gofmt`.

If you want to avoid integrating with `gopls`, and are OK with the overhead of
calling `goimports` from scratch on each save, you should be able to call both
tools; for example, `goimports file.go && gofumpt file.go`.

### Contributing

Issues and pull requests are welcome! Please open an issue to discuss a feature
before sending a pull request.

We also use the `#gofumpt` channel over at the
[Gophers Slack](https://invite.slack.golangbridge.org/) to chat.

When reporting a formatting bug, insert a `//gofumpt:diagnose` comment.
The comment will be rewritten to include useful debugging information.
For instance:

```
$ cat f.go
package p

//gofumpt:diagnose
$ gofumpt f.go
package p

//gofumpt:diagnose v0.1.1-0.20211103104632-bdfa3b02e50a -lang=go1.16
```

### License

Note that much of the code is copied from Go's `gofmt` command. You can tell
which files originate from the Go repository from their copyright headers. Their
license file is `LICENSE.google`.

`gofumpt`'s original source files are also under the 3-clause BSD license, with
the separate file `LICENSE`.
