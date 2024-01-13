package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var MonkeyC = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "MonkeyC",
		Aliases:   []string{"monkeyc"},
		Filenames: []string{"*.mc"},
		MimeTypes: []string{"text/x-monkeyc"},
	},
	monkeyCRules,
))

func monkeyCRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`\n`, Text, nil},
			{`//(\n|[\w\W]*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil},
			{`/(\\\n)?[*][\w\W]*`, CommentMultiline, nil},
			{`:[a-zA-Z_][\w_\.]*`, StringSymbol, nil},
			{`[{}\[\]\(\),;:\.]`, Punctuation, nil},
			{`[&~\|\^!+\-*\/%=?]`, Operator, nil},
			{`=>|[+-]=|&&|\|\||>>|<<|[<>]=?|[!=]=`, Operator, nil},
			{`\b(and|or|instanceof|has|extends|new)`, OperatorWord, nil},
			{Words(``, `\b`, `NaN`, `null`, `true`, `false`), KeywordConstant, nil},
			{`(using)((?:\s|\\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`(class)((?:\s|\\\\s)+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`(function)((?:\s|\\\\s)+)`, ByGroups(KeywordDeclaration, Text), Push("function")},
			{`(module)((?:\s|\\\\s)+)`, ByGroups(KeywordDeclaration, Text), Push("module")},
			{`\b(if|else|for|switch|case|while|break|continue|default|do|try|catch|finally|return|throw|extends|function)\b`, Keyword, nil},
			{`\b(const|enum|hidden|public|protected|private|static)\b`, KeywordType, nil},
			{`\bvar\b`, KeywordDeclaration, nil},
			{`\b(Activity(Monitor|Recording)?|Ant(Plus)?|Application|Attention|Background|Communications|Cryptography|FitContributor|Graphics|Gregorian|Lang|Math|Media|Persisted(Content|Locations)|Position|Properties|Sensor(History|Logging)?|Storage|StringUtil|System|Test|Time(r)?|Toybox|UserProfile|WatchUi|Rez|Drawables|Strings|Fonts|method)\b`, NameBuiltin, nil},
			{`\b(me|self|\$)\b`, NameBuiltinPseudo, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^''])*'`, LiteralStringSingle, nil},
			{`-?(0x[0-9a-fA-F]+l?)`, NumberHex, nil},
			{`-?([0-9]+(\.[0-9]+[df]?|[df]))\b`, NumberFloat, nil},
			{`-?([0-9]+l?)`, NumberInteger, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"import": {
			{`([a-zA-Z_][\w_\.]*)(?:(\s+)(as)(\s+)([a-zA-Z_][\w_]*))?`, ByGroups(NameNamespace, Text, KeywordNamespace, Text, NameNamespace), nil},
			Default(Pop(1)),
		},
		"class": {
			{`([a-zA-Z_][\w_\.]*)(?:(\s+)(extends)(\s+)([a-zA-Z_][\w_\.]*))?`, ByGroups(NameClass, Text, KeywordDeclaration, Text, NameClass), nil},
			Default(Pop(1)),
		},
		"function": {
			{`initialize`, NameFunctionMagic, nil},
			{`[a-zA-Z_][\w_\.]*`, NameFunction, nil},
			Default(Pop(1)),
		},
		"module": {
			{`[a-zA-Z_][\w_\.]*`, NameNamespace, nil},
			Default(Pop(1)),
		},
	}
}
