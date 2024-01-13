package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Termcap lexer.
var Termcap = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Termcap",
		Aliases:   []string{"termcap"},
		Filenames: []string{"termcap", "termcap.src"},
		MimeTypes: []string{},
	},
	termcapRules,
))

func termcapRules() Rules {
	return Rules{
		"root": {
			{`^#.*$`, Comment, nil},
			{`^[^\s#:|]+`, NameTag, Push("names")},
		},
		"names": {
			{`\n`, Text, Pop(1)},
			{`:`, Punctuation, Push("defs")},
			{`\|`, Punctuation, nil},
			{`[^:|]+`, NameAttribute, nil},
		},
		"defs": {
			{`\\\n[ \t]*`, Text, nil},
			{`\n[ \t]*`, Text, Pop(2)},
			{`(#)([0-9]+)`, ByGroups(Operator, LiteralNumber), nil},
			{`=`, Operator, Push("data")},
			{`:`, Punctuation, nil},
			{`[^\s:=#]+`, NameClass, nil},
		},
		"data": {
			{`\\072`, Literal, nil},
			{`:`, Punctuation, Pop(1)},
			{`[^:\\]+`, Literal, nil},
			{`.`, Literal, nil},
		},
	}
}
