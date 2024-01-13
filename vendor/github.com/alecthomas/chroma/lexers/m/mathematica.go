package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Mathematica lexer.
var Mathematica = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Mathematica",
		Aliases:   []string{"mathematica", "mma", "nb"},
		Filenames: []string{"*.nb", "*.cdf", "*.nbp", "*.ma"},
		MimeTypes: []string{"application/mathematica", "application/vnd.wolfram.mathematica", "application/vnd.wolfram.mathematica.package", "application/vnd.wolfram.cdf"},
	},
	mathematicaRules,
))

func mathematicaRules() Rules {
	return Rules{
		"root": {
			{`(?s)\(\*.*?\*\)`, Comment, nil},
			{"([a-zA-Z]+[A-Za-z0-9]*`)", NameNamespace, nil},
			{`([A-Za-z0-9]*_+[A-Za-z0-9]*)`, NameVariable, nil},
			{`#\d*`, NameVariable, nil},
			{`([a-zA-Z]+[a-zA-Z0-9]*)`, Name, nil},
			{`-?\d+\.\d*`, LiteralNumberFloat, nil},
			{`-?\d*\.\d+`, LiteralNumberFloat, nil},
			{`-?\d+`, LiteralNumberInteger, nil},
			{Words(``, ``, `;;`, `=`, `=.`, `!===`, `:=`, `->`, `:>`, `/.`, `+`, `-`, `*`, `/`, `^`, `&&`, `||`, `!`, `<>`, `|`, `/;`, `?`, `@`, `//`, `/@`, `@@`, `@@@`, `~~`, `===`, `&`, `<`, `>`, `<=`, `>=`), Operator, nil},
			{Words(``, ``, `,`, `;`, `(`, `)`, `[`, `]`, `{`, `}`), Punctuation, nil},
			{`".*?"`, LiteralString, nil},
			{`\s+`, TextWhitespace, nil},
		},
	}
}
