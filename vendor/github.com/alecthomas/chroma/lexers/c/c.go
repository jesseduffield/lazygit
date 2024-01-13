package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// C lexer.
var C = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "C",
		Aliases:   []string{"c"},
		Filenames: []string{"*.c", "*.h", "*.idc", "*.x[bp]m"},
		MimeTypes: []string{"text/x-chdr", "text/x-csrc", "image/x-xbitmap", "image/x-xpixmap"},
		EnsureNL:  true,
	},
	cRules,
))

func cRules() Rules {
	return Rules{
		"whitespace": {
			{`^#if\s+0`, CommentPreproc, Push("if0")},
			{`^#`, CommentPreproc, Push("macro")},
			{`^(\s*(?:/[*].*?[*]/\s*)?)(#if\s+0)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("if0")},
			{`^(\s*(?:/[*].*?[*]/\s*)?)(#)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("macro")},
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`//(\n|[\w\W]*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil},
			{`/(\\\n)?[*][\w\W]*`, CommentMultiline, nil},
		},
		"statements": {
			{`(L?)(")`, ByGroups(LiteralStringAffix, LiteralString), Push("string")},
			{`(L?)(')(\\.|\\[0-7]{1,3}|\\x[a-fA-F0-9]{1,2}|[^\\\'\n])(')`, ByGroups(LiteralStringAffix, LiteralStringChar, LiteralStringChar, LiteralStringChar), nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[LlUu]*`, LiteralNumberFloat, nil},
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+[LlUu]*`, LiteralNumberHex, nil},
			{`0[0-7]+[LlUu]*`, LiteralNumberOct, nil},
			{`\d+[LlUu]*`, LiteralNumberInteger, nil},
			{`\*/`, Error, nil},
			{`[~!%^&*+=|?:<>/-]`, Operator, nil},
			{`[()\[\],.]`, Punctuation, nil},
			{Words(``, `\b`, `asm`, `auto`, `break`, `case`, `const`, `continue`, `default`, `do`, `else`, `enum`, `extern`, `for`, `goto`, `if`, `register`, `restricted`, `return`, `sizeof`, `static`, `struct`, `switch`, `typedef`, `union`, `volatile`, `while`), Keyword, nil},
			{`(bool|int|long|float|short|double|char((8|16|32)_t)?|unsigned|signed|void|u?int(_fast|_least|)(8|16|32|64)_t)\b`, KeywordType, nil},
			{Words(``, `\b`, `inline`, `_inline`, `__inline`, `naked`, `restrict`, `thread`, `typename`), KeywordReserved, nil},
			{`(__m(128i|128d|128|64))\b`, KeywordReserved, nil},
			{Words(`__`, `\b`, `asm`, `int8`, `based`, `except`, `int16`, `stdcall`, `cdecl`, `fastcall`, `int32`, `declspec`, `finally`, `int64`, `try`, `leave`, `wchar_t`, `w64`, `unaligned`, `raise`, `noop`, `identifier`, `forceinline`, `assume`), KeywordReserved, nil},
			{`(true|false|NULL)\b`, NameBuiltin, nil},
			{`([a-zA-Z_]\w*)(\s*)(:)(?!:)`, ByGroups(NameLabel, Text, Punctuation), nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"root": {
			Include("whitespace"),
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;{]*)(\{)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), Push("function")},
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;]*)(;)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), nil},
			Default(Push("statement")),
		},
		"statement": {
			Include("whitespace"),
			Include("statements"),
			{`[{}]`, Punctuation, nil},
			{`;`, Punctuation, Pop(1)},
		},
		"function": {
			Include("whitespace"),
			Include("statements"),
			{`;`, Punctuation, nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\([\\abfnrtv"\']|x[a-fA-F0-9]{2,4}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|[0-7]{1,3})`, LiteralStringEscape, nil},
			{`[^\\"\n]+`, LiteralString, nil},
			{`\\\n`, LiteralString, nil},
			{`\\`, LiteralString, nil},
		},
		"macro": {
			{`(include)(\s*(?:/[*].*?[*]/\s*)?)([^\n]+)`, ByGroups(CommentPreproc, Text, CommentPreprocFile), nil},
			{`[^/\n]+`, CommentPreproc, nil},
			{`/[*](.|\n)*?[*]/`, CommentMultiline, nil},
			{`//.*?\n`, CommentSingle, Pop(1)},
			{`/`, CommentPreproc, nil},
			{`(?<=\\)\n`, CommentPreproc, nil},
			{`\n`, CommentPreproc, Pop(1)},
		},
		"if0": {
			{`^\s*#if.*?(?<!\\)\n`, CommentPreproc, Push()},
			{`^\s*#el(?:se|if).*\n`, CommentPreproc, Pop(1)},
			{`^\s*#endif.*?(?<!\\)\n`, CommentPreproc, Pop(1)},
			{`.*?\n`, Comment, nil},
		},
	}
}
