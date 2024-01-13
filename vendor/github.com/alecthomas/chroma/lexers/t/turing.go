package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Turing lexer.
var Turing = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Turing",
		Aliases:   []string{"turing"},
		Filenames: []string{"*.turing", "*.tu"},
		MimeTypes: []string{"text/x-turing"},
	},
	turingRules,
))

func turingRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`%(.*?)\n`, CommentSingle, nil},
			{`/(\\\n)?[*](.|\n)*?[*](\\\n)?/`, CommentMultiline, nil},
			{`(var|fcn|function|proc|procedure|process|class|end|record|type|begin|case|loop|for|const|union|monitor|module|handler)\b`, KeywordDeclaration, nil},
			{`(all|asm|assert|bind|bits|body|break|by|cheat|checked|close|condition|decreasing|def|deferred|else|elsif|exit|export|external|flexible|fork|forward|free|get|if|implement|import|include|inherit|init|invariant|label|new|objectclass|of|opaque|open|packed|pause|pervasive|post|pre|priority|put|quit|read|register|result|seek|self|set|signal|skip|tag|tell|then|timeout|to|unchecked|unqualified|wait|when|write)\b`, Keyword, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{Words(``, `\b`, `addressint`, `array`, `boolean`, `char`, `int`, `int1`, `int2`, `int4`, `int8`, `nat`, `nat1`, `nat2`, `nat4`, `nat8`, `pointer`, `real`, `real4`, `real8`, `string`, `enum`), KeywordType, nil},
			{`\d+i`, LiteralNumber, nil},
			{`\d+\.\d*([Ee][-+]\d+)?i`, LiteralNumber, nil},
			{`\.\d+([Ee][-+]\d+)?i`, LiteralNumber, nil},
			{`\d+[Ee][-+]\d+i`, LiteralNumber, nil},
			{`\d+(\.\d+[eE][+\-]?\d+|\.\d*|[eE][+\-]?\d+)`, LiteralNumberFloat, nil},
			{`\.\d+([eE][+\-]?\d+)?`, LiteralNumberFloat, nil},
			{`0[0-7]+`, LiteralNumberOct, nil},
			{`0[xX][0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`(0|[1-9][0-9]*)`, LiteralNumberInteger, nil},
			{`(div|mod|rem|\*\*|=|<|>|>=|<=|not=|not|and|or|xor|=>|in|shl|shr|->|~|~=|~in|&|:=|\.\.|[\^+\-*/&#])`, Operator, nil},
			{`'(\\['"\\abfnrtv]|\\x[0-9a-fA-F]{2}|\\[0-7]{1,3}|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}|[^\\])'`, LiteralStringChar, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`[()\[\]{}.,:]`, Punctuation, nil},
			{`[^\W\d]\w*`, NameOther, nil},
		},
	}
}
