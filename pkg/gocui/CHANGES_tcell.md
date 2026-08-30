# Change from termbox to tcell

Original GOCUI was written on top of [termbox](https://github.com/nsf/termbox-go) package. This document describes changes which were done to be able to use to [tcell/v2](https://github.com/gdamore/tcell) package.

## Attribute color

Attribute type represents a terminal attribute like color and font effects. Color and font effects can be combined using bitwise OR (`|`).

In `termbox` colors were represented by range 1 to 256. `0` was default color which uses the terminal default setting.

In `tcell` colors can be represented in 24bit, and all of them starts from 0. Valid colors have special flag which gives them real value starting from 4294967296. `0` is a default similart to `termbox`.
The change to support all these colors was made in a way, that original colors from 1 to 256 are backward compatible and if user has color specified as
`Attribute(ansicolor+1)` without the valid color flag, it will be translated to `tcell` color by subtracting 1 and making the color valid by adding the flag. This should ensure backward compatibility.

All the color constants are the same with different underlying values. From user perspective, this should be fine unless some arithmetic is done with it. For example `ColorBlack` was `1` in original version but is `4294967296` in new version.

GOCUI provides a few helper functions which could be used to get the real color value or to create a color attribute.

- `(a Attribute).Hex()` - returns `int32` value of the color represented as `Red << 16 | Green << 8 | Blue`
- `(a Attribute).RGB()` - returns 3 `int32` values for red, green and blue color.
- `GetColor(string)` - creates `Attribute` from color passed as a string. This can be hex value or color name (W3C name).
- `Get256Color(int32)` - creates `Attribute` from color number (ANSI colors).
- `GetRGBColor(int32)` - creates `Attribute` from color number created the same way as `Hex()` function returns.
- `NewRGBColor(int32, int32, int32)` - creates `Attribute` from color numbers for red, green and blue values.

## Attribute font effect

There were 3 attributes for font effect, `AttrBold`, `AttrUnderline` and `AttrReverse`.

`tcell` supports more attributes, so they were added. All of these attributes have different values from before. However they can be used in the same way as before.

All the font effect attributes:
- `AttrBold`
- `AttrBlink`
- `AttrReverse`
- `AttrUnderline`
- `AttrDim`
- `AttrItalic`
- `AttrStrikeThrough`

## OutputMode

`OutputMode` in `termbox` was used to translate colors into the correct range. So for example in `OutputGrayscale` you had colors from 1 - 24 all representing gray colors in range 232 - 255, and white and black color.

`tcell` colors are 24bit and they are translated by the library into the color which can be read by terminal.

The original translation from `termbox` was included in GOCUI to be backward compatible. This is enabled in all the original modes: `OutputNormal`, `Output216`, `OutputGrayscale` and `Output256`.

`OutputTrue` is a new mode. It is recomended, because in this mode GOCUI doesn't do any kind of translation of the colors and pass them directly to `tcell`. If user wants to use true color in terminal and this mode doesn't work, it might be because of the terminal setup. `tcell` has a documentation what needs to be done, but in short `COLORTERM=truecolor` environment variable should help (see [_examples/colorstrue.go](./_examples/colorstrue.go)). Other way would be to have `TERM` environment variable having value with suffix `-truecolor`. To disable true color set `TCELL_TRUECOLOR=disable`.

## Keybinding

`termbox` had different way of handling input from terminal than `tcell`. This leads to some adjustement on how the keys are represented.
In general, all the keys in GOCUI should be presented from before, but the underlying values might be different. This could lead to some problems if a user uses different parser to create the `Key` for the keybinding. If using GOCUI parser, everything should be ok.

Mouse is handled differently in `tcell`, but translation was done to keep it in the same way as it was before. However this was harder to test due to different behaviour across the platforms, so if anything is missing or not working, please report.
