package k

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Kotlin lexer.
var Kotlin = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Kotlin",
		Aliases:   []string{"kotlin"},
		Filenames: []string{"*.kt"},
		MimeTypes: []string{"text/x-kotlin"},
		DotAll:    true,
	},
	kotlinRules,
))

func kotlinRules() Rules {
	const kotlinIdentifier = "(?:[_\\p{L}][\\p{L}\\p{N}]*|`@?[_\\p{L}][\\p{L}\\p{N}]+`)"

	return Rules{
		"root": {
			{`^\s*\[.*?\]`, NameAttribute, nil},
			{`[^\S\n]+`, Text, nil},
			{`\\\n`, Text, nil},
			{`//[^\n]*\n?`, CommentSingle, nil},
			{`/[*].*?[*]/`, CommentMultiline, nil},
			{`\n`, Text, nil},
			{`!==|!in|!is|===`, Operator, nil},
			{`%=|&&|\*=|\+\+|\+=|--|-=|->|\.\.|\/=|::|<=|==|>=|!!|!=|\|\||\?[:.]`, Operator, nil},
			{`[~!%^&*()+=|\[\]:;,.<>\/?-]`, Punctuation, nil},
			{`[{}]`, Punctuation, nil},
			{`"""`, LiteralString, Push("rawstring")},
			{`"`, LiteralStringDouble, Push("string")},
			{`(')(\\u[0-9a-fA-F]{4})(')`, ByGroups(LiteralStringChar, LiteralStringEscape, LiteralStringChar), nil},
			{`'\\.'|'[^\\]'`, LiteralStringChar, nil},
			{`0[xX][0-9a-fA-F]+[Uu]?[Ll]?|[0-9]+(\.[0-9]*)?([eE][+-][0-9]+)?[fF]?[Uu]?[Ll]?`, LiteralNumber, nil},
			{`(companion)(\s+)(object)`, ByGroups(Keyword, Text, Keyword), nil},
			{`(class|interface|object)(\s+)`, ByGroups(Keyword, Text), Push("class")},
			{`(package|import)(\s+)`, ByGroups(Keyword, Text), Push("package")},
			{`(val|var)(\s+)`, ByGroups(Keyword, Text), Push("property")},
			{`(fun)(\s+)`, ByGroups(Keyword, Text), Push("function")},
			{`(abstract|actual|annotation|as|as\?|break|by|catch|class|companion|const|constructor|continue|crossinline|data|delegate|do|dynamic|else|enum|expect|external|false|field|file|final|finally|for|fun|get|if|import|in|infix|init|inline|inner|interface|internal|is|it|lateinit|noinline|null|object|open|operator|out|override|package|param|private|property|protected|public|receiver|reified|return|sealed|set|setparam|super|suspend|tailrec|this|throw|true|try|typealias|typeof|val|var|vararg|when|where|while)\b`, Keyword, nil},
			{`@` + kotlinIdentifier, NameDecorator, nil},
			{kotlinIdentifier, Name, nil},
		},
		"package": {
			{`\S+`, NameNamespace, Pop(1)},
		},
		"class": {
			// \x60 is the back tick character (`)
			{`\x60[^\x60]+?\x60`, NameClass, Pop(1)},
			{kotlinIdentifier, NameClass, Pop(1)},
		},
		"property": {
			{`\x60[^\x60]+?\x60`, NameProperty, Pop(1)},
			{kotlinIdentifier, NameProperty, Pop(1)},
		},
		"generics-specification": {
			{`<`, Punctuation, Push("generics-specification")}, // required for generics inside generics e.g. <T : List<Int> >
			{`>`, Punctuation, Pop(1)},
			{`[,:*?]`, Punctuation, nil},
			{`(in|out|reified)`, Keyword, nil},
			{`\x60[^\x60]+?\x60`, NameClass, nil},
			{kotlinIdentifier, NameClass, nil},
			{`\s+`, Text, nil},
		},
		"function": {
			{`<`, Punctuation, Push("generics-specification")},
			{`\x60[^\x60]+?\x60`, NameFunction, Pop(1)},
			{kotlinIdentifier, NameFunction, Pop(1)},
			{`\s+`, Text, nil},
		},
		"rawstring": {
			// raw strings don't allow character escaping
			{`"""`, LiteralString, Pop(1)},
			{`(?:[^$"]+|\"{1,2}[^"])+`, LiteralString, nil},
			Include("string-interpol"),
			// remaining dollar signs are just a string
			{`\$`, LiteralString, nil},
		},
		"string": {
			{`\\[tbnr'"\\\$]`, LiteralStringEscape, nil},
			{`\\u[0-9a-fA-F]{4}`, LiteralStringEscape, nil},
			{`"`, LiteralStringDouble, Pop(1)},
			Include("string-interpol"),
			{`[^\n\\"$]+`, LiteralStringDouble, nil},
			// remaining dollar signs are just a string
			{`\$`, LiteralStringDouble, nil},
		},
		"string-interpol": {
			{`\$` + kotlinIdentifier, LiteralStringInterpol, nil},
			{`\${[^}\n]*}`, LiteralStringInterpol, nil},
		},
	}
}
