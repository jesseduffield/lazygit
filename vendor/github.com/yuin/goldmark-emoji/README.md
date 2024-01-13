goldmark-emoji
=========================

[![GoDev][godev-image]][godev-url]

[godev-image]: https://pkg.go.dev/badge/github.com/yuin/goldmark-emoji
[godev-url]: https://pkg.go.dev/github.com/yuin/goldmark-emoji

goldmark-emoji is an extension for the [goldmark](http://github.com/yuin/goldmark) 
that parses `:joy:` style emojis.

Installation
--------------------

```
go get github.com/yuin/goldmark-emoji
```

Usage
--------------------

```go
import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark-emoji"
    "github.com/yuin/goldmark-emoji/definition"
)

func main() {
    markdown := goldmark.New(
        goldmark.WithExtensions(
            emoji.Emoji,
        ),
    )
    source := `
    Joy :joy:
    `
    var buf bytes.Buffer
    if err := markdown.Convert([]byte(source), &buf); err != nil {
        panic(err)
    }
    fmt.Print(buf.String())
}
```

See `emoji_test.go` for detailed usage.

### Options

Options for the extension

| Option | Description |
| ------ | ----------- |
| `WithEmojis` | Definition of emojis. This defaults to github emoji set |
| `WithRenderingMethod` | `Entity` : renders as HTML entities, `Twemoji` : renders as an img tag that uses [twemoji](https://github.com/twitter/twemoji), `Func` : renders using a go function |
| `WithTwemojiTemplate` | Twemoji img tag printf template |
| `WithRendererFunc` | renders by a go function |



License
--------------------
MIT

Author
--------------------
Yusuke Inuzuka

