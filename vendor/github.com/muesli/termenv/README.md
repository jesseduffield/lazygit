<p align="center">
    <img src="https://stuff.charm.sh/termenv.png" width="480" alt="termenv Logo">
    <br />
    <a href="https://github.com/muesli/termenv/releases"><img src="https://img.shields.io/github/release/muesli/termenv.svg" alt="Latest Release"></a>
    <a href="https://godoc.org/github.com/muesli/termenv"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/muesli/termenv/actions"><img src="https://github.com/muesli/termenv/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://coveralls.io/github/muesli/termenv?branch=master"><img src="https://coveralls.io/repos/github/muesli/termenv/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/muesli/termenv"><img src="https://goreportcard.com/badge/muesli/termenv" alt="Go ReportCard"></a>
    <br />
    <img src="https://github.com/muesli/termenv/raw/master/examples/hello-world/hello-world.png" alt="Example terminal output">
</p>

`termenv` lets you safely use advanced styling options on the terminal. It
gathers information about the terminal environment in terms of its ANSI & color
support and offers you convenient methods to colorize and style your output,
without you having to deal with all kinds of weird ANSI escape sequences and
color conversions.

## Features

- RGB/TrueColor support
- Detects the supported color range of your terminal
- Automatically converts colors to the best matching, available colors
- Terminal theme (light/dark) detection
- Chainable syntax
- Nested styles

## Installation

```bash
go get github.com/muesli/termenv
```

## Usage

```go
output := termenv.NewOutput(os.Stdout)
```

`termenv` queries the terminal's capabilities it is running in, so you can
safely use advanced features, like RGB colors or ANSI styles. `output.Profile`
returns the supported profile:

- `termenv.Ascii` - no ANSI support detected, ASCII only
- `termenv.ANSI` - 16 color ANSI support
- `termenv.ANSI256` - Extended 256 color ANSI support
- `termenv.TrueColor` - RGB/TrueColor support

Alternatively, you can use `termenv.EnvColorProfile` which evaluates the
terminal like `ColorProfile`, but also respects the `NO_COLOR` and
`CLICOLOR_FORCE` environment variables.

You can also query the terminal for its color scheme, so you know whether your
app is running in a light- or dark-themed environment:

```go
// Returns terminal's foreground color
color := output.ForegroundColor()

// Returns terminal's background color
color := output.BackgroundColor()

// Returns whether terminal uses a dark-ish background
darkTheme := output.HasDarkBackground()
```

### Manual Profile Selection

If you don't want to rely on the automatic detection, you can manually select
the profile you want to use:

```go
output := termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
```

## Colors

`termenv` supports multiple color profiles: Ascii (black & white only),
ANSI (16 colors), ANSI Extended (256 colors), and TrueColor (24-bit RGB). Colors
will automatically be degraded to the best matching available color in the
desired profile:

`TrueColor` => `ANSI 256 Colors` => `ANSI 16 Colors` => `Ascii`

```go
s := output.String("Hello World")

// Supports hex values
// Will automatically degrade colors on terminals not supporting RGB
s.Foreground(output.Color("#abcdef"))
// but also supports ANSI colors (0-255)
s.Background(output.Color("69"))
// ...or the color.Color interface
s.Foreground(output.FromColor(color.RGBA{255, 128, 0, 255}))

// Combine fore- & background colors
s.Foreground(output.Color("#ffffff")).Background(output.Color("#0000ff"))

// Supports the fmt.Stringer interface
fmt.Println(s)
```

## Styles

You can use a chainable syntax to compose your own styles:

```go
s := output.String("foobar")

// Text styles
s.Bold()
s.Faint()
s.Italic()
s.CrossOut()
s.Underline()
s.Overline()

// Reverse swaps current fore- & background colors
s.Reverse()

// Blinking text
s.Blink()

// Combine multiple options
s.Bold().Underline()
```

## Template Helpers

`termenv` provides a set of helper functions to style your Go templates:

```go
// load template helpers
f := output.TemplateFuncs()
tpl := template.New("tpl").Funcs(f)

// apply bold style in a template
bold := `{{ Bold "Hello World" }}`

// examples for colorized templates
col := `{{ Color "#ff0000" "#0000ff" "Red on Blue" }}`
fg := `{{ Foreground "#ff0000" "Red Foreground" }}`
bg := `{{ Background "#0000ff" "Blue Background" }}`

// wrap styles
wrap := `{{ Bold (Underline "Hello World") }}`

// parse and render
tpl, err = tpl.Parse(bold)

var buf bytes.Buffer
tpl.Execute(&buf, nil)
fmt.Println(&buf)
```

Other available helper functions are: `Faint`, `Italic`, `CrossOut`,
`Underline`, `Overline`, `Reverse`, and `Blink`.

## Positioning

```go
// Move the cursor to a given position
output.MoveCursor(row, column)

// Save the cursor position
output.SaveCursorPosition()

// Restore a saved cursor position
output.RestoreCursorPosition()

// Move the cursor up a given number of lines
output.CursorUp(n)

// Move the cursor down a given number of lines
output.CursorDown(n)

// Move the cursor up a given number of lines
output.CursorForward(n)

// Move the cursor backwards a given number of cells
output.CursorBack(n)

// Move the cursor down a given number of lines and place it at the beginning
// of the line
output.CursorNextLine(n)

// Move the cursor up a given number of lines and place it at the beginning of
// the line
output.CursorPrevLine(n)
```

## Screen

```go
// Reset the terminal to its default style, removing any active styles
output.Reset()

// RestoreScreen restores a previously saved screen state
output.RestoreScreen()

// SaveScreen saves the screen state
output.SaveScreen()

// Switch to the altscreen. The former view can be restored with ExitAltScreen()
output.AltScreen()

// Exit the altscreen and return to the former terminal view
output.ExitAltScreen()

// Clear the visible portion of the terminal
output.ClearScreen()

// Clear the current line
output.ClearLine()

// Clear a given number of lines
output.ClearLines(n)

// Set the scrolling region of the terminal
output.ChangeScrollingRegion(top, bottom)

// Insert the given number of lines at the top of the scrollable region, pushing
// lines below down
output.InsertLines(n)

// Delete the given number of lines, pulling any lines in the scrollable region
// below up
output.DeleteLines(n)
```

## Session

```go
// SetWindowTitle sets the terminal window title
output.SetWindowTitle(title)

// SetForegroundColor sets the default foreground color
output.SetForegroundColor(color)

// SetBackgroundColor sets the default background color
output.SetBackgroundColor(color)

// SetCursorColor sets the cursor color
output.SetCursorColor(color)

// Hide the cursor
output.HideCursor()

// Show the cursor
output.ShowCursor()
```

## Mouse

```go
// Enable X10 mouse mode, only button press events are sent
output.EnableMousePress()

// Disable X10 mouse mode
output.DisableMousePress()

// Enable Mouse Tracking mode
output.EnableMouse()

// Disable Mouse Tracking mode
output.DisableMouse()

// Enable Hilite Mouse Tracking mode
output.EnableMouseHilite()

// Disable Hilite Mouse Tracking mode
output.DisableMouseHilite()

// Enable Cell Motion Mouse Tracking mode
output.EnableMouseCellMotion()

// Disable Cell Motion Mouse Tracking mode
output.DisableMouseCellMotion()

// Enable All Motion Mouse mode
output.EnableMouseAllMotion()

// Disable All Motion Mouse mode
output.DisableMouseAllMotion()
```

## Bracketed Paste

```go
// Enables bracketed paste mode
termenv.EnableBracketedPaste()

// Disables bracketed paste mode
termenv.DisableBracketedPaste()
```

## Optional Feature Support

| Terminal         | Alt Screen | Query Color Scheme | Query Cursor Position | Set Window Title | Change Cursor Color | Change Default Foreground Setting | Change Default Background Setting | Copy (OSC52) | Hyperlinks (OSC8) | Bracketed Paste |
| ---------------- | :--------: | :----------------: | :-------------------: | :--------------: | :-----------------: | :-------------------------------: | :-------------------------------: | :----------: | :---------------: | :-------------: |
| alacritty        |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |  ❌[^alacritty]   |        ✅       |
| foot             |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |        ✅         |        ✅       |
| kitty            |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |        ✅         |        ✅       |
| Konsole          |     ✅     |         ✅         |          ✅           |        ✅        |         ❌          |                ✅                 |                ✅                 | ❌[^konsole] |        ✅         |        ✅       |
| rxvt             |     ✅     |         ❌         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ❌      |        ❌         |        ✅       |
| urxvt            |     ✅     |         ❌         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |  ✅[^urxvt]  |        ❌         |        ✅       |
| screen           |     ✅     |      ⛔[^mux]      |          ✅           |        ✅        |         ❌          |                ❌                 |                ✅                 |      ✅      |    ❌[^screen]    |        ❌       |
| st               |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |        ❌         |        ✅       |
| tmux             |     ✅     |      ⛔[^mux]      |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |     ❌[^tmux]     |        ✅       |
| vte-based[^vte]  |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ❌                 |   ❌[^vte]   |        ✅         |        ✅       |
| wezterm          |     ✅     |         ✅         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |        ✅         |        ✅       |
| xterm            |     ✅     |         ✅         |          ✅           |        ✅        |         ❌          |                ❌                 |                ❌                 |      ✅      |        ❌         |        ✅       |
| Linux Console    |     ✅     |         ❌         |          ✅           |        ⛔        |         ❌          |                ❌                 |                ❌                 |      ⛔      |        ⛔         |        ❌       |
| Apple Terminal   |     ✅     |         ✅         |          ✅           |        ✅        |         ❌          |                ✅                 |                ✅                 |  ✅[^apple]  |        ❌         |        ✅       |
| iTerm            |     ✅     |         ✅         |          ✅           |        ✅        |         ❌          |                ❌                 |                ❌                 |      ✅      |        ✅         |        ✅       |
| Windows cmd      |     ✅     |         ❌         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ❌      |        ❌         |        ❌       |
| Windows Terminal |     ✅     |         ❌         |          ✅           |        ✅        |         ✅          |                ✅                 |                ✅                 |      ✅      |        ✅         |        ✅       |

[^vte]: This covers all vte-based terminals, including Gnome Terminal, guake, Pantheon Terminal, Terminator, Tilix, XFCE Terminal. OSC52 is not supported, see [issue#2495](https://gitlab.gnome.org/GNOME/vte/-/issues/2495).
[^mux]: Unavailable as multiplexers (like tmux or screen) can be connected to multiple terminals (with different color settings) at the same time.
[^urxvt]: Workaround for urxvt not supporting OSC52. See [this](https://unix.stackexchange.com/a/629485) for more information.
[^konsole]: OSC52 is not supported, for more info see [bug#372116](https://bugs.kde.org/show_bug.cgi?id=372116).
[^apple]: OSC52 works with a [workaround](https://github.com/roy2220/osc52pty).
[^tmux]: OSC8 is not supported, for more info see [issue#911](https://github.com/tmux/tmux/issues/911).
[^screen]: OSC8 is not supported, for more info see [bug#50952](https://savannah.gnu.org/bugs/index.php?50952).
[^alacritty]: OSC8 is not supported, for more info see [issue#922](https://github.com/alacritty/alacritty/issues/922).

You can help improve this list! Check out [how to](ansi_compat.md) and open an issue or pull request.

### Color Support

- 24-bit (RGB): alacritty, foot, iTerm, kitty, Konsole, st, tmux, vte-based, wezterm, Windows Terminal
- 8-bit (256): rxvt, screen, xterm, Apple Terminal
- 4-bit (16): Linux Console

## Platform Support

`termenv` works on Unix systems (like Linux, macOS, or BSD) and Windows. While
terminal applications on Unix support ANSI styling out-of-the-box, on Windows
you need to enable ANSI processing in your application first:

```go
    restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
    if err != nil {
        panic(err)
    }
    defer restoreConsole()
```

The above code is safe to include on non-Windows systems or when os.Stdout does
not refer to a terminal (e.g. in tests).


## Color Chart

![ANSI color chart](https://github.com/muesli/termenv/raw/master/examples/color-chart/color-chart.png)

You can find the source code used to create this chart in `termenv`'s examples.

## Related Projects

- [reflow](https://github.com/muesli/reflow) - ANSI-aware text operations
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - style definitions for nice terminal layouts 👄
- [ansi](https://github.com/muesli/ansi) - ANSI sequence helpers

## termenv in the Wild

Need some inspiration or just want to see how others are using `termenv`? Check
out these projects:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - a powerful little TUI framework 🏗
- [Glamour](https://github.com/charmbracelet/glamour) - stylesheet-based markdown rendering for your CLI apps 💇🏻‍♀️
- [Glow](https://github.com/charmbracelet/glow) - a markdown renderer for the command-line 💅🏻
- [duf](https://github.com/muesli/duf) - Disk Usage/Free Utility - a better 'df' alternative
- [gitty](https://github.com/muesli/gitty) - contextual information about your git projects
- [slides](https://github.com/maaslalani/slides) - terminal-based presentation tool

## Feedback

Got some feedback or suggestions? Please open an issue or drop me a note!

- [Twitter](https://twitter.com/mueslix)
- [The Fediverse](https://mastodon.social/@fribbledom)

## License

[MIT](https://github.com/muesli/termenv/raw/master/LICENSE)
