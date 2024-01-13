package b

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Ballerina lexer.
var Ballerina = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Ballerina",
		Aliases:   []string{"ballerina"},
		Filenames: []string{"*.bal"},
		MimeTypes: []string{"text/x-ballerina"},
		DotAll:    true,
	},
	ballerinaRules,
))

func ballerinaRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`(break|catch|continue|done|else|finally|foreach|forever|fork|if|lock|match|return|throw|transaction|try|while)\b`, Keyword, nil},
			{`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},
			{`@[^\W\d][\w.]*`, NameDecorator, nil},
			{`(annotation|bind|but|endpoint|error|function|object|private|public|returns|service|type|var|with|worker)\b`, KeywordDeclaration, nil},
			{`(boolean|byte|decimal|float|int|json|map|nil|record|string|table|xml)\b`, KeywordType, nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(import)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
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
		"import": {
			{`[\w.]+`, NameNamespace, Pop(1)},
		},
	}
}
