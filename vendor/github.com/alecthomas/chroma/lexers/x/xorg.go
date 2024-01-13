package x

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Xorg lexer.
var Xorg = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Xorg",
		Aliases:   []string{"xorg.conf"},
		Filenames: []string{"xorg.conf"},
		MimeTypes: []string{},
	},
	xorgRules,
))

func xorgRules() Rules {
	return Rules{
		"root": {
			{`\s+`, TextWhitespace, nil},
			{`#.*$`, Comment, nil},
			{`((|Sub)Section)(\s+)("\w+")`, ByGroups(KeywordNamespace, LiteralStringEscape, TextWhitespace, LiteralStringEscape), nil},
			{`(End(|Sub)Section)`, KeywordNamespace, nil},
			{`(\w+)(\s+)([^\n#]+)`, ByGroups(NameKeyword, TextWhitespace, LiteralString), nil},
		},
	}
}
