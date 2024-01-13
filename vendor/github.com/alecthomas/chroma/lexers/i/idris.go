package i

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Idris lexer.
var Idris = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Idris",
		Aliases:   []string{"idris", "idr"},
		Filenames: []string{"*.idr"},
		MimeTypes: []string{"text/x-idris"},
	},
	idrisRules,
))

func idrisRules() Rules {
	return Rules{
		"root": {
			{`^(\s*)(%lib|link|flag|include|hide|freeze|access|default|logging|dynamic|name|error_handlers|language)`, ByGroups(Text, KeywordReserved), nil},
			{`(\s*)(--(?![!#$%&*+./<=>?@^|_~:\\]).*?)$`, ByGroups(Text, CommentSingle), nil},
			{`(\s*)(\|{3}.*?)$`, ByGroups(Text, CommentSingle), nil},
			{`(\s*)(\{-)`, ByGroups(Text, CommentMultiline), Push("comment")},
			{`^(\s*)([^\s(){}]+)(\s*)(:)(\s*)`, ByGroups(Text, NameFunction, Text, OperatorWord, Text), nil},
			{`\b(case|class|data|default|using|do|else|if|in|infix[lr]?|instance|rewrite|auto|namespace|codata|mutual|private|public|abstract|total|partial|let|proof|of|then|static|where|_|with|pattern|term|syntax|prefix|postulate|parameters|record|dsl|impossible|implicit|tactics|intros|intro|compute|refine|exact|trivial)(?!\')\b`, KeywordReserved, nil},
			{`(import|module)(\s+)`, ByGroups(KeywordReserved, Text), Push("module")},
			{`('')?[A-Z][\w\']*`, KeywordType, nil},
			{`[a-z][\w\']*`, Text, nil},
			{`(<-|::|->|=>|=)`, OperatorWord, nil},
			{`([(){}\[\]:!#$%&*+.\\/<=>?@^|~-]+)`, OperatorWord, nil},
			{`\d+[eE][+-]?\d+`, LiteralNumberFloat, nil},
			{`\d+\.\d+([eE][+-]?\d+)?`, LiteralNumberFloat, nil},
			{`0[xX][\da-fA-F]+`, LiteralNumberHex, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`'`, LiteralStringChar, Push("character")},
			{`"`, LiteralString, Push("string")},
			{`[^\s(){}]+`, Text, nil},
			{`\s+?`, Text, nil},
		},
		"module": {
			{`\s+`, Text, nil},
			{`([A-Z][\w.]*)(\s+)(\()`, ByGroups(NameNamespace, Text, Punctuation), Push("funclist")},
			{`[A-Z][\w.]*`, NameNamespace, Pop(1)},
		},
		"funclist": {
			{`\s+`, Text, nil},
			{`[A-Z]\w*`, KeywordType, nil},
			{`(_[\w\']+|[a-z][\w\']*)`, NameFunction, nil},
			{`--.*$`, CommentSingle, nil},
			{`\{-`, CommentMultiline, Push("comment")},
			{`,`, Punctuation, nil},
			{`[:!#$%&*+.\\/<=>?@^|~-]+`, Operator, nil},
			{`\(`, Punctuation, Push("funclist", "funclist")},
			{`\)`, Punctuation, Pop(2)},
		},
		"comment": {
			{`[^-{}]+`, CommentMultiline, nil},
			{`\{-`, CommentMultiline, Push()},
			{`-\}`, CommentMultiline, Pop(1)},
			{`[-{}]`, CommentMultiline, nil},
		},
		"character": {
			{`[^\\']`, LiteralStringChar, nil},
			{`\\`, LiteralStringEscape, Push("escape")},
			{`'`, LiteralStringChar, Pop(1)},
		},
		"string": {
			{`[^\\"]+`, LiteralString, nil},
			{`\\`, LiteralStringEscape, Push("escape")},
			{`"`, LiteralString, Pop(1)},
		},
		"escape": {
			{`[abfnrtv"\'&\\]`, LiteralStringEscape, Pop(1)},
			{`\^[][A-Z@^_]`, LiteralStringEscape, Pop(1)},
			{`NUL|SOH|[SE]TX|EOT|ENQ|ACK|BEL|BS|HT|LF|VT|FF|CR|S[OI]|DLE|DC[1-4]|NAK|SYN|ETB|CAN|EM|SUB|ESC|[FGRU]S|SP|DEL`, LiteralStringEscape, Pop(1)},
			{`o[0-7]+`, LiteralStringEscape, Pop(1)},
			{`x[\da-fA-F]+`, LiteralStringEscape, Pop(1)},
			{`\d+`, LiteralStringEscape, Pop(1)},
			{`\s+\\`, LiteralStringEscape, Pop(1)},
		},
	}
}
