package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// ANTLR lexer.
var ANTLR = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ANTLR",
		Aliases:   []string{"antlr"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	antlrRules,
))

func antlrRules() Rules {
	return Rules{
		"whitespace": {
			{`\s+`, TextWhitespace, nil},
		},
		"comments": {
			{`//.*$`, Comment, nil},
			{`/\*(.|\n)*?\*/`, Comment, nil},
		},
		"root": {
			Include("whitespace"),
			Include("comments"),
			{`(lexer|parser|tree)?(\s*)(grammar\b)(\s*)([A-Za-z]\w*)(;)`, ByGroups(Keyword, TextWhitespace, Keyword, TextWhitespace, NameClass, Punctuation), nil},
			{`options\b`, Keyword, Push("options")},
			{`tokens\b`, Keyword, Push("tokens")},
			{`(scope)(\s*)([A-Za-z]\w*)(\s*)(\{)`, ByGroups(Keyword, TextWhitespace, NameVariable, TextWhitespace, Punctuation), Push("action")},
			{`(catch|finally)\b`, Keyword, Push("exception")},
			{`(@[A-Za-z]\w*)(\s*)(::)?(\s*)([A-Za-z]\w*)(\s*)(\{)`, ByGroups(NameLabel, TextWhitespace, Punctuation, TextWhitespace, NameLabel, TextWhitespace, Punctuation), Push("action")},
			{`((?:protected|private|public|fragment)\b)?(\s*)([A-Za-z]\w*)(!)?`, ByGroups(Keyword, TextWhitespace, NameLabel, Punctuation), Push("rule-alts", "rule-prelims")},
		},
		"exception": {
			{`\n`, TextWhitespace, Pop(1)},
			{`\s`, TextWhitespace, nil},
			Include("comments"),
			{`\[`, Punctuation, Push("nested-arg-action")},
			{`\{`, Punctuation, Push("action")},
		},
		"rule-prelims": {
			Include("whitespace"),
			Include("comments"),
			{`returns\b`, Keyword, nil},
			{`\[`, Punctuation, Push("nested-arg-action")},
			{`\{`, Punctuation, Push("action")},
			{`(throws)(\s+)([A-Za-z]\w*)`, ByGroups(Keyword, TextWhitespace, NameLabel), nil},
			{`(,)(\s*)([A-Za-z]\w*)`, ByGroups(Punctuation, TextWhitespace, NameLabel), nil},
			{`options\b`, Keyword, Push("options")},
			{`(scope)(\s+)(\{)`, ByGroups(Keyword, TextWhitespace, Punctuation), Push("action")},
			{`(scope)(\s+)([A-Za-z]\w*)(\s*)(;)`, ByGroups(Keyword, TextWhitespace, NameLabel, TextWhitespace, Punctuation), nil},
			{`(@[A-Za-z]\w*)(\s*)(\{)`, ByGroups(NameLabel, TextWhitespace, Punctuation), Push("action")},
			{`:`, Punctuation, Pop(1)},
		},
		"rule-alts": {
			Include("whitespace"),
			Include("comments"),
			{`options\b`, Keyword, Push("options")},
			{`:`, Punctuation, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralString, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`<<([^>]|>[^>])>>`, LiteralString, nil},
			{`\$?[A-Z_]\w*`, NameConstant, nil},
			{`\$?[a-z_]\w*`, NameVariable, nil},
			{`(\+|\||->|=>|=|\(|\)|\.\.|\.|\?|\*|\^|!|\#|~)`, Operator, nil},
			{`,`, Punctuation, nil},
			{`\[`, Punctuation, Push("nested-arg-action")},
			{`\{`, Punctuation, Push("action")},
			{`;`, Punctuation, Pop(1)},
		},
		"tokens": {
			Include("whitespace"),
			Include("comments"),
			{`\{`, Punctuation, nil},
			{`([A-Z]\w*)(\s*)(=)?(\s*)(\'(?:\\\\|\\\'|[^\']*)\')?(\s*)(;)`, ByGroups(NameLabel, TextWhitespace, Punctuation, TextWhitespace, LiteralString, TextWhitespace, Punctuation), nil},
			{`\}`, Punctuation, Pop(1)},
		},
		"options": {
			Include("whitespace"),
			Include("comments"),
			{`\{`, Punctuation, nil},
			{`([A-Za-z]\w*)(\s*)(=)(\s*)([A-Za-z]\w*|\'(?:\\\\|\\\'|[^\']*)\'|[0-9]+|\*)(\s*)(;)`, ByGroups(NameVariable, TextWhitespace, Punctuation, TextWhitespace, Text, TextWhitespace, Punctuation), nil},
			{`\}`, Punctuation, Pop(1)},
		},
		"action": {
			{`([^${}\'"/\\]+|"(\\\\|\\"|[^"])*"|'(\\\\|\\'|[^'])*'|//.*$\n?|/\*(.|\n)*?\*/|/(?!\*)(\\\\|\\/|[^/])*/|\\(?!%)|/)+`, Other, nil},
			{`(\\)(%)`, ByGroups(Punctuation, Other), nil},
			{`(\$[a-zA-Z]+)(\.?)(text|value)?`, ByGroups(NameVariable, Punctuation, NameProperty), nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
		},
		"nested-arg-action": {
			{`([^$\[\]\'"/]+|"(\\\\|\\"|[^"])*"|'(\\\\|\\'|[^'])*'|//.*$\n?|/\*(.|\n)*?\*/|/(?!\*)(\\\\|\\/|[^/])*/|/)+`, Other, nil},
			{`\[`, Punctuation, Push()},
			{`\]`, Punctuation, Pop(1)},
			{`(\$[a-zA-Z]+)(\.?)(text|value)?`, ByGroups(NameVariable, Punctuation, NameProperty), nil},
			{`(\\\\|\\\]|\\\[|[^\[\]])+`, Other, nil},
		},
	}
}
