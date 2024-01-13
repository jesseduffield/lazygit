package b

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Blitzbasic lexer.
var Blitzbasic = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "BlitzBasic",
		Aliases:         []string{"blitzbasic", "b3d", "bplus"},
		Filenames:       []string{"*.bb", "*.decls"},
		MimeTypes:       []string{"text/x-bb"},
		CaseInsensitive: true,
	},
	blitzbasicRules,
))

func blitzbasicRules() Rules {
	return Rules{
		"root": {
			{`[ \t]+`, Text, nil},
			{`;.*?\n`, CommentSingle, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`[0-9]+\.[0-9]*(?!\.)`, LiteralNumberFloat, nil},
			{`\.[0-9]+(?!\.)`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`\$[0-9a-f]+`, LiteralNumberHex, nil},
			{`\%[10]+`, LiteralNumberBin, nil},
			{Words(`\b`, `\b`, `Shl`, `Shr`, `Sar`, `Mod`, `Or`, `And`, `Not`, `Abs`, `Sgn`, `Handle`, `Int`, `Float`, `Str`, `First`, `Last`, `Before`, `After`), Operator, nil},
			{`([+\-*/~=<>^])`, Operator, nil},
			{`[(),:\[\]\\]`, Punctuation, nil},
			{`\.([ \t]*)([a-z]\w*)`, NameLabel, nil},
			{`\b(New)\b([ \t]+)([a-z]\w*)`, ByGroups(KeywordReserved, Text, NameClass), nil},
			{`\b(Gosub|Goto)\b([ \t]+)([a-z]\w*)`, ByGroups(KeywordReserved, Text, NameLabel), nil},
			{`\b(Object)\b([ \t]*)([.])([ \t]*)([a-z]\w*)\b`, ByGroups(Operator, Text, Punctuation, Text, NameClass), nil},
			{`\b([a-z]\w*)(?:([ \t]*)(@{1,2}|[#$%])|([ \t]*)([.])([ \t]*)(?:([a-z]\w*)))?\b([ \t]*)(\()`, ByGroups(NameFunction, Text, KeywordType, Text, Punctuation, Text, NameClass, Text, Punctuation), nil},
			{`\b(Function)\b([ \t]+)([a-z]\w*)(?:([ \t]*)(@{1,2}|[#$%])|([ \t]*)([.])([ \t]*)(?:([a-z]\w*)))?`, ByGroups(KeywordReserved, Text, NameFunction, Text, KeywordType, Text, Punctuation, Text, NameClass), nil},
			{`\b(Type)([ \t]+)([a-z]\w*)`, ByGroups(KeywordReserved, Text, NameClass), nil},
			{`\b(Pi|True|False|Null)\b`, KeywordConstant, nil},
			{`\b(Local|Global|Const|Field|Dim)\b`, KeywordDeclaration, nil},
			{Words(`\b`, `\b`, `End`, `Return`, `Exit`, `Chr`, `Len`, `Asc`, `New`, `Delete`, `Insert`, `Include`, `Function`, `Type`, `If`, `Then`, `Else`, `ElseIf`, `EndIf`, `For`, `To`, `Next`, `Step`, `Each`, `While`, `Wend`, `Repeat`, `Until`, `Forever`, `Select`, `Case`, `Default`, `Goto`, `Gosub`, `Data`, `Read`, `Restore`), KeywordReserved, nil},
			{`([a-z]\w*)(?:([ \t]*)(@{1,2}|[#$%])|([ \t]*)([.])([ \t]*)(?:([a-z]\w*)))?`, ByGroups(NameVariable, Text, KeywordType, Text, Punctuation, Text, NameClass), nil},
		},
		"string": {
			{`""`, LiteralStringDouble, nil},
			{`"C?`, LiteralStringDouble, Pop(1)},
			{`[^"]+`, LiteralStringDouble, nil},
		},
	}
}
