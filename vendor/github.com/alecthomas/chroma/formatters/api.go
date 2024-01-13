package formatters

import (
	"io"
	"sort"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/formatters/svg"
)

var (
	// NoOp formatter.
	NoOp = Register("noop", chroma.FormatterFunc(func(w io.Writer, s *chroma.Style, iterator chroma.Iterator) error {
		for t := iterator(); t != chroma.EOF; t = iterator() {
			if _, err := io.WriteString(w, t.Value); err != nil {
				return err
			}
		}
		return nil
	}))
	// Default HTML formatter outputs self-contained HTML.
	htmlFull = Register("html", html.New(html.Standalone(true), html.WithClasses(true))) // nolint
	SVG      = Register("svg", svg.New(svg.EmbedFont("Liberation Mono", svg.FontLiberationMono, svg.WOFF)))
)

// Fallback formatter.
var Fallback = NoOp

// Registry of Formatters.
var Registry = map[string]chroma.Formatter{}

// Names of registered formatters.
func Names() []string {
	out := []string{}
	for name := range Registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Get formatter by name.
//
// If the given formatter is not found, the Fallback formatter will be returned.
func Get(name string) chroma.Formatter {
	if f, ok := Registry[name]; ok {
		return f
	}
	return Fallback
}

// Register a named formatter.
func Register(name string, formatter chroma.Formatter) chroma.Formatter {
	Registry[name] = formatter
	return formatter
}
