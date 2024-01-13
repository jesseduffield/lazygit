package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Coffeescript lexer.
var Coffeescript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:         "CoffeeScript",
		Aliases:      []string{"coffee-script", "coffeescript", "coffee"},
		Filenames:    []string{"*.coffee"},
		MimeTypes:    []string{"text/coffeescript"},
		NotMultiline: true,
		DotAll:       true,
	},
	coffeescriptRules,
))

func coffeescriptRules() Rules {
	return Rules{
		"commentsandwhitespace": {
			{`\s+`, Text, nil},
			{`###[^#].*?###`, CommentMultiline, nil},
			{`#(?!##[^#]).*?\n`, CommentSingle, nil},
		},
		"multilineregex": {
			{`[^/#]+`, LiteralStringRegex, nil},
			{`///([gim]+\b|\B)`, LiteralStringRegex, Pop(1)},
			{`#\{`, LiteralStringInterpol, Push("interpoling_string")},
			{`[/#]`, LiteralStringRegex, nil},
		},
		"slashstartsregex": {
			Include("commentsandwhitespace"),
			{`///`, LiteralStringRegex, Push("#pop", "multilineregex")},
			{`/(?! )(\\.|[^[/\\\n]|\[(\\.|[^\]\\\n])*])+/([gim]+\b|\B)`, LiteralStringRegex, Pop(1)},
			{`/`, Operator, nil},
			Default(Pop(1)),
		},
		"root": {
			Include("commentsandwhitespace"),
			{`^(?=\s|/)`, Text, Push("slashstartsregex")},
			{"\\+\\+|~|&&|\\band\\b|\\bor\\b|\\bis\\b|\\bisnt\\b|\\bnot\\b|\\?|:|\\|\\||\\\\(?=\\n)|(<<|>>>?|==?(?!>)|!=?|=(?!>)|-(?!>)|[<>+*`%&\\|\\^/])=?", Operator, Push("slashstartsregex")},
			{`(?:\([^()]*\))?\s*[=-]>`, NameFunction, Push("slashstartsregex")},
			{`[{(\[;,]`, Punctuation, Push("slashstartsregex")},
			{`[})\].]`, Punctuation, nil},
			{`(?<![.$])(for|own|in|of|while|until|loop|break|return|continue|switch|when|then|if|unless|else|throw|try|catch|finally|new|delete|typeof|instanceof|super|extends|this|class|by)\b`, Keyword, Push("slashstartsregex")},
			{`(?<![.$])(true|false|yes|no|on|off|null|NaN|Infinity|undefined)\b`, KeywordConstant, nil},
			{`(Array|Boolean|Date|Error|Function|Math|netscape|Number|Object|Packages|RegExp|String|sun|decodeURI|decodeURIComponent|encodeURI|encodeURIComponent|eval|isFinite|isNaN|parseFloat|parseInt|document|window)\b`, NameBuiltin, nil},
			{`[$a-zA-Z_][\w.:$]*\s*[:=]\s`, NameVariable, Push("slashstartsregex")},
			{`@[$a-zA-Z_][\w.:$]*\s*[:=]\s`, NameVariableInstance, Push("slashstartsregex")},
			{`@`, NameOther, Push("slashstartsregex")},
			{`@?[$a-zA-Z_][\w$]*`, NameOther, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"""`, LiteralString, Push("tdqs")},
			{`'''`, LiteralString, Push("tsqs")},
			{`"`, LiteralString, Push("dqs")},
			{`'`, LiteralString, Push("sqs")},
		},
		"strings": {
			{`[^#\\\'"]+`, LiteralString, nil},
		},
		"interpoling_string": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			Include("root"),
		},
		"dqs": {
			{`"`, LiteralString, Pop(1)},
			{`\\.|\'`, LiteralString, nil},
			{`#\{`, LiteralStringInterpol, Push("interpoling_string")},
			{`#`, LiteralString, nil},
			Include("strings"),
		},
		"sqs": {
			{`'`, LiteralString, Pop(1)},
			{`#|\\.|"`, LiteralString, nil},
			Include("strings"),
		},
		"tdqs": {
			{`"""`, LiteralString, Pop(1)},
			{`\\.|\'|"`, LiteralString, nil},
			{`#\{`, LiteralStringInterpol, Push("interpoling_string")},
			{`#`, LiteralString, nil},
			Include("strings"),
		},
		"tsqs": {
			{`'''`, LiteralString, Pop(1)},
			{`#|\\.|\'|"`, LiteralString, nil},
			Include("strings"),
		},
	}
}
