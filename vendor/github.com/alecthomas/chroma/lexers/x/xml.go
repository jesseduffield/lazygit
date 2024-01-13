package x

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// XML lexer.
var XML = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "XML",
		Aliases:   []string{"xml"},
		Filenames: []string{"*.xml", "*.xsl", "*.rss", "*.xslt", "*.xsd", "*.wsdl", "*.wsf", "*.svg", "*.csproj", "*.vcxproj", "*.fsproj"},
		MimeTypes: []string{"text/xml", "application/xml", "image/svg+xml", "application/rss+xml", "application/atom+xml"},
		DotAll:    true,
	},
	xmlRules,
))

func xmlRules() Rules {
	return Rules{
		"root": {
			{`[^<&]+`, Text, nil},
			{`&\S*?;`, NameEntity, nil},
			{`\<\!\[CDATA\[.*?\]\]\>`, CommentPreproc, nil},
			{`<!--`, Comment, Push("comment")},
			{`<\?.*?\?>`, CommentPreproc, nil},
			{`<![^>]*>`, CommentPreproc, nil},
			{`<\s*[\w:.-]+`, NameTag, Push("tag")},
			{`<\s*/\s*[\w:.-]+\s*>`, NameTag, nil},
		},
		"comment": {
			{`[^-]+`, Comment, nil},
			{`-->`, Comment, Pop(1)},
			{`-`, Comment, nil},
		},
		"tag": {
			{`\s+`, Text, nil},
			{`[\w.:-]+\s*=`, NameAttribute, Push("attr")},
			{`/?\s*>`, NameTag, Pop(1)},
		},
		"attr": {
			{`\s+`, Text, nil},
			{`".*?"`, LiteralString, Pop(1)},
			{`'.*?'`, LiteralString, Pop(1)},
			{`[^\s>]+`, LiteralString, Pop(1)},
		},
	}
}
