package n

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Nasm lexer.
var Nasm = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "NASM",
		Aliases:         []string{"nasm"},
		Filenames:       []string{"*.asm", "*.ASM"},
		MimeTypes:       []string{"text/x-nasm"},
		CaseInsensitive: true,
	},
	nasmRules,
))

func nasmRules() Rules {
	return Rules{
		"root": {
			{`^\s*%`, CommentPreproc, Push("preproc")},
			Include("whitespace"),
			{`[a-z$._?][\w$.?#@~]*:`, NameLabel, nil},
			{`([a-z$._?][\w$.?#@~]*)(\s+)(equ)`, ByGroups(NameConstant, KeywordDeclaration, KeywordDeclaration), Push("instruction-args")},
			{`BITS|USE16|USE32|SECTION|SEGMENT|ABSOLUTE|EXTERN|GLOBAL|ORG|ALIGN|STRUC|ENDSTRUC|COMMON|CPU|GROUP|UPPERCASE|IMPORT|EXPORT|LIBRARY|MODULE`, Keyword, Push("instruction-args")},
			{`(?:res|d)[bwdqt]|times`, KeywordDeclaration, Push("instruction-args")},
			{`[a-z$._?][\w$.?#@~]*`, NameFunction, Push("instruction-args")},
			{`[\r\n]+`, Text, nil},
		},
		"instruction-args": {
			{"\"(\\\\\"|[^\"\\n])*\"|'(\\\\'|[^'\\n])*'|`(\\\\`|[^`\\n])*`", LiteralString, nil},
			{`(?:0x[0-9a-f]+|$0[0-9a-f]*|[0-9]+[0-9a-f]*h)`, LiteralNumberHex, nil},
			{`[0-7]+q`, LiteralNumberOct, nil},
			{`[01]+b`, LiteralNumberBin, nil},
			{`[0-9]+\.e?[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			Include("punctuation"),
			{`r[0-9][0-5]?[bwd]|[a-d][lh]|[er]?[a-d]x|[er]?[sb]p|[er]?[sd]i|[c-gs]s|st[0-7]|mm[0-7]|cr[0-4]|dr[0-367]|tr[3-7]`, NameBuiltin, nil},
			{`[a-z$._?][\w$.?#@~]*`, NameVariable, nil},
			{`[\r\n]+`, Text, Pop(1)},
			Include("whitespace"),
		},
		"preproc": {
			{`[^;\n]+`, CommentPreproc, nil},
			{`;.*?\n`, CommentSingle, Pop(1)},
			{`\n`, CommentPreproc, Pop(1)},
		},
		"whitespace": {
			{`\n`, Text, nil},
			{`[ \t]+`, Text, nil},
			{`;.*`, CommentSingle, nil},
		},
		"punctuation": {
			{`[,():\[\]]+`, Punctuation, nil},
			{`[&|^<>+*/%~-]+`, Operator, nil},
			{`[$]+`, KeywordConstant, nil},
			{`seg|wrt|strict`, OperatorWord, nil},
			{`byte|[dq]?word`, KeywordType, nil},
		},
	}
}
