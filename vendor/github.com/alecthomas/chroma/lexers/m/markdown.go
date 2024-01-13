package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/h"
	"github.com/alecthomas/chroma/lexers/internal"
)

// Markdown lexer.
var Markdown = internal.Register(DelegatingLexer(h.HTML, MustNewLazyLexer(
	&Config{
		Name:      "markdown",
		Aliases:   []string{"md", "mkd"},
		Filenames: []string{"*.md", "*.mkd", "*.markdown"},
		MimeTypes: []string{"text/x-markdown"},
	},
	markdownRules,
)))

func markdownRules() Rules {
	return Rules{
		"root": {
			{`^(#[^#].+\n)`, ByGroups(GenericHeading), nil},
			{`^(#{2,6}.+\n)`, ByGroups(GenericSubheading), nil},
			{`^(\s*)([*-] )(\[[ xX]\])( .+\n)`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			{`^(\s*)([*-])(\s)(.+\n)`, ByGroups(Text, Keyword, Text, UsingSelf("inline")), nil},
			{`^(\s*)([0-9]+\.)( .+\n)`, ByGroups(Text, Keyword, UsingSelf("inline")), nil},
			{`^(\s*>\s)(.+\n)`, ByGroups(Keyword, GenericEmph), nil},
			{"^(```\\n)([\\w\\W]*?)(^```$)", ByGroups(String, Text, String), nil},
			{
				"^(```)(\\w+)(\\n)([\\w\\W]*?)(^```$)",
				UsingByGroup(
					internal.Get,
					2, 4,
					String, String, String, Text, String,
				),
				nil,
			},
			Include("inline"),
		},
		"inline": {
			{`\\.`, Text, nil},
			{`(\s)(\*|_)((?:(?!\2).)*)(\2)((?=\W|\n))`, ByGroups(Text, GenericEmph, GenericEmph, GenericEmph, Text), nil},
			{`(\s)((\*\*|__).*?)\3((?=\W|\n))`, ByGroups(Text, GenericStrong, GenericStrong, Text), nil},
			{`(\s)(~~[^~]+~~)((?=\W|\n))`, ByGroups(Text, GenericDeleted, Text), nil},
			{"`[^`]+`", LiteralStringBacktick, nil},
			{`[@#][\w/:]+`, NameEntity, nil},
			{`(!?\[)([^]]+)(\])(\()([^)]+)(\))`, ByGroups(Text, NameTag, Text, Text, NameAttribute, Text), nil},
			{`[^\\\s]+`, Other, nil},
			{`.|\n`, Other, nil},
		},
	}
}
