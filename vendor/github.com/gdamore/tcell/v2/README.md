<img src="logos/tcell.png" style="float: right"/>

# Tcell

_Tcell_ is a _Go_ package that provides a cell based view for text terminals, like _XTerm_.
It was inspired by _termbox_, but includes many additional improvements.

[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://stand-with-ukraine.pp.ua)
[![Linux](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/linux.yml?branch=main&logoColor=grey&logo=linux&label=)](https://github.com/gdamore/tcell/actions/workflows/linux.yml)
[![Windows](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/windows.yml?branch=main&logoColor=grey&label=Windows)](https://github.com/gdamore/tcell/actions/workflows/windows.yml)
[![Web Assembly](https://img.shields.io/github/actions/workflow/status/gdamore/tcell/webasm.yml?branch=main&logoColor=grey&logo=webassembly&label=)](https://github.com/gdamore/tcell/actions/workflows/webasm.yml)
[![Apache License](https://img.shields.io/github/license/gdamore/tcell.svg?logoColor=silver&logo=opensourceinitiative&color=blue&label=)](https://github.com/gdamore/tcell/blob/master/LICENSE)
[![Docs](https://img.shields.io/badge/godoc-reference-blue.svg?label=&logo=go)](https://pkg.go.dev/github.com/gdamore/tcell/v2)
[![Discord](https://img.shields.io/discord/639503822733180969?label=&logo=discord)](https://discord.gg/urTTxDN)
[![Coverage](https://img.shields.io/codecov/c/github/gdamore/tcell?logoColor=grey&logo=codecov&label=)](https://codecov.io/gh/gdamore/tcell)
[![Go Report Card](https://goreportcard.com/badge/github.com/gdamore/tcell/v2)](https://goreportcard.com/report/github.com/gdamore/tcell/v2)
[![Latest Release](https://img.shields.io/github/v/release/gdamore/tcell.svg?logo=github&label=)](https://github.com/gdamore/tcell/releases)

Please see [here](UKRAINE.md) for an important message for the people of Russia.

NOTE: This is version 2 of _Tcell_. There are breaking changes relative to version 1.
Version 1.x remains available using the import `github.com/gdamore/tcell`.

## Tutorial

A brief, and still somewhat rough, [tutorial](TUTORIAL.md) is available.

## Examples

A number of example are posted up on our [Gallery](https://github.com/gdamore/tcell/wikis/Gallery/).

Let us know if you want to add your masterpiece to the list!

## Pure Go Terminfo Database

_Tcell_ includes a full parser and expander for terminfo capability strings,
so that it can avoid hard coding escape strings for formatting. It also favors
portability, and includes support for all POSIX systems.

The database is also flexible & extensible, and can be modified by either running
a program to build the entire database, or an entry for just a single terminal.

## More Portable

_Tcell_ is portable to a wide variety of systems, and is pure Go, without
any need for CGO.
_Tcell_ is believed to work with mainstream systems officially supported by golang.

Following the Go support policy, _Tcell_ officially only supports the current ("stable") version of go,
and the version immediately prior ("oldstable").  This policy is necessary to make sure that we can
update dependencies to pick up security fixes and new features, and it allows us to adopt changes
(such as library and language features) that are only supported in newer versions of Go.

## No Async IO

_Tcell_ is able to operate without requiring `SIGIO` signals (unlike _termbox_),
or asynchronous I/O, and can instead use standard Go file objects and Go routines.
This means it should be safe, especially for
use with programs that use exec, or otherwise need to manipulate the tty streams.
This model is also much closer to idiomatic Go, leading to fewer surprises.

## Rich Unicode & non-Unicode support

_Tcell_ includes enhanced support for Unicode, including wide characters and
combining characters, provided your terminal can support them.
Note that
Windows terminals generally don't support the full Unicode repertoire.

It will also convert to and from Unicode locales, so that the program
can work with UTF-8 internally, and get reasonable output in other locales.
_Tcell_ tries hard to convert to native characters on both input and output.
On output _Tcell_ even makes use of the alternate character set to facilitate
drawing certain characters.

## More Function Keys

_Tcell_ also has richer support for a larger number of special keys that some
terminals can send.

## Better Color Handling

_Tcell_ will respect your terminal's color space as specified within your terminfo entries.
For example attempts to emit color sequences on VT100 terminals
won't result in unintended consequences.

_Tcell_ maps 16 colors down to 8, for terminals that need it.
(The upper 8 colors are just brighter versions of the lower 8.)

## Better Mouse Support

_Tcell_ supports enhanced mouse tracking mode, so your application can receive
regular mouse motion events, and wheel events, if your terminal supports it.

## _Termbox_ Compatibility

A compatibility layer for _termbox_ is provided in the `compat` directory.
To use it, try importing `github.com/gdamore/tcell/termbox` instead.
Most _termbox-go_ programs will probably work without further modification.

## Working With Unicode

Internally _Tcell_ uses UTF-8, just like Go.
However, _Tcell_ understands how to
convert to and from other character sets, using the capabilities of
the `golang.org/x/text/encoding` packages.
Your application must supply
them, as the full set of the most common ones bloats the program by about 2 MB.
If you're lazy, and want them all anyway, see the `encoding` sub-directory.

## Wide & Combining Characters

The `SetContent()` API takes a primary rune, and an optional list of combining runes.
If any of the runes is a wide (East Asian) rune occupying two cells,
then the library will skip output from the following cell. Care must be
taken in the application to avoid explicitly attempting to set content in the
next cell, otherwise the results are undefined. (Normally the wide character
is displayed, and the other character is not; do not depend on that behavior.)

## Colors

_Tcell_ assumes the ANSI/XTerm color model, including the 256 color map that
XTerm uses when it supports 256 colors. The terminfo guidance will be
honored, with respect to the number of colors supported. Also, only
terminals which expose ANSI style `setaf` and `setab` will support color;
if you have a color terminal that only has `setf` and `setb`, please submit
a ticket.

## 24-bit Color

_Tcell_ _supports 24-bit color!_ (That is, if your terminal can support it.)

NOTE: Technically the approach of using 24-bit RGB values for color is more
accurately described as "direct color", but most people use the term "true color".
We follow the (inaccurate) common convention.

There are a few ways you can enable (or disable) true color.

- For many terminals, we can detect it automatically if your terminal
  includes the `RGB` or `Tc` capabilities (or rather it did when the database
  was updated.)

- You can force this one by setting the `COLORTERM` environment variable to
  `24-bit`, `truecolor` or `24bit`. This is the same method used
  by most other terminal applications that support 24-bit color.

- If you set your `TERM` environment variable to a value with the suffix `-truecolor`
  then 24-bit color compatible with XTerm and ECMA-48 will be assumed.
  (This feature is deprecated.
  It is recommended to use one of other methods listed above.)

- You can disable 24-bit color by setting `TCELL_TRUECOLOR=disable` in your
  environment.

When using TrueColor, programs will display the colors that the programmer
intended, overriding any "`themes`" you may have set in your terminal
emulator. (For some cases, accurate color fidelity is more important
than respecting themes. For other cases, such as typical text apps that
only use a few colors, its more desirable to respect the themes that
the user has established.)

## Performance

Reasonable attempts have been made to minimize sending data to terminals,
avoiding repeated sequences or drawing the same cell on refresh updates.

## Terminfo

(Not relevant for Windows users.)

The Terminfo implementation operates with a built-in database.
This should satisfy most users. However, it can also (on systems
with ncurses installed), dynamically parse the output from `infocmp`
for terminals it does not already know about.

See the `terminfo/` directory for more information about generating
new entries for the built-in database.

_Tcell_ requires that the terminal support the `cup` mode of cursor addressing.
Ancient terminals without the ability to position the cursor directly
are not supported.
This is unlikely to be a problem; such terminals have not been mass-produced
since the early 1970s.

## Mouse Support

Mouse support is detected via the `kmous` terminfo variable, however,
enablement/disablement and decoding mouse events is done using hard coded
sequences based on the XTerm X11 model. All popular
terminals with mouse tracking support this model. (Full terminfo support
is not possible as terminfo sequences are not defined.)

On Windows, the mouse works normally.

Mouse wheel buttons on various terminals are known to work, but the support
in terminal emulators, as well as support for various buttons and
live mouse tracking, varies widely.
Modern _xterm_, macOS _Terminal_, and _iTerm_ all work well.

## Bracketed Paste

Terminals that appear to support the XTerm mouse model also can support
bracketed paste, for applications that opt-in. See `EnablePaste()` for details.

## Testability

There is a `SimulationScreen`, that can be used to simulate a real screen
for automated testing. The supplied tests do this. The simulation contains
event delivery, screen resizing support, and capabilities to inject events
and examine "`physical`" screen contents.

## Platforms

### POSIX (Linux, FreeBSD, macOS, Solaris, etc.)

Everything works using pure Go on mainstream platforms. Some more esoteric
platforms (e.g., AIX) may need to be added. Pull requests are welcome!

### Windows

Windows console mode applications are supported.

Modern console applications like ConEmu and the Windows Terminal,
support all the good features (resize, mouse tracking, etc.)

### WASM

WASM is supported, but needs additional setup detailed in [README-wasm](README-wasm.md).

### Plan9 and its variants

Plan 9 is supported on a limited basis. The Plan 9 backend opens `/dev/cons` for I/O, enables raw mode by writing `rawon`/`rawoff` to `/dev/consctl`, watches `/dev/wctl` for resize notifications, and then constructs a **terminfo-backed** `Screen` (so `NewScreen` works as on other platforms). Typical usage is inside `vt(1)` with `TERM=vt100`. Expect **monochrome text** and **no mouse reporting** under stock `vt(1)` (it generally does not emit ANSI color or xterm mouse sequences). If a Plan 9 terminal supplies ANSI color escape sequences and xterm-style mouse reporting, color can be picked up via **terminfo** and mouse support could be added by wiring those sequences into the Plan 9 TTY path; contributions that improve terminal detection and broaden feature support are welcome.

### Commercial Support

_Tcell_ is absolutely free, but if you want to obtain commercial, professional support, there are options.

- [TideLift](https://tidelift.com/) subscriptions include support for _Tcell_, as well as many other open source packages.
- [Staysail Systems Inc.](mailto:info@staysail.tech) offers direct support, and custom development around _Tcell_ on an hourly basis.
