package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Tasm lexer.
var Tasm = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "TASM",
		Aliases:         []string{"tasm"},
		Filenames:       []string{"*.asm", "*.ASM", "*.tasm"},
		MimeTypes:       []string{"text/x-tasm"},
		CaseInsensitive: true,
	},
	tasmRules,
))

func tasmRules() Rules {
	return Rules{
		"root": {
			{`^\s*%`, CommentPreproc, Push("preproc")},
			Include("whitespace"),
			{`[@a-z$._?][\w$.?#@~]*:`, NameLabel, nil},
			{`BITS|USE16|USE32|SECTION|SEGMENT|ABSOLUTE|EXTERN|GLOBAL|ORG|ALIGN|STRUC|ENDSTRUC|ENDS|COMMON|CPU|GROUP|UPPERCASE|INCLUDE|EXPORT|LIBRARY|MODULE|PROC|ENDP|USES|ARG|DATASEG|UDATASEG|END|IDEAL|P386|MODEL|ASSUME|CODESEG|SIZE`, Keyword, Push("instruction-args")},
			{`([@a-z$._?][\w$.?#@~]*)(\s+)(db|dd|dw|T[A-Z][a-z]+)`, ByGroups(NameConstant, KeywordDeclaration, KeywordDeclaration), Push("instruction-args")},
			{`(?:res|d)[bwdqt]|times`, KeywordDeclaration, Push("instruction-args")},
			{`[@a-z$._?][\w$.?#@~]*`, NameFunction, Push("instruction-args")},
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
			{`[@a-z$._?][\w$.?#@~]*`, NameVariable, nil},
			{`(\\\s*)(;.*)([\r\n])`, ByGroups(Text, CommentSingle, Text), nil},
			{`[\r\n]+`, Text, Pop(1)},
			Include("whitespace"),
		},
		"preproc": {
			{`[^;\n]+`, CommentPreproc, nil},
			{`;.*?\n`, CommentSingle, Pop(1)},
			{`\n`, CommentPreproc, Pop(1)},
		},
		"whitespace": {
			{`[\n\r]`, Text, nil},
			{`\\[\n\r]`, Text, nil},
			{`[ \t]+`, Text, nil},
			{`;.*`, CommentSingle, nil},
		},
		"punctuation": {
			{`[,():\[\]]+`, Punctuation, nil},
			{`[&|^<>+*=/%~-]+`, Operator, nil},
			{`[$]+`, KeywordConstant, nil},
			{`seg|wrt|strict`, OperatorWord, nil},
			{`byte|[dq]?word`, KeywordType, nil},
		},
	}
}
