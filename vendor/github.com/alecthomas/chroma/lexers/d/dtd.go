package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Dtd lexer.
var Dtd = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "DTD",
		Aliases:   []string{"dtd"},
		Filenames: []string{"*.dtd"},
		MimeTypes: []string{"application/xml-dtd"},
		DotAll:    true,
	},
	dtdRules,
))

func dtdRules() Rules {
	return Rules{
		"root": {
			Include("common"),
			{`(<!ELEMENT)(\s+)(\S+)`, ByGroups(Keyword, Text, NameTag), Push("element")},
			{`(<!ATTLIST)(\s+)(\S+)`, ByGroups(Keyword, Text, NameTag), Push("attlist")},
			{`(<!ENTITY)(\s+)(\S+)`, ByGroups(Keyword, Text, NameEntity), Push("entity")},
			{`(<!NOTATION)(\s+)(\S+)`, ByGroups(Keyword, Text, NameTag), Push("notation")},
			{`(<!\[)([^\[\s]+)(\s*)(\[)`, ByGroups(Keyword, NameEntity, Text, Keyword), nil},
			{`(<!DOCTYPE)(\s+)([^>\s]+)`, ByGroups(Keyword, Text, NameTag), nil},
			{`PUBLIC|SYSTEM`, KeywordConstant, nil},
			{`[\[\]>]`, Keyword, nil},
		},
		"common": {
			{`\s+`, Text, nil},
			{`(%|&)[^;]*;`, NameEntity, nil},
			{`<!--`, Comment, Push("comment")},
			{`[(|)*,?+]`, Operator, nil},
			{`"[^"]*"`, LiteralStringDouble, nil},
			{`\'[^\']*\'`, LiteralStringSingle, nil},
		},
		"comment": {
			{`[^-]+`, Comment, nil},
			{`-->`, Comment, Pop(1)},
			{`-`, Comment, nil},
		},
		"element": {
			Include("common"),
			{`EMPTY|ANY|#PCDATA`, KeywordConstant, nil},
			{`[^>\s|()?+*,]+`, NameTag, nil},
			{`>`, Keyword, Pop(1)},
		},
		"attlist": {
			Include("common"),
			{`CDATA|IDREFS|IDREF|ID|NMTOKENS|NMTOKEN|ENTITIES|ENTITY|NOTATION`, KeywordConstant, nil},
			{`#REQUIRED|#IMPLIED|#FIXED`, KeywordConstant, nil},
			{`xml:space|xml:lang`, KeywordReserved, nil},
			{`[^>\s|()?+*,]+`, NameAttribute, nil},
			{`>`, Keyword, Pop(1)},
		},
		"entity": {
			Include("common"),
			{`SYSTEM|PUBLIC|NDATA`, KeywordConstant, nil},
			{`[^>\s|()?+*,]+`, NameEntity, nil},
			{`>`, Keyword, Pop(1)},
		},
		"notation": {
			Include("common"),
			{`SYSTEM|PUBLIC`, KeywordConstant, nil},
			{`[^>\s|()?+*,]+`, NameAttribute, nil},
			{`>`, Keyword, Pop(1)},
		},
	}
}
