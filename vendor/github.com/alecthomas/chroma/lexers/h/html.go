package h

import (
	. "github.com/alecthomas/chroma"          // nolint
	. "github.com/alecthomas/chroma/lexers/c" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	. "github.com/alecthomas/chroma/lexers/j" // nolint
)

// HTML lexer.
var HTML = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "HTML",
		Aliases:         []string{"html"},
		Filenames:       []string{"*.html", "*.htm", "*.xhtml", "*.xslt"},
		MimeTypes:       []string{"text/html", "application/xhtml+xml"},
		NotMultiline:    true,
		DotAll:          true,
		CaseInsensitive: true,
	},
	htmlRules,
))

func htmlRules() Rules {
	return Rules{
		"root": {
			{`[^<&]+`, Text, nil},
			{`&\S*?;`, NameEntity, nil},
			{`\<\!\[CDATA\[.*?\]\]\>`, CommentPreproc, nil},
			{`<!--`, Comment, Push("comment")},
			{`<\?.*?\?>`, CommentPreproc, nil},
			{`<![^>]*>`, CommentPreproc, nil},
			{`(<)(\s*)(script)(\s*)`, ByGroups(Punctuation, Text, NameTag, Text), Push("script-content", "tag")},
			{`(<)(\s*)(style)(\s*)`, ByGroups(Punctuation, Text, NameTag, Text), Push("style-content", "tag")},
			{`(<)(\s*)([\w:.-]+)`, ByGroups(Punctuation, Text, NameTag), Push("tag")},
			{`(<)(\s*)(/)(\s*)([\w:.-]+)(\s*)(>)`, ByGroups(Punctuation, Text, Punctuation, Text, NameTag, Text, Punctuation), nil},
		},
		"comment": {
			{`[^-]+`, Comment, nil},
			{`-->`, Comment, Pop(1)},
			{`-`, Comment, nil},
		},
		"tag": {
			{`\s+`, Text, nil},
			{`([\w:-]+\s*)(=)(\s*)`, ByGroups(NameAttribute, Operator, Text), Push("attr")},
			{`[\w:-]+`, NameAttribute, nil},
			{`(/?)(\s*)(>)`, ByGroups(Punctuation, Text, Punctuation), Pop(1)},
		},
		"script-content": {
			{`(<)(\s*)(/)(\s*)(script)(\s*)(>)`, ByGroups(Punctuation, Text, Punctuation, Text, NameTag, Text, Punctuation), Pop(1)},
			{`.+?(?=<\s*/\s*script\s*>)`, Using(Javascript), nil},
		},
		"style-content": {
			{`(<)(\s*)(/)(\s*)(style)(\s*)(>)`, ByGroups(Punctuation, Text, Punctuation, Text, NameTag, Text, Punctuation), Pop(1)},
			{`.+?(?=<\s*/\s*style\s*>)`, Using(CSS), nil},
		},
		"attr": {
			{`".*?"`, LiteralString, Pop(1)},
			{`'.*?'`, LiteralString, Pop(1)},
			{`[^\s>]+`, LiteralString, Pop(1)},
		},
	}
}
