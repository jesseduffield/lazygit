package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Angular2 lexer.
var Angular2 = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Angular2",
		Aliases:   []string{"ng2"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	angular2Rules,
))

func angular2Rules() Rules {
	return Rules{
		"root": {
			{`[^{([*#]+`, Other, nil},
			{`(\{\{)(\s*)`, ByGroups(CommentPreproc, Text), Push("ngExpression")},
			{`([([]+)([\w:.-]+)([\])]+)(\s*)(=)(\s*)`, ByGroups(Punctuation, NameAttribute, Punctuation, Text, Operator, Text), Push("attr")},
			{`([([]+)([\w:.-]+)([\])]+)(\s*)`, ByGroups(Punctuation, NameAttribute, Punctuation, Text), nil},
			{`([*#])([\w:.-]+)(\s*)(=)(\s*)`, ByGroups(Punctuation, NameAttribute, Punctuation, Operator), Push("attr")},
			{`([*#])([\w:.-]+)(\s*)`, ByGroups(Punctuation, NameAttribute, Punctuation), nil},
		},
		"ngExpression": {
			{`\s+(\|\s+)?`, Text, nil},
			{`\}\}`, CommentPreproc, Pop(1)},
			{`:?(true|false)`, LiteralStringBoolean, nil},
			{`:?"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`:?'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[0-9](\.[0-9]*)?(eE[+-][0-9])?[flFLdD]?|0[xX][0-9a-fA-F]+[Ll]?`, LiteralNumber, nil},
			{`[a-zA-Z][\w-]*(\(.*\))?`, NameVariable, nil},
			{`\.[\w-]+(\(.*\))?`, NameVariable, nil},
			{`(\?)(\s*)([^}\s]+)(\s*)(:)(\s*)([^}\s]+)(\s*)`, ByGroups(Operator, Text, LiteralString, Text, Operator, Text, LiteralString, Text), nil},
		},
		"attr": {
			{`".*?"`, LiteralString, Pop(1)},
			{`'.*?'`, LiteralString, Pop(1)},
			{`[^\s>]+`, LiteralString, Pop(1)},
		},
	}
}
