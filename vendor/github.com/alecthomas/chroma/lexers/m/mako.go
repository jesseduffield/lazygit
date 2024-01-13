package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	. "github.com/alecthomas/chroma/lexers/p" // nolint
)

// Mako lexer.
var Mako = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Mako",
		Aliases:   []string{"mako"},
		Filenames: []string{"*.mao"},
		MimeTypes: []string{"application/x-mako"},
	},
	makoRules,
))

func makoRules() Rules {
	return Rules{
		"root": {
			{`(\s*)(%)(\s*end(?:\w+))(\n|\Z)`, ByGroups(Text, CommentPreproc, Keyword, Other), nil},
			{`(\s*)(%)([^\n]*)(\n|\Z)`, ByGroups(Text, CommentPreproc, Using(Python), Other), nil},
			{`(\s*)(##[^\n]*)(\n|\Z)`, ByGroups(Text, CommentPreproc, Other), nil},
			{`(?s)<%doc>.*?</%doc>`, CommentPreproc, nil},
			{`(<%)([\w.:]+)`, ByGroups(CommentPreproc, NameBuiltin), Push("tag")},
			{`(</%)([\w.:]+)(>)`, ByGroups(CommentPreproc, NameBuiltin, CommentPreproc), nil},
			{`<%(?=([\w.:]+))`, CommentPreproc, Push("ondeftags")},
			{`(<%(?:!?))(.*?)(%>)(?s)`, ByGroups(CommentPreproc, Using(Python), CommentPreproc), nil},
			{`(\$\{)(.*?)(\})`, ByGroups(CommentPreproc, Using(Python), CommentPreproc), nil},
			{`(?sx)
                (.+?)                # anything, followed by:
                (?:
                 (?<=\n)(?=%|\#\#) | # an eval or comment line
                 (?=\#\*) |          # multiline comment
                 (?=</?%) |          # a python block
                                     # call start or end
                 (?=\$\{) |          # a substitution
                 (?<=\n)(?=\s*%) |
                                     # - don't consume
                 (\\\n) |            # an escaped newline
                 \Z                  # end of string
                )
            `, ByGroups(Other, Operator), nil},
			{`\s+`, Text, nil},
		},
		"ondeftags": {
			{`<%`, CommentPreproc, nil},
			{`(?<=<%)(include|inherit|namespace|page)`, NameBuiltin, nil},
			Include("tag"),
		},
		"tag": {
			{`((?:\w+)\s*=)(\s*)(".*?")`, ByGroups(NameAttribute, Text, LiteralString), nil},
			{`/?\s*>`, CommentPreproc, Pop(1)},
			{`\s+`, Text, nil},
		},
		"attr": {
			{`".*?"`, LiteralString, Pop(1)},
			{`'.*?'`, LiteralString, Pop(1)},
			{`[^\s>]+`, LiteralString, Pop(1)},
		},
	}
}
