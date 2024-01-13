package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Gas lexer.
var Gas = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "GAS",
		Aliases:   []string{"gas", "asm"},
		Filenames: []string{"*.s", "*.S"},
		MimeTypes: []string{"text/x-gas"},
	},
	gasRules,
))

func gasRules() Rules {
	return Rules{
		"root": {
			Include("whitespace"),
			{`(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+):`, NameLabel, nil},
			{`\.(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameAttribute, Push("directive-args")},
			{`lock|rep(n?z)?|data\d+`, NameAttribute, nil},
			{`(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameFunction, Push("instruction-args")},
			{`[\r\n]+`, Text, nil},
		},
		"directive-args": {
			{`(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameConstant, nil},
			{`"(\\"|[^"])*"`, LiteralString, nil},
			{`@(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameAttribute, nil},
			{`(?:0[xX][a-zA-Z0-9]+|\d+)`, LiteralNumberInteger, nil},
			{`[\r\n]+`, Text, Pop(1)},
			Include("punctuation"),
			Include("whitespace"),
		},
		"instruction-args": {
			{`([a-z0-9]+)( )(<)((?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+))(>)`, ByGroups(LiteralNumberHex, Text, Punctuation, NameConstant, Punctuation), nil},
			{`([a-z0-9]+)( )(<)((?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+))([-+])((?:0[xX][a-zA-Z0-9]+|\d+))(>)`, ByGroups(LiteralNumberHex, Text, Punctuation, NameConstant, Punctuation, LiteralNumberInteger, Punctuation), nil},
			{`(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameConstant, nil},
			{`(?:0[xX][a-zA-Z0-9]+|\d+)`, LiteralNumberInteger, nil},
			{`%(?:[a-zA-Z$_][\w$.@-]*|\.[\w$.@-]+)`, NameVariable, nil},
			{`$(?:0[xX][a-zA-Z0-9]+|\d+)`, LiteralNumberInteger, nil},
			{`$'(.|\\')'`, LiteralStringChar, nil},
			{`[\r\n]+`, Text, Pop(1)},
			Include("punctuation"),
			Include("whitespace"),
		},
		"whitespace": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`[;#].*?\n`, Comment, nil},
		},
		"punctuation": {
			{`[-*,.()\[\]!:]+`, Punctuation, nil},
		},
	}
}
