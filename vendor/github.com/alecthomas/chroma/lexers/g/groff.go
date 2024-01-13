package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Groff lexer.
var Groff = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Groff",
		Aliases:   []string{"groff", "nroff", "man"},
		Filenames: []string{"*.[1-9]", "*.1p", "*.3pm", "*.man"},
		MimeTypes: []string{"application/x-troff", "text/troff"},
	},
	func() Rules {
		return Rules{
			"root": {
				{`(\.)(\w+)`, ByGroups(Text, Keyword), Push("request")},
				{`\.`, Punctuation, Push("request")},
				{`[^\\\n]+`, Text, Push("textline")},
				Default(Push("textline")),
			},
			"textline": {
				Include("escapes"),
				{`[^\\\n]+`, Text, nil},
				{`\n`, Text, Pop(1)},
			},
			"escapes": {
				{`\\"[^\n]*`, Comment, nil},
				{`\\[fn]\w`, LiteralStringEscape, nil},
				{`\\\(.{2}`, LiteralStringEscape, nil},
				{`\\.\[.*\]`, LiteralStringEscape, nil},
				{`\\.`, LiteralStringEscape, nil},
				{`\\\n`, Text, Push("request")},
			},
			"request": {
				{`\n`, Text, Pop(1)},
				Include("escapes"),
				{`"[^\n"]+"`, LiteralStringDouble, nil},
				{`\d+`, LiteralNumber, nil},
				{`\S+`, LiteralString, nil},
				{`\s+`, Text, nil},
			},
		}
	},
))
