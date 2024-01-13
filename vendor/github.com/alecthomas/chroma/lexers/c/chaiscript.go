package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Chaiscript lexer.
var Chaiscript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ChaiScript",
		Aliases:   []string{"chai", "chaiscript"},
		Filenames: []string{"*.chai"},
		MimeTypes: []string{"text/x-chaiscript", "application/x-chaiscript"},
		DotAll:    true,
	},
	chaiscriptRules,
))

func chaiscriptRules() Rules {
	return Rules{
		"commentsandwhitespace": {
			{`\s+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`^\#.*?\n`, CommentSingle, nil},
		},
		"slashstartsregex": {
			Include("commentsandwhitespace"),
			{`/(\\.|[^[/\\\n]|\[(\\.|[^\]\\\n])*])+/([gim]+\b|\B)`, LiteralStringRegex, Pop(1)},
			{`(?=/)`, Text, Push("#pop", "badregex")},
			Default(Pop(1)),
		},
		"badregex": {
			{`\n`, Text, Pop(1)},
		},
		"root": {
			Include("commentsandwhitespace"),
			{`\n`, Text, nil},
			{`[^\S\n]+`, Text, nil},
			{`\+\+|--|~|&&|\?|:|\|\||\\(?=\n)|\.\.(<<|>>>?|==?|!=?|[-<>+*%&|^/])=?`, Operator, Push("slashstartsregex")},
			{`[{(\[;,]`, Punctuation, Push("slashstartsregex")},
			{`[})\].]`, Punctuation, nil},
			{`[=+\-*/]`, Operator, nil},
			{`(for|in|while|do|break|return|continue|if|else|throw|try|catch)\b`, Keyword, Push("slashstartsregex")},
			{`(var)\b`, KeywordDeclaration, Push("slashstartsregex")},
			{`(attr|def|fun)\b`, KeywordReserved, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{`(eval|throw)\b`, NameBuiltin, nil},
			{"`\\S+`", NameBuiltin, nil},
			{`[$a-zA-Z_]\w*`, NameOther, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"`, LiteralStringDouble, Push("dqstring")},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
		},
		"dqstring": {
			{`\$\{[^"}]+?\}`, LiteralStringInterpol, nil},
			{`\$`, LiteralStringDouble, nil},
			{`\\\\`, LiteralStringDouble, nil},
			{`\\"`, LiteralStringDouble, nil},
			{`[^\\"$]+`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
	}
}
