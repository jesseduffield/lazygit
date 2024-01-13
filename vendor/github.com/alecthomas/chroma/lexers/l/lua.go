package l

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Lua lexer.
var Lua = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Lua",
		Aliases:   []string{"lua"},
		Filenames: []string{"*.lua", "*.wlua"},
		MimeTypes: []string{"text/x-lua", "application/x-lua"},
	},
	luaRules,
))

func luaRules() Rules {
	return Rules{
		"root": {
			{`#!.*`, CommentPreproc, nil},
			Default(Push("base")),
		},
		"ws": {
			{`(?:--\[(=*)\[[\w\W]*?\](\1)\])`, CommentMultiline, nil},
			{`(?:--.*$)`, CommentSingle, nil},
			{`(?:\s+)`, Text, nil},
		},
		"base": {
			Include("ws"),
			{`(?i)0x[\da-f]*(\.[\da-f]*)?(p[+-]?\d+)?`, LiteralNumberHex, nil},
			{`(?i)(\d*\.\d+|\d+\.\d*)(e[+-]?\d+)?`, LiteralNumberFloat, nil},
			{`(?i)\d+e[+-]?\d+`, LiteralNumberFloat, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`(?s)\[(=*)\[.*?\]\1\]`, LiteralString, nil},
			{`::`, Punctuation, Push("label")},
			{`\.{3}`, Punctuation, nil},
			{`[=<>|~&+\-*/%#^]+|\.\.`, Operator, nil},
			{`[\[\]{}().,:;]`, Punctuation, nil},
			{`(and|or|not)\b`, OperatorWord, nil},
			{`(break|do|else|elseif|end|for|if|in|repeat|return|then|until|while)\b`, KeywordReserved, nil},
			{`goto\b`, KeywordReserved, Push("goto")},
			{`(local)\b`, KeywordDeclaration, nil},
			{`(true|false|nil)\b`, KeywordConstant, nil},
			{`(function)\b`, KeywordReserved, Push("funcname")},
			{`[A-Za-z_]\w*(\.[A-Za-z_]\w*)?`, Name, nil},
			{`'`, LiteralStringSingle, Combined("stringescape", "sqs")},
			{`"`, LiteralStringDouble, Combined("stringescape", "dqs")},
		},
		"funcname": {
			Include("ws"),
			{`[.:]`, Punctuation, nil},
			{`(?:[^\W\d]\w*)(?=(?:(?:--\[(=*)\[[\w\W]*?\](\2)\])|(?:--.*$)|(?:\s+))*[.:])`, NameClass, nil},
			{`(?:[^\W\d]\w*)`, NameFunction, Pop(1)},
			{`\(`, Punctuation, Pop(1)},
		},
		"goto": {
			Include("ws"),
			{`(?:[^\W\d]\w*)`, NameLabel, Pop(1)},
		},
		"label": {
			Include("ws"),
			{`::`, Punctuation, Pop(1)},
			{`(?:[^\W\d]\w*)`, NameLabel, nil},
		},
		"stringescape": {
			{`\\([abfnrtv\\"\']|[\r\n]{1,2}|z\s*|x[0-9a-fA-F]{2}|\d{1,3}|u\{[0-9a-fA-F]+\})`, LiteralStringEscape, nil},
		},
		"sqs": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`[^\\']+`, LiteralStringSingle, nil},
		},
		"dqs": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`[^\\"]+`, LiteralStringDouble, nil},
		},
	}
}
