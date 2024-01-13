package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Twig lexer.
var Twig = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Twig",
		Aliases:   []string{"twig"},
		Filenames: []string{},
		MimeTypes: []string{"application/x-twig"},
		DotAll:    true,
	},
	twigRules,
))

func twigRules() Rules {
	return Rules{
		"root": {
			{`[^{]+`, Other, nil},
			{`\{\{`, CommentPreproc, Push("var")},
			{`\{\#.*?\#\}`, Comment, nil},
			{`(\{%)(-?\s*)(raw)(\s*-?)(%\})(.*?)(\{%)(-?\s*)(endraw)(\s*-?)(%\})`, ByGroups(CommentPreproc, Text, Keyword, Text, CommentPreproc, Other, CommentPreproc, Text, Keyword, Text, CommentPreproc), nil},
			{`(\{%)(-?\s*)(verbatim)(\s*-?)(%\})(.*?)(\{%)(-?\s*)(endverbatim)(\s*-?)(%\})`, ByGroups(CommentPreproc, Text, Keyword, Text, CommentPreproc, Other, CommentPreproc, Text, Keyword, Text, CommentPreproc), nil},
			{`(\{%)(-?\s*)(filter)(\s+)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w-]|[^\x00-\x7f])*)`, ByGroups(CommentPreproc, Text, Keyword, Text, NameFunction), Push("tag")},
			{`(\{%)(-?\s*)([a-zA-Z_]\w*)`, ByGroups(CommentPreproc, Text, Keyword), Push("tag")},
			{`\{`, Other, nil},
		},
		"varnames": {
			{`(\|)(\s*)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w-]|[^\x00-\x7f])*)`, ByGroups(Operator, Text, NameFunction), nil},
			{`(is)(\s+)(not)?(\s*)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w-]|[^\x00-\x7f])*)`, ByGroups(Keyword, Text, Keyword, Text, NameFunction), nil},
			{`(?i)(true|false|none|null)\b`, KeywordPseudo, nil},
			{`(in|not|and|b-and|or|b-or|b-xor|isif|elseif|else|importconstant|defined|divisibleby|empty|even|iterable|odd|sameasmatches|starts\s+with|ends\s+with)\b`, Keyword, nil},
			{`(loop|block|parent)\b`, NameBuiltin, nil},
			{`(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w-]|[^\x00-\x7f])*`, NameVariable, nil},
			{`\.(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w-]|[^\x00-\x7f])*`, NameVariable, nil},
			{`\.[0-9]+`, LiteralNumber, nil},
			{`:?"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`:?'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`([{}()\[\]+\-*/,:~%]|\.\.|\?|:|\*\*|\/\/|!=|[><=]=?)`, Operator, nil},
			{`[0-9](\.[0-9]*)?(eE[+-][0-9])?[flFLdD]?|0[xX][0-9a-fA-F]+[Ll]?`, LiteralNumber, nil},
		},
		"var": {
			{`\s+`, Text, nil},
			{`(-?)(\}\})`, ByGroups(Text, CommentPreproc), Pop(1)},
			Include("varnames"),
		},
		"tag": {
			{`\s+`, Text, nil},
			{`(-?)(%\})`, ByGroups(Text, CommentPreproc), Pop(1)},
			Include("varnames"),
			{`.`, Punctuation, nil},
		},
	}
}
