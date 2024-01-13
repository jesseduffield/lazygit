package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var ArmAsm = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ArmAsm",
		Aliases:   []string{"armasm"},
		EnsureNL:  true,
		Filenames: []string{"*.s", "*.S"},
		MimeTypes: []string{"text/x-armasm", "text/x-asm"},
	},
	armasmRules,
))

func armasmRules() Rules {
	return Rules{
		"commentsandwhitespace": {
			{`\s+`, Text, nil},
			{`[@;].*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
		},
		"literal": {
			// Binary
			{`0b[01]+`, NumberBin, Pop(1)},
			// Hex
			{`0x\w{1,8}`, NumberHex, Pop(1)},
			// Octal
			{`0\d+`, NumberOct, Pop(1)},
			// Float
			{`\d+?\.\d+?`, NumberFloat, Pop(1)},
			// Integer
			{`\d+`, NumberInteger, Pop(1)},
			// String
			{`(")(.+)(")`, ByGroups(Punctuation, StringDouble, Punctuation), Pop(1)},
			// Char
			{`(')(.{1}|\\.{1})(')`, ByGroups(Punctuation, StringChar, Punctuation), Pop(1)},
		},
		"opcode": {
			// Escape at line end
			{`\n`, Text, Pop(1)},
			// Comment
			{`(@|;).*\n`, CommentSingle, Pop(1)},
			// Whitespace
			{`(\s+|,)`, Text, nil},
			// Register by number
			{`[rapcfxwbhsdqv]\d{1,2}`, NameClass, nil},
			// Address by hex
			{`=0x\w+`, ByGroups(Text, NameLabel), nil},
			// Pseudo address by label
			{`(=)(\w+)`, ByGroups(Text, NameLabel), nil},
			// Immediate
			{`#`, Text, Push("literal")},
		},
		"root": {
			Include("commentsandwhitespace"),
			// Directive with optional param
			{`(\.\w+)([ \t]+\w+\s+?)?`, ByGroups(KeywordNamespace, NameLabel), nil},
			// Label with data
			{`(\w+)(:)(\s+\.\w+\s+)`, ByGroups(NameLabel, Punctuation, KeywordNamespace), Push("literal")},
			// Label
			{`(\w+)(:)`, ByGroups(NameLabel, Punctuation), nil},
			// Syscall Op
			{`svc\s+\w+`, NameNamespace, nil},
			// Opcode
			{`[a-zA-Z]+`, Text, Push("opcode")},
		},
	}
}
