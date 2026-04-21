# CLI Color

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/color?style=flat-square)
[![Actions Status](https://github.com/gookit/color/workflows/action-tests/badge.svg)](https://github.com/gookit/color/actions)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7fef8d74c1d64afc99ce0f2c6d3f8af1)](https://www.codacy.com/gh/gookit/color/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gookit/color&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://pkg.go.dev/badge/github.com/gookit/color.svg)](https://pkg.go.dev/github.com/gookit/color?tab=overview)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/color)](https://github.com/gookit/color)
[![Coverage Status](https://coveralls.io/repos/github/gookit/color/badge.svg?branch=master)](https://coveralls.io/github/gookit/color?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/color)](https://goreportcard.com/report/github.com/gookit/color)

Golang下的命令行色彩使用库, 拥有丰富的色彩(16/256/True)渲染输出，通用的API方法，兼容Windows系统

> **[EN README](README.md)**

基本颜色预览：

![basic-color](_examples/images/basic-color2.png)

现在，256色和RGB色彩也已经支持windows CMD和PowerShell中工作：

![color-on-cmd-pwsh](_examples/images/color-on-cmd-pwsh.jpg)

## 功能特色

  - 使用简单方便
  - 支持丰富的颜色输出, 16色(4bit)，256色(8bit)，RGB色彩(24bit, RGB)
    - 16色(4bit)是最常用和支持最广的，支持Windows `cmd.exe`
    - 自 `v1.2.4` 起 **256色(8bit)，RGB色彩(24bit)均支持Windows CMD和PowerShell终端**
    - 请查看 [this gist](https://gist.github.com/XVilka/8346728) 了解支持RGB色彩的终端
  - 支持转换 `HEX` `HSL` 等为RGB色彩
  - 提供通用的API方法：`Print` `Printf` `Println` `Sprint` `Sprintf`
  - 同时支持html标签式的颜色渲染，除了使用内置标签，同时支持自定义颜色属性
    - 例如: `this an <green>message</> <fg=red;bg=blue>text</>` 标签内部文本将会渲染对应色彩
    - 自定义颜色属性: 支持使用16色彩名称，256色彩值，rgb色彩值以及hex色彩值
  - 基础色彩: `Bold` `Black` `White` `Gray` `Red` `Green` `Yellow` `Blue` `Magenta` `Cyan`
  - 扩展风格: `Info` `Note` `Light` `Error` `Danger` `Notice` `Success` `Comment` `Primary` `Warning` `Question` `Secondary`
  - 支持通过设置环境变量 `NO_COLOR` 来禁用色彩，或者使用 `FORCE_COLOR` 来强制使用色彩渲染.
  - 支持 Rgb, 256, 16 色彩之间的互相转换
  - 支持Linux、Mac，同时兼容Windows系统环境

## GoDoc

[godoc for github](https://pkg.go.dev/github.com/gookit/color)

## 安装

```bash
go get github.com/gookit/color
```

## 快速开始

如下，引入当前包就可以快速的使用

```go
package main

import (
	"fmt"
	
	"github.com/gookit/color"
)

func main() {
	// 简单快速的使用，跟 fmt.Print* 类似
	color.Redp("Simple to use color")
	color.Redln("Simple to use color")
	color.Greenp("Simple to use color\n")
	color.Cyanln("Simple to use color")
	color.Yellowln("Simple to use color")

	// 简单快速的使用，跟 fmt.Print* 类似
	color.Red.Println("Simple to use color")
	color.Green.Print("Simple to use color\n")
	color.Cyan.Printf("Simple to use %s\n", "color")
	color.Yellow.Printf("Simple to use %s\n", "color")

	// use like func
	red := color.FgRed.Render
	green := color.FgGreen.Render
	fmt.Printf("%s line %s library\n", red("Command"), green("color"))

	// 自定义颜色
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")

	// 也可以:
	color.Style{color.FgCyan, color.OpBold}.Println("custom color style")
	
	// internal style:
	color.Info.Println("message")
	color.Warn.Println("message")
	color.Error.Println("message")
	
	// 使用内置颜色标签
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")
	// 自定义标签: 支持使用16色彩名称，256色彩值，rgb色彩值以及hex色彩值
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

> 运行 demo: `go run ./_examples/demo.go`

![colored-out](_examples/images/color-demo.jpg)

## 基础颜色(16-color)

提供通用的API方法：`Print` `Printf` `Println` `Sprint` `Sprintf`

> 支持在windows `cmd.exe`  `powerShell` 等终端使用

```go
color.Bold.Println("bold message")
color.Black.Println("bold message")
color.White.Println("bold message")
// ...

// Only use foreground color
color.FgCyan.Printf("Simple to use %s\n", "color")
// Only use background color
color.BgRed.Printf("Simple to use %s\n", "color")
```

> 运行demo: `go run ./_examples/color_16.go`

![basic-color](_examples/images/basic-color.png)

### 构建风格

```go
// 仅设置前景色
color.FgCyan.Printf("Simple to use %s\n", "color")
// 仅设置背景色
color.BgRed.Printf("Simple to use %s\n", "color")

// 完全自定义: 前景色 背景色 选项
style := color.New(color.FgWhite, color.BgBlack, color.OpBold)
style.Println("custom color style")

// 也可以:
color.Style{color.FgCyan, color.OpBold}.Println("custom color style")
```

直接设置控制台属性：

```go
// 设置console颜色
color.Set(color.FgCyan)

// 输出信息
fmt.Print("message")

// 重置console颜色
color.Reset()
```

> 当然，color已经内置丰富的色彩风格支持

### 扩展风格方法 

提供通用的API方法：`Print` `Printf` `Println` `Sprint` `Sprintf`

> 支持在windows `cmd.exe`  `powerShell` 等终端使用

基础使用：

```go
// print message
color.Info.Println("Info message")
color.Note.Println("Note message")
color.Notice.Println("Notice message")
// ...
```

Run demo: `go run ./_examples/theme_basic.go`

![theme-basic](_examples/images/theme-basic.png)

**简约提示风格**

```go
color.Info.Tips("Info tips message")
color.Notice.Tips("Notice tips message")
color.Error.Tips("Error tips message")
// ...
```

Run demo: `go run ./_examples/theme_tips.go`

![theme-tips](_examples/images/theme-tips.png)

**着重提示风格**

```go
color.Info.Prompt("Info prompt message")
color.Error.Prompt("Error prompt message")
color.Danger.Prompt("Danger prompt message")
```

Run demo: `go run ./_examples/theme_prompt.go`

![theme-prompt](_examples/images/theme-prompt.png)

**强调提示风格**

```go
color.Warn.Block("Warn block message")
color.Debug.Block("Debug block message")
color.Question.Block("Question block message")
```

Run demo: `go run ./_examples/theme_block.go`

![theme-block](_examples/images/theme-block.png)

## 256 色彩使用

> 256色彩在 `v1.2.4` 后支持Windows CMD,PowerShell 环境

### 使用前景或后景色
 
- `color.C256(val uint8, isBg ...bool) Color256`

```go
c := color.C256(132) // fg color
c.Println("message")
c.Printf("format %s", "message")

c := color.C256(132, true) // bg color
c.Println("message")
c.Printf("format %s", "message")
```

### 使用256 色彩风格

> 可同时设置前景和背景色
 
- `color.S256(fgAndBg ...uint8) *Style256`

```go
s := color.S256(32, 203)
s.Println("message")
s.Printf("format %s", "message")
```

可以同时添加选项设置:

```go
s := color.S256(32, 203)
s.SetOpts(color.Opts{color.OpBold})

s.Println("style with options")
s.Printf("style with %s\n", "options")
```

> 运行 demo: `go run ./_examples/color_256.go`

![color-tags](_examples/images/color-256.png)

## RGB/True色彩使用

> RGB色彩在 `v1.2.4` 后支持 Windows `CMD`, `PowerShell` 环境

**效果预览:**

> 运行 demo: `Run demo: go run ./_examples/color_rgb.go`

![color-rgb](_examples/images/color-rgb.png)

代码示例：

```go
color.RGB(30, 144, 255).Println("message. use RGB number")

color.HEX("#1976D2").Println("blue-darken")
color.HEX("#D50000", true).Println("red-accent. use HEX style")

color.RGBStyleFromString("213,0,0").Println("red-accent. use RGB number")
color.HEXStyle("eee", "D50000").Println("deep-purple color")
```

### 使用前景或后景色 

- `color.RGB(r, g, b uint8, isBg ...bool) RGBColor`

```go
c := color.RGB(30,144,255) // fg color
c.Println("message")
c.Printf("format %s", "message")

c := color.RGB(30,144,255, true) // bg color
c.Println("message")
c.Printf("format %s", "message")
```

- `color.HEX(hex string, isBg ...bool) RGBColor` 从16进制颜色创建

```go
c := color.HEX("ccc") // 也可以写为: "cccccc" "#cccccc"
c.Println("message")
c.Printf("format %s", "message")

c = color.HEX("aabbcc", true) // as bg color
c.Println("message")
c.Printf("format %s", "message")
```

### 使用RGB风格

> TIP: 可同时设置前景和背景色

- `color.NewRGBStyle(fg RGBColor, bg ...RGBColor) *RGBStyle`

```go
s := color.NewRGBStyle(RGB(20, 144, 234), RGB(234, 78, 23))
s.Println("message")
s.Printf("format %s", "message")
```

- `color.HEXStyle(fg string, bg ...string) *RGBStyle` 从16进制颜色创建

```go
s := color.HEXStyle("11aa23", "eee")
s.Println("message")
s.Printf("format %s", "message")
```

- 可以同时添加选项设置:

```go
s := color.HEXStyle("11aa23", "eee")
s.SetOpts(color.Opts{color.OpBold})

s.Println("style with options")
s.Printf("style with %s\n", "options")
```

## 使用颜色标签

`Print,Printf,Println` 等方法支持自动解析并渲染 HTML 风格的颜色标签

> **支持** 在windows `cmd.exe` `PowerShell` 使用

简单示例:

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

输出效果, 示例代码请看 [_examples/demo_tag.go](_examples/demo_tag.go):

![demo_tag](_examples/images/demo_tag.png)

**颜色标签格式:**

- 直接使用内置风格标签: `<TAG_NAME>CONTENT</>` e.g: `<info>message</>`
- 自定义标签属性: `<fg=VALUE;bg=VALUE;op=VALUES>CONTENT</>` e.g: `<fg=167;bg=232>wel</>`

使用内置的颜色标签，可以非常方便简单的构建自己需要的任何格式

> 同时支持自定义颜色属性: 支持使用16色彩名称，256色彩值，rgb色彩值以及hex色彩值

```go
// 使用内置的 color tag
color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")
color.Println("<suc>hello</>")
color.Println("<error>hello</>")
color.Println("<warning>hello</>")

// 自定义颜色属性
color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")

// 自定义颜色属性: 支持使用16色彩名称，256色彩值，rgb色彩值以及hex色彩值
color.Println("<fg=11aa23>he</><bg=120,35,156>llo</>, <fg=167;bg=232>wel</><fg=red>come</>")
```

### 自定义标签属性

标签属性格式:

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

### 内置标签

内置标签请参见变量 `colorTags` 定义, 源文件 [color_tag.go](color_tag.go)

```go
// use style tag
color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")
color.Println("<suc>hello</>")
color.Println("<error>hello</>")
```

> 运行 demo: `go run ./_examples/color_tag.go`

![color-tags](_examples/images/color-tags.png)

**使用 `color.Tag` 包装标签**:

可以使用通用的输出API方法,给后面输出的文本信息加上给定的颜色风格标签

```go
// set a style tag
color.Tag("info").Print("info style text")
color.Tag("info").Printf("%s style text", "info")
color.Tag("info").Println("info style text")
```

## 颜色转换

支持 Rgb, 256, 16 色彩之间的互相转换 `Rgb <=> 256 <=> 16`

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

### 颜色转换方法

`color` 内置了许多颜色转换工具方法

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

**转换为 `RGBColor`**:

- `func RGBFromSlice(rgb []uint8, isBg ...bool) RGBColor`
- `func RGBFromString(rgb string, isBg ...bool) RGBColor`
- `func HEX(hex string, isBg ...bool) RGBColor`
- `func HSL(h, s, l float64, isBg ...bool) RGBColor`
- `func HSLInt(h, s, l int, isBg ...bool) RGBColor`

## 工具方法参考

一些有用的工具方法参考

- `Disable()` 禁用颜色渲染输出
- `SetOutput(io.Writer)` 自定义设置渲染后的彩色文本输出位置
- `ForceOpenColor()` 强制开启颜色渲染
- `ClearCode(str string) string` Use for clear color codes
- `Colors2code(colors ...Color) string` Convert colors to code. return like "32;45;3"
- `ClearTag(s string) string` clear all color html-tag for a string
- `IsConsole(w io.Writer)` Determine whether w is one of stderr, stdout, stdin
- 更多请查看文档 https://pkg.go.dev/github.com/gookit/color

### 检测支持的颜色级别

`color` 会自动检查当前环境支持的颜色级别

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

- `func SupportColor() bool` 当前环境是否支持色彩输出
- `func Support256Color() bool` 当前环境是否支持256色彩输出
- `func SupportTrueColor() bool` 当前环境是否支持(RGB)True色彩输出
- `func TermColorLevel() Level` 获取当前支持的颜色级别

## 使用Color的项目

看看这些使用了 https://github.com/gookit/color 的项目:

- https://github.com/Delta456/box-cli-maker Make Highly Customized Boxes for your CLI
- https://github.com/flipped-aurora/gin-vue-admin 基于gin+vue搭建的（中）后台系统框架
- https://github.com/JanDeDobbeleer/oh-my-posh A prompt theme engine for any shell.
- https://github.com/jesseduffield/lazygit Simple terminal UI for git commands
- https://github.com/olivia-ai/olivia 💁‍♀️Your new best friend powered by an artificial neural network
- https://github.com/pterm/pterm PTerm is a modern Go module to beautify console output. Featuring charts, progressbars, tables, trees, etc.
- https://github.com/securego/gosec Golang security checker
- https://github.com/TNK-Studio/lazykube ⎈ The lazier way to manage kubernetes.
- [+ See More](https://pkg.go.dev/github.com/gookit/color?tab=importedby)

## Gookit 工具包

  - [gookit/ini](https://github.com/gookit/ini) INI配置读取管理，支持多文件加载，数据覆盖合并, 解析ENV变量, 解析变量引用
  - [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
  - [gookit/gcli](https://github.com/gookit/gcli) Go的命令行应用，工具库，运行CLI命令，支持命令行色彩，用户交互，进度显示，数据格式化显示
  - [gookit/slog](https://github.com/gookit/slog) 简洁易扩展的go日志库
  - [gookit/event](https://github.com/gookit/event) Go实现的轻量级的事件管理、调度程序库, 支持设置监听器的优先级, 支持对一组事件进行监听
  - [gookit/cache](https://github.com/gookit/cache) 通用的缓存使用包装库，通过包装各种常用的驱动，来提供统一的使用API
  - [gookit/config](https://github.com/gookit/config) Go应用配置管理，支持多种格式（JSON, YAML, TOML, INI, HCL, ENV, Flags），多文件加载，远程文件加载，数据合并
  - [gookit/color](https://github.com/gookit/color) CLI 控制台颜色渲染工具库, 拥有简洁的使用API，支持16色，256色，RGB色彩渲染输出
  - [gookit/filter](https://github.com/gookit/filter) 提供对Golang数据的过滤，净化，转换
  - [gookit/validate](https://github.com/gookit/validate) Go通用的数据验证与过滤库，使用简单，内置大部分常用验证、过滤器
  - [gookit/goutil](https://github.com/gookit/goutil) Go 的一些工具函数，格式化，特殊处理，常用信息获取等
  - 更多请查看 https://github.com/gookit

## 参考项目

  - [inhere/console](https://github.com/inhere/php-console)
  - [muesli/termenv](https://github.com/muesli/termenv)
  - [xo/terminfo](https://github.com/xo/terminfo)
  - [beego/bee](https://github.com/beego/bee)
  - [issue9/term](https://github.com/issue9/term)
  - [ANSI转义序列](https://zh.wikipedia.org/wiki/ANSI转义序列)
  - [Standard ANSI color map](https://conemu.github.io/en/AnsiEscapeCodes.html#Standard_ANSI_color_map)
  - [Terminal Colors](https://gist.github.com/XVilka/8346728)

## License

MIT
