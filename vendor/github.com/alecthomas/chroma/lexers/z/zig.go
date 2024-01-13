package z

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Zig lexer.
var Zig = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Zig",
		Aliases:   []string{"zig"},
		Filenames: []string{"*.zig"},
		MimeTypes: []string{"text/zig"},
	},
	zigRules,
))

func zigRules() Rules {
	return Rules{
		"root": {
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`//.*?\n`, CommentSingle, nil},
			{Words(``, `\b`, `break`, `return`, `continue`, `asm`, `defer`, `errdefer`, `unreachable`, `try`, `catch`, `async`, `await`, `suspend`, `resume`, `cancel`), Keyword, nil},
			{Words(``, `\b`, `const`, `var`, `extern`, `packed`, `export`, `pub`, `noalias`, `inline`, `comptime`, `nakedcc`, `stdcallcc`, `volatile`, `allowzero`, `align`, `linksection`, `threadlocal`), KeywordReserved, nil},
			{Words(``, `\b`, `struct`, `enum`, `union`, `error`), Keyword, nil},
			{Words(``, `\b`, `while`, `for`), Keyword, nil},
			{Words(``, `\b`, `bool`, `f16`, `f32`, `f64`, `f128`, `void`, `noreturn`, `type`, `anyerror`, `promise`, `i0`, `u0`, `isize`, `usize`, `comptime_int`, `comptime_float`, `c_short`, `c_ushort`, `c_int`, `c_uint`, `c_long`, `c_ulong`, `c_longlong`, `c_ulonglong`, `c_longdouble`, `c_voidi8`, `u8`, `i16`, `u16`, `i32`, `u32`, `i64`, `u64`, `i128`, `u128`), KeywordType, nil},
			{Words(``, `\b`, `true`, `false`, `null`, `undefined`), KeywordConstant, nil},
			{Words(``, `\b`, `if`, `else`, `switch`, `and`, `or`, `orelse`), Keyword, nil},
			{Words(``, `\b`, `fn`, `usingnamespace`, `test`), Keyword, nil},
			{`0x[0-9a-fA-F]+\.[0-9a-fA-F]+([pP][\-+]?[0-9a-fA-F]+)?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+\.?[pP][\-+]?[0-9a-fA-F]+`, LiteralNumberFloat, nil},
			{`[0-9]+\.[0-9]+([eE][-+]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`[0-9]+\.?[eE][-+]?[0-9]+`, LiteralNumberFloat, nil},
			{`0b(?:_?[01])+`, LiteralNumberBin, nil},
			{`0o(?:_?[0-7])+`, LiteralNumberOct, nil},
			{`0x(?:_?[0-9a-fA-F])+`, LiteralNumberHex, nil},
			{`(?:_?[0-9])+`, LiteralNumberInteger, nil},
			{`@[a-zA-Z_]\w*`, NameBuiltin, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
			{`\'\\\'\'`, LiteralStringEscape, nil},
			{`\'\\(|x[a-fA-F0-9]{2}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{6}|[nr\\t\'"])\'`, LiteralStringEscape, nil},
			{`\'[^\\\']\'`, LiteralString, nil},
			{`\\\\[^\n]*`, LiteralStringHeredoc, nil},
			{`c\\\\[^\n]*`, LiteralStringHeredoc, nil},
			{`c?"`, LiteralString, Push("string")},
			{`[+%=><|^!?/\-*&~:]`, Operator, nil},
			{`[{}()\[\],.;]`, Punctuation, nil},
		},
		"string": {
			{`\\(x[a-fA-F0-9]{2}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{6}|[nr\\t\'"])`, LiteralStringEscape, nil},
			{`[^\\"\n]+`, LiteralString, nil},
			{`"`, LiteralString, Pop(1)},
		},
	}
}
