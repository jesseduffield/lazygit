package o

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Objective-C lexer.
var ObjectiveC = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Objective-C",
		Aliases:   []string{"objective-c", "objectivec", "obj-c", "objc"},
		Filenames: []string{"*.m", "*.h"},
		MimeTypes: []string{"text/x-objective-c"},
	},
	objectiveCRules,
))

func objectiveCRules() Rules {
	return Rules{
		"statements": {
			{`@"`, LiteralString, Push("string")},
			{`@(YES|NO)`, LiteralNumber, nil},
			{`@'(\\.|\\[0-7]{1,3}|\\x[a-fA-F0-9]{1,2}|[^\\\'\n])'`, LiteralStringChar, nil},
			{`@(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[lL]?`, LiteralNumberFloat, nil},
			{`@(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`@0x[0-9a-fA-F]+[Ll]?`, LiteralNumberHex, nil},
			{`@0[0-7]+[Ll]?`, LiteralNumberOct, nil},
			{`@\d+[Ll]?`, LiteralNumberInteger, nil},
			{`@\(`, Literal, Push("literal_number")},
			{`@\[`, Literal, Push("literal_array")},
			{`@\{`, Literal, Push("literal_dictionary")},
			{Words(``, `\b`, `@selector`, `@private`, `@protected`, `@public`, `@encode`, `@synchronized`, `@try`, `@throw`, `@catch`, `@finally`, `@end`, `@property`, `@synthesize`, `__bridge`, `__bridge_transfer`, `__autoreleasing`, `__block`, `__weak`, `__strong`, `weak`, `strong`, `copy`, `retain`, `assign`, `unsafe_unretained`, `atomic`, `nonatomic`, `readonly`, `readwrite`, `setter`, `getter`, `typeof`, `in`, `out`, `inout`, `release`, `class`, `@dynamic`, `@optional`, `@required`, `@autoreleasepool`), Keyword, nil},
			{Words(``, `\b`, `id`, `instancetype`, `Class`, `IMP`, `SEL`, `BOOL`, `IBOutlet`, `IBAction`, `unichar`), KeywordType, nil},
			{`@(true|false|YES|NO)\n`, NameBuiltin, nil},
			{`(YES|NO|nil|self|super)\b`, NameBuiltin, nil},
			{`(Boolean|UInt8|SInt8|UInt16|SInt16|UInt32|SInt32)\b`, KeywordType, nil},
			{`(TRUE|FALSE)\b`, NameBuiltin, nil},
			{`(@interface|@implementation)(\s+)`, ByGroups(Keyword, Text), Push("#pop", "oc_classname")},
			{`(@class|@protocol)(\s+)`, ByGroups(Keyword, Text), Push("#pop", "oc_forward_classname")},
			{`@`, Punctuation, nil},
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
			{`(bool|int|long|float|short|double|char|unsigned|signed|void)\b`, KeywordType, nil},
			{Words(``, `\b`, `inline`, `_inline`, `__inline`, `naked`, `restrict`, `thread`, `typename`), KeywordReserved, nil},
			{`(__m(128i|128d|128|64))\b`, KeywordReserved, nil},
			{Words(`__`, `\b`, `asm`, `int8`, `based`, `except`, `int16`, `stdcall`, `cdecl`, `fastcall`, `int32`, `declspec`, `finally`, `int64`, `try`, `leave`, `wchar_t`, `w64`, `unaligned`, `raise`, `noop`, `identifier`, `forceinline`, `assume`), KeywordReserved, nil},
			{`(true|false|NULL)\b`, NameBuiltin, nil},
			{`([a-zA-Z_]\w*)(\s*)(:)(?!:)`, ByGroups(NameLabel, Text, Punctuation), nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"oc_classname": {
			{`([a-zA-Z$_][\w$]*)(\s*:\s*)([a-zA-Z$_][\w$]*)?(\s*)(\{)`, ByGroups(NameClass, Text, NameClass, Text, Punctuation), Push("#pop", "oc_ivars")},
			{`([a-zA-Z$_][\w$]*)(\s*:\s*)([a-zA-Z$_][\w$]*)?`, ByGroups(NameClass, Text, NameClass), Pop(1)},
			{`([a-zA-Z$_][\w$]*)(\s*)(\([a-zA-Z$_][\w$]*\))(\s*)(\{)`, ByGroups(NameClass, Text, NameLabel, Text, Punctuation), Push("#pop", "oc_ivars")},
			{`([a-zA-Z$_][\w$]*)(\s*)(\([a-zA-Z$_][\w$]*\))`, ByGroups(NameClass, Text, NameLabel), Pop(1)},
			{`([a-zA-Z$_][\w$]*)(\s*)(\{)`, ByGroups(NameClass, Text, Punctuation), Push("#pop", "oc_ivars")},
			{`([a-zA-Z$_][\w$]*)`, NameClass, Pop(1)},
		},
		"oc_forward_classname": {
			{`([a-zA-Z$_][\w$]*)(\s*,\s*)`, ByGroups(NameClass, Text), Push("oc_forward_classname")},
			{`([a-zA-Z$_][\w$]*)(\s*;?)`, ByGroups(NameClass, Text), Pop(1)},
		},
		"oc_ivars": {
			Include("whitespace"),
			Include("statements"),
			{`;`, Punctuation, nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
		},
		"root": {
			{`^([-+])(\s*)(\(.*?\))?(\s*)([a-zA-Z$_][\w$]*:?)`, ByGroups(Punctuation, Text, UsingSelf("root"), Text, NameFunction), Push("method")},
			Include("whitespace"),
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;{]*)(\{)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), Push("function")},
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;]*)(;)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), nil},
			Default(Push("statement")),
		},
		"method": {
			Include("whitespace"),
			{`,`, Punctuation, nil},
			{`\.\.\.`, Punctuation, nil},
			{`(\(.*?\))(\s*)([a-zA-Z$_][\w$]*)`, ByGroups(UsingSelf("root"), Text, NameVariable), nil},
			{`[a-zA-Z$_][\w$]*:`, NameFunction, nil},
			{`;`, Punctuation, Pop(1)},
			{`\{`, Punctuation, Push("function")},
			Default(Pop(1)),
		},
		"literal_number": {
			{`\(`, Punctuation, Push("literal_number_inner")},
			{`\)`, Literal, Pop(1)},
			Include("statement"),
		},
		"literal_number_inner": {
			{`\(`, Punctuation, Push()},
			{`\)`, Punctuation, Pop(1)},
			Include("statement"),
		},
		"literal_array": {
			{`\[`, Punctuation, Push("literal_array_inner")},
			{`\]`, Literal, Pop(1)},
			Include("statement"),
		},
		"literal_array_inner": {
			{`\[`, Punctuation, Push()},
			{`\]`, Punctuation, Pop(1)},
			Include("statement"),
		},
		"literal_dictionary": {
			{`\}`, Literal, Pop(1)},
			Include("statement"),
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
