package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Modula-2 lexer.
var Modula2 = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Modula-2",
		Aliases:   []string{"modula2", "m2"},
		Filenames: []string{"*.def", "*.mod"},
		MimeTypes: []string{"text/x-modula2"},
		DotAll:    true,
	},
	modula2Rules,
))

func modula2Rules() Rules {
	return Rules{
		"whitespace": {
			{`\n+`, Text, nil},
			{`\s+`, Text, nil},
		},
		"dialecttags": {
			{`\(\*!m2pim\*\)`, CommentSpecial, nil},
			{`\(\*!m2iso\*\)`, CommentSpecial, nil},
			{`\(\*!m2r10\*\)`, CommentSpecial, nil},
			{`\(\*!objm2\*\)`, CommentSpecial, nil},
			{`\(\*!m2iso\+aglet\*\)`, CommentSpecial, nil},
			{`\(\*!m2pim\+gm2\*\)`, CommentSpecial, nil},
			{`\(\*!m2iso\+p1\*\)`, CommentSpecial, nil},
			{`\(\*!m2iso\+xds\*\)`, CommentSpecial, nil},
		},
		"identifiers": {
			{`([a-zA-Z_$][\w$]*)`, Name, nil},
		},
		"prefixed_number_literals": {
			{`0b[01]+(\'[01]+)*`, LiteralNumberBin, nil},
			{`0[ux][0-9A-F]+(\'[0-9A-F]+)*`, LiteralNumberHex, nil},
		},
		"plain_number_literals": {
			{`[0-9]+(\'[0-9]+)*\.[0-9]+(\'[0-9]+)*[eE][+-]?[0-9]+(\'[0-9]+)*`, LiteralNumberFloat, nil},
			{`[0-9]+(\'[0-9]+)*\.[0-9]+(\'[0-9]+)*`, LiteralNumberFloat, nil},
			{`[0-9]+(\'[0-9]+)*`, LiteralNumberInteger, nil},
		},
		"suffixed_number_literals": {
			{`[0-7]+B`, LiteralNumberOct, nil},
			{`[0-7]+C`, LiteralNumberOct, nil},
			{`[0-9A-F]+H`, LiteralNumberHex, nil},
		},
		"string_literals": {
			{`'(\\\\|\\'|[^'])*'`, LiteralString, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
		},
		"digraph_operators": {
			{`\*\.`, Operator, nil},
			{`\+>`, Operator, nil},
			{`<>`, Operator, nil},
			{`<=`, Operator, nil},
			{`>=`, Operator, nil},
			{`==`, Operator, nil},
			{`::`, Operator, nil},
			{`:=`, Operator, nil},
			{`\+\+`, Operator, nil},
			{`--`, Operator, nil},
		},
		"unigraph_operators": {
			{`[+-]`, Operator, nil},
			{`[*/]`, Operator, nil},
			{`\\`, Operator, nil},
			{`[=#<>]`, Operator, nil},
			{`\^`, Operator, nil},
			{`@`, Operator, nil},
			{`&`, Operator, nil},
			{`~`, Operator, nil},
			{"`", Operator, nil},
		},
		"digraph_punctuation": {
			{`\.\.`, Punctuation, nil},
			{`<<`, Punctuation, nil},
			{`>>`, Punctuation, nil},
			{`->`, Punctuation, nil},
			{`\|#`, Punctuation, nil},
			{`##`, Punctuation, nil},
			{`\|\*`, Punctuation, nil},
		},
		"unigraph_punctuation": {
			{`[()\[\]{},.:;|]`, Punctuation, nil},
			{`!`, Punctuation, nil},
			{`\?`, Punctuation, nil},
		},
		"comments": {
			{`^//.*?\n`, CommentSingle, nil},
			{`\(\*([^$].*?)\*\)`, CommentMultiline, nil},
			{`/\*(.*?)\*/`, CommentMultiline, nil},
		},
		"pragmas": {
			{`<\*.*?\*>`, CommentPreproc, nil},
			{`\(\*\$.*?\*\)`, CommentPreproc, nil},
		},
		"root": {
			Include("whitespace"),
			Include("dialecttags"),
			Include("pragmas"),
			Include("comments"),
			Include("identifiers"),
			Include("suffixed_number_literals"),
			Include("prefixed_number_literals"),
			Include("plain_number_literals"),
			Include("string_literals"),
			Include("digraph_punctuation"),
			Include("digraph_operators"),
			Include("unigraph_punctuation"),
			Include("unigraph_operators"),
		},
	}
}
