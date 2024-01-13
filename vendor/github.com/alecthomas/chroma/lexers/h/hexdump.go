package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Hexdump lexer.
var Hexdump = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Hexdump",
		Aliases:   []string{"hexdump"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	hexdumpRules,
))

func hexdumpRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			Include("offset"),
			{`([0-9A-Ha-h]{2})(\-)([0-9A-Ha-h]{2})`, ByGroups(LiteralNumberHex, Punctuation, LiteralNumberHex), nil},
			{`[0-9A-Ha-h]{2}`, LiteralNumberHex, nil},
			{`(\s{2,3})(\>)(.{16})(\<)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), Push("bracket-strings")},
			{`(\s{2,3})(\|)(.{16})(\|)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), Push("piped-strings")},
			{`(\s{2,3})(\>)(.{1,15})(\<)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), nil},
			{`(\s{2,3})(\|)(.{1,15})(\|)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), nil},
			{`(\s{2,3})(.{1,15})$`, ByGroups(Text, LiteralString), nil},
			{`(\s{2,3})(.{16}|.{20})$`, ByGroups(Text, LiteralString), Push("nonpiped-strings")},
			{`\s`, Text, nil},
			{`^\*`, Punctuation, nil},
		},
		"offset": {
			{`^([0-9A-Ha-h]+)(:)`, ByGroups(NameLabel, Punctuation), Push("offset-mode")},
			{`^[0-9A-Ha-h]+`, NameLabel, nil},
		},
		"offset-mode": {
			{`\s`, Text, Pop(1)},
			{`[0-9A-Ha-h]+`, NameLabel, nil},
			{`:`, Punctuation, nil},
		},
		"piped-strings": {
			{`\n`, Text, nil},
			Include("offset"),
			{`[0-9A-Ha-h]{2}`, LiteralNumberHex, nil},
			{`(\s{2,3})(\|)(.{1,16})(\|)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), nil},
			{`\s`, Text, nil},
			{`^\*`, Punctuation, nil},
		},
		"bracket-strings": {
			{`\n`, Text, nil},
			Include("offset"),
			{`[0-9A-Ha-h]{2}`, LiteralNumberHex, nil},
			{`(\s{2,3})(\>)(.{1,16})(\<)$`, ByGroups(Text, Punctuation, LiteralString, Punctuation), nil},
			{`\s`, Text, nil},
			{`^\*`, Punctuation, nil},
		},
		"nonpiped-strings": {
			{`\n`, Text, nil},
			Include("offset"),
			{`([0-9A-Ha-h]{2})(\-)([0-9A-Ha-h]{2})`, ByGroups(LiteralNumberHex, Punctuation, LiteralNumberHex), nil},
			{`[0-9A-Ha-h]{2}`, LiteralNumberHex, nil},
			{`(\s{19,})(.{1,20}?)$`, ByGroups(Text, LiteralString), nil},
			{`(\s{2,3})(.{1,20})$`, ByGroups(Text, LiteralString), nil},
			{`\s`, Text, nil},
			{`^\*`, Punctuation, nil},
		},
	}
}
