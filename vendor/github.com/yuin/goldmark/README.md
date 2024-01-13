goldmark
==========================================

[![https://pkg.go.dev/github.com/yuin/goldmark](https://pkg.go.dev/badge/github.com/yuin/goldmark.svg)](https://pkg.go.dev/github.com/yuin/goldmark)
[![https://github.com/yuin/goldmark/actions?query=workflow:test](https://github.com/yuin/goldmark/workflows/test/badge.svg?branch=master&event=push)](https://github.com/yuin/goldmark/actions?query=workflow:test)
[![https://coveralls.io/github/yuin/goldmark](https://coveralls.io/repos/github/yuin/goldmark/badge.svg?branch=master)](https://coveralls.io/github/yuin/goldmark)
[![https://goreportcard.com/report/github.com/yuin/goldmark](https://goreportcard.com/badge/github.com/yuin/goldmark)](https://goreportcard.com/report/github.com/yuin/goldmark)

> A Markdown parser written in Go. Easy to extend, standards-compliant, well-structured.

goldmark is compliant with CommonMark 0.30.

Motivation
----------------------
I needed a Markdown parser for Go that satisfies the following requirements:

- Easy to extend.
    - Markdown is poor in document expressions compared to other light markup languages such as reStructuredText.
    - We have extensions to the Markdown syntax, e.g. PHP Markdown Extra, GitHub Flavored Markdown.
- Standards-compliant.
    - Markdown has many dialects.
    - GitHub-Flavored Markdown is widely used and is based upon CommonMark, effectively mooting the question of whether or not CommonMark is an ideal specification.
        - CommonMark is complicated and hard to implement.
- Well-structured.
    - AST-based; preserves source position of nodes.
- Written in pure Go.

[golang-commonmark](https://gitlab.com/golang-commonmark/markdown) may be a good choice, but it seems to be a copy of [markdown-it](https://github.com/markdown-it).

[blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) is a fast and widely-used implementation, but is not CommonMark-compliant and cannot be extended from outside of the package, since its AST uses structs instead of interfaces.

Furthermore, its behavior differs from other implementations in some cases, especially regarding lists: [Deep nested lists don't output correctly #329](https://github.com/russross/blackfriday/issues/329), [List block cannot have a second line #244](https://github.com/russross/blackfriday/issues/244), etc.

This behavior sometimes causes problems. If you migrate your Markdown text from GitHub to blackfriday-based wikis, many lists will immediately be broken.

As mentioned above, CommonMark is complicated and hard to implement, so Markdown parsers based on CommonMark are few and far between.

Features
----------------------

- **Standards-compliant.**  goldmark is fully compliant with the latest [CommonMark](https://commonmark.org/) specification.
- **Extensible.**  Do you want to add a `@username` mention syntax to Markdown?
  You can easily do so in goldmark. You can add your AST nodes,
  parsers for block-level elements, parsers for inline-level elements,
  transformers for paragraphs, transformers for the whole AST structure, and
  renderers.
- **Performance.**  goldmark's performance is on par with that of cmark,
  the CommonMark reference implementation written in C.
- **Robust.**  goldmark is tested with `go test --fuzz`.
- **Built-in extensions.**  goldmark ships with common extensions like tables, strikethrough,
  task lists, and definition lists.
- **Depends only on standard libraries.**

Installation
----------------------
```bash
$ go get github.com/yuin/goldmark
```


Usage
----------------------
Import packages:

```go
import (
    "bytes"
    "github.com/yuin/goldmark"
)
```


Convert Markdown documents with the CommonMark-compliant mode:

```go
var buf bytes.Buffer
if err := goldmark.Convert(source, &buf); err != nil {
  panic(err)
}
```

With options
------------------------------

```go
var buf bytes.Buffer
if err := goldmark.Convert(source, &buf, parser.WithContext(ctx)); err != nil {
  panic(err)
}
```

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `parser.WithContext` | A `parser.Context` | Context for the parsing phase. |

Context options
----------------------

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `parser.WithIDs` | A `parser.IDs` | `IDs` allows you to change logics that are related to element id(ex: Auto heading id generation). |


Custom parser and renderer
--------------------------
```go
import (
    "bytes"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/yuin/goldmark/parser"
    "github.com/yuin/goldmark/renderer/html"
)

md := goldmark.New(
          goldmark.WithExtensions(extension.GFM),
          goldmark.WithParserOptions(
              parser.WithAutoHeadingID(),
          ),
          goldmark.WithRendererOptions(
              html.WithHardWraps(),
              html.WithXHTML(),
          ),
      )
var buf bytes.Buffer
if err := md.Convert(source, &buf); err != nil {
    panic(err)
}
```

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `goldmark.WithParser` | `parser.Parser`  | This option must be passed before `goldmark.WithParserOptions` and `goldmark.WithExtensions` |
| `goldmark.WithRenderer` | `renderer.Renderer`  | This option must be passed before `goldmark.WithRendererOptions` and `goldmark.WithExtensions`  |
| `goldmark.WithParserOptions` | `...parser.Option`  |  |
| `goldmark.WithRendererOptions` | `...renderer.Option` |  |
| `goldmark.WithExtensions` | `...goldmark.Extender`  |  |

Parser and Renderer options
------------------------------

### Parser options

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `parser.WithBlockParsers` | A `util.PrioritizedSlice` whose elements are `parser.BlockParser` | Parsers for parsing block level elements. |
| `parser.WithInlineParsers` | A `util.PrioritizedSlice` whose elements are `parser.InlineParser` | Parsers for parsing inline level elements. |
| `parser.WithParagraphTransformers` | A `util.PrioritizedSlice` whose elements are `parser.ParagraphTransformer` | Transformers for transforming paragraph nodes. |
| `parser.WithASTTransformers` | A `util.PrioritizedSlice` whose elements are `parser.ASTTransformer` | Transformers for transforming an AST. |
| `parser.WithAutoHeadingID` | `-` | Enables auto heading ids. |
| `parser.WithAttribute` | `-` | Enables custom attributes. Currently only headings supports attributes. |

### HTML Renderer options

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `html.WithWriter` | `html.Writer` | `html.Writer` for writing contents to an `io.Writer`. |
| `html.WithHardWraps` | `-` | Render newlines as `<br>`.|
| `html.WithXHTML` | `-` | Render as XHTML. |
| `html.WithUnsafe` | `-` | By default, goldmark does not render raw HTML or potentially dangerous links. With this option, goldmark renders such content as written. |

### Built-in extensions

- `extension.Table`
    - [GitHub Flavored Markdown: Tables](https://github.github.com/gfm/#tables-extension-)
- `extension.Strikethrough`
    - [GitHub Flavored Markdown: Strikethrough](https://github.github.com/gfm/#strikethrough-extension-)
- `extension.Linkify`
    - [GitHub Flavored Markdown: Autolinks](https://github.github.com/gfm/#autolinks-extension-)
- `extension.TaskList`
    - [GitHub Flavored Markdown: Task list items](https://github.github.com/gfm/#task-list-items-extension-)
- `extension.GFM`
    - This extension enables Table, Strikethrough, Linkify and TaskList.
    - This extension does not filter tags defined in [6.11: Disallowed Raw HTML (extension)](https://github.github.com/gfm/#disallowed-raw-html-extension-).
    If you need to filter HTML tags, see [Security](#security).
    - If you need to parse github emojis, you can use [goldmark-emoji](https://github.com/yuin/goldmark-emoji) extension.
- `extension.DefinitionList`
    - [PHP Markdown Extra: Definition lists](https://michelf.ca/projects/php-markdown/extra/#def-list)
- `extension.Footnote`
    - [PHP Markdown Extra: Footnotes](https://michelf.ca/projects/php-markdown/extra/#footnotes)
- `extension.Typographer`
    - This extension substitutes punctuations with typographic entities like [smartypants](https://daringfireball.net/projects/smartypants/).
- `extension.CJK`
    - This extension is a shortcut for CJK related functionalities.

### Attributes
The `parser.WithAttribute` option allows you to define attributes on some elements.

Currently only headings support attributes.

**Attributes are being discussed in the
[CommonMark forum](https://talk.commonmark.org/t/consistent-attribute-syntax/272).
This syntax may possibly change in the future.**


#### Headings

```
## heading ## {#id .className attrName=attrValue class="class1 class2"}

## heading {#id .className attrName=attrValue class="class1 class2"}
```

```
heading {#id .className attrName=attrValue}
============
```

### Table extension
The Table extension implements [Table(extension)](https://github.github.com/gfm/#tables-extension-), as
defined in [GitHub Flavored Markdown Spec](https://github.github.com/gfm/).

Specs are defined for XHTML, so specs use some deprecated attributes for HTML5.

You can override alignment rendering method via options.

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `extension.WithTableCellAlignMethod` | `extension.TableCellAlignMethod` | Option indicates how are table cells aligned. |

### Typographer extension

The Typographer extension translates plain ASCII punctuation characters into typographic-punctuation HTML entities.

Default substitutions are:

| Punctuation | Default entity |
| ------------ | ---------- |
| `'`           | `&lsquo;`, `&rsquo;` |
| `"`           | `&ldquo;`, `&rdquo;` |
| `--`       | `&ndash;` |
| `---`      | `&mdash;` |
| `...`      | `&hellip;` |
| `<<`       | `&laquo;` |
| `>>`       | `&raquo;` |

You can override the default substitutions via `extensions.WithTypographicSubstitutions`:

```go
markdown := goldmark.New(
    goldmark.WithExtensions(
        extension.NewTypographer(
            extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
                extension.LeftSingleQuote:  []byte("&sbquo;"),
                extension.RightSingleQuote: nil, // nil disables a substitution
            }),
        ),
    ),
)
```

### Linkify extension

The Linkify extension implements [Autolinks(extension)](https://github.github.com/gfm/#autolinks-extension-), as
defined in [GitHub Flavored Markdown Spec](https://github.github.com/gfm/).

Since the spec does not define details about URLs, there are numerous ambiguous cases.

You can override autolinking patterns via options.

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `extension.WithLinkifyAllowedProtocols` | `[][]byte` | List of allowed protocols such as `[][]byte{ []byte("http:") }` |
| `extension.WithLinkifyURLRegexp` | `*regexp.Regexp` | Regexp that defines URLs, including protocols |
| `extension.WithLinkifyWWWRegexp` | `*regexp.Regexp` | Regexp that defines URL starting with `www.`. This pattern corresponds to [the extended www autolink](https://github.github.com/gfm/#extended-www-autolink) |
| `extension.WithLinkifyEmailRegexp` | `*regexp.Regexp` | Regexp that defines email addresses` |

Example, using [xurls](https://github.com/mvdan/xurls):

```go
import "mvdan.cc/xurls/v2"

markdown := goldmark.New(
    goldmark.WithRendererOptions(
        html.WithXHTML(),
        html.WithUnsafe(),
    ),
    goldmark.WithExtensions(
        extension.NewLinkify(
            extension.WithLinkifyAllowedProtocols([][]byte{
                []byte("http:"),
                []byte("https:"),
            }),
            extension.WithLinkifyURLRegexp(
                xurls.Strict,
            ),
        ),
    ),
)
```

### Footnotes extension

The Footnote extension implements [PHP Markdown Extra: Footnotes](https://michelf.ca/projects/php-markdown/extra/#footnotes).

This extension has some options:

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `extension.WithFootnoteIDPrefix` | `[]byte` |  a prefix for the id attributes.|
| `extension.WithFootnoteIDPrefixFunction` | `func(gast.Node) []byte` |  a function that determines the id attribute for given Node.|
| `extension.WithFootnoteLinkTitle` | `[]byte` |  an optional title attribute for footnote links.|
| `extension.WithFootnoteBacklinkTitle` | `[]byte` |  an optional title attribute for footnote backlinks. |
| `extension.WithFootnoteLinkClass` | `[]byte` |  a class for footnote links. This defaults to `footnote-ref`. |
| `extension.WithFootnoteBacklinkClass` | `[]byte` |  a class for footnote backlinks. This defaults to `footnote-backref`. |
| `extension.WithFootnoteBacklinkHTML` | `[]byte` |  a class for footnote backlinks. This defaults to `&#x21a9;&#xfe0e;`. |

Some options can have special substitutions. Occurrences of “^^” in the string will be replaced by the corresponding footnote number in the HTML output. Occurrences of “%%” will be replaced by a number for the reference (footnotes can have multiple references).

`extension.WithFootnoteIDPrefix` and `extension.WithFootnoteIDPrefixFunction` are useful if you have multiple Markdown documents displayed inside one HTML document to avoid footnote ids to clash each other.

`extension.WithFootnoteIDPrefix` sets fixed id prefix, so you may write codes like the following:

```go
for _, path := range files {
    source := readAll(path)
    prefix := getPrefix(path)

    markdown := goldmark.New(
        goldmark.WithExtensions(
            NewFootnote(
                WithFootnoteIDPrefix([]byte(path)),
            ),
        ),
    )
    var b bytes.Buffer
    err := markdown.Convert(source, &b)
    if err != nil {
        t.Error(err.Error())
    }
}
```

`extension.WithFootnoteIDPrefixFunction` determines an id prefix by calling given function, so you may write codes like the following:

```go
markdown := goldmark.New(
    goldmark.WithExtensions(
        NewFootnote(
                WithFootnoteIDPrefixFunction(func(n gast.Node) []byte {
                    v, ok := n.OwnerDocument().Meta()["footnote-prefix"]
                    if ok {
                        return util.StringToReadOnlyBytes(v.(string))
                    }
                    return nil
                }),
        ),
    ),
)

for _, path := range files {
    source := readAll(path)
    var b bytes.Buffer

    doc := markdown.Parser().Parse(text.NewReader(source))
    doc.Meta()["footnote-prefix"] = getPrefix(path)
    err := markdown.Renderer().Render(&b, source, doc)
}
```

You can use [goldmark-meta](https://github.com/yuin/goldmark-meta) to define a id prefix in the markdown document:


```markdown
---
title: document title
slug: article1
footnote-prefix: article1
---

# My article

```

### CJK extension
CommonMark gives compatibilities a high priority and original markdown was designed by westerners. So CommonMark lacks considerations for languages like CJK.

This extension provides additional options for CJK users.

| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| `extension.WithEastAsianLineBreaks` | `-` | Soft line breaks are rendered as a newline. Some asian users will see it as an unnecessary space. With this option, soft line breaks between east asian wide characters will be ignored. |
| `extension.WithEscapedSpace` | `-` | Without spaces around an emphasis started with east asian punctuations, it is not interpreted as an emphasis(as defined in CommonMark spec). With this option, you can avoid this inconvenient behavior by putting 'not rendered' spaces around an emphasis like `太郎は\ **「こんにちわ」**\ といった`. |

 
Security
--------------------
By default, goldmark does not render raw HTML or potentially-dangerous URLs.
If you need to gain more control over untrusted contents, it is recommended that you
use an HTML sanitizer such as [bluemonday](https://github.com/microcosm-cc/bluemonday).

Benchmark
--------------------
You can run this benchmark in the `_benchmark` directory.

### against other golang libraries

blackfriday v2 seems to be the fastest, but as it is not CommonMark compliant, its performance cannot be directly compared to that of the CommonMark-compliant libraries.

goldmark, meanwhile, builds a clean, extensible AST structure, achieves full compliance with
CommonMark, and consumes less memory, all while being reasonably fast.

- MBP 2019 13″(i5, 16GB), Go1.17

```
BenchmarkMarkdown/Blackfriday-v2-8                   302           3743747 ns/op         3290445 B/op      20050 allocs/op
BenchmarkMarkdown/GoldMark-8                         280           4200974 ns/op         2559738 B/op      13435 allocs/op
BenchmarkMarkdown/CommonMark-8                       226           5283686 ns/op         2702490 B/op      20792 allocs/op
BenchmarkMarkdown/Lute-8                              12          92652857 ns/op        10602649 B/op      40555 allocs/op
BenchmarkMarkdown/GoMarkdown-8                        13          81380167 ns/op         2245002 B/op      22889 allocs/op
```

### against cmark (CommonMark reference implementation written in C)

- MBP 2019 13″(i5, 16GB), Go1.17

```
----------- cmark -----------
file: _data.md
iteration: 50
average: 0.0044073057 sec
------- goldmark -------
file: _data.md
iteration: 50
average: 0.0041611990 sec
```

As you can see, goldmark's performance is on par with cmark's.

Extensions
--------------------

- [goldmark-meta](https://github.com/yuin/goldmark-meta): A YAML metadata
  extension for the goldmark Markdown parser.
- [goldmark-highlighting](https://github.com/yuin/goldmark-highlighting): A syntax-highlighting extension
  for the goldmark markdown parser.
- [goldmark-emoji](https://github.com/yuin/goldmark-emoji): An emoji
  extension for the goldmark Markdown parser.
- [goldmark-mathjax](https://github.com/litao91/goldmark-mathjax): Mathjax support for the goldmark markdown parser
- [goldmark-pdf](https://github.com/stephenafamo/goldmark-pdf): A PDF renderer that can be passed to `goldmark.WithRenderer()`.
- [goldmark-hashtag](https://github.com/abhinav/goldmark-hashtag): Adds support for `#hashtag`-based tagging to goldmark.
- [goldmark-wikilink](https://github.com/abhinav/goldmark-wikilink): Adds support for `[[wiki]]`-style links to goldmark.
- [goldmark-toc](https://github.com/abhinav/goldmark-toc): Adds support for generating tables-of-contents for goldmark documents.
- [goldmark-mermaid](https://github.com/abhinav/goldmark-mermaid): Adds support for rendering [Mermaid](https://mermaid-js.github.io/mermaid/) diagrams in goldmark documents.
- [goldmark-pikchr](https://github.com/jchenry/goldmark-pikchr): Adds support for rendering [Pikchr](https://pikchr.org/home/doc/trunk/homepage.md) diagrams in goldmark documents.
- [goldmark-embed](https://github.com/13rac1/goldmark-embed): Adds support for rendering embeds from YouTube links.
- [goldmark-latex](https://github.com/soypat/goldmark-latex): A $\LaTeX$ renderer that can be passed to `goldmark.WithRenderer()`.


goldmark internal(for extension developers)
----------------------------------------------
### Overview
goldmark's Markdown processing is outlined in the diagram below.

```
            <Markdown in []byte, parser.Context>
                           |
                           V
            +-------- parser.Parser ---------------------------
            | 1. Parse block elements into AST
            |   1. If a parsed block is a paragraph, apply 
            |      ast.ParagraphTransformer
            | 2. Traverse AST and parse blocks.
            |   1. Process delimiters(emphasis) at the end of
            |      block parsing
            | 3. Apply parser.ASTTransformers to AST
                           |
                           V
                      <ast.Node>
                           |
                           V
            +------- renderer.Renderer ------------------------
            | 1. Traverse AST and apply renderer.NodeRenderer
            |    corespond to the node type

                           |
                           V
                        <Output>
```

### Parsing
Markdown documents are read through `text.Reader` interface.

AST nodes do not have concrete text. AST nodes have segment information of the documents, represented by `text.Segment` .

`text.Segment` has 3 attributes: `Start`, `End`, `Padding` .

(TBC)

**TODO**

See `extension` directory for examples of extensions.

Summary:

1. Define AST Node as a struct in which `ast.BaseBlock` or `ast.BaseInline` is embedded.
2. Write a parser that implements `parser.BlockParser` or `parser.InlineParser`.
3. Write a renderer that implements `renderer.NodeRenderer`.
4. Define your goldmark extension that implements `goldmark.Extender`.


Donation
--------------------
BTC: 1NEDSyUmo4SMTDP83JJQSWi1MvQUGGNMZB

License
--------------------
MIT

Author
--------------------
Yusuke Inuzuka
