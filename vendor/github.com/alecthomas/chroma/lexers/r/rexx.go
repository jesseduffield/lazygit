package r

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Rexx lexer.
var Rexx = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Rexx",
		Aliases:         []string{"rexx", "arexx"},
		Filenames:       []string{"*.rexx", "*.rex", "*.rx", "*.arexx"},
		MimeTypes:       []string{"text/x-rexx"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	rexxRules,
))

func rexxRules() Rules {
	return Rules{
		"root": {
			{`\s`, TextWhitespace, nil},
			{`/\*`, CommentMultiline, Push("comment")},
			{`"`, LiteralString, Push("string_double")},
			{`'`, LiteralString, Push("string_single")},
			{`[0-9]+(\.[0-9]+)?(e[+-]?[0-9])?`, LiteralNumber, nil},
			{`([a-z_]\w*)(\s*)(:)(\s*)(procedure)\b`, ByGroups(NameFunction, TextWhitespace, Operator, TextWhitespace, KeywordDeclaration), nil},
			{`([a-z_]\w*)(\s*)(:)`, ByGroups(NameLabel, TextWhitespace, Operator), nil},
			Include("function"),
			Include("keyword"),
			Include("operator"),
			{`[a-z_]\w*`, Text, nil},
		},
		"function": {
			{Words(``, `(\s*)(\()`, `abbrev`, `abs`, `address`, `arg`, `b2x`, `bitand`, `bitor`, `bitxor`, `c2d`, `c2x`, `center`, `charin`, `charout`, `chars`, `compare`, `condition`, `copies`, `d2c`, `d2x`, `datatype`, `date`, `delstr`, `delword`, `digits`, `errortext`, `form`, `format`, `fuzz`, `insert`, `lastpos`, `left`, `length`, `linein`, `lineout`, `lines`, `max`, `min`, `overlay`, `pos`, `queued`, `random`, `reverse`, `right`, `sign`, `sourceline`, `space`, `stream`, `strip`, `substr`, `subword`, `symbol`, `time`, `trace`, `translate`, `trunc`, `value`, `verify`, `word`, `wordindex`, `wordlength`, `wordpos`, `words`, `x2b`, `x2c`, `x2d`, `xrange`), ByGroups(NameBuiltin, TextWhitespace, Operator), nil},
		},
		"keyword": {
			{`(address|arg|by|call|do|drop|else|end|exit|for|forever|if|interpret|iterate|leave|nop|numeric|off|on|options|parse|pull|push|queue|return|say|select|signal|to|then|trace|until|while)\b`, KeywordReserved, nil},
		},
		"operator": {
			{`(-|//|/|\(|\)|\*\*|\*|\\<<|\\<|\\==|\\=|\\>>|\\>|\\|\|\||\||&&|&|%|\+|<<=|<<|<=|<>|<|==|=|><|>=|>>=|>>|>|¬<<|¬<|¬==|¬=|¬>>|¬>|¬|\.|,)`, Operator, nil},
		},
		"string_double": {
			{`[^"\n]+`, LiteralString, nil},
			{`""`, LiteralString, nil},
			{`"`, LiteralString, Pop(1)},
			{`\n`, Text, Pop(1)},
		},
		"string_single": {
			{`[^\'\n]`, LiteralString, nil},
			{`\'\'`, LiteralString, nil},
			{`\'`, LiteralString, Pop(1)},
			{`\n`, Text, Pop(1)},
		},
		"comment": {
			{`[^*]+`, CommentMultiline, nil},
			{`\*/`, CommentMultiline, Pop(1)},
			{`\*`, CommentMultiline, nil},
		},
	}
}
