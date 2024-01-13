package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Snobol lexer.
var Snobol = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Snobol",
		Aliases:   []string{"snobol"},
		Filenames: []string{"*.snobol"},
		MimeTypes: []string{"text/x-snobol"},
	},
	snobolRules,
))

func snobolRules() Rules {
	return Rules{
		"root": {
			{`\*.*\n`, Comment, nil},
			{`[+.] `, Punctuation, Push("statement")},
			{`-.*\n`, Comment, nil},
			{`END\s*\n`, NameLabel, Push("heredoc")},
			{`[A-Za-z$][\w$]*`, NameLabel, Push("statement")},
			{`\s+`, Text, Push("statement")},
		},
		"statement": {
			{`\s*\n`, Text, Pop(1)},
			{`\s+`, Text, nil},
			{`(?<=[^\w.])(LT|LE|EQ|NE|GE|GT|INTEGER|IDENT|DIFFER|LGT|SIZE|REPLACE|TRIM|DUPL|REMDR|DATE|TIME|EVAL|APPLY|OPSYN|LOAD|UNLOAD|LEN|SPAN|BREAK|ANY|NOTANY|TAB|RTAB|REM|POS|RPOS|FAIL|FENCE|ABORT|ARB|ARBNO|BAL|SUCCEED|INPUT|OUTPUT|TERMINAL)(?=[^\w.])`, NameBuiltin, nil},
			{`[A-Za-z][\w.]*`, Name, nil},
			{`\*\*|[?$.!%*/#+\-@|&\\=]`, Operator, nil},
			{`"[^"]*"`, LiteralString, nil},
			{`'[^']*'`, LiteralString, nil},
			{`[0-9]+(?=[^.EeDd])`, LiteralNumberInteger, nil},
			{`[0-9]+(\.[0-9]*)?([EDed][-+]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`:`, Punctuation, Push("goto")},
			{`[()<>,;]`, Punctuation, nil},
		},
		"goto": {
			{`\s*\n`, Text, Pop(2)},
			{`\s+`, Text, nil},
			{`F|S`, Keyword, nil},
			{`(\()([A-Za-z][\w.]*)(\))`, ByGroups(Punctuation, NameLabel, Punctuation), nil},
		},
		"heredoc": {
			{`.*\n`, LiteralStringHeredoc, nil},
		},
	}
}
