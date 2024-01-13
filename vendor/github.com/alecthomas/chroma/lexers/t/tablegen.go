package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// TableGen lexer.
var Tablegen = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "TableGen",
		Aliases:   []string{"tablegen"},
		Filenames: []string{"*.td"},
		MimeTypes: []string{"text/x-tablegen"},
	},
	tablegenRules,
))

func tablegenRules() Rules {
	return Rules{
		"root": {
			Include("macro"),
			Include("whitespace"),
			{`c?"[^"]*?"`, LiteralString, nil},
			Include("keyword"),
			{`\$[_a-zA-Z][_\w]*`, NameVariable, nil},
			{`\d*[_a-zA-Z][_\w]*`, NameVariable, nil},
			{`\[\{[\w\W]*?\}\]`, LiteralString, nil},
			{`[+-]?\d+|0x[\da-fA-F]+|0b[01]+`, LiteralNumber, nil},
			{`[=<>{}\[\]()*.,!:;]`, Punctuation, nil},
		},
		"macro": {
			{`(#include\s+)("[^"]*")`, ByGroups(CommentPreproc, LiteralString), nil},
			{`^\s*#(ifdef|ifndef)\s+[_\w][_\w\d]*`, CommentPreproc, nil},
			{`^\s*#define\s+[_\w][_\w\d]*`, CommentPreproc, nil},
			{`^\s*#endif`, CommentPreproc, nil},
		},
		"whitespace": {
			{`(\n|\s)+`, Text, nil},
			{`//.*?\n`, Comment, nil},
		},
		"keyword": {
			{Words(``, `\b`, `bit`, `bits`, `class`, `code`, `dag`, `def`, `defm`, `field`, `foreach`, `in`, `int`, `let`, `list`, `multiclass`, `string`), Keyword, nil},
		},
	}
}
