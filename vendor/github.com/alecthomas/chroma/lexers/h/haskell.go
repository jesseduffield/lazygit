package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Haskell lexer.
var Haskell = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Haskell",
		Aliases:   []string{"haskell", "hs"},
		Filenames: []string{"*.hs"},
		MimeTypes: []string{"text/x-haskell"},
	},
	haskellRules,
))

func haskellRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`--(?![!#$%&*+./<=>?@^|_~:\\]).*?$`, CommentSingle, nil},
			{`\{-`, CommentMultiline, Push("comment")},
			{`\bimport\b`, KeywordReserved, Push("import")},
			{`\bmodule\b`, KeywordReserved, Push("module")},
			{`\berror\b`, NameException, nil},
			{`\b(case|class|data|default|deriving|do|else|family|if|in|infix[lr]?|instance|let|newtype|of|then|type|where|_)(?!\')\b`, KeywordReserved, nil},
			{`'[^\\]'`, LiteralStringChar, nil},
			{`^[_\p{Ll}][\w\']*`, NameFunction, nil},
			{`'?[_\p{Ll}][\w']*`, Name, nil},
			{`('')?[\p{Lu}][\w\']*`, KeywordType, nil},
			{`(')[\p{Lu}][\w\']*`, KeywordType, nil},
			{`(')\[[^\]]*\]`, KeywordType, nil},
			{`(')\([^)]*\)`, KeywordType, nil},
			{`\\(?![:!#$%&*+.\\/<=>?@^|~-]+)`, NameFunction, nil},
			{`(<-|::|->|=>|=)(?![:!#$%&*+.\\/<=>?@^|~-]+)`, OperatorWord, nil},
			{`:[:!#$%&*+.\\/<=>?@^|~-]*`, KeywordType, nil},
			{`[:!#$%&*+.\\/<=>?@^|~-]+`, Operator, nil},
			{`\d+[eE][+-]?\d+`, LiteralNumberFloat, nil},
			{`\d+\.\d+([eE][+-]?\d+)?`, LiteralNumberFloat, nil},
			{`0[oO][0-7]+`, LiteralNumberOct, nil},
			{`0[xX][\da-fA-F]+`, LiteralNumberHex, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`'`, LiteralStringChar, Push("character")},
			{`"`, LiteralString, Push("string")},
			{`\[\]`, KeywordType, nil},
			{`\(\)`, NameBuiltin, nil},
			{"[][(),;`{}]", Punctuation, nil},
		},
		"import": {
			{`\s+`, Text, nil},
			{`"`, LiteralString, Push("string")},
			{`\)`, Punctuation, Pop(1)},
			{`qualified\b`, Keyword, nil},
			{`([\p{Lu}][\w.]*)(\s+)(as)(\s+)([\p{Lu}][\w.]*)`, ByGroups(NameNamespace, Text, Keyword, Text, Name), Pop(1)},
			{`([\p{Lu}][\w.]*)(\s+)(hiding)(\s+)(\()`, ByGroups(NameNamespace, Text, Keyword, Text, Punctuation), Push("funclist")},
			{`([\p{Lu}][\w.]*)(\s+)(\()`, ByGroups(NameNamespace, Text, Punctuation), Push("funclist")},
			{`[\w.]+`, NameNamespace, Pop(1)},
		},
		"module": {
			{`\s+`, Text, nil},
			{`([\p{Lu}][\w.]*)(\s+)(\()`, ByGroups(NameNamespace, Text, Punctuation), Push("funclist")},
			{`[\p{Lu}][\w.]*`, NameNamespace, Pop(1)},
		},
		"funclist": {
			{`\s+`, Text, nil},
			{`[\p{Lu}]\w*`, KeywordType, nil},
			{`(_[\w\']+|[\p{Ll}][\w\']*)`, NameFunction, nil},
			{`--(?![!#$%&*+./<=>?@^|_~:\\]).*?$`, CommentSingle, nil},
			{`\{-`, CommentMultiline, Push("comment")},
			{`,`, Punctuation, nil},
			{`[:!#$%&*+.\\/<=>?@^|~-]+`, Operator, nil},
			{`\(`, Punctuation, Push("funclist", "funclist")},
			{`\)`, Punctuation, Pop(2)},
		},
		"comment": {
			{`[^-{}]+`, CommentMultiline, nil},
			{`\{-`, CommentMultiline, Push()},
			{`-\}`, CommentMultiline, Pop(1)},
			{`[-{}]`, CommentMultiline, nil},
		},
		"character": {
			{`[^\\']'`, LiteralStringChar, Pop(1)},
			{`\\`, LiteralStringEscape, Push("escape")},
			{`'`, LiteralStringChar, Pop(1)},
		},
		"string": {
			{`[^\\"]+`, LiteralString, nil},
			{`\\`, LiteralStringEscape, Push("escape")},
			{`"`, LiteralString, Pop(1)},
		},
		"escape": {
			{`[abfnrtv"\'&\\]`, LiteralStringEscape, Pop(1)},
			{`\^[][\p{Lu}@^_]`, LiteralStringEscape, Pop(1)},
			{`NUL|SOH|[SE]TX|EOT|ENQ|ACK|BEL|BS|HT|LF|VT|FF|CR|S[OI]|DLE|DC[1-4]|NAK|SYN|ETB|CAN|EM|SUB|ESC|[FGRU]S|SP|DEL`, LiteralStringEscape, Pop(1)},
			{`o[0-7]+`, LiteralStringEscape, Pop(1)},
			{`x[\da-fA-F]+`, LiteralStringEscape, Pop(1)},
			{`\d+`, LiteralStringEscape, Pop(1)},
			{`\s+\\`, LiteralStringEscape, Pop(1)},
		},
	}
}
