package n

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Newspeak lexer.
var Newspeak = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Newspeak",
		Aliases:   []string{"newspeak"},
		Filenames: []string{"*.ns2"},
		MimeTypes: []string{"text/x-newspeak"},
	},
	newspeakRules,
))

func newspeakRules() Rules {
	return Rules{
		"root": {
			{`\b(Newsqueak2)\b`, KeywordDeclaration, nil},
			{`'[^']*'`, LiteralString, nil},
			{`\b(class)(\s+)(\w+)(\s*)`, ByGroups(KeywordDeclaration, Text, NameClass, Text), nil},
			{`\b(mixin|self|super|private|public|protected|nil|true|false)\b`, Keyword, nil},
			{`(\w+\:)(\s*)([a-zA-Z_]\w+)`, ByGroups(NameFunction, Text, NameVariable), nil},
			{`(\w+)(\s*)(=)`, ByGroups(NameAttribute, Text, Operator), nil},
			{`<\w+>`, CommentSpecial, nil},
			Include("expressionstat"),
			Include("whitespace"),
		},
		"expressionstat": {
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`:\w+`, NameVariable, nil},
			{`(\w+)(::)`, ByGroups(NameVariable, Operator), nil},
			{`\w+:`, NameFunction, nil},
			{`\w+`, NameVariable, nil},
			{`\(|\)`, Punctuation, nil},
			{`\[|\]`, Punctuation, nil},
			{`\{|\}`, Punctuation, nil},
			{`(\^|\+|\/|~|\*|<|>|=|@|%|\||&|\?|!|,|-|:)`, Operator, nil},
			{`\.|;`, Punctuation, nil},
			Include("whitespace"),
			Include("literals"),
		},
		"literals": {
			{`\$.`, LiteralString, nil},
			{`'[^']*'`, LiteralString, nil},
			{`#'[^']*'`, LiteralStringSymbol, nil},
			{`#\w+:?`, LiteralStringSymbol, nil},
			{`#(\+|\/|~|\*|<|>|=|@|%|\||&|\?|!|,|-)+`, LiteralStringSymbol, nil},
		},
		"whitespace": {
			{`\s+`, Text, nil},
			{`"[^"]*"`, Comment, nil},
		},
	}
}
