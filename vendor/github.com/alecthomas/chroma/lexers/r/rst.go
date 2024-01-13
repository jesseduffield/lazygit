package r

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Restructuredtext lexer.
var Restructuredtext = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "reStructuredText",
		Aliases:   []string{"rst", "rest", "restructuredtext"},
		Filenames: []string{"*.rst", "*.rest"},
		MimeTypes: []string{"text/x-rst", "text/prs.fallenstein.rst"},
	},
	restructuredtextRules,
))

func restructuredtextRules() Rules {
	return Rules{
		"root": {
			{"^(=+|-+|`+|:+|\\.+|\\'+|\"+|~+|\\^+|_+|\\*+|\\++|#+)([ \\t]*\\n)(.+)(\\n)(\\1)(\\n)", ByGroups(GenericHeading, Text, GenericHeading, Text, GenericHeading, Text), nil},
			{"^(\\S.*)(\\n)(={3,}|-{3,}|`{3,}|:{3,}|\\.{3,}|\\'{3,}|\"{3,}|~{3,}|\\^{3,}|_{3,}|\\*{3,}|\\+{3,}|#{3,})(\\n)", ByGroups(GenericHeading, Text, GenericHeading, Text), nil},
			{`^(\s*)([-*+])( .+\n(?:\1  .+\n)*)`, ByGroups(Text, LiteralNumber, UsingSelf("inline")), nil},
			{`^(\s*)([0-9#ivxlcmIVXLCM]+\.)( .+\n(?:\1  .+\n)*)`, ByGroups(Text, LiteralNumber, UsingSelf("inline")), nil},
			{`^(\s*)(\(?[0-9#ivxlcmIVXLCM]+\))( .+\n(?:\1  .+\n)*)`, ByGroups(Text, LiteralNumber, UsingSelf("inline")), nil},
			{`^(\s*)([A-Z]+\.)( .+\n(?:\1  .+\n)+)`, ByGroups(Text, LiteralNumber, UsingSelf("inline")), nil},
			{`^(\s*)(\(?[A-Za-z]+\))( .+\n(?:\1  .+\n)+)`, ByGroups(Text, LiteralNumber, UsingSelf("inline")), nil},
			{`^(\s*)(\|)( .+\n(?:\|  .+\n)*)`, ByGroups(Text, Operator, UsingSelf("inline")), nil},
			{`^( *\.\.)(\s*)((?:source)?code(?:-block)?)(::)([ \t]*)([^\n]+)(\n[ \t]*\n)([ \t]+)(.*)(\n)((?:(?:\8.*|)\n)+)`, EmitterFunc(rstCodeBlock), nil},
			{`^( *\.\.)(\s*)([\w:-]+?)(::)(?:([ \t]*)(.*))`, ByGroups(Punctuation, Text, OperatorWord, Punctuation, Text, UsingSelf("inline")), nil},
			{`^( *\.\.)(\s*)(_(?:[^:\\]|\\.)+:)(.*?)$`, ByGroups(Punctuation, Text, NameTag, UsingSelf("inline")), nil},
			{`^( *\.\.)(\s*)(\[.+\])(.*?)$`, ByGroups(Punctuation, Text, NameTag, UsingSelf("inline")), nil},
			{`^( *\.\.)(\s*)(\|.+\|)(\s*)([\w:-]+?)(::)(?:([ \t]*)(.*))`, ByGroups(Punctuation, Text, NameTag, Text, OperatorWord, Punctuation, Text, UsingSelf("inline")), nil},
			{`^ *\.\..*(\n( +.*\n|\n)+)?`, CommentPreproc, nil},
			{`^( *)(:[a-zA-Z-]+:)(\s*)$`, ByGroups(Text, NameClass, Text), nil},
			{`^( *)(:.*?:)([ \t]+)(.*?)$`, ByGroups(Text, NameClass, Text, NameFunction), nil},
			{`^(\S.*(?<!::)\n)((?:(?: +.*)\n)+)`, ByGroups(UsingSelf("inline"), UsingSelf("inline")), nil},
			{`(::)(\n[ \t]*\n)([ \t]+)(.*)(\n)((?:(?:\3.*|)\n)+)`, ByGroups(LiteralStringEscape, Text, LiteralString, LiteralString, Text, LiteralString), nil},
			Include("inline"),
		},
		"inline": {
			{`\\.`, Text, nil},
			{"``", LiteralString, Push("literal")},
			{"(`.+?)(<.+?>)(`__?)", ByGroups(LiteralString, LiteralStringInterpol, LiteralString), nil},
			{"`.+?`__?", LiteralString, nil},
			{"(`.+?`)(:[a-zA-Z0-9:-]+?:)?", ByGroups(NameVariable, NameAttribute), nil},
			{"(:[a-zA-Z0-9:-]+?:)(`.+?`)", ByGroups(NameAttribute, NameVariable), nil},
			{`\*\*.+?\*\*`, GenericStrong, nil},
			{`\*.+?\*`, GenericEmph, nil},
			{`\[.*?\]_`, LiteralString, nil},
			{`<.+?>`, NameTag, nil},
			{"[^\\\\\\n\\[*`:]+", Text, nil},
			{`.`, Text, nil},
		},
		"literal": {
			{"[^`]+", LiteralString, nil},
			{"``((?=$)|(?=[-/:.,; \\n\\x00\\\u2010\\\u2011\\\u2012\\\u2013\\\u2014\\\u00a0\\'\\\"\\)\\]\\}\\>\\\u2019\\\u201d\\\u00bb\\!\\?]))", LiteralString, Pop(1)},
			{"`", LiteralString, nil},
		},
	}
}

func rstCodeBlock(groups []string, state *LexerState) Iterator {
	iterators := []Iterator{}
	tokens := []Token{
		{Punctuation, groups[1]},
		{Text, groups[2]},
		{OperatorWord, groups[3]},
		{Punctuation, groups[4]},
		{Text, groups[5]},
		{Keyword, groups[6]},
		{Text, groups[7]},
	}
	code := strings.Join(groups[8:], "")
	lexer := internal.Get(groups[6])
	if lexer == nil {
		tokens = append(tokens, Token{String, code})
		iterators = append(iterators, Literator(tokens...))
	} else {
		sub, err := lexer.Tokenise(nil, code)
		if err != nil {
			panic(err)
		}
		iterators = append(iterators, Literator(tokens...), sub)
	}
	return Concaterator(iterators...)
}
