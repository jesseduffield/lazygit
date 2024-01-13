package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Actionscript 3 lexer.
var Actionscript3 = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ActionScript 3",
		Aliases:   []string{"as3", "actionscript3"},
		Filenames: []string{"*.as"},
		MimeTypes: []string{"application/x-actionscript3", "text/x-actionscript3", "text/actionscript3"},
		DotAll:    true,
	},
	actionscript3Rules,
))

func actionscript3Rules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`(function\s+)([$a-zA-Z_]\w*)(\s*)(\()`, ByGroups(KeywordDeclaration, NameFunction, Text, Operator), Push("funcparams")},
			{`(var|const)(\s+)([$a-zA-Z_]\w*)(\s*)(:)(\s*)([$a-zA-Z_]\w*(?:\.<\w+>)?)`, ByGroups(KeywordDeclaration, Text, Name, Text, Punctuation, Text, KeywordType), nil},
			{`(import|package)(\s+)((?:[$a-zA-Z_]\w*|\.)+)(\s*)`, ByGroups(Keyword, Text, NameNamespace, Text), nil},
			{`(new)(\s+)([$a-zA-Z_]\w*(?:\.<\w+>)?)(\s*)(\()`, ByGroups(Keyword, Text, KeywordType, Text, Operator), nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`/(\\\\|\\/|[^\n])*/[gisx]*`, LiteralStringRegex, nil},
			{`(\.)([$a-zA-Z_]\w*)`, ByGroups(Operator, NameAttribute), nil},
			{`(case|default|for|each|in|while|do|break|return|continue|if|else|throw|try|catch|with|new|typeof|arguments|instanceof|this|switch|import|include|as|is)\b`, Keyword, nil},
			{`(class|public|final|internal|native|override|private|protected|static|import|extends|implements|interface|intrinsic|return|super|dynamic|function|const|get|namespace|package|set)\b`, KeywordDeclaration, nil},
			{`(true|false|null|NaN|Infinity|-Infinity|undefined|void)\b`, KeywordConstant, nil},
			{`(decodeURI|decodeURIComponent|encodeURI|escape|eval|isFinite|isNaN|isXMLName|clearInterval|fscommand|getTimer|getURL|getVersion|isFinite|parseFloat|parseInt|setInterval|trace|updateAfterEvent|unescape)\b`, NameFunction, nil},
			{`[$a-zA-Z_]\w*`, Name, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[~^*!%&<>|+=:;,/?\\{}\[\]().-]+`, Operator, nil},
		},
		"funcparams": {
			{`\s+`, Text, nil},
			{`(\s*)(\.\.\.)?([$a-zA-Z_]\w*)(\s*)(:)(\s*)([$a-zA-Z_]\w*(?:\.<\w+>)?|\*)(\s*)`, ByGroups(Text, Punctuation, Name, Text, Operator, Text, KeywordType, Text), Push("defval")},
			{`\)`, Operator, Push("type")},
		},
		"type": {
			{`(\s*)(:)(\s*)([$a-zA-Z_]\w*(?:\.<\w+>)?|\*)`, ByGroups(Text, Operator, Text, KeywordType), Pop(2)},
			{`\s+`, Text, Pop(2)},
			Default(Pop(2)),
		},
		"defval": {
			{`(=)(\s*)([^(),]+)(\s*)(,?)`, ByGroups(Operator, Text, UsingSelf("root"), Text, Operator), Pop(1)},
			{`,`, Operator, Pop(1)},
			Default(Pop(1)),
		},
	}
}
