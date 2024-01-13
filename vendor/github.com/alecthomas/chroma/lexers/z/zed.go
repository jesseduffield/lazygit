package z

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Zed lexer.
var Zed = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Zed",
		Aliases:   []string{"zed"},
		Filenames: []string{"*.zed"},
		MimeTypes: []string{"text/zed"},
	},
	zedRules,
).SetAnalyser(func(text string) float32 {
	if strings.Contains(text, "definition ") && strings.Contains(text, "relation ") && strings.Contains(text, "permission ") {
		return 0.9
	}
	if strings.Contains(text, "definition ") {
		return 0.5
	}
	if strings.Contains(text, "relation ") {
		return 0.5
	}
	if strings.Contains(text, "permission ") {
		return 0.25
	}
	return 0.0
}))

func zedRules() Rules {
	return Rules{
		"root": {
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil},
			{`/(\\\n)?[*][\w\W]*`, CommentMultiline, nil},
			{Words(``, `\b`, `definition`), KeywordType, nil},
			{Words(``, `\b`, `relation`), KeywordNamespace, nil},
			{Words(``, `\b`, `permission`), KeywordDeclaration, nil},
			{`[a-zA-Z_]\w*/`, NameNamespace, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
			{`#[a-zA-Z_]\w*`, NameVariable, nil},
			{`[+%=><|^!?/\-*&~:]`, Operator, nil},
			{`[{}()\[\],.;]`, Punctuation, nil},
		},
	}
}
