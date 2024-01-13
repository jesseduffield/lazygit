package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Pony lexer.
var Pony = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Pony",
		Aliases:   []string{"pony"},
		Filenames: []string{"*.pony"},
		MimeTypes: []string{},
	},
	ponyRules,
))

func ponyRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`[^\S\n]+`, Text, nil},
			{`//.*\n`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("nested_comment")},
			{`"""(?:.|\n)*?"""`, LiteralStringDoc, nil},
			{`"`, LiteralString, Push("string")},
			{`\'.*\'`, LiteralStringChar, nil},
			{`=>|[]{}:().~;,|&!^?[]`, Punctuation, nil},
			{Words(``, `\b`, `addressof`, `and`, `as`, `consume`, `digestof`, `is`, `isnt`, `not`, `or`), OperatorWord, nil},
			{`!=|==|<<|>>|[-+/*%=<>]`, Operator, nil},
			{Words(``, `\b`, `box`, `break`, `compile_error`, `compile_intrinsic`, `continue`, `do`, `else`, `elseif`, `embed`, `end`, `error`, `for`, `if`, `ifdef`, `in`, `iso`, `lambda`, `let`, `match`, `object`, `recover`, `ref`, `repeat`, `return`, `tag`, `then`, `this`, `trn`, `try`, `until`, `use`, `var`, `val`, `where`, `while`, `with`, `#any`, `#read`, `#send`, `#share`), Keyword, nil},
			{`(actor|class|struct|primitive|interface|trait|type)((?:\s)+)`, ByGroups(Keyword, Text), Push("typename")},
			{`(new|fun|be)((?:\s)+)`, ByGroups(Keyword, Text), Push("methodname")},
			{Words(``, `\b`, `U8`, `U16`, `U32`, `U64`, `ULong`, `USize`, `U128`, `Unsigned`, `Stringable`, `String`, `StringBytes`, `StringRunes`, `InputNotify`, `InputStream`, `Stdin`, `ByteSeq`, `ByteSeqIter`, `OutStream`, `StdStream`, `SourceLoc`, `I8`, `I16`, `I32`, `I64`, `ILong`, `ISize`, `I128`, `Signed`, `Seq`, `RuntimeOptions`, `Real`, `Integer`, `SignedInteger`, `UnsignedInteger`, `FloatingPoint`, `Number`, `Int`, `ReadSeq`, `ReadElement`, `Pointer`, `Platform`, `NullablePointer`, `None`, `Iterator`, `F32`, `F64`, `Float`, `Env`, `DoNotOptimise`, `DisposableActor`, `Less`, `Equal`, `Greater`, `Compare`, `HasEq`, `Equatable`, `Comparable`, `Bool`, `AsioEventID`, `AsioEventNotify`, `AsioEvent`, `Array`, `ArrayKeys`, `ArrayValues`, `ArrayPairs`, `Any`, `AmbientAuth`), KeywordType, nil},
			{`_?[A-Z]\w*`, NameClass, nil},
			{`string\(\)`, NameOther, nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`(true|false)\b`, Keyword, nil},
			{`_\d*`, Name, nil},
			{`_?[a-z][\w\'_]*`, Name, nil},
		},
		"typename": {
			{`(iso|trn|ref|val|box|tag)?((?:\s)*)(_?[A-Z]\w*)`, ByGroups(Keyword, Text, NameClass), Pop(1)},
		},
		"methodname": {
			{`(iso|trn|ref|val|box|tag)?((?:\s)*)(_?[a-z]\w*)`, ByGroups(Keyword, Text, NameFunction), Pop(1)},
		},
		"nested_comment": {
			{`[^*/]+`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push()},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\"`, LiteralString, nil},
			{`[^\\"]+`, LiteralString, nil},
		},
	}
}
