package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Prolog lexer.
var Prolog = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Prolog",
		Aliases:   []string{"prolog"},
		Filenames: []string{"*.ecl", "*.prolog", "*.pro", "*.pl"},
		MimeTypes: []string{"text/x-prolog"},
	},
	prologRules,
))

func prologRules() Rules {
	return Rules{
		"root": {
			{`/\*`, CommentMultiline, Push("nested-comment")},
			{`%.*`, CommentSingle, nil},
			{`0\'.`, LiteralStringChar, nil},
			{`0b[01]+`, LiteralNumberBin, nil},
			{`0o[0-7]+`, LiteralNumberOct, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`\d\d?\'[a-zA-Z0-9]+`, LiteralNumberInteger, nil},
			{`(\d+\.\d*|\d*\.\d+)([eE][+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`[\[\](){}|.,;!]`, Punctuation, nil},
			{`:-|-->`, Punctuation, nil},
			{`"(?:\\x[0-9a-fA-F]+\\|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}|\\[0-7]+\\|\\["\nabcefnrstv]|[^\\"])*"`, LiteralStringDouble, nil},
			{`'(?:''|[^'])*'`, LiteralStringAtom, nil},
			{`is\b`, Operator, nil},
			{`(<|>|=<|>=|==|=:=|=|/|//|\*|\+|-)(?=\s|[a-zA-Z0-9\[])`, Operator, nil},
			{`(mod|div|not)\b`, Operator, nil},
			{`_`, Keyword, nil},
			{`([a-z]+)(:)`, ByGroups(NameNamespace, Punctuation), nil},
			{`([a-zÀ-῿぀-퟿-￯][\w$À-῿぀-퟿-￯]*)(\s*)(:-|-->)`, ByGroups(NameFunction, Text, Operator), nil},
			{`([a-zÀ-῿぀-퟿-￯][\w$À-῿぀-퟿-￯]*)(\s*)(\()`, ByGroups(NameFunction, Text, Punctuation), nil},
			{`[a-zÀ-῿぀-퟿-￯][\w$À-῿぀-퟿-￯]*`, LiteralStringAtom, nil},
			{`[#&*+\-./:<=>?@\\^~¡-¿‐-〿]+`, LiteralStringAtom, nil},
			{`[A-Z_]\w*`, NameVariable, nil},
			{`\s+|[ -‏￰-￾￯]`, Text, nil},
		},
		"nested-comment": {
			{`\*/`, CommentMultiline, Pop(1)},
			{`/\*`, CommentMultiline, Push()},
			{`[^*/]+`, CommentMultiline, nil},
			{`[*/]`, CommentMultiline, nil},
		},
	}
}
