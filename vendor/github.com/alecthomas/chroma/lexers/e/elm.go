package e

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Elm lexer.
var Elm = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Elm",
		Aliases:   []string{"elm"},
		Filenames: []string{"*.elm"},
		MimeTypes: []string{"text/x-elm"},
	},
	elmRules,
))

func elmRules() Rules {
	return Rules{
		"root": {
			{`\{-`, CommentMultiline, Push("comment")},
			{`--.*`, CommentSingle, nil},
			{`\s+`, Text, nil},
			{`"`, LiteralString, Push("doublequote")},
			{`^\s*module\s*`, KeywordNamespace, Push("imports")},
			{`^\s*import\s*`, KeywordNamespace, Push("imports")},
			{`\[glsl\|.*`, NameEntity, Push("shader")},
			{Words(``, `\b`, `alias`, `as`, `case`, `else`, `if`, `import`, `in`, `let`, `module`, `of`, `port`, `then`, `type`, `where`), KeywordReserved, nil},
			{`[A-Z]\w*`, KeywordType, nil},
			{`^main `, KeywordReserved, nil},
			{Words(`\(`, `\)`, `~`, `||`, `|>`, `|`, "`", `^`, `\`, `'`, `>>`, `>=`, `>`, `==`, `=`, `<~`, `<|`, `<=`, `<<`, `<-`, `<`, `::`, `:`, `/=`, `//`, `/`, `..`, `.`, `->`, `-`, `++`, `+`, `*`, `&&`, `%`), NameFunction, nil},
			{Words(``, ``, `~`, `||`, `|>`, `|`, "`", `^`, `\`, `'`, `>>`, `>=`, `>`, `==`, `=`, `<~`, `<|`, `<=`, `<<`, `<-`, `<`, `::`, `:`, `/=`, `//`, `/`, `..`, `.`, `->`, `-`, `++`, `+`, `*`, `&&`, `%`), NameFunction, nil},
			Include("numbers"),
			{`[a-z_][a-zA-Z_\']*`, NameVariable, nil},
			{`[,()\[\]{}]`, Punctuation, nil},
		},
		"comment": {
			{`-(?!\})`, CommentMultiline, nil},
			{`\{-`, CommentMultiline, Push("comment")},
			{`[^-}]`, CommentMultiline, nil},
			{`-\}`, CommentMultiline, Pop(1)},
		},
		"doublequote": {
			{`\\u[0-9a-fA-F]{4}`, LiteralStringEscape, nil},
			{`\\[nrfvb\\"]`, LiteralStringEscape, nil},
			{`[^"]`, LiteralString, nil},
			{`"`, LiteralString, Pop(1)},
		},
		"imports": {
			{`\w+(\.\w+)*`, NameClass, Pop(1)},
		},
		"numbers": {
			{`_?\d+\.(?=\d+)`, LiteralNumberFloat, nil},
			{`_?\d+`, LiteralNumberInteger, nil},
		},
		"shader": {
			{`\|(?!\])`, NameEntity, nil},
			{`\|\]`, NameEntity, Pop(1)},
			{`.*\n`, NameEntity, nil},
		},
	}
}
