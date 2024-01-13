package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// ProtocolBuffer lexer.
var ProtocolBuffer = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Protocol Buffer",
		Aliases:   []string{"protobuf", "proto"},
		Filenames: []string{"*.proto"},
		MimeTypes: []string{},
	},
	protocolBufferRules,
))

func protocolBufferRules() Rules {
	return Rules{
		"root": {
			{`[ \t]+`, Text, nil},
			{`[,;{}\[\]()<>]`, Punctuation, nil},
			{`/(\\\n)?/(\n|(.|\n)*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?\*(.|\n)*?\*(\\\n)?/`, CommentMultiline, nil},
			{Words(`\b`, `\b`, `import`, `option`, `optional`, `required`, `repeated`, `default`, `packed`, `ctype`, `extensions`, `to`, `max`, `rpc`, `returns`, `oneof`), Keyword, nil},
			{Words(``, `\b`, `int32`, `int64`, `uint32`, `uint64`, `sint32`, `sint64`, `fixed32`, `fixed64`, `sfixed32`, `sfixed64`, `float`, `double`, `bool`, `string`, `bytes`), KeywordType, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{`(package)(\s+)`, ByGroups(KeywordNamespace, Text), Push("package")},
			{`(message|extend)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("message")},
			{`(enum|group|service)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("type")},
			{`\".*?\"`, LiteralString, nil},
			{`\'.*?\'`, LiteralString, nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[LlUu]*`, LiteralNumberFloat, nil},
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`(\-?(inf|nan))\b`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+[LlUu]*`, LiteralNumberHex, nil},
			{`0[0-7]+[LlUu]*`, LiteralNumberOct, nil},
			{`\d+[LlUu]*`, LiteralNumberInteger, nil},
			{`[+-=]`, Operator, nil},
			{`([a-zA-Z_][\w.]*)([ \t]*)(=)`, ByGroups(Name, Text, Operator), nil},
			{`[a-zA-Z_][\w.]*`, Name, nil},
		},
		"package": {
			{`[a-zA-Z_]\w*`, NameNamespace, Pop(1)},
			Default(Pop(1)),
		},
		"message": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
			Default(Pop(1)),
		},
		"type": {
			{`[a-zA-Z_]\w*`, Name, Pop(1)},
			Default(Pop(1)),
		},
	}
}
