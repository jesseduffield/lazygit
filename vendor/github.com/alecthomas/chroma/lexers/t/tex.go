package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Tex lexer.
var TeX = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "TeX",
		Aliases:   []string{"tex", "latex"},
		Filenames: []string{"*.tex", "*.aux", "*.toc"},
		MimeTypes: []string{"text/x-tex", "text/x-latex"},
	},
	texRules,
))

func texRules() Rules {
	return Rules{
		"general": {
			{`%.*?\n`, Comment, nil},
			{`[{}]`, NameBuiltin, nil},
			{`[&_^]`, NameBuiltin, nil},
		},
		"root": {
			{`\\\[`, LiteralStringBacktick, Push("displaymath")},
			{`\\\(`, LiteralString, Push("inlinemath")},
			{`\$\$`, LiteralStringBacktick, Push("displaymath")},
			{`\$`, LiteralString, Push("inlinemath")},
			{`\\([a-zA-Z]+|.)`, Keyword, Push("command")},
			{`\\$`, Keyword, nil},
			Include("general"),
			{`[^\\$%&_^{}]+`, Text, nil},
		},
		"math": {
			{`\\([a-zA-Z]+|.)`, NameVariable, nil},
			Include("general"),
			{`[0-9]+`, LiteralNumber, nil},
			{`[-=!+*/()\[\]]`, Operator, nil},
			{`[^=!+*/()\[\]\\$%&_^{}0-9-]+`, NameBuiltin, nil},
		},
		"inlinemath": {
			{`\\\)`, LiteralString, Pop(1)},
			{`\$`, LiteralString, Pop(1)},
			Include("math"),
		},
		"displaymath": {
			{`\\\]`, LiteralString, Pop(1)},
			{`\$\$`, LiteralString, Pop(1)},
			{`\$`, NameBuiltin, nil},
			Include("math"),
		},
		"command": {
			{`\[.*?\]`, NameAttribute, nil},
			{`\*`, Keyword, nil},
			Default(Pop(1)),
		},
	}
}
