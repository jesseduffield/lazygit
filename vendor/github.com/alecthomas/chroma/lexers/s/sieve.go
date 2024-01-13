package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Sieve lexer.
var Sieve = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Sieve",
		Aliases:   []string{"sieve"},
		Filenames: []string{"*.siv", "*.sieve"},
		MimeTypes: []string{},
	},
	func() Rules {
		return Rules{
			"root": {
				{`\s+`, Text, nil},
				{`[();,{}\[\]]`, Punctuation, nil},
				{`(?i)require`, KeywordNamespace, nil},
				{`(?i)(:)(addresses|all|contains|content|create|copy|comparator|count|days|detail|domain|fcc|flags|from|handle|importance|is|localpart|length|lowerfirst|lower|matches|message|mime|options|over|percent|quotewildcard|raw|regex|specialuse|subject|text|under|upperfirst|upper|value)`, ByGroups(NameTag, NameTag), nil},
				{`(?i)(address|addflag|allof|anyof|body|discard|elsif|else|envelope|ereject|exists|false|fileinto|if|hasflag|header|keep|notify_method_capability|notify|not|redirect|reject|removeflag|setflag|size|spamtest|stop|string|true|vacation|virustest)`, NameBuiltin, nil},
				{`(?i)set`, KeywordDeclaration, nil},
				{`([0-9.]+)([kmgKMG])?`, ByGroups(LiteralNumber, LiteralNumber), nil},
				{`#.*$`, CommentSingle, nil},
				{`/\*.*\*/`, CommentMultiline, nil},
				{`"[^"]*?"`, LiteralString, nil},
				{`text:`, NameTag, Push("text")},
			},
			"text": {
				{`[^.].*?\n`, LiteralString, nil},
				{`^\.`, Punctuation, Pop(1)},
			},
		}
	},
))
