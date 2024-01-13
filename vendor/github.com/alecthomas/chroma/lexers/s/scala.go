package s

import (
	"fmt"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Scala lexer.
var Scala = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Scala",
		Aliases:   []string{"scala"},
		Filenames: []string{"*.scala"},
		MimeTypes: []string{"text/x-scala"},
		DotAll:    true,
	},
	scalaRules,
))

func scalaRules() Rules {
	var (
		scalaOp     = "[-~\\^\\*!%&\\\\<>\\|+=:/?@\xa6-\xa7\xa9\xac\xae\xb0-\xb1\xb6\xd7\xf7\u03f6\u0482\u0606-\u0608\u060e-\u060f\u06e9\u06fd-\u06fe\u07f6\u09fa\u0b70\u0bf3-\u0bf8\u0bfa\u0c7f\u0cf1-\u0cf2\u0d79\u0f01-\u0f03\u0f13-\u0f17\u0f1a-\u0f1f\u0f34\u0f36\u0f38\u0fbe-\u0fc5\u0fc7-\u0fcf\u109e-\u109f\u1360\u1390-\u1399\u1940\u19e0-\u19ff\u1b61-\u1b6a\u1b74-\u1b7c\u2044\u2052\u207a-\u207c\u208a-\u208c\u2100-\u2101\u2103-\u2106\u2108-\u2109\u2114\u2116-\u2118\u211e-\u2123\u2125\u2127\u2129\u212e\u213a-\u213b\u2140-\u2144\u214a-\u214d\u214f\u2190-\u2328\u232b-\u244a\u249c-\u24e9\u2500-\u2767\u2794-\u27c4\u27c7-\u27e5\u27f0-\u2982\u2999-\u29d7\u29dc-\u29fb\u29fe-\u2b54\u2ce5-\u2cea\u2e80-\u2ffb\u3004\u3012-\u3013\u3020\u3036-\u3037\u303e-\u303f\u3190-\u3191\u3196-\u319f\u31c0-\u31e3\u3200-\u321e\u322a-\u3250\u3260-\u327f\u328a-\u32b0\u32c0-\u33ff\u4dc0-\u4dff\ua490-\ua4c6\ua828-\ua82b\ufb29\ufdfd\ufe62\ufe64-\ufe66\uff0b\uff1c-\uff1e\uff5c\uff5e\uffe2\uffe4\uffe8-\uffee\ufffc-\ufffd]+"
		scalaUpper  = `[\\$_\p{Lu}]`
		scalaLetter = `[\\$_\p{L}]`
		scalaIDRest = fmt.Sprintf(`%s(?:%s|[0-9])*(?:(?<=_)%s)?`, scalaLetter, scalaLetter, scalaOp)
	)

	return Rules{
		"root": {
			{`(class|trait|object)(\s+)`, ByGroups(Keyword, Text), Push("class")},
			{`[^\S\n]+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("comment")},
			{`@` + scalaIDRest, NameDecorator, nil},
			{`(abstract|ca(?:se|tch)|d(?:ef|o)|e(?:lse|xtends)|f(?:inal(?:ly)?|or(?:Some)?)|i(?:f|mplicit)|lazy|match|new|override|pr(?:ivate|otected)|re(?:quires|turn)|s(?:ealed|uper)|t(?:h(?:is|row)|ry)|va[lr]|w(?:hile|ith)|yield)\b|(<[%:-]|=>|>:|[#=@_⇒←])(\b|(?=\s)|$)`, Keyword, nil},
			{`:(?!` + scalaOp + `%s)`, Keyword, Push("type")},
			{fmt.Sprintf("%s%s\\b", scalaUpper, scalaIDRest), NameClass, nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(import|package)(\s+)`, ByGroups(Keyword, Text), Push("import")},
			{`(type)(\s+)`, ByGroups(Keyword, Text), Push("type")},
			{`""".*?"""(?!")`, LiteralString, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'\\.'|'[^\\]'|'\\u[0-9a-fA-F]{4}'`, LiteralStringChar, nil},
			{"'" + scalaIDRest, TextSymbol, nil},
			{`[fs]"""`, LiteralString, Push("interptriplestring")},
			{`[fs]"`, LiteralString, Push("interpstring")},
			{`raw"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{scalaIDRest, Name, nil},
			{"`[^`]+`", Name, nil},
			{`\[`, Operator, Push("typeparam")},
			{`[(){};,.#]`, Operator, nil},
			{scalaOp, Operator, nil},
			{`([0-9][0-9]*\.[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?[fFdD]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+L?`, LiteralNumberInteger, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{fmt.Sprintf("(%s|%s|`[^`]+`)(\\s*)(\\[)", scalaIDRest, scalaOp), ByGroups(NameClass, Text, Operator), Push("typeparam")},
			{`\s+`, Text, nil},
			{`\{`, Operator, Pop(1)},
			{`\(`, Operator, Pop(1)},
			{`//.*?\n`, CommentSingle, Pop(1)},
			{fmt.Sprintf("%s|%s|`[^`]+`", scalaIDRest, scalaOp), NameClass, Pop(1)},
		},
		"type": {
			{`\s+`, Text, nil},
			{`<[%:]|>:|[#_]|forSome|type`, Keyword, nil},
			{`([,);}]|=>|=|⇒)(\s*)`, ByGroups(Operator, Text), Pop(1)},
			{`[({]`, Operator, Push()},
			{fmt.Sprintf("((?:%s|%s|`[^`]+`)(?:\\.(?:%s|%s|`[^`]+`))*)(\\s*)(\\[)", scalaIDRest, scalaOp, scalaIDRest, scalaOp), ByGroups(KeywordType, Text, Operator), Push("#pop", "typeparam")},
			{fmt.Sprintf("((?:%s|%s|`[^`]+`)(?:\\.(?:%s|%s|`[^`]+`))*)(\\s*)$", scalaIDRest, scalaOp, scalaIDRest, scalaOp), ByGroups(KeywordType, Text), Pop(1)},
			{`//.*?\n`, CommentSingle, Pop(1)},
			{fmt.Sprintf("\\.|%s|%s|`[^`]+`", scalaIDRest, scalaOp), KeywordType, nil},
		},
		"typeparam": {
			{`[\s,]+`, Text, nil},
			{`<[%:]|=>|>:|[#_⇒]|forSome|type`, Keyword, nil},
			{`([\])}])`, Operator, Pop(1)},
			{`[(\[{]`, Operator, Push()},
			{fmt.Sprintf("\\.|%s|%s|`[^`]+`", scalaIDRest, scalaOp), KeywordType, nil},
		},
		"comment": {
			{`[^/*]+`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push()},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
		"import": {
			{fmt.Sprintf("(%s|\\.)+", scalaIDRest), NameNamespace, Pop(1)},
		},
		"interpstringcommon": {
			{`[^"$\\]+`, LiteralString, nil},
			{`\$\$`, LiteralString, nil},
			{`\$` + scalaLetter + `(?:` + scalaLetter + `|\d)*`, LiteralStringInterpol, nil},
			{`\$\{`, LiteralStringInterpol, Push("interpbrace")},
			{`\\.`, LiteralString, nil},
		},
		"interptriplestring": {
			{`"""(?!")`, LiteralString, Pop(1)},
			{`"`, LiteralString, nil},
			Include("interpstringcommon"),
		},
		"interpstring": {
			{`"`, LiteralString, Pop(1)},
			Include("interpstringcommon"),
		},
		"interpbrace": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			{`\{`, LiteralStringInterpol, Push()},
			Include("root"),
		},
	}
}
