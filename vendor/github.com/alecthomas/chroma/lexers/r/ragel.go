package r

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Ragel lexer.
var Ragel = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Ragel",
		Aliases:   []string{"ragel"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	ragelRules,
))

func ragelRules() Rules {
	return Rules{
		"whitespace": {
			{`\s+`, TextWhitespace, nil},
		},
		"comments": {
			{`\#.*$`, Comment, nil},
		},
		"keywords": {
			{`(access|action|alphtype)\b`, Keyword, nil},
			{`(getkey|write|machine|include)\b`, Keyword, nil},
			{`(any|ascii|extend|alpha|digit|alnum|lower|upper)\b`, Keyword, nil},
			{`(xdigit|cntrl|graph|print|punct|space|zlen|empty)\b`, Keyword, nil},
		},
		"numbers": {
			{`0x[0-9A-Fa-f]+`, LiteralNumberHex, nil},
			{`[+-]?[0-9]+`, LiteralNumberInteger, nil},
		},
		"literals": {
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralString, nil},
			{`\[(\\\\|\\\]|[^\]])*\]`, LiteralString, nil},
			{`/(?!\*)(\\\\|\\/|[^/])*/`, LiteralStringRegex, nil},
		},
		"identifiers": {
			{`[a-zA-Z_]\w*`, NameVariable, nil},
		},
		"operators": {
			{`,`, Operator, nil},
			{`\||&|--?`, Operator, nil},
			{`\.|<:|:>>?`, Operator, nil},
			{`:`, Operator, nil},
			{`->`, Operator, nil},
			{`(>|\$|%|<|@|<>)(/|eof\b)`, Operator, nil},
			{`(>|\$|%|<|@|<>)(!|err\b)`, Operator, nil},
			{`(>|\$|%|<|@|<>)(\^|lerr\b)`, Operator, nil},
			{`(>|\$|%|<|@|<>)(~|to\b)`, Operator, nil},
			{`(>|\$|%|<|@|<>)(\*|from\b)`, Operator, nil},
			{`>|@|\$|%`, Operator, nil},
			{`\*|\?|\+|\{[0-9]*,[0-9]*\}`, Operator, nil},
			{`!|\^`, Operator, nil},
			{`\(|\)`, Operator, nil},
		},
		"root": {
			Include("literals"),
			Include("whitespace"),
			Include("comments"),
			Include("keywords"),
			Include("numbers"),
			Include("identifiers"),
			Include("operators"),
			{`\{`, Punctuation, Push("host")},
			{`=`, Operator, nil},
			{`;`, Punctuation, nil},
		},
		"host": {
			{`([^{}\'"/#]+|[^\\]\\[{}]|"(\\\\|\\"|[^"])*"|'(\\\\|\\'|[^'])*'|//.*$\n?|/\*(.|\n)*?\*/|\#.*$\n?|/(?!\*)(\\\\|\\/|[^/])*/|/)+`, Other, nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
		},
	}
}
