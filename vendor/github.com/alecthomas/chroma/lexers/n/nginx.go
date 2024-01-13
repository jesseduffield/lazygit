package n

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Nginx Configuration File lexer.
var Nginx = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Nginx configuration file",
		Aliases:   []string{"nginx"},
		Filenames: []string{"nginx.conf"},
		MimeTypes: []string{"text/x-nginx-conf"},
	},
	nginxRules,
))

func nginxRules() Rules {
	return Rules{
		"root": {
			{`(include)(\s+)([^\s;]+)`, ByGroups(Keyword, Text, Name), nil},
			{`[^\s;#]+`, Keyword, Push("stmt")},
			Include("base"),
		},
		"block": {
			{`\}`, Punctuation, Pop(2)},
			{`[^\s;#]+`, KeywordNamespace, Push("stmt")},
			Include("base"),
		},
		"stmt": {
			{`\{`, Punctuation, Push("block")},
			{`;`, Punctuation, Pop(1)},
			Include("base"),
		},
		"base": {
			{`#.*\n`, CommentSingle, nil},
			{`on|off`, NameConstant, nil},
			{`\$[^\s;#()]+`, NameVariable, nil},
			{`([a-z0-9.-]+)(:)([0-9]+)`, ByGroups(Name, Punctuation, LiteralNumberInteger), nil},
			{`[a-z-]+/[a-z-+]+`, LiteralString, nil},
			{`[0-9]+[km]?\b`, LiteralNumberInteger, nil},
			{`(~)(\s*)([^\s{]+)`, ByGroups(Punctuation, Text, LiteralStringRegex), nil},
			{`[:=~]`, Punctuation, nil},
			{`[^\s;#{}$]+`, LiteralString, nil},
			{`/[^\s;#]*`, Name, nil},
			{`\s+`, Text, nil},
			{`[$;]`, Text, nil},
		},
	}
}
