package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Handlebars lexer.
var Handlebars = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Handlebars",
		Aliases:   []string{"handlebars", "hbs"},
		Filenames: []string{"*.handlebars", "*.hbs"},
		MimeTypes: []string{},
	},
	handlebarsRules,
))

func handlebarsRules() Rules {
	return Rules{
		"root": {
			{`[^{]+`, Other, nil},
			{`\{\{!.*\}\}`, Comment, nil},
			{`(\{\{\{)(\s*)`, ByGroups(CommentSpecial, Text), Push("tag")},
			{`(\{\{)(\s*)`, ByGroups(CommentPreproc, Text), Push("tag")},
		},
		"tag": {
			{`\s+`, Text, nil},
			{`\}\}\}`, CommentSpecial, Pop(1)},
			{`\}\}`, CommentPreproc, Pop(1)},
			{`([#/]*)(each|if|unless|else|with|log|in(?:line)?)`, ByGroups(Keyword, Keyword), nil},
			{`#\*inline`, Keyword, nil},
			{`([#/])([\w-]+)`, ByGroups(NameFunction, NameFunction), nil},
			{`([\w-]+)(=)`, ByGroups(NameAttribute, Operator), nil},
			{`(>)(\s*)(@partial-block)`, ByGroups(Keyword, Text, Keyword), nil},
			{`(#?>)(\s*)([\w-]+)`, ByGroups(Keyword, Text, NameVariable), nil},
			{`(>)(\s*)(\()`, ByGroups(Keyword, Text, Punctuation), Push("dynamic-partial")},
			Include("generic"),
		},
		"dynamic-partial": {
			{`\s+`, Text, nil},
			{`\)`, Punctuation, Pop(1)},
			{`(lookup)(\s+)(\.|this)(\s+)`, ByGroups(Keyword, Text, NameVariable, Text), nil},
			{`(lookup)(\s+)(\S+)`, ByGroups(Keyword, Text, UsingSelf("variable")), nil},
			{`[\w-]+`, NameFunction, nil},
			Include("generic"),
		},
		"variable": {
			{`[a-zA-Z][\w-]*`, NameVariable, nil},
			{`\.[\w-]+`, NameVariable, nil},
			{`(this\/|\.\/|(\.\.\/)+)[\w-]+`, NameVariable, nil},
		},
		"generic": {
			Include("variable"),
			{`:?"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`:?'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[0-9](\.[0-9]*)?(eE[+-][0-9])?[flFLdD]?|0[xX][0-9a-fA-F]+[Ll]?`, LiteralNumber, nil},
		},
	}
}
