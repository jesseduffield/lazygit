package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Awk lexer.
var Awk = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Awk",
		Aliases:   []string{"awk", "gawk", "mawk", "nawk"},
		Filenames: []string{"*.awk"},
		MimeTypes: []string{"application/x-awk"},
	},
	awkRules,
))

func awkRules() Rules {
	return Rules{
		"commentsandwhitespace": {
			{`\s+`, Text, nil},
			{`#.*$`, CommentSingle, nil},
		},
		"slashstartsregex": {
			Include("commentsandwhitespace"),
			{`/(\\.|[^[/\\\n]|\[(\\.|[^\]\\\n])*])+/\B`, LiteralStringRegex, Pop(1)},
			{`(?=/)`, Text, Push("#pop", "badregex")},
			Default(Pop(1)),
		},
		"badregex": {
			{`\n`, Text, Pop(1)},
		},
		"root": {
			{`^(?=\s|/)`, Text, Push("slashstartsregex")},
			Include("commentsandwhitespace"),
			{`\+\+|--|\|\||&&|in\b|\$|!?~|\|&|(\*\*|[-<>+*%\^/!=|])=?`, Operator, Push("slashstartsregex")},
			{`[{(\[;,]`, Punctuation, Push("slashstartsregex")},
			{`[})\].]`, Punctuation, nil},
			{`(break|continue|do|while|exit|for|if|else|return|switch|case|default)\b`, Keyword, Push("slashstartsregex")},
			{`function\b`, KeywordDeclaration, Push("slashstartsregex")},
			{`(atan2|cos|exp|int|log|rand|sin|sqrt|srand|gensub|gsub|index|length|match|split|patsplit|sprintf|sub|substr|tolower|toupper|close|fflush|getline|next(file)|print|printf|strftime|systime|mktime|delete|system|strtonum|and|compl|lshift|or|rshift|asorti?|isarray|bindtextdomain|dcn?gettext|@(include|load|namespace))\b`, KeywordReserved, nil},
			{`(ARGC|ARGIND|ARGV|BEGIN(FILE)?|BINMODE|CONVFMT|ENVIRON|END(FILE)?|ERRNO|FIELDWIDTHS|FILENAME|FNR|FPAT|FS|IGNORECASE|LINT|NF|NR|OFMT|OFS|ORS|PROCINFO|RLENGTH|RS|RSTART|RT|SUBSEP|TEXTDOMAIN)\b`, NameBuiltin, nil},
			{`[@$a-zA-Z_]\w*`, NameOther, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
		},
	}
}
