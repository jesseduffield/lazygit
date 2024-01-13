package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Dylan lexer.
var Dylan = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Dylan",
		Aliases:         []string{"dylan"},
		Filenames:       []string{"*.dylan", "*.dyl", "*.intr"},
		MimeTypes:       []string{"text/x-dylan"},
		CaseInsensitive: true,
	},
	func() Rules {
		return Rules{
			"root": {
				{`\s+`, Whitespace, nil},
				{`//.*?\n`, CommentSingle, nil},
				{`([a-z0-9-]+:)([ \t]*)(.*(?:\n[ \t].+)*)`, ByGroups(NameAttribute, Whitespace, LiteralString), nil},
				Default(Push("code")),
			},
			"code": {
				{`\s+`, Whitespace, nil},
				{`//.*?\n`, CommentSingle, nil},
				{`/\*`, CommentMultiline, Push("comment")},
				{`"`, LiteralString, Push("string")},
				{`'(\\.|\\[0-7]{1,3}|\\x[a-f0-9]{1,2}|[^\\\'\n])'`, LiteralStringChar, nil},
				{`#b[01]+`, LiteralNumberBin, nil},
				{`#o[0-7]+`, LiteralNumberOct, nil},
				{`[-+]?(\d*\.\d+([ed][-+]?\d+)?|\d+(\.\d*)?e[-+]?\d+)`, LiteralNumberFloat, nil},
				{`[-+]?\d+`, LiteralNumberInteger, nil},
				{`#x[0-9a-f]+`, LiteralNumberHex, nil},

				{`(\?\\?)([\w!&*<>|^$%@+~?/=-]+)(:)(token|name|variable|expression|body|case-body|\*)`,
					ByGroups(Operator, NameVariable, Operator, NameBuiltin), nil},
				{`(\?)(:)(token|name|variable|expression|body|case-body|\*)`,
					ByGroups(Operator, Operator, NameVariable), nil},
				{`(\?\\?)([\w!&*<>|^$%@+~?/=-]+)`, ByGroups(Operator, NameVariable), nil},

				{`(=>|::|#\(|#\[|##|\?\?|\?=|\?|[(){}\[\],.;])`, Punctuation, nil},
				{`:=`, Operator, nil},
				{`#[tf]`, Literal, nil},
				{`#"`, LiteralStringSymbol, Push("symbol")},
				{`#[a-z0-9-]+`, Keyword, nil},
				{`#(all-keys|include|key|next|rest)`, Keyword, nil},
				{`[\w!&*<>|^$%@+~?/=-]+:`, KeywordConstant, nil},
				{`<[\w!&*<>|^$%@+~?/=-]+>`, NameClass, nil},
				{`\*[\w!&*<>|^$%@+~?/=-]+\*`, NameVariableGlobal, nil},
				{`\$[\w!&*<>|^$%@+~?/=-]+`, NameConstant, nil},
				{`(let|method|function)([ \t]+)([\w!&*<>|^$%@+~?/=-]+)`, ByGroups(NameBuiltin, Whitespace, NameVariable), nil},
				{`(error|signal|return|break)`, NameException, nil},
				{`(\\?)([\w!&*<>|^$%@+~?/=-]+)`, ByGroups(Operator, Name), nil},
			},
			"comment": {
				{`[^*/]`, CommentMultiline, nil},
				{`/\*`, CommentMultiline, Push()},
				{`\*/`, CommentMultiline, Pop(1)},
				{`[*/]`, CommentMultiline, nil},
			},
			"symbol": {
				{`"`, LiteralStringSymbol, Pop(1)},
				{`[^\\"]+`, LiteralStringSymbol, nil},
			},
			"string": {
				{`"`, LiteralString, Pop(1)},
				{`\\([\\abfnrtv"\']|x[a-f0-9]{2,4}|[0-7]{1,3})`, LiteralStringEscape, nil},
				{`[^\\"\n]+`, LiteralString, nil},
				{`\\\n`, LiteralString, nil},
				{`\\`, LiteralString, nil},
			},
		}
	},
))
