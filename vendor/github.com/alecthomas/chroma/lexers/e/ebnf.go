package e

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Ebnf lexer.
var Ebnf = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "EBNF",
		Aliases:   []string{"ebnf"},
		Filenames: []string{"*.ebnf"},
		MimeTypes: []string{"text/x-ebnf"},
	},
	ebnfRules,
))

func ebnfRules() Rules {
	return Rules{
		"root": {
			Include("whitespace"),
			Include("comment_start"),
			Include("identifier"),
			{`=`, Operator, Push("production")},
		},
		"production": {
			Include("whitespace"),
			Include("comment_start"),
			Include("identifier"),
			{`"[^"]*"`, LiteralStringDouble, nil},
			{`'[^']*'`, LiteralStringSingle, nil},
			{`(\?[^?]*\?)`, NameEntity, nil},
			{`[\[\]{}(),|]`, Punctuation, nil},
			{`-`, Operator, nil},
			{`;`, Punctuation, Pop(1)},
			{`\.`, Punctuation, Pop(1)},
		},
		"whitespace": {
			{`\s+`, Text, nil},
		},
		"comment_start": {
			{`\(\*`, CommentMultiline, Push("comment")},
		},
		"comment": {
			{`[^*)]`, CommentMultiline, nil},
			Include("comment_start"),
			{`\*\)`, CommentMultiline, Pop(1)},
			{`[*)]`, CommentMultiline, nil},
		},
		"identifier": {
			{`([a-zA-Z][\w \-]*)`, Keyword, nil},
		},
	}
}
