package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Cfengine3 lexer.
var Cfengine3 = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "CFEngine3",
		Aliases:   []string{"cfengine3", "cf3"},
		Filenames: []string{"*.cf"},
		MimeTypes: []string{},
	},
	cfengine3Rules,
))

func cfengine3Rules() Rules {
	return Rules{
		"root": {
			{`#.*?\n`, Comment, nil},
			{`(body)(\s+)(\S+)(\s+)(control)`, ByGroups(Keyword, Text, Keyword, Text, Keyword), nil},
			{`(body|bundle)(\s+)(\S+)(\s+)(\w+)(\()`, ByGroups(Keyword, Text, Keyword, Text, NameFunction, Punctuation), Push("arglist")},
			{`(body|bundle)(\s+)(\S+)(\s+)(\w+)`, ByGroups(Keyword, Text, Keyword, Text, NameFunction), nil},
			{`(")([^"]+)(")(\s+)(string|slist|int|real)(\s*)(=>)(\s*)`, ByGroups(Punctuation, NameVariable, Punctuation, Text, KeywordType, Text, Operator, Text), nil},
			{`(\S+)(\s*)(=>)(\s*)`, ByGroups(KeywordReserved, Text, Operator, Text), nil},
			{`"`, LiteralString, Push("string")},
			{`(\w+)(\()`, ByGroups(NameFunction, Punctuation), nil},
			{`([\w.!&|()]+)(::)`, ByGroups(NameClass, Punctuation), nil},
			{`(\w+)(:)`, ByGroups(KeywordDeclaration, Punctuation), nil},
			{`@[{(][^)}]+[})]`, NameVariable, nil},
			{`[(){},;]`, Punctuation, nil},
			{`=>`, Operator, nil},
			{`->`, Operator, nil},
			{`\d+\.\d+`, LiteralNumberFloat, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`\w+`, NameFunction, nil},
			{`\s+`, Text, nil},
		},
		"string": {
			{`\$[{(]`, LiteralStringInterpol, Push("interpol")},
			{`\\.`, LiteralStringEscape, nil},
			{`"`, LiteralString, Pop(1)},
			{`\n`, LiteralString, nil},
			{`.`, LiteralString, nil},
		},
		"interpol": {
			{`\$[{(]`, LiteralStringInterpol, Push()},
			{`[})]`, LiteralStringInterpol, Pop(1)},
			{`[^${()}]+`, LiteralStringInterpol, nil},
		},
		"arglist": {
			{`\)`, Punctuation, Pop(1)},
			{`,`, Punctuation, nil},
			{`\w+`, NameVariable, nil},
			{`\s+`, Text, nil},
		},
	}
}
