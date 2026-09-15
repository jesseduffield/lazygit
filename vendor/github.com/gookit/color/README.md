# CLI Color

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/color?style=flat-square)
[![Actions Status](https://github.com/gookit/color/workflows/action-tests/badge.svg)](https://github.com/gookit/color/actions)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7fef8d74c1d64afc99ce0f2c6d3f8af1)](https://www.codacy.com/gh/gookit/color/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gookit/color&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://pkg.go.dev/badge/github.com/gookit/color.svg)](https://pkg.go.dev/github.com/gookit/color?tab=overview)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/color)](https://github.com/gookit/color)
[![Coverage Status](https://coveralls.io/repos/github/gookit/color/badge.svg?branch=master)](https://coveralls.io/github/gookit/color?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/color)](https://goreportcard.com/report/github.com/gookit/color)

A command-line color library with 16/256/True color support, universal API methods and Windows support.

> **[中文说明](README.zh-CN.md)**

Basic color preview:

![basic-color](_examples/images/basic-color2.png)

Now, 256 colors and RGB colors have also been supported to work in Windows CMD and PowerShell:

![color-on-cmd-pwsh](_examples/images/color-on-cmd-pwsh.jpg)

## Features

  - Simple to use, zero dependencies
  - Supports rich color output: 16-color (4-bit), 256-color (8-bit), true color (24-bit, RGB)
    - 16-color output is the most commonly used and most widely supported, working on any Windows version
    - Since `v1.2.4` **the 256-color (8-bit), true color (24-bit) support windows CMD and PowerShell**
    - See [this gist](https://gist.github.com/XVilka/8346728) for information on true color support
  - Support converts `HEX` `HSL` value to RGB color
  - Generic API methods: `Print`, `Printf`, `Println`, `Sprint`, `Sprintf`
  - Supports HTML tag-style color rendering, such as `<green>message</> <fg=red;bg=blue>text</>`.
    - In addition to using built-in tags, it also supports custom color attributes
    - Custom color attributes support the use of 16 color names, 256 color values, rgb color values and hex color values
    - Support working on Windows `cmd` and `powerShell` terminal
  - Basic colors: `Bold`, `Black`, `White`, `Gray`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`
  - Additional styles: `Info`, `Note`, `Light`, `Error`, `Danger`, `Notice`, `Success`, `Comment`, `Primary`, `Warning`, `Question`, `Secondary`
  - Support by set `NO_COLOR` for disable color or use `FORCE_COLOR` for force open color render.
  - Support Rgb, 256, 16 color conversion

## GoDoc

See [godoc for github](https://pkg.go.dev/github.com/gookit/color)

## Install

```bash
go get github.com/gookit/color
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/gookit/color"
)

func main() {
	// quick use package func
	color.Redp("Simple to use color")
	color.Redln("Simple to use color")
	color.Greenp("Simple to use color\n")
	color.Cyanln("Simple to use color")
	color.Yellowln("Simple to use color")

	// quick use like fmt.Print*
	color.Red.Println("Simple to use color")
	color.Green.Print("Simple to use color\n")
	color.Cyan.Printf("Simple to use %s\n", "color")
	color.Yellow.Printf("Simple to use %s\n", "color")

	// use like func
	red := color.FgRed.Render
	green := color.FgGreen.Render
	fmt.Printf("%s line %s library\n", red("Command"), green("color"))

	// custom color
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")

	// can also:
	color.Style{color.FgCyan, color.OpBold}.Println("custom color style")

	// internal theme/style:
	color.Info.Tips("message")
	color.Info.Prompt("message")
	color.Info.Println("message")
	color.Warn.Println("message")
	color.Error.Println("message")

	// use style tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")
	// Custom label attr: Supports the use of 16 color names, 256 color values, rgb color values and hex color values
	color.Println("<fg=11aa23>he</><bg=120,35,156>llo</>, <fg=167;bg=232>wel</><fg=red>come</>")

	// apply a style tag
	color.Tag("info").Println("info style text")

	// prompt message
	color.Info.Prompt("prompt style message")
	color.Warn.Prompt("prompt style message")

	// tips message
	color.Info.Tips("tips style message")
	color.Warn.Tips("tips style message")
}
```

Run demo: `go run ./_examples/demo.go`

![colored-out](_examples/images/color-demo.jpg)

## Basic/16 color

Supported on any Windows version. Provide generic API methods: `Print`, `Printf`, `Println`, `Sprint`, `Sprintf`

```go
color.Bold.Println("bold message")
color.Cyan.Println("yellow message")
color.Yellow.Println("yellow message")
color.Magenta.Println("yellow message")

// Only use foreground color
color.FgCyan.Printf("Simple to use %s\n", "color")
// Only use background color
color.BgRed.Printf("Simple to use %s\n", "color")
```

Run demo: `go run ./_examples/color_16.go`

![basic-color](_examples/images/basic-color.png)

### Custom build color

```go
// Full custom: foreground, background, option
myStyle := color.New(color.FgWhite, color.BgBlack, color.OpBold)
myStyle.Println("custom color style")

// can also:
color.Style{color.FgCyan, color.OpBold}.Println("custom color style")
```

custom set console settings:

```go
// set console color
color.Set(color.FgCyan)

// print message
fmt.Print("message")

// reset console settings
color.Reset()
```

### Additional styles

provide generic API methods: `Print`, `Printf`, `Println`, `Sprint`, `Sprintf`

print message use defined style:

```go
color.Info.Println("Info message")
color.Notice.Println("Notice message")
color.Error.Println("Error message")
// ...
```

Run demo: `go run ./_examples/theme_basic.go`

![theme-basic](_examples/images/theme-basic.png)

**Tips style**

```go
color.Info.Tips("Info tips message")
color.Notice.Tips("Notice tips message")
color.Error.Tips("Error tips message")
color.Secondary.Tips("Secondary tips message")
```

Run demo: `go run ./_examples/theme_tips.go`

![theme-tips](_examples/images/theme-tips.png)

**Prompt Style**

```go
color.Info.Prompt("Info prompt message")
color.Notice.Prompt("Notice prompt message")
color.Error.Prompt("Error prompt message")
// ...
```

Run demo: `go run ./_examples/theme_prompt.go`

![theme-prompt](_examples/images/theme-prompt.png)

**Block Style**

```go
color.Danger.Block("Danger block message")
color.Warn.Block("Warn block message")
// ...
```

Run demo: `go run ./_examples/theme_block.go`

![theme-block](_examples/images/theme-block.png)

## 256-color usage

> 256 colors support Windows CMD, PowerShell environment after `v1.2.4`

### Set the foreground or background color

- `color.C256(val uint8, isBg ...bool) Color256`

```go
c := color.C256(132) // fg color
c.Println("message")
c.Printf("format %s", "message")

c := color.C256(132, true) // bg color
c.Println("message")
c.Printf("format %s", "message")
```

### 256-color style

Can be used to set foreground and background colors at the same time.

- `S256(fgAndBg ...uint8) *Style256`

```go
s := color.S256(32, 203)
s.Println("message")
s.Printf("format %s", "message")
```

with options:

```go
s := color.S256(32, 203)
s.SetOpts(color.Opts{color.OpBold})

s.Println("style with options")
s.Printf("style with %s\n", "options")
```

Run demo: `go run ./_examples/color_256.go`

![color-tags](_examples/images/color-256.png)

## RGB/True color

> RGB colors support Windows `CMD`, `PowerShell` environment after `v1.2.4`

**Preview:**

> Run demo: `Run demo: go run ./_examples/color_rgb.go`

![color-rgb](_examples/images/color-rgb.png)

example:

```go
color.RGB(30, 144, 255).Println("message. use RGB number")

color.HEX("#1976D2").Println("blue-darken")
color.HEX("#D50000", true).Println("red-accent. use HEX style")

color.RGBStyleFromString("213,0,0").Println("red-accent. use RGB number")
color.HEXStyle("eee", "D50000").Println("deep-purple color")
```

### Set the foreground or background color

- `color.RGB(r, g, b uint8, isBg ...bool) RGBColor`

```go
c := color.RGB(30,144,255) // fg color
c.Println("message")
c.Printf("format %s", "message")

c := color.RGB(30,144,255, true) // bg color
c.Println("message")
c.Printf("format %s", "message")
```

Create a style from an hexadecimal color string:

- `color.HEX(hex string, isBg ...bool) RGBColor`

```go
c := color.HEX("ccc") // can also: "cccccc" "#cccccc"
c.Println("message")
c.Printf("format %s", "message")

c = color.HEX("aabbcc", true) // as bg color
c.Println("message")
c.Printf("format %s", "message")
```

### RGB color style

Can be used to set the foreground and background colors at the same time.

- `color.NewRGBStyle(fg RGBColor, bg ...RGBColor) *RGBStyle`

```go
s := color.NewRGBStyle(RGB(20, 144, 234), RGB(234, 78, 23))
s.Println("message")
s.Printf("format %s", "message")
```

Create a style from an hexadecimal color string:

- `color.HEXStyle(fg string, bg ...string) *RGBStyle`

```go
s := color.HEXStyle("11aa23", "eee")
s.Println("message")
s.Printf("format %s", "message")
```

with options:

```go
s := color.HEXStyle("11aa23", "eee")
s.SetOpts(color.Opts{color.OpBold})

s.Println("style with options")
s.Printf("style with %s\n", "options")
```

## HTML-like tag usage

`Print,Printf,Println` functions support auto parse and render color tags.

```go
	text := `
  <mga1>gookit/color:</>
     A <green>command-line</> 
     <cyan>color library</> with <fg=167;bg=232>256-color</>
     and <fg=11aa23;op=bold>True-color</> support,
     <fg=mga;op=i>universal API</> methods
     and <cyan>Windows</> support.
`
	color.Print(text)
```

Preview, code please see [_examples/demo_tag.go](_examples/demo_tag.go):

![demo_tag](_examples/images/demo_tag.png)

**Tag formats:**

- Use built in tags: `<TAG_NAME>CONTENT</>` e.g: `<info>message</>`
- Custom tag attributes: `<fg=VALUE;bg=VALUE;op=VALUES>CONTENT</>` e.g: `<fg=167;bg=232>wel</>`

> **Supported** on Windows `cmd.exe` `PowerShell`.

Examples:

```go
// use style tag
color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")
color.Println("<suc>hello</>")
color.Println("<error>hello</>")
color.Println("<warning>hello</>")

// custom color attributes
color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")

// Custom label attr: Supports the use of 16 color names, 256 color values, rgb color values and hex color values
color.Println("<fg=11aa23>he</><bg=120,35,156>llo</>, <fg=167;bg=232>wel</><fg=red>come</>")
```
 
### Tag attributes

tag attributes format:

```text
attr format:
 // VALUE please see var: FgColors, BgColors, AllOptions
 "fg=VALUE;bg=VALUE;op=VALUE"

16 color:
 "fg=yellow"
 "bg=red"
 "op=bold,underscore" // option is allow multi value
 "fg=white;bg=blue;op=bold"
 "fg=white;op=bold,underscore"

256 color:
 "fg=167"
 "fg=167;bg=23"
 "fg=167;bg=23;op=bold"
 
True color:
 // hex
 "fg=fc1cac"
 "fg=fc1cac;bg=c2c3c4"
 // r,g,b
 "fg=23,45,214"
 "fg=23,45,214;bg=109,99,88"
```

> tag attributes parse please see `func ParseCodeFromAttr()`

### Built-in tags

Built-in tags please see var `colorTags` in [color_tag.go](color_tag.go)

```go
// use style tag
color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")
color.Println("<suc>hello</>")
color.Println("<error>hello</>")
```

Run demo: `go run ./_examples/color_tag.go`

![color-tags](_examples/images/color-tags.png)

**Use `color.Tag` build message**:

```go
// set a style tag
color.Tag("info").Print("info style text")
color.Tag("info").Printf("%s style text", "info")
color.Tag("info").Println("info style text")
```

## Color convert

Supports conversion between Rgb, 256, 16 colors, `Rgb <=> 256 <=> 16`

```go
basic := color.Red
basic.Println("basic color")

c256 := color.Red.C256()
c256.Println("256 color")
c256.C16().Println("basic color")

rgb := color.Red.RGB()
rgb.Println("rgb color")
rgb.C256().Println("256 color")
```

### Convert utils

`color` has many built-in color conversion utility functions.

```go
func Basic2hex(val uint8) string

func Bg2Fg(val uint8) uint8
func Fg2Bg(val uint8) uint8

func C256ToRgb(val uint8) (rgb []uint8)
func C256ToRgbV1(val uint8) (rgb []uint8)

func Hex2basic(hex string, asBg ...bool) uint8
func Hex2rgb(hex string) []int
func HexToRGB(hex string) []int
func HexToRgb(hex string) (rgb []int)

func HslIntToRgb(h, s, l int) (rgb []uint8)
func HslToRgb(h, s, l float64) (rgb []uint8)
func HsvToRgb(h, s, v int) (rgb []uint8)

func Rgb2ansi(r, g, b uint8, isBg bool) uint8
func Rgb2basic(r, g, b uint8, isBg bool) uint8
func Rgb2hex(rgb []int) string
func Rgb2short(r, g, b uint8) uint8
func RgbTo256(r, g, b uint8) uint8
func RgbTo256Table() map[string]uint8
func RgbToAnsi(r, g, b uint8, isBg bool) uint8
func RgbToHex(rgb []int) string
func RgbToHsl(r, g, b uint8) []float64
func RgbToHslInt(r, g, b uint8) []int
```

**Convert to `RGBColor`**:

- `func RGBFromSlice(rgb []uint8, isBg ...bool) RGBColor`
- `func RGBFromString(rgb string, isBg ...bool) RGBColor`
- `func HEX(hex string, isBg ...bool) RGBColor`
- `func HSL(h, s, l float64, isBg ...bool) RGBColor`
- `func HSLInt(h, s, l int, isBg ...bool) RGBColor`

## Util functions

There are some useful functions reference

- `Disable()` disable color render
- `SetOutput(io.Writer)` custom set the colored text output writer
- `ForceOpenColor()` force open color render
- `Colors2code(colors ...Color) string` Convert colors to code. return like "32;45;3"
- `ClearCode(str string) string` Use for clear color codes
- `ClearTag(s string) string` clear all color html-tag for a string
- `IsConsole(w io.Writer)` Determine whether w is one of stderr, stdout, stdin

> More useful func please see https://pkg.go.dev/github.com/gookit/color

### Detect color level

`color` automatically checks the color levels supported by the current environment.

```go
// Level is the color level supported by a terminal.
type Level = terminfo.ColorLevel

// terminal color available level alias of the terminfo.ColorLevel*
const (
	LevelNo  = terminfo.ColorLevelNone     // not support color.
	Level16  = terminfo.ColorLevelBasic    // basic - 3/4 bit color supported
	Level256 = terminfo.ColorLevelHundreds // hundreds - 8-bit color supported
	LevelRgb = terminfo.ColorLevelMillions // millions - (24 bit)true color supported
)
```

- `func SupportColor() bool` Whether the current environment supports color output
- `func Support256Color() bool` Whether the current environment supports 256-color output
- `func SupportTrueColor() bool` Whether the current environment supports (RGB)True-color output
- `func TermColorLevel() Level` Get the currently supported color level


## Projects using color

Check out these projects, which use https://github.com/gookit/color :

- https://github.com/Delta456/box-cli-maker Make Highly Customized Boxes for your CLI
- https://github.com/flipped-aurora/gin-vue-admin 基于gin+vue搭建的（中）后台系统框架
- https://github.com/JanDeDobbeleer/oh-my-posh A prompt theme engine for any shell.
- https://github.com/jesseduffield/lazygit Simple terminal UI for git commands
- https://github.com/olivia-ai/olivia 💁‍♀️Your new best friend powered by an artificial neural network  
- https://github.com/pterm/pterm PTerm is a modern Go module to beautify console output. Featuring charts, progressbars, tables, trees, etc.
- https://github.com/securego/gosec Golang security checker
- https://github.com/TNK-Studio/lazykube ⎈ The lazier way to manage kubernetes.
- [+ See More](https://pkg.go.dev/github.com/gookit/color?tab=importedby)

## Gookit packages

  - [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
  - [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
  - [gookit/gcli](https://github.com/gookit/gcli) build CLI application, tool library, running CLI commands
  - [gookit/slog](https://github.com/gookit/slog) Concise and extensible go log library
  - [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
  - [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
  - [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
  - [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
  - [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
  - [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
  - [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
  - More, please see https://github.com/gookit

## See also

  - [inhere/console](https://github.com/inhere/php-console)
  - [xo/terminfo](https://github.com/xo/terminfo)
  - [beego/bee](https://github.com/beego/bee)
  - [issue9/term](https://github.com/issue9/term)
  - [muesli/termenv](https://github.com/muesli/termenv)
  - [ANSI escape code](https://en.wikipedia.org/wiki/ANSI_escape_code)
  - [Standard ANSI color map](https://conemu.github.io/en/AnsiEscapeCodes.html#Standard_ANSI_color_map)
  - [Terminal Colors](https://gist.github.com/XVilka/8346728)

## License

[MIT](/LICENSE)
