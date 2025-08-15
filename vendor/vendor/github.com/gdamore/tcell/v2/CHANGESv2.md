## Breaking Changes in _Tcell_ v2

A number of changes were made to _Tcell_ for version two, and some of these are breaking.

### Import Path

The import path for tcell has changed to `github.com/gdamore/tcell/v2` to reflect a new major version.

### Style Is Not Numeric

The type `Style` has changed to a structure, to allow us to add additional data such as flags for color setting,
more attribute bits, and so forth.
Applications that relied on this being a number will need to be updated to use the accessor methods.

### Mouse Event Changes

The middle mouse button was reported as button 2 on Linux, but as button 3 on Windows,
and the right mouse button was reported the reverse way.
_Tcell_ now always reports the right mouse button as button 2, and the middle button as button 3.
To help make this clearer, new symbols `ButtonPrimary`, `ButtonSecondary`, and
`ButtonMiddle` are provided.
(Note that which button is right vs. left may be impacted by user preferences.
Usually the left button will be considered the Primary, and the right will be the Secondary.)
Applications may need to adjust their handling of mouse buttons 2 and 3 accordingly.

### Terminals Removed

A number of terminals have been removed.
These are mostly ancient definitions unlikely to be used by anyone, such as `adm3a`.

### High Number Function Keys

Historically terminfo reported function keys with modifiers set as a different
function key altogether.  For example, Shift-F1 was reported as F13 on XTerm.
_Tcell_ now prefers to report these using the base key (such as F1) with modifiers added.
This works on XTerm and VTE based emulators, but some emulators may not support this.
The new behavior more closely aligns with behavior on Windows platforms.

## New Features in _Tcell_ v2

These features are not breaking, but are introduced in version 2.

### Improved Modifier Support

For terminals that appear to behave like the venerable XTerm, _tcell_
automatically adds modifier reporting for ALT, CTRL, SHIFT, and META keys
when the terminal reports them.

### Better Support for Palettes (Themes)

When using a color by its name or palette entry, _Tcell_ now tries to
use that palette entry as is; this should avoid some inconsistency and respect
terminal themes correctly.

When true fidelity to RGB values is needed, the new `TrueColor()` API can be used
to create a direct color, which bypasses the palette altogether.

### Automatic TrueColor Detection

For some terminals, if the `Tc` or `RGB` properties are present in terminfo,
_Tcell_ will automatically assume the terminal supports 24-bit color.

### ColorReset

A new color value, `ColorReset` can be used on the foreground or background
to reset the color the default used by the terminal.

### tmux Support

_Tcell_ now has improved support for tmux, when the `$TERM` variable is set to "tmux".

### Strikethrough Support

_Tcell_ has support for strikethrough when the terminal supports it, using the new `StrikeThrough()` API.

### Bracketed Paste Support

_Tcell_ provides the long requested capability to discriminate paste event by using the
bracketed-paste capability present in some terminals.  This is automatically available on
terminals that support XTerm style mouse handling, but applications must opt-in to this
by using the new `EnablePaste()` function.  A new `EventPaste` type of event will be
delivered when starting and finishing a paste operation.