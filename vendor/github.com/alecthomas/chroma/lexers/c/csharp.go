package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// CSharp lexer.
var CSharp = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "C#",
		Aliases:   []string{"csharp", "c#"},
		Filenames: []string{"*.cs"},
		MimeTypes: []string{"text/x-csharp"},
		DotAll:    true,
		EnsureNL:  true,
	},
	cSharpRules,
))

func cSharpRules() Rules {
	return Rules{
		"root": {
			{`^\s*\[.*?\]`, NameAttribute, nil},
			{`[^\S\n]+`, Text, nil},
			{`\\\n`, Text, nil},
			{`///[^\n\r]+`, CommentSpecial, nil},
			{`//[^\n\r]+`, CommentSingle, nil},
			{`/[*].*?[*]/`, CommentMultiline, nil},
			{`\n`, Text, nil},
			{`[~!%^&*()+=|\[\]:;,.<>/?-]`, Punctuation, nil},
			{`[{}]`, Punctuation, nil},
			{`@"(""|[^"])*"`, LiteralString, nil},
			{`\$@?"(""|[^"])*"`, LiteralString, nil},
			{`"(\\\\|\\"|[^"\n])*["\n]`, LiteralString, nil},
			{`'\\.'|'[^\\]'`, LiteralStringChar, nil},
			{`0[xX][0-9a-fA-F]+[Ll]?|[0-9_](\.[0-9]*)?([eE][+-]?[0-9]+)?[flFLdD]?`, LiteralNumber, nil},
			{`#[ \t]*(if|endif|else|elif|define|undef|line|error|warning|region|endregion|pragma|nullable)\b[^\n\r]+`, CommentPreproc, nil},
			{`\b(extern)(\s+)(alias)\b`, ByGroups(Keyword, Text, Keyword), nil},
			{`(abstract|as|async|await|base|break|by|case|catch|checked|const|continue|default|delegate|do|else|enum|event|explicit|extern|false|finally|fixed|for|foreach|goto|if|implicit|in|init|internal|is|let|lock|new|null|on|operator|out|override|params|private|protected|public|readonly|ref|return|sealed|sizeof|stackalloc|static|switch|this|throw|true|try|typeof|unchecked|unsafe|virtual|void|while|get|set|new|partial|yield|add|remove|value|alias|ascending|descending|from|group|into|orderby|select|thenby|where|join|equals)\b`, Keyword, nil},
			{`(global)(::)`, ByGroups(Keyword, Punctuation), nil},
			{`(bool|byte|char|decimal|double|dynamic|float|int|long|object|sbyte|short|string|uint|ulong|ushort|var)\b\??`, KeywordType, nil},
			{`(class|struct|record|interface)(\s+)`, ByGroups(Keyword, Text), Push("class")},
			{`(namespace|using)(\s+)`, ByGroups(Keyword, Text), Push("namespace")},
			{`@?[_a-zA-Z]\w*`, Name, nil},
		},
		"class": {
			{`@?[_a-zA-Z]\w*`, NameClass, Pop(1)},
			Default(Pop(1)),
		},
		"namespace": {
			{`(?=\()`, Text, Pop(1)},
			{`(@?[_a-zA-Z]\w*|\.)+`, NameNamespace, Pop(1)},
		},
	}
}
