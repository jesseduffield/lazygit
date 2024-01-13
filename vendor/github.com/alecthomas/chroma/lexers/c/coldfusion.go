package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Cfstatement lexer.
var Cfstatement = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "cfstatement",
		Aliases:         []string{"cfs"},
		Filenames:       []string{},
		MimeTypes:       []string{},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	cfstatementRules,
))

func cfstatementRules() Rules {
	return Rules{
		"root": {
			{`//.*?\n`, CommentSingle, nil},
			{`/\*(?:.|\n)*?\*/`, CommentMultiline, nil},
			{`\+\+|--`, Operator, nil},
			{`[-+*/^&=!]`, Operator, nil},
			{`<=|>=|<|>|==`, Operator, nil},
			{`mod\b`, Operator, nil},
			{`(eq|lt|gt|lte|gte|not|is|and|or)\b`, Operator, nil},
			{`\|\||&&`, Operator, nil},
			{`\?`, Operator, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`'.*?'`, LiteralStringSingle, nil},
			{`\d+`, LiteralNumber, nil},
			{`(if|else|len|var|xml|default|break|switch|component|property|function|do|try|catch|in|continue|for|return|while|required|any|array|binary|boolean|component|date|guid|numeric|query|string|struct|uuid|case)\b`, Keyword, nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(application|session|client|cookie|super|this|variables|arguments)\b`, NameConstant, nil},
			{`([a-z_$][\w.]*)(\s*)(\()`, ByGroups(NameFunction, Text, Punctuation), nil},
			{`[a-z_$][\w.]*`, NameVariable, nil},
			{`[()\[\]{};:,.\\]`, Punctuation, nil},
			{`\s+`, Text, nil},
		},
		"string": {
			{`""`, LiteralStringDouble, nil},
			{`#.+?#`, LiteralStringInterpol, nil},
			{`[^"#]+`, LiteralStringDouble, nil},
			{`#`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
	}
}
