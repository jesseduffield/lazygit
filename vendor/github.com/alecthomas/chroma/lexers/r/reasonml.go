package r

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Reasonml lexer.
var Reasonml = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ReasonML",
		Aliases:   []string{"reason", "reasonml"},
		Filenames: []string{"*.re", "*.rei"},
		MimeTypes: []string{"text/x-reasonml"},
	},
	reasonmlRules,
))

func reasonmlRules() Rules {
	return Rules{
		"escape-sequence": {
			{`\\[\\"\'ntbr]`, LiteralStringEscape, nil},
			{`\\[0-9]{3}`, LiteralStringEscape, nil},
			{`\\x[0-9a-fA-F]{2}`, LiteralStringEscape, nil},
		},
		"root": {
			{`\s+`, Text, nil},
			{`false|true|\(\)|\[\]`, NameBuiltinPseudo, nil},
			{`\b([A-Z][\w\']*)(?=\s*\.)`, NameNamespace, Push("dotted")},
			{`\b([A-Z][\w\']*)`, NameClass, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`\/\*(?![\/])`, CommentMultiline, Push("comment")},
			{`\b(as|assert|begin|class|constraint|do|done|downto|else|end|exception|external|false|for|fun|esfun|function|functor|if|in|include|inherit|initializer|lazy|let|switch|module|pub|mutable|new|nonrec|object|of|open|pri|rec|sig|struct|then|to|true|try|type|val|virtual|when|while|with)\b`, Keyword, nil},
			{"(~|\\}|\\|]|\\||\\|\\||\\{<|\\{|`|_|]|\\[\\||\\[>|\\[<|\\[|\\?\\?|\\?|>\\}|>]|>|=|<-|<|;;|;|:>|:=|::|:|\\.\\.\\.|\\.\\.|\\.|=>|-\\.|-|,|\\+|\\*|\\)|\\(|&&|&|#|!=)", OperatorWord, nil},
			{`([=<>@^|&+\*/$%-]|[!?~])?[!$%&*+\./:<=>?@^|~-]`, Operator, nil},
			{`\b(and|asr|land|lor|lsl|lsr|lxor|mod|or)\b`, OperatorWord, nil},
			{`\b(unit|int|float|bool|string|char|list|array)\b`, KeywordType, nil},
			{`[^\W\d][\w']*`, Name, nil},
			{`-?\d[\d_]*(.[\d_]*)?([eE][+\-]?\d[\d_]*)`, LiteralNumberFloat, nil},
			{`0[xX][\da-fA-F][\da-fA-F_]*`, LiteralNumberHex, nil},
			{`0[oO][0-7][0-7_]*`, LiteralNumberOct, nil},
			{`0[bB][01][01_]*`, LiteralNumberBin, nil},
			{`\d[\d_]*`, LiteralNumberInteger, nil},
			{`'(?:(\\[\\\"'ntbr ])|(\\[0-9]{3})|(\\x[0-9a-fA-F]{2}))'`, LiteralStringChar, nil},
			{`'.'`, LiteralStringChar, nil},
			{`'`, Keyword, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`[~?][a-z][\w\']*:`, NameVariable, nil},
		},
		"comment": {
			{`[^\/*]+`, CommentMultiline, nil},
			{`\/\*`, CommentMultiline, Push()},
			{`\*\/`, CommentMultiline, Pop(1)},
			{`[\*]`, CommentMultiline, nil},
		},
		"string": {
			{`[^\\"]+`, LiteralStringDouble, nil},
			Include("escape-sequence"),
			{`\\\n`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"dotted": {
			{`\s+`, Text, nil},
			{`\.`, Punctuation, nil},
			{`[A-Z][\w\']*(?=\s*\.)`, NameNamespace, nil},
			{`[A-Z][\w\']*`, NameClass, Pop(1)},
			{`[a-z_][\w\']*`, Name, Pop(1)},
			Default(Pop(1)),
		},
	}
}
