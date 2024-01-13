package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Dart lexer.
var Dart = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Dart",
		Aliases:   []string{"dart"},
		Filenames: []string{"*.dart"},
		MimeTypes: []string{"text/x-dart"},
		DotAll:    true,
	},
	dartRules,
))

func dartRules() Rules {
	return Rules{
		"root": {
			Include("string_literal"),
			{`#!(.*?)$`, CommentPreproc, nil},
			{`\b(import|export)\b`, Keyword, Push("import_decl")},
			{`\b(library|source|part of|part)\b`, Keyword, nil},
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`\b(class)\b(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`\b(assert|break|case|catch|continue|default|do|else|finally|for|if|in|is|new|return|super|switch|this|throw|try|while)\b`, Keyword, nil},
			{`\b(abstract|async|await|const|extends|factory|final|get|implements|native|operator|set|static|sync|typedef|var|with|yield)\b`, KeywordDeclaration, nil},
			{`\b(bool|double|dynamic|int|num|Object|String|void)\b`, KeywordType, nil},
			{`\b(false|null|true)\b`, KeywordConstant, nil},
			{`[~!%^&*+=|?:<>/-]|as\b`, Operator, nil},
			{`[a-zA-Z_$]\w*:`, NameLabel, nil},
			{`[a-zA-Z_$]\w*`, Name, nil},
			{`[(){}\[\],.;]`, Punctuation, nil},
			{`0[xX][0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`\d+(\.\d*)?([eE][+-]?\d+)?`, LiteralNumber, nil},
			{`\.\d+([eE][+-]?\d+)?`, LiteralNumber, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`[a-zA-Z_$]\w*`, NameClass, Pop(1)},
		},
		"import_decl": {
			Include("string_literal"),
			{`\s+`, Text, nil},
			{`\b(as|show|hide)\b`, Keyword, nil},
			{`[a-zA-Z_$]\w*`, Name, nil},
			{`\,`, Punctuation, nil},
			{`\;`, Punctuation, Pop(1)},
		},
		"string_literal": {
			{`r"""([\w\W]*?)"""`, LiteralStringDouble, nil},
			{`r'''([\w\W]*?)'''`, LiteralStringSingle, nil},
			{`r"(.*?)"`, LiteralStringDouble, nil},
			{`r'(.*?)'`, LiteralStringSingle, nil},
			{`"""`, LiteralStringDouble, Push("string_double_multiline")},
			{`'''`, LiteralStringSingle, Push("string_single_multiline")},
			{`"`, LiteralStringDouble, Push("string_double")},
			{`'`, LiteralStringSingle, Push("string_single")},
		},
		"string_common": {
			{`\\(x[0-9A-Fa-f]{2}|u[0-9A-Fa-f]{4}|u\{[0-9A-Fa-f]*\}|[a-z'\"$\\])`, LiteralStringEscape, nil},
			{`(\$)([a-zA-Z_]\w*)`, ByGroups(LiteralStringInterpol, Name), nil},
			{`(\$\{)(.*?)(\})`, ByGroups(LiteralStringInterpol, UsingSelf("root"), LiteralStringInterpol), nil},
		},
		"string_double": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`[^"$\\\n]+`, LiteralStringDouble, nil},
			Include("string_common"),
			{`\$+`, LiteralStringDouble, nil},
		},
		"string_double_multiline": {
			{`"""`, LiteralStringDouble, Pop(1)},
			{`[^"$\\]+`, LiteralStringDouble, nil},
			Include("string_common"),
			{`(\$|\")+`, LiteralStringDouble, nil},
		},
		"string_single": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`[^'$\\\n]+`, LiteralStringSingle, nil},
			Include("string_common"),
			{`\$+`, LiteralStringSingle, nil},
		},
		"string_single_multiline": {
			{`'''`, LiteralStringSingle, Pop(1)},
			{`[^\'$\\]+`, LiteralStringSingle, nil},
			Include("string_common"),
			{`(\$|\')+`, LiteralStringSingle, nil},
		},
	}
}
