package j

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Java lexer.
var Java = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Java",
		Aliases:   []string{"java"},
		Filenames: []string{"*.java"},
		MimeTypes: []string{"text/x-java"},
		DotAll:    true,
		EnsureNL:  true,
	},
	javaRules,
))

func javaRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`(assert|break|case|catch|continue|default|do|else|finally|for|if|goto|instanceof|new|return|switch|this|throw|try|while)\b`, Keyword, nil},
			{`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},
			{`@[^\W\d][\w.]*`, NameDecorator, nil},
			{`(abstract|const|enum|extends|final|implements|native|private|protected|public|static|strictfp|super|synchronized|throws|transient|volatile)\b`, KeywordDeclaration, nil},
			{`(boolean|byte|char|double|float|int|long|short|void)\b`, KeywordType, nil},
			{`(package)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(class|interface)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`(import(?:\s+static)?)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'\\.'|'[^\\]'|'\\u[0-9a-fA-F]{4}'`, LiteralStringChar, nil},
			{`(\.)((?:[^\W\d]|\$)[\w$]*)`, ByGroups(Operator, NameAttribute), nil},
			{`^\s*([^\W\d]|\$)[\w$]*:`, NameLabel, nil},
			{`([^\W\d]|\$)[\w$]*`, Name, nil},
			{`([0-9][0-9_]*\.([0-9][0-9_]*)?|\.[0-9][0-9_]*)([eE][+\-]?[0-9][0-9_]*)?[fFdD]?|[0-9][eE][+\-]?[0-9][0-9_]*[fFdD]?|[0-9]([eE][+\-]?[0-9][0-9_]*)?[fFdD]|0[xX]([0-9a-fA-F][0-9a-fA-F_]*\.?|([0-9a-fA-F][0-9a-fA-F_]*)?\.[0-9a-fA-F][0-9a-fA-F_]*)[pP][+\-]?[0-9][0-9_]*[fFdD]?`, LiteralNumberFloat, nil},
			{`0[xX][0-9a-fA-F][0-9a-fA-F_]*[lL]?`, LiteralNumberHex, nil},
			{`0[bB][01][01_]*[lL]?`, LiteralNumberBin, nil},
			{`0[0-7_]+[lL]?`, LiteralNumberOct, nil},
			{`0|[1-9][0-9_]*[lL]?`, LiteralNumberInteger, nil},
			{`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`([^\W\d]|\$)[\w$]*`, NameClass, Pop(1)},
		},
		"import": {
			{`[\w.]+\*?`, NameNamespace, Pop(1)},
		},
	}
}
