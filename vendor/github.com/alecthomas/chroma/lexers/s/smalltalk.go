package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Smalltalk lexer.
var Smalltalk = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Smalltalk",
		Aliases:   []string{"smalltalk", "squeak", "st"},
		Filenames: []string{"*.st"},
		MimeTypes: []string{"text/x-smalltalk"},
	},
	smalltalkRules,
))

func smalltalkRules() Rules {
	return Rules{
		"root": {
			{`(<)(\w+:)(.*?)(>)`, ByGroups(Text, Keyword, Text, Text), nil},
			Include("squeak fileout"),
			Include("whitespaces"),
			Include("method definition"),
			{`(\|)([\w\s]*)(\|)`, ByGroups(Operator, NameVariable, Operator), nil},
			Include("objects"),
			{`\^|:=|_`, Operator, nil},
			{`[\]({}.;!]`, Text, nil},
		},
		"method definition": {
			{`([a-zA-Z]+\w*:)(\s*)(\w+)`, ByGroups(NameFunction, Text, NameVariable), nil},
			{`^(\b[a-zA-Z]+\w*\b)(\s*)$`, ByGroups(NameFunction, Text), nil},
			{`^([-+*/\\~<>=|&!?,@%]+)(\s*)(\w+)(\s*)$`, ByGroups(NameFunction, Text, NameVariable, Text), nil},
		},
		"blockvariables": {
			Include("whitespaces"),
			{`(:)(\s*)(\w+)`, ByGroups(Operator, Text, NameVariable), nil},
			{`\|`, Operator, Pop(1)},
			Default(Pop(1)),
		},
		"literals": {
			{`'(''|[^'])*'`, LiteralString, Push("afterobject")},
			{`\$.`, LiteralStringChar, Push("afterobject")},
			{`#\(`, LiteralStringSymbol, Push("parenth")},
			{`\)`, Text, Push("afterobject")},
			{`(\d+r)?-?\d+(\.\d+)?(e-?\d+)?`, LiteralNumber, Push("afterobject")},
		},
		"_parenth_helper": {
			Include("whitespaces"),
			{`(\d+r)?-?\d+(\.\d+)?(e-?\d+)?`, LiteralNumber, nil},
			{`[-+*/\\~<>=|&#!?,@%\w:]+`, LiteralStringSymbol, nil},
			{`'(''|[^'])*'`, LiteralString, nil},
			{`\$.`, LiteralStringChar, nil},
			{`#*\(`, LiteralStringSymbol, Push("inner_parenth")},
		},
		"parenth": {
			{`\)`, LiteralStringSymbol, Push("root", "afterobject")},
			Include("_parenth_helper"),
		},
		"inner_parenth": {
			{`\)`, LiteralStringSymbol, Pop(1)},
			Include("_parenth_helper"),
		},
		"whitespaces": {
			{`\s+`, Text, nil},
			{`"(""|[^"])*"`, Comment, nil},
		},
		"objects": {
			{`\[`, Text, Push("blockvariables")},
			{`\]`, Text, Push("afterobject")},
			{`\b(self|super|true|false|nil|thisContext)\b`, NameBuiltinPseudo, Push("afterobject")},
			{`\b[A-Z]\w*(?!:)\b`, NameClass, Push("afterobject")},
			{`\b[a-z]\w*(?!:)\b`, NameVariable, Push("afterobject")},
			{`#("(""|[^"])*"|[-+*/\\~<>=|&!?,@%]+|[\w:]+)`, LiteralStringSymbol, Push("afterobject")},
			Include("literals"),
		},
		"afterobject": {
			{`! !$`, Keyword, Pop(1)},
			Include("whitespaces"),
			{`\b(ifTrue:|ifFalse:|whileTrue:|whileFalse:|timesRepeat:)`, NameBuiltin, Pop(1)},
			{`\b(new\b(?!:))`, NameBuiltin, nil},
			{`:=|_`, Operator, Pop(1)},
			{`\b[a-zA-Z]+\w*:`, NameFunction, Pop(1)},
			{`\b[a-zA-Z]+\w*`, NameFunction, nil},
			{`\w+:?|[-+*/\\~<>=|&!?,@%]+`, NameFunction, Pop(1)},
			{`\.`, Punctuation, Pop(1)},
			{`;`, Punctuation, nil},
			{`[\])}]`, Text, nil},
			{`[\[({]`, Text, Pop(1)},
		},
		"squeak fileout": {
			{`^"(""|[^"])*"!`, Keyword, nil},
			{`^'(''|[^'])*'!`, Keyword, nil},
			{`^(!)(\w+)( commentStamp: )(.*?)( prior: .*?!\n)(.*?)(!)`, ByGroups(Keyword, NameClass, Keyword, LiteralString, Keyword, Text, Keyword), nil},
			{`^(!)(\w+(?: class)?)( methodsFor: )('(?:''|[^'])*')(.*?!)`, ByGroups(Keyword, NameClass, Keyword, LiteralString, Keyword), nil},
			{`^(\w+)( subclass: )(#\w+)(\s+instanceVariableNames: )(.*?)(\s+classVariableNames: )(.*?)(\s+poolDictionaries: )(.*?)(\s+category: )(.*?)(!)`, ByGroups(NameClass, Keyword, LiteralStringSymbol, Keyword, LiteralString, Keyword, LiteralString, Keyword, LiteralString, Keyword, LiteralString, Keyword), nil},
			{`^(\w+(?: class)?)(\s+instanceVariableNames: )(.*?)(!)`, ByGroups(NameClass, Keyword, LiteralString, Keyword), nil},
			{`(!\n)(\].*)(! !)$`, ByGroups(Keyword, Text, Keyword), nil},
			{`! !$`, Keyword, nil},
		},
	}
}
