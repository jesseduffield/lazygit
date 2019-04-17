go-gitconfig 
====

[![GitHub release](http://img.shields.io/github/release/tcnksm/go-gitconfig.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/544ee33aea87f6374f001483.svg?style=flat-square)][wercker]
[![Coveralls](http://img.shields.io/coveralls/tcnksm/go-gitconfig.svg?style=flat-square)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[release]: https://github.com/tcnksm/go-gitconfig/releases
[wercker]: https://app.wercker.com/project/bykey/89c5a6e50a0daceec971ff5ce210164a
[coveralls]: https://coveralls.io/r/tcnksm/go-gitconfig
[license]: https://github.com/tcnksm/go-gitconfig/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/tcnksm/go-gitconfig


`go-gitconfig` is a pacakge to use `gitconfig` values in Golang.

Sometimes you want to extract username or its email address **implicitly** in your tool.
Now most of developer use `git`, so we can use its configuration variables. `go-gitconfig` is for that.

`go-gitconfig` is very small, so it may not be included what you want to use.
If you want to use more git specific variable, check [Other](##VS).

## Usage

If you want to use git user name defined in `~/.gitconfig`: 

```go
username, err := gitconfig.Username()
```

Or git user email defined in `~/.gitconfig`: 

```go
email, err := gitconfig.Email()
```

Or, if you want to extract origin url of current project (from `.git/config`):

```go
url, err := gitconfig.OriginURL()
```

You can also extract value by key:

```go
editor, err := gitconfig.Global("core.editor")
```

```go
remote, err := gitconfig.Local("branch.master.remote")
```

See more details in document at [https://godoc.org/github.com/tcnksm/go-gitconfig](https://godoc.org/github.com/tcnksm/go-gitconfig). 

## Install

To install, use `go get`:

```bash
$ go get -d github.com/tcnksm/go-gitconfig
```

## VS.

- [speedata/gogit](https://github.com/speedata/gogit)
- [libgit2/git2go](https://github.com/libgit2/git2go)

These packages have many features to use git from golang. `go-gitconfig` is very simple alternative and focus to extract information from gitconfig. `go-gitconfig` is used in [tcnksm/ghr](https://github.com/tcnksm/ghr). 

## Contribution

1. Fork ([https://github.com/tcnksm/go-gitconfig/fork](https://github.com/tcnksm/go-gitconfig/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

## Author

[tcnksm](https://github.com/tcnksm)
