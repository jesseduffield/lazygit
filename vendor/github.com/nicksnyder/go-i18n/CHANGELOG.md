# Changelog

Major version changes are documented in the changelog.

To see the documentation for minor or patch version, [view the release notes](https://github.com/nicksnyder/go-i18n/releases).

## v2

### Motivation

The first commit to this project was January 2012 (go1 had not yet been released) and v1.0.0 was tagged June 2015 (go1.4).
This project has evolved with the Go ecosystem since then in a backwards compatible way,
but there is a growing list of issues and warts that cannot be addressed without breaking compatiblity.

v2 is rewrite of the API from first principals to make it more idiomatic Go, and to resolve a backlog of issues: https://github.com/nicksnyder/go-i18n/milestone/1

### Improvements

* Use `golang.org/x/text/language` to get standardized behavior for language matching (https://github.com/nicksnyder/go-i18n/issues/30, https://github.com/nicksnyder/go-i18n/issues/44, https://github.com/nicksnyder/go-i18n/issues/76)
* Remove global state so that the race detector does not complain when downstream projects run tests that depend on go-i18n in parallel (https://github.com/nicksnyder/go-i18n/issues/82)
* Automatically extract messages from Go source code (https://github.com/nicksnyder/go-i18n/issues/64)
* Provide clearer documentation and examples (https://github.com/nicksnyder/go-i18n/issues/27)
* Reduce complexity of file format for simple translations (https://github.com/nicksnyder/go-i18n/issues/85)
* Support descriptions for messages (https://github.com/nicksnyder/go-i18n/issues/8)
* Support custom template delimiters (https://github.com/nicksnyder/go-i18n/issues/88)

### Upgrading from v1

The i18n package in v2 is completely different than v1.
Refer to the [documentation](https://godoc.org/github.com/nicksnyder/go-i18n/v2/i18n) and [README](https://github.com/nicksnyder/go-i18n/blob/master/README.md) for guidance.

The goi18n command has similarities and differences:

* `goi18n merge` has a new implementation but accomplishes the same task.
* `goi18n extract` extracts messages from Go source files.
* `goi18n constants` no longer exists. Prefer to extract messages directly from Go source files.

v2 makes changes to the canonical message file format, but you can use v1 message files with v2. Message files will be converted to the new format the first time they are processed by the new `goi18n merge` command.

v2 requires Go 1.9 or newer.
