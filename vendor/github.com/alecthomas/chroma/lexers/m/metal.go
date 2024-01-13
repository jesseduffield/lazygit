package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Metal lexer.
var Metal = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Metal",
		Aliases:   []string{"metal"},
		Filenames: []string{"*.metal"},
		MimeTypes: []string{"text/x-metal"},
		EnsureNL:  true,
	},
	metalRules,
))

func metalRules() Rules {
	return Rules{
		"statements": {
			{Words(``, `\b`, `namespace`, `operator`, `template`, `this`, `using`, `constexpr`), Keyword, nil},
			{`(enum)\b(\s+)(class)\b(\s*)`, ByGroups(Keyword, Text, Keyword, Text), Push("classname")},
			{`(class|struct|enum|union)\b(\s*)`, ByGroups(Keyword, Text), Push("classname")},
			{`\[\[.+\]\]`, NameAttribute, nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[LlUu]*`, LiteralNumberFloat, nil},
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`0[xX]([0-9A-Fa-f]('?[0-9A-Fa-f]+)*)[LlUu]*`, LiteralNumberHex, nil},
			{`0('?[0-7]+)+[LlUu]*`, LiteralNumberOct, nil},
			{`0[Bb][01]('?[01]+)*[LlUu]*`, LiteralNumberBin, nil},
			{`[0-9]('?[0-9]+)*[LlUu]*`, LiteralNumberInteger, nil},
			{`\*/`, Error, nil},
			{`[~!%^&*+=|?:<>/-]`, Operator, nil},
			{`[()\[\],.]`, Punctuation, nil},
			{Words(``, `\b`, `break`, `case`, `const`, `continue`, `do`, `else`, `enum`, `extern`, `for`, `if`, `return`, `sizeof`, `static`, `struct`, `switch`, `typedef`, `union`, `while`), Keyword, nil},
			{`(bool|float|half|long|ptrdiff_t|size_t|unsigned|u?char|u?int((8|16|32|64)_t)?|u?short)\b`, KeywordType, nil},
			{`(bool|float|half|u?(char|int|long|short))(2|3|4)\b`, KeywordType, nil},
			{`packed_(float|half|long|u?(char|int|short))(2|3|4)\b`, KeywordType, nil},
			{`(float|half)(2|3|4)x(2|3|4)\b`, KeywordType, nil},
			{`atomic_u?int\b`, KeywordType, nil},
			{`(rg?(8|16)(u|s)norm|rgba(8|16)(u|s)norm|srgba8unorm|rgb10a2|rg11b10f|rgb9e5)\b`, KeywordType, nil},
			{`(array|depth(2d|cube)(_array)?|depth2d_ms(_array)?|sampler|texture_buffer|texture(1|2)d(_array)?|texture2d_ms(_array)?|texture3d|texturecube(_array)?|uniform|visible_function_table)\b`, KeywordType, nil},
			{`(true|false|NULL)\b`, NameBuiltin, nil},
			{Words(``, `\b`, `device`, `constant`, `ray_data`, `thread`, `threadgroup`, `threadgroup_imageblock`), Keyword, nil},
			{`([a-zA-Z_]\w*)(\s*)(:)(?!:)`, ByGroups(NameLabel, Text, Punctuation), nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"root": {
			Include("whitespace"),
			{`(fragment|kernel|vertex)?((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;{]*)(\{)`, ByGroups(Keyword, UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), Push("function")},
			{`(fragment|kernel|vertex)?((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;]*)(;)`, ByGroups(Keyword, UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), nil},
			Default(Push("statement")),
		},
		"classname": {
			{`(\[\[.+\]\])(\s*)`, ByGroups(NameAttribute, Text), nil},
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
			{`\s*(?=[>{])`, Text, Pop(1)},
		},
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
		"statement": {
			Include("whitespace"),
			Include("statements"),
			{`[{]`, Punctuation, Push("root")},
			{`[;}]`, Punctuation, Pop(1)},
		},
		"function": {
			Include("whitespace"),
			Include("statements"),
			{`;`, Punctuation, nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
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
