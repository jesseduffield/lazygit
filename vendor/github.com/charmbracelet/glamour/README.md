# Glamour

<p>
    <img src="https://stuff.charm.sh/glamour/glamour-github-header.png" width="245" alt="Glamour Title Treatment"><br>
    <a href="https://github.com/charmbracelet/glamour/releases"><img src="https://img.shields.io/github/release/charmbracelet/glamour.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/glamour?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/glamour/actions"><img src="https://github.com/charmbracelet/glamour/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://coveralls.io/github/charmbracelet/glamour?branch=master"><img src="https://coveralls.io/repos/github/charmbracelet/glamour/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/charmbracelet/glamour"><img src="https://goreportcard.com/badge/charmbracelet/glamour" alt="Go ReportCard"></a>
</p>

Stylesheet-based markdown rendering for your CLI apps.

![Glamour dark style example](https://stuff.charm.sh/glamour/glamour-example.png)

`glamour` lets you render [markdown](https://en.wikipedia.org/wiki/Markdown)
documents & templates on [ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code)
compatible terminals. You can create your own stylesheet or simply use one of
the stylish defaults.

## Usage

```go
import "github.com/charmbracelet/glamour"

in := `# Hello World

This is a simple example of Markdown rendering with Glamour!
Check out the [other examples](https://github.com/charmbracelet/glamour/tree/master/examples) too.

Bye!
`

out, err := glamour.Render(in, "dark")
fmt.Print(out)
```

<img src="https://github.com/charmbracelet/glamour/raw/master/examples/helloworld/helloworld.png" width="600" alt="Hello World example">

### Custom Renderer

```go
import "github.com/charmbracelet/glamour"

r, _ := glamour.NewTermRenderer(
    // detect background color and pick either the default dark or light theme
    glamour.WithAutoStyle(),
    // wrap output at specific width (default is 80)
    glamour.WithWordWrap(40),
)

out, err := r.Render(in)
fmt.Print(out)
```

## Styles

You can find all available default styles in our [gallery](https://github.com/charmbracelet/glamour/tree/master/styles/gallery).
Want to create your own style? [Learn how!](https://github.com/charmbracelet/glamour/tree/master/styles)

There are a few options for using a custom style:
1. Call `glamour.Render(inputText, "desiredStyle")`
1. Set the `GLAMOUR_STYLE` environment variable to your desired default style or a file location for a style and call `glamour.RenderWithEnvironmentConfig(inputText)`
1. Set the `GLAMOUR_STYLE` environment variable and pass `glamour.WithEnvironmentConfig()` to your custom renderer

## Glamourous Projects

Check out these projects, which use `glamour`:
- [Glow](https://github.com/charmbracelet/glow), a markdown renderer for
the command-line.
- [GitHub CLI](https://github.com/cli/cli), GitHub’s official command line tool.
- [GLab](https://github.com/profclems/glab), an open source GitLab command line tool.
- [Gitea CLI](https://gitea.com/gitea/tea), Gitea's official command line tool.
- [Meteor](https://github.com/odpf/meteor), an easy-to-use, plugin-driven metadata collection framework.

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.social/@charmcli)
* [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/glamour/raw/master/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
