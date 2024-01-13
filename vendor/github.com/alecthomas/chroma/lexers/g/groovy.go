package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Groovy lexer.
var Groovy = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Groovy",
		Aliases:   []string{"groovy"},
		Filenames: []string{"*.groovy", "*.gradle"},
		MimeTypes: []string{"text/x-groovy"},
		DotAll:    true,
	},
	groovyRules,
))

func groovyRules() Rules {
	return Rules{
		"root": {
			{`#!(.*?)$`, CommentPreproc, Push("base")},
			Default(Push("base")),
		},
		"base": {
			{`^(\s*(?:[a-zA-Z_][\w.\[\]]*\s+)+?)([a-zA-Z_]\w*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`@[a-zA-Z_][\w.]*`, NameDecorator, nil},
			{`(as|assert|break|case|catch|continue|default|do|else|finally|for|if|in|goto|instanceof|new|return|switch|this|throw|try|while|in|as)\b`, Keyword, nil},
			{`(abstract|const|enum|extends|final|implements|native|private|protected|public|static|strictfp|super|synchronized|throws|transient|volatile)\b`, KeywordDeclaration, nil},
			{`(def|boolean|byte|char|double|float|int|long|short|void)\b`, KeywordType, nil},
			{`(package)(\s+)`, ByGroups(KeywordNamespace, Text), nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(class|interface)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`(import)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`""".*?"""`, LiteralStringDouble, nil},
			{`'''.*?'''`, LiteralStringSingle, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`\$/((?!/\$).)*/\$`, LiteralString, nil},
			{`/(\\\\|\\"|[^/])*/`, LiteralString, nil},
			{`'\\.'|'[^\\]'|'\\u[0-9a-fA-F]{4}'`, LiteralStringChar, nil},
			{`(\.)([a-zA-Z_]\w*)`, ByGroups(Operator, NameAttribute), nil},
			{`[a-zA-Z_]\w*:`, NameLabel, nil},
			{`[a-zA-Z_$]\w*`, Name, nil},
			{`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+L?`, LiteralNumberInteger, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
		},
		"import": {
			{`[\w.]+\*?`, NameNamespace, Pop(1)},
		},
	}
}
