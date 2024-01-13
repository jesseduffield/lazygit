package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Ceylon lexer.
var Ceylon = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Ceylon",
		Aliases:   []string{"ceylon"},
		Filenames: []string{"*.ceylon"},
		MimeTypes: []string{"text/x-ceylon"},
		DotAll:    true,
	},
	ceylonRules,
))

func ceylonRules() Rules {
	return Rules{
		"root": {
			{`^(\s*(?:[a-zA-Z_][\w.\[\]]*\s+)+?)([a-zA-Z_]\w*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("comment")},
			{`(shared|abstract|formal|default|actual|variable|deprecated|small|late|literal|doc|by|see|throws|optional|license|tagged|final|native|annotation|sealed)\b`, NameDecorator, nil},
			{`(break|case|catch|continue|else|finally|for|in|if|return|switch|this|throw|try|while|is|exists|dynamic|nonempty|then|outer|assert|let)\b`, Keyword, nil},
			{`(abstracts|extends|satisfies|super|given|of|out|assign)\b`, KeywordDeclaration, nil},
			{`(function|value|void|new)\b`, KeywordType, nil},
			{`(assembly|module|package)(\s+)`, ByGroups(KeywordNamespace, Text), nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(class|interface|object|alias)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`(import)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'\\.'|'[^\\]'|'\\\{#[0-9a-fA-F]{4}\}'`, LiteralStringChar, nil},
			{"\".*``.*``.*\"", LiteralStringInterpol, nil},
			{`(\.)([a-z_]\w*)`, ByGroups(Operator, NameAttribute), nil},
			{`[a-zA-Z_]\w*:`, NameLabel, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
			{`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator, nil},
			{`\d{1,3}(_\d{3})+\.\d{1,3}(_\d{3})+[kMGTPmunpf]?`, LiteralNumberFloat, nil},
			{`\d{1,3}(_\d{3})+\.[0-9]+([eE][+-]?[0-9]+)?[kMGTPmunpf]?`, LiteralNumberFloat, nil},
			{`[0-9][0-9]*\.\d{1,3}(_\d{3})+[kMGTPmunpf]?`, LiteralNumberFloat, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][+-]?[0-9]+)?[kMGTPmunpf]?`, LiteralNumberFloat, nil},
			{`#([0-9a-fA-F]{4})(_[0-9a-fA-F]{4})+`, LiteralNumberHex, nil},
			{`#[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`\$([01]{4})(_[01]{4})+`, LiteralNumberBin, nil},
			{`\$[01]+`, LiteralNumberBin, nil},
			{`\d{1,3}(_\d{3})+[kMGTP]?`, LiteralNumberInteger, nil},
			{`[0-9]+[kMGTP]?`, LiteralNumberInteger, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`[A-Za-z_]\w*`, NameClass, Pop(1)},
		},
		"import": {
			{`[a-z][\w.]*`, NameNamespace, Pop(1)},
		},
		"comment": {
			{`[^*/]`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push()},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
	}
}
