package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// PowerQuery lexer.
var PowerQuery = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "PowerQuery",
		Aliases:         []string{"powerquery", "pq"},
		Filenames:       []string{"*.pq"},
		MimeTypes:       []string{"text/x-powerquery"},
		DotAll:          true,
		CaseInsensitive: true,
	},
	powerqueryRules,
))

func powerqueryRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`(and|as|each|else|error|false|if|in|is|let|meta|not|null|or|otherwise|section|shared|then|true|try|type)\b`, Keyword, nil},
			{`(#binary|#date|#datetime|#datetimezone|#duration|#infinity|#nan|#sections|#shared|#table|#time)\b`, KeywordType, nil},
			{`(([a-zA-Z]|_)[\w|._]*|#"[^"]+")`, Name, nil},
			{`0[xX][0-9a-fA-F][0-9a-fA-F_]*[lL]?`, LiteralNumberHex, nil},
			{`([0-9]+\.[0-9]+|\.[0-9]+)([eE][0-9]+)?`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`[\(\)\[\]\{\}]`, Punctuation, nil},
			{`\.\.|\.\.\.|=>|<=|>=|<>|[@!?,;=<>\+\-\*\/&]`, Operator, nil},
		},
	}
}
