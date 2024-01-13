package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Django/Jinja lexer.
var DjangoJinja = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Django/Jinja",
		Aliases:   []string{"django", "jinja"},
		Filenames: []string{},
		MimeTypes: []string{"application/x-django-templating", "application/x-jinja"},
		DotAll:    true,
	},
	djangoJinjaRules,
))

func djangoJinjaRules() Rules {
	return Rules{
		"root": {
			{`[^{]+`, Other, nil},
			{`\{\{`, CommentPreproc, Push("var")},
			{`\{[*#].*?[*#]\}`, Comment, nil},
			{`(\{%)(-?\s*)(comment)(\s*-?)(%\})(.*?)(\{%)(-?\s*)(endcomment)(\s*-?)(%\})`, ByGroups(CommentPreproc, Text, Keyword, Text, CommentPreproc, Comment, CommentPreproc, Text, Keyword, Text, CommentPreproc), nil},
			{`(\{%)(-?\s*)(raw)(\s*-?)(%\})(.*?)(\{%)(-?\s*)(endraw)(\s*-?)(%\})`, ByGroups(CommentPreproc, Text, Keyword, Text, CommentPreproc, Text, CommentPreproc, Text, Keyword, Text, CommentPreproc), nil},
			{`(\{%)(-?\s*)(filter)(\s+)([a-zA-Z_]\w*)`, ByGroups(CommentPreproc, Text, Keyword, Text, NameFunction), Push("block")},
			{`(\{%)(-?\s*)([a-zA-Z_]\w*)`, ByGroups(CommentPreproc, Text, Keyword), Push("block")},
			{`\{`, Other, nil},
		},
		"varnames": {
			{`(\|)(\s*)([a-zA-Z_]\w*)`, ByGroups(Operator, Text, NameFunction), nil},
			{`(is)(\s+)(not)?(\s+)?([a-zA-Z_]\w*)`, ByGroups(Keyword, Text, Keyword, Text, NameFunction), nil},
			{`(_|true|false|none|True|False|None)\b`, KeywordPseudo, nil},
			{`(in|as|reversed|recursive|not|and|or|is|if|else|import|with(?:(?:out)?\s*context)?|scoped|ignore\s+missing)\b`, Keyword, nil},
			{`(loop|block|super|forloop)\b`, NameBuiltin, nil},
			{`[a-zA-Z_][\w-]*`, NameVariable, nil},
			{`\.\w+`, NameVariable, nil},
			{`:?"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`:?'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`([{}()\[\]+\-*/,:~]|[><=]=?)`, Operator, nil},
			{`[0-9](\.[0-9]*)?(eE[+-][0-9])?[flFLdD]?|0[xX][0-9a-fA-F]+[Ll]?`, LiteralNumber, nil},
		},
		"var": {
			{`\s+`, Text, nil},
			{`(-?)(\}\})`, ByGroups(Text, CommentPreproc), Pop(1)},
			Include("varnames"),
		},
		"block": {
			{`\s+`, Text, nil},
			{`(-?)(%\})`, ByGroups(Text, CommentPreproc), Pop(1)},
			Include("varnames"),
			{`.`, Punctuation, nil},
		},
	}
}
