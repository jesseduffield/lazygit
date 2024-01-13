package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	. "github.com/alecthomas/chroma/lexers/p"
)

// Genshi Text lexer.
var GenshiText = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Genshi Text",
		Aliases:   []string{"genshitext"},
		Filenames: []string{},
		MimeTypes: []string{"application/x-genshi-text", "text/x-genshi"},
	},
	genshiTextRules,
))

func genshiTextRules() Rules {
	return Rules{
		"root": {
			{`[^#$\s]+`, Other, nil},
			{`^(\s*)(##.*)$`, ByGroups(Text, Comment), nil},
			{`^(\s*)(#)`, ByGroups(Text, CommentPreproc), Push("directive")},
			Include("variable"),
			{`[#$\s]`, Other, nil},
		},
		"directive": {
			{`\n`, Text, Pop(1)},
			{`(?:def|for|if)\s+.*`, Using(Python), Pop(1)},
			{`(choose|when|with)([^\S\n]+)(.*)`, ByGroups(Keyword, Text, Using(Python)), Pop(1)},
			{`(choose|otherwise)\b`, Keyword, Pop(1)},
			{`(end\w*)([^\S\n]*)(.*)`, ByGroups(Keyword, Text, Comment), Pop(1)},
		},
		"variable": {
			{`(?<!\$)(\$\{)(.+?)(\})`, ByGroups(CommentPreproc, Using(Python), CommentPreproc), nil},
			{`(?<!\$)(\$)([a-zA-Z_][\w.]*)`, NameVariable, nil},
		},
	}
}

// Html+Genshi lexer.
var GenshiHTMLTemplate = internal.Register(MustNewLazyLexer(
	&Config{
		Name:         "Genshi HTML",
		Aliases:      []string{"html+genshi", "html+kid"},
		Filenames:    []string{},
		MimeTypes:    []string{"text/html+genshi"},
		NotMultiline: true,
		DotAll:       true,
	},
	genshiMarkupRules,
))

// Genshi lexer.
var Genshi = internal.Register(MustNewLazyLexer(
	&Config{
		Name:         "Genshi",
		Aliases:      []string{"genshi", "kid", "xml+genshi", "xml+kid"},
		Filenames:    []string{"*.kid"},
		MimeTypes:    []string{"application/x-genshi", "application/x-kid"},
		NotMultiline: true,
		DotAll:       true,
	},
	genshiMarkupRules,
))

func genshiMarkupRules() Rules {
	return Rules{
		"root": {
			{`[^<$]+`, Other, nil},
			{`(<\?python)(.*?)(\?>)`, ByGroups(CommentPreproc, Using(Python), CommentPreproc), nil},
			{`<\s*(script|style)\s*.*?>.*?<\s*/\1\s*>`, Other, nil},
			{`<\s*py:[a-zA-Z0-9]+`, NameTag, Push("pytag")},
			{`<\s*[a-zA-Z0-9:.]+`, NameTag, Push("tag")},
			Include("variable"),
			{`[<$]`, Other, nil},
		},
		"pytag": {
			{`\s+`, Text, nil},
			{`[\w:-]+\s*=`, NameAttribute, Push("pyattr")},
			{`/?\s*>`, NameTag, Pop(1)},
		},
		"pyattr": {
			{`(")(.*?)(")`, ByGroups(LiteralString, Using(Python), LiteralString), Pop(1)},
			{`(')(.*?)(')`, ByGroups(LiteralString, Using(Python), LiteralString), Pop(1)},
			{`[^\s>]+`, LiteralString, Pop(1)},
		},
		"tag": {
			{`\s+`, Text, nil},
			{`py:[\w-]+\s*=`, NameAttribute, Push("pyattr")},
			{`[\w:-]+\s*=`, NameAttribute, Push("attr")},
			{`/?\s*>`, NameTag, Pop(1)},
		},
		"attr": {
			{`"`, LiteralString, Push("attr-dstring")},
			{`'`, LiteralString, Push("attr-sstring")},
			{`[^\s>]*`, LiteralString, Pop(1)},
		},
		"attr-dstring": {
			{`"`, LiteralString, Pop(1)},
			Include("strings"),
			{`'`, LiteralString, nil},
		},
		"attr-sstring": {
			{`'`, LiteralString, Pop(1)},
			Include("strings"),
			{`'`, LiteralString, nil},
		},
		"strings": {
			{`[^"'$]+`, LiteralString, nil},
			Include("variable"),
		},
		"variable": {
			{`(?<!\$)(\$\{)(.+?)(\})`, ByGroups(CommentPreproc, Using(Python), CommentPreproc), nil},
			{`(?<!\$)(\$)([a-zA-Z_][\w\.]*)`, NameVariable, nil},
		},
	}
}
