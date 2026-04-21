## Breaking Changes in _Tcell_ v3

There are a number of changes in _Tcell_ v3, mostly aimed at simplifying things for
applications, but some also intended to reduce the burden for support.  Every application
will need at least some changes, but it is expected that those changes will be small,
possibly even mechanical in nature.

### Cell and Contents APIs

In order to improve support for multi-rune grapheme clusters, and to provide an
experience that reduces friction when using it, some APIs have been removed, and
newer APIs exist in their place.

- `SetCell` and `SetContents` are removed.  Use `Put` instead.
- `GetContents` is removed. Use `Get` instead.

### Events (PostEvent, PollEvent, ChannelEvents)

The event channel is now directly exposed via `EventQ`, and events may be read from or written
directly to the channel in the standard Go fashion.  This should help applications that want to
integrate into `select` statements (e.g. for timed key presses).

The `ChannelEvents`, `PollEvent`, `PostEvent`, and `PostEventWait` functions are removed, as
applications can now just access the event channel directly.

### Key Event Changes

`EventKey` now carries a string for `KeyRune` instead of a single rune.
As a result the old `Rune` method for `EventKey` is replaced by `Str`.
The main difference for most users will be that `Str` returns a string, and most
of the time that string will consist of only a single rune. However, it is possible
now to inject synthetic key strokes consisting of multi-rune grapheme clusters.

Additionally the following special keys are removed, as they are delivered
instead as `KeyRune` with the relevant rune, and the `ModCtrl` modifier:
`KeyCtrlSpace`, `KeyCtrlLeftSq`, `KeyCtrlRightSq`, `KeyCtrlBackslash`, and
`KeyCtrlUnderscore`.

Note that `KeyRune` will never have a `ModShift` applied unless it is applied
also with other modifiers.

The `KeyCtrlA` through `KeyCtrlZ` keys are delivered, but will also carry
the associated lower case rune (e.g. "a", "b", etc.) and `ModCtrl`.

The `KeyBackspace2` key is no longer delivered, but is converted to
`KeyBackspace`. (This resolves some inconsistency around e.g. CTRL-H vs DELETE.)

### Termbox Compatibility Removed

The `termbox` compatibility package is removed. Few applications were using it,
and the compatibility was imperfect. Also the package had limited support for many
newer features. Further, _Termbox_ itself is no longer being maintained.
Applications that still need this should keep using _Tcell_ v2.

### Terminfo Removed

The Terminfo subsystem has been removed entirely.
Essentially the old terminfo based design has long proved to be inferior for modern terminal
applications, and has not kept up with newer terminal features such as 24-bit color,
different mouse reporting modes, bracketed paste, advanced text styling, and so forth.

As part of this, we're removing the parsed terminfo logic entirely.  It turns out that pretty much
all of the terminal logic can be consolidated to just a few classes of terminals with substantial
overlap.

A consequence of this is that support for some legacy terminals that are either functionally
extinct (such as _hpterm_) or unlikely to be found outside of a museum (such as VT52, Wyse50, or
anything produced more than 40 years ago.)

Note that VT100 and later will work in emulation, and VT220 and later physical terminals should still work. 
VT100 physical terminals may not work, as the padding delays that existed for them are removed.
Those delays hurt emulations that do not need them, and existed only to accommodate limitations found on the
physical hardware from the 1970s.

Note that we still examine `$TERM` when appropriate, but if the value is not one we recognize,
then we will assume something reasonably capable and compatible at some level with _xterm_ or
at least ECMA-48.

### Color, Attributes, Etc. Bit Sizes

The `Color` type is now only 32-bits, which should save some memory on large terminal windows.
The `AttrMask` type is now only 16-bits, and the `UnderlineStyle` is now 8 bits.
All these lead to further savings in the memory per cell.

### Underline

`AttrUnderline` is gone.  It was not sufficient to describe styled and colored underlines.

### Removed Capability Queries

Deprecated APIs `HasKey`, `HasMouse`, and `CanDisplay` are removed.
These functions weren't reliable and served no useful purpose.

### Windows Console API

`NewConsoleScreen` is removed as is support for Windows console mode.

Instead this uses the more modern Windows VT modes.
As a consequence, this means that _Tcell_ on Windows requires at least Winows 10 build 1703 (the Creators Update).
If you are using a version of Windows 10 older than that, you should really upgrade for _many_ reasons, not just
because _Tcell_ doesn't support it anymore.

### InputProcessor is no longer Public

This structure, and the associated `NewInputProcessor` function, were made public incorrectly.
They are not part of our public API going forward, and are now private symbols.
