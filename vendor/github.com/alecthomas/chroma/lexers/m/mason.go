package m

import (
	. "github.com/alecthomas/chroma"          // nolint
	. "github.com/alecthomas/chroma/lexers/h" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	. "github.com/alecthomas/chroma/lexers/p" // nolint
)

// Mason lexer.
var Mason = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Mason",
		Aliases:   []string{"mason"},
		Filenames: []string{"*.m", "*.mhtml", "*.mc", "*.mi", "autohandler", "dhandler"},
		MimeTypes: []string{"application/x-mason"},
		Priority:  0.1,
	},
	masonRules,
))

func masonRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`(<%doc>)(.*?)(</%doc>)(?s)`, ByGroups(NameTag, CommentMultiline, NameTag), nil},
			{`(<%(?:def|method))(\s*)(.*?)(>)(.*?)(</%\2\s*>)(?s)`, ByGroups(NameTag, Text, NameFunction, NameTag, UsingSelf("root"), NameTag), nil},
			{`(<%\w+)(.*?)(>)(.*?)(</%\2\s*>)(?s)`, ByGroups(NameTag, NameFunction, NameTag, Using(Perl), NameTag), nil},
			{`(<&[^|])(.*?)(,.*?)?(&>)(?s)`, ByGroups(NameTag, NameFunction, Using(Perl), NameTag), nil},
			{`(<&\|)(.*?)(,.*?)?(&>)(?s)`, ByGroups(NameTag, NameFunction, Using(Perl), NameTag), nil},
			{`</&>`, NameTag, nil},
			{`(<%!?)(.*?)(%>)(?s)`, ByGroups(NameTag, Using(Perl), NameTag), nil},
			{`(?<=^)#[^\n]*(\n|\Z)`, Comment, nil},
			{`(?<=^)(%)([^\n]*)(\n|\Z)`, ByGroups(NameTag, Using(Perl), Other), nil},
			{`(?sx)
                 (.+?)               # anything, followed by:
                 (?:
                  (?<=\n)(?=[%#]) |  # an eval or comment line
                  (?=</?[%&]) |      # a substitution or block or
                                     # call start or end
                                     # - don't consume
                  (\\\n) |           # an escaped newline
                  \Z                 # end of string
                 )`, ByGroups(Using(HTML), Operator), nil},
		},
	}
}
