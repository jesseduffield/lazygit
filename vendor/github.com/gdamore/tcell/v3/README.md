<img src="logos/tcell.png" style="float: right"/>

# Tcell

_Tcell_ is a _Go_ package that provides a cell based view for text terminals, like _XTerm_.
It was inspired by _termbox_, but includes many additional improvements.

[![Stand With Ukraine](logos/ukraine.svg)](https://stand-with-ukraine.pp.ua)
[![Docs](https://img.shields.io/badge/godoc-reference-blue.svg?label=&logo=go)](https://pkg.go.dev/github.com/gdamore/tcell/v3)
[![Linux](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/linux.yml?branch=main&logoColor=grey&logo=linux&label=)](https://github.com/gdamore/tcell/actions/workflows/linux.yml)
[![macOS](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/macos.yml?branch=main&logoColor=grey&logo=apple&label=)](https://github.com/gdamore/tcell/actions/workflows/macos.yml)
[![Windows](https://custom-icon-badges.demolab.com/github/actions/workflow/status/gdamore/tcell/windows.yml?branch=main&logoColor=grey&logo=windows10&label=)](https://github.com/gdamore/tcell/actions/workflows/windows.yml)
[![Web Assembly](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/webasm.yml?branch=main&logoColor=grey&logo=webassembly&label=)](https://github.com/gdamore/tcell/actions/workflows/webasm.yml)
[![Coverage](https://img.shields.io/codecov/c/github/gdamore/tcell?logoColor=grey&logo=codecov&label=)](https://codecov.io/gh/gdamore/tcell)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat&label=&logo=go&logoColor=grey)](https://goreportcard.com/report/github.com/gdamore/tcell/v3)
[![Discord](https://img.shields.io/discord/639503822733180969?label=&logo=discord)](https://discord.gg/urTTxDN)
[![Latest Release](https://img.shields.io/github/v/release/gdamore/tcell.svg?logo=github&label=)](https://github.com/gdamore/tcell/releases)

> [!NOTE]
> This is version 3 of _Tcell_.
> There are breaking changes relative to versions 1 and 2.
> [Version 2](https://github.com/gdamore/tcell/tree/v2) remains available using the import `github.com/gdamore/tcell/v2`.
> [Version 1](https://github.com/gdamore/tcell/tree/v1) Version 1.x remains available using the import `github.com/gdamore/tcell`, but is
> unmaintained and should not be used.

## Tutorial

A brief, and still somewhat rough, [tutorial](TUTORIAL.md) is available.

## Examples

A number of example are posted up on our [Gallery](https://github.com/gdamore/tcell/wikis/Gallery/).
That's a wiki, and please do submit updates if you have something you want to showcase.

There are also demonstration programs in the `./demos` directory, as well as some in `./_demos`.

## More Portable

_Tcell_ is portable to a wide variety of systems, and is pure Go, without any need for CGO.
_Tcell_ works with mainstream systems officially supported by golang.

Following the Go support policy, _Tcell_ officially only supports the current ("stable") version of go,
and the version immediately prior ("oldstable").  This policy is necessary to make sure that we can
update dependencies to pick up security fixes and new features, and it allows us to adopt changes
(such as library and language features) that are only supported in newer versions of Go.

## Rich Unicode & non-Unicode support

_Tcell_ includes enhanced support for Unicode, including wide characters and
grapheme clusters, provided your terminal can support them.

It will also convert to and from Unicode locales, so that the program
can work with UTF-8 internally, and get reasonable output in other locales.
_Tcell_ tries hard to convert to native characters on both input and output.
On output _Tcell_ even makes use of the alternate character set to facilitate
drawing certain characters.

## Better Keyboard Support

_Tcell_ also has richer support for a larger number of special keys that some
terminals can send. On modern terminal emulators we can also support a rich set of
modifiers, and can discriminate between e.g. CTRL-I and TAB.  (This does require
the terminal emulator to support one of the modern keyboard protocols.)

## Better Mouse Support

_Tcell_ supports enhanced mouse tracking mode, so your application can receive
regular mouse motion events, click-drag, and wheel events, if your terminal supports it.

## Working With Unicode

Internally _Tcell_ uses UTF-8, just like Go.
However, _Tcell_ understands how to
convert to and from other character sets, using the capabilities of
the `golang.org/x/text/encoding` packages.
Your application must supply
them, as the full set of the most common ones bloats the program by about 2 MB.
If you're lazy, and want them all anyway, see the `encoding` sub-directory.

## Wide & Combining Characters

The `Put()` API takes a string, which should be legal UTF-8, and displays
the first grapheme cluster (which may composed of multiple runes).
It returns the actual width displayed, which can be used to advance the column positiion
for the next display grapheme.  Alternatively, `PutStr()` or `PutStrStyled()`
can be used to display a single line of text (which will be clipped at the
edge of the screen).

If a second character is displayed immediately in the cell adjacent to a
wide character (offset by one instead of by two), then the results are undefined.

## Colors

_Tcell_ assumes the ANSI/XTerm color palette for up to 256 colors, although terminals
such as legacy ANSI terminals may only support 8 colors.

## 24-bit Color

_Tcell_ _supports 24-bit color!_ (That is, if your terminal can support it.)

There are a few ways you can enable (or disable) 24-bit color.

- You can force this one by setting the `COLORTERM` environment variable to
  `truecolor`. This environment variable is frequently set by terminal emulators
  that support 24-bit color.

- On Windows, 24-bit color support is assumed. (All modern Windows terminal emulators support it.)

- If you set your `TERM` environment variable to a value with the suffix `-truecolor`
  or `-direct`, then 24-bit color compatible with XTerm and ECMA-48 will be assumed.

- You can disable 24-bit color by setting `TCELL_TRUECOLOR=disable` in your
  environment.

When using 24-bit color, programs will display the colors that the programmer
intended, overriding any "`themes`" the user may have set in their terminal
emulator. (For some cases, accurate color fidelity is more important
than respecting themes. For other cases, such as typical text apps that
only use a few colors, its more desirable to respect the themes that
the user has established.)

## Performance

Reasonable attempts have been made to minimize sending data to terminals,
avoiding repeated sequences or drawing the same cell on refresh updates.

## Mouse Support

Mouse tracking, buttons, and even wheel mice are supported on most terminal
emulators, as well as Windows.

## Bracketed Paste

Terminals that support support it, can use bracketed paste.
See `EnablePaste()` for details.

## Breaking Changes in v3

There are a number of changes in _Tcell_ version 3, which break compatibility with
version 2 and version 1. 

Your application will almost certainly need some minor updates to work with version 3.

Please see the [CHANGESv3](CHANGESv3.md) document for a list.

## Platforms

### POSIX (Linux, FreeBSD, macOS, Solaris, etc.)

Everything works using pure Go on mainstream platforms.
Esoteric platforms (e.g. zOS or AIX) are supported on a best-effort
only basis. Pull requests to fix any issues found are welcome!

### Windows

Modern Windows is supported.  Please see the [README-windows](README-windows.md)
document for much more detailed information.

### WASM

WASM is supported, but needs additional setup detailed in [README-wasm](README-wasm.md).

### Plan 9

Plan 9 is supported on a best-effort basis.  Please see the [README-plan9](README-plan9.md)
document for more information.

### Commercial Support

_Tcell_ is absolutely free, but if you want to obtain commercial, professional support, there are options.

- [TideLift](https://tidelift.com/) subscriptions include support for _Tcell_, as well as many other open source packages.
- [Staysail Systems Inc.](mailto:info@staysail.tech) offers direct support, and custom development around _Tcell_ on an hourly basis.
