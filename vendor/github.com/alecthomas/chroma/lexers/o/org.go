package o

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Org mode lexer.
var Org = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Org Mode",
		Aliases:   []string{"org", "orgmode"},
		Filenames: []string{"*.org"},
		MimeTypes: []string{"text/org"}, // https://lists.gnu.org/r/emacs-orgmode/2017-09/msg00087.html
	},
	orgRules,
))

func orgRules() Rules {
	return Rules{
		"root": {
			{`^# .*$`, Comment, nil},
			// Headings
			{`^(\*)( COMMENT)( .*)$`, ByGroups(GenericHeading, NameEntity, GenericStrong), nil},
			{`^(\*\*+)( COMMENT)( .*)$`, ByGroups(GenericSubheading, NameEntity, Text), nil},
			{`^(\*)( DONE)( .*)$`, ByGroups(GenericHeading, LiteralStringRegex, GenericStrong), nil},
			{`^(\*\*+)( DONE)( .*)$`, ByGroups(GenericSubheading, LiteralStringRegex, Text), nil},
			{`^(\*)( TODO)( .*)$`, ByGroups(GenericHeading, Error, GenericStrong), nil},
			{`^(\*\*+)( TODO)( .*)$`, ByGroups(GenericSubheading, Error, Text), nil},
			{`^(\*)( .+?)( :[a-zA-Z0-9_@:]+:)$`, ByGroups(GenericHeading, GenericStrong, GenericEmph), nil}, // Level 1 heading with tags
			{`^(\*)( .+)$`, ByGroups(GenericHeading, GenericStrong), nil},                                   // // Level 1 heading with NO tags
			{`^(\*\*+)( .+?)( :[a-zA-Z0-9_@:]+:)$`, ByGroups(GenericSubheading, Text, GenericEmph), nil},    // Level 2+ heading with tags
			{`^(\*\*+)( .+)$`, ByGroups(GenericSubheading, Text), nil},                                      // Level 2+ heading with NO tags
			// Checkbox lists
			{`^( *)([+-] )(\[[ X]\])( .+)$`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			{`^( +)(\* )(\[[ X]\])( .+)$`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			// Definition lists
			{`^( *)([+-] )([^ \n]+ ::)( .+)$`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			{`^( +)(\* )([^ \n]+ ::)( .+)$`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			// Unordered lists
			{`^( *)([+-] )(.+)$`, ByGroups(Text, Keyword, UsingSelf("inline")), nil},
			{`^( +)(\* )(.+)$`, ByGroups(Text, Keyword, UsingSelf("inline")), nil},
			// Ordered lists
			{`^( *)([0-9]+[.)])( \[@[0-9]+\])( .+)$`, ByGroups(Text, Keyword, GenericEmph, UsingSelf("inline")), nil},
			{`^( *)([0-9]+[.)])( .+)$`, ByGroups(Text, Keyword, UsingSelf("inline")), nil},
			// Dynamic Blocks
			{`(?i)^( *#\+begin: )([^ ]+)([\w\W]*?\n)([\w\W]*?)(^ *#\+end: *$)`, ByGroups(Comment, CommentSpecial, Comment, UsingSelf("inline"), Comment), nil},
			// Blocks
			// - Comment Blocks
			{`(?i)^( *#\+begin_comment *\n)([\w\W]*?)(^ *#\+end_comment *$)`, ByGroups(Comment, Comment, Comment), nil},
			// - Src Blocks
			{
				`(?i)^( *#\+begin_src )([^ \n]+)(.*?\n)([\w\W]*?)(^ *#\+end_src *$)`,
				UsingByGroup(
					internal.Get,
					2, 4,
					Comment, CommentSpecial, Comment, Text, Comment,
				),
				nil,
			},
			// - Export Blocks
			{
				`(?i)^( *#\+begin_export )(\w+)( *\n)([\w\W]*?)(^ *#\+end_export *$)`,
				UsingByGroup(
					internal.Get,
					2, 4,
					Comment, CommentSpecial, Text, Text, Comment,
				),
				nil,
			},
			// - Org Special, Example, Verse, etc. Blocks
			{`(?i)^( *#\+begin_)(\w+)( *\n)([\w\W]*?)(^ *#\+end_\2)( *$)`, ByGroups(Comment, Comment, Text, Text, Comment, Text), nil},
			// Keywords
			{`^(#\+\w+)(:.*)$`, ByGroups(CommentSpecial, Comment), nil}, // Other Org keywords like #+title
			// Properties and Drawers
			{`(?i)^( *:\w+: *\n)([\w\W]*?)(^ *:end: *$)`, ByGroups(Comment, CommentSpecial, Comment), nil},
			// Line break operator
			{`^(.*)(\\\\)$`, ByGroups(UsingSelf("inline"), Operator), nil},
			// Deadline/Scheduled
			{`(?i)^( *(?:DEADLINE|SCHEDULED): )(<[^<>]+?> *)$`, ByGroups(Comment, CommentSpecial), nil}, // DEADLINE/SCHEDULED: <datestamp>
			// DONE state CLOSED
			{`(?i)^( *CLOSED: )(\[[^][]+?\] *)$`, ByGroups(Comment, CommentSpecial), nil}, // CLOSED: [datestamp]
			// All other lines
			Include("inline"),
		},
		"inline": {
			{`(\s)*(\*[^ \n*][^*]+?[^ \n*]\*)((?=\W|\n|$))`, ByGroups(Text, GenericStrong, Text), nil},                          // Bold
			{`(\s)*(/[^/]+?/)((?=\W|\n|$))`, ByGroups(Text, GenericEmph, Text), nil},                                            // Italic
			{`(\s)*(=[^\n=]+?=)((?=\W|\n|$))`, ByGroups(Text, NameClass, Text), nil},                                            // Verbatim
			{`(\s)*(~[^\n~]+?~)((?=\W|\n|$))`, ByGroups(Text, NameClass, Text), nil},                                            // Code
			{`(\s)*(\+[^+]+?\+)((?=\W|\n|$))`, ByGroups(Text, GenericDeleted, Text), nil},                                       // Strikethrough
			{`(\s)*(_[^_]+?_)((?=\W|\n|$))`, ByGroups(Text, GenericUnderline, Text), nil},                                       // Underline
			{`(<)([^<>]+?)(>)`, ByGroups(Text, String, Text), nil},                                                              // <datestamp>
			{`[{]{3}[^}]+[}]{3}`, NameBuiltin, nil},                                                                             // {{{macro(foo,1)}}}
			{`([^[])(\[fn:)([^]]+?)(\])([^]])`, ByGroups(Text, NameBuiltinPseudo, LiteralString, NameBuiltinPseudo, Text), nil}, // [fn:1]
			// Links
			{`(\[\[)([^][]+?)(\]\[)([^][]+)(\]\])`, ByGroups(Text, NameAttribute, Text, NameTag, Text), nil}, // [[link][descr]]
			{`(\[\[)([^][]+?)(\]\])`, ByGroups(Text, NameAttribute, Text), nil},                              // [[link]]
			{`(<<)([^<>]+?)(>>)`, ByGroups(Text, NameAttribute, Text), nil},                                  // <<targetlink>>
			// Tables
			{`^( *)(\|[ -].*?[ -]\|)$`, ByGroups(Text, String), nil},
			// Blank lines, newlines
			{`\n`, Text, nil},
			// Any other text
			{`.`, Text, nil},
		},
	}
}
