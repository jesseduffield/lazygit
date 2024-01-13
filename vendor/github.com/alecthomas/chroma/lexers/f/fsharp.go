package f

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Fsharp lexer.
var Fsharp = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "FSharp",
		Aliases:   []string{"fsharp"},
		Filenames: []string{"*.fs", "*.fsi"},
		MimeTypes: []string{"text/x-fsharp"},
	},
	fsharpRules,
))

func fsharpRules() Rules {
	return Rules{
		"escape-sequence": {
			{`\\[\\"\'ntbrafv]`, LiteralStringEscape, nil},
			{`\\[0-9]{3}`, LiteralStringEscape, nil},
			{`\\u[0-9a-fA-F]{4}`, LiteralStringEscape, nil},
			{`\\U[0-9a-fA-F]{8}`, LiteralStringEscape, nil},
		},
		"root": {
			{`\s+`, Text, nil},
			{`\(\)|\[\]`, NameBuiltinPseudo, nil},
			{`\b(?<!\.)([A-Z][\w\']*)(?=\s*\.)`, NameNamespace, Push("dotted")},
			{`\b([A-Z][\w\']*)`, Name, nil},
			{`///.*?\n`, LiteralStringDoc, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`\(\*(?!\))`, Comment, Push("comment")},
			{`@"`, LiteralString, Push("lstring")},
			{`"""`, LiteralString, Push("tqs")},
			{`"`, LiteralString, Push("string")},
			{`\b(open|module)(\s+)([\w.]+)`, ByGroups(Keyword, Text, NameNamespace), nil},
			{`\b(let!?)(\s+)(\w+)`, ByGroups(Keyword, Text, NameVariable), nil},
			{`\b(type)(\s+)(\w+)`, ByGroups(Keyword, Text, NameClass), nil},
			{`\b(member|override)(\s+)(\w+)(\.)(\w+)`, ByGroups(Keyword, Text, Name, Punctuation, NameFunction), nil},
			{`\b(abstract|as|assert|base|begin|class|default|delegate|do!|do|done|downcast|downto|elif|else|end|exception|extern|false|finally|for|function|fun|global|if|inherit|inline|interface|internal|in|lazy|let!|let|match|member|module|mutable|namespace|new|null|of|open|override|private|public|rec|return!|return|select|static|struct|then|to|true|try|type|upcast|use!|use|val|void|when|while|with|yield!|yield|atomic|break|checked|component|const|constraint|constructor|continue|eager|event|external|fixed|functor|include|method|mixin|object|parallel|process|protected|pure|sealed|tailcall|trait|virtual|volatile)\b`, Keyword, nil},
			{"``([^`\\n\\r\\t]|`[^`\\n\\r\\t])+``", Name, nil},
			{"(!=|#|&&|&|\\(|\\)|\\*|\\+|,|-\\.|->|-|\\.\\.|\\.|::|:=|:>|:|;;|;|<-|<\\]|<|>\\]|>|\\?\\?|\\?|\\[<|\\[\\||\\[|\\]|_|`|\\{|\\|\\]|\\||\\}|~|<@@|<@|=|@>|@@>)", Operator, nil},
			{`([=<>@^|&+\*/$%-]|[!?~])?[!$%&*+\./:<=>?@^|~-]`, Operator, nil},
			{`\b(and|or|not)\b`, OperatorWord, nil},
			{`\b(sbyte|byte|char|nativeint|unativeint|float32|single|float|double|int8|uint8|int16|uint16|int32|uint32|int64|uint64|decimal|unit|bool|string|list|exn|obj|enum)\b`, KeywordType, nil},
			{`#[ \t]*(if|endif|else|line|nowarn|light|\d+)\b.*?\n`, CommentPreproc, nil},
			{`[^\W\d][\w']*`, Name, nil},
			{`\d[\d_]*[uU]?[yslLnQRZINGmM]?`, LiteralNumberInteger, nil},
			{`0[xX][\da-fA-F][\da-fA-F_]*[uU]?[yslLn]?[fF]?`, LiteralNumberHex, nil},
			{`0[oO][0-7][0-7_]*[uU]?[yslLn]?`, LiteralNumberOct, nil},
			{`0[bB][01][01_]*[uU]?[yslLn]?`, LiteralNumberBin, nil},
			{`-?\d[\d_]*(.[\d_]*)?([eE][+\-]?\d[\d_]*)[fFmM]?`, LiteralNumberFloat, nil},
			{`'(?:(\\[\\\"'ntbr ])|(\\[0-9]{3})|(\\x[0-9a-fA-F]{2}))'B?`, LiteralStringChar, nil},
			{`'.'`, LiteralStringChar, nil},
			{`'`, Keyword, nil},
			{`@?"`, LiteralStringDouble, Push("string")},
			{`[~?][a-z][\w\']*:`, NameVariable, nil},
		},
		"dotted": {
			{`\s+`, Text, nil},
			{`\.`, Punctuation, nil},
			{`[A-Z][\w\']*(?=\s*\.)`, NameNamespace, nil},
			{`[A-Z][\w\']*`, Name, Pop(1)},
			{`[a-z_][\w\']*`, Name, Pop(1)},
			Default(Pop(1)),
		},
		"comment": {
			{`[^(*)@"]+`, Comment, nil},
			{`\(\*`, Comment, Push()},
			{`\*\)`, Comment, Pop(1)},
			{`@"`, LiteralString, Push("lstring")},
			{`"""`, LiteralString, Push("tqs")},
			{`"`, LiteralString, Push("string")},
			{`[(*)@]`, Comment, nil},
		},
		"string": {
			{`[^\\"]+`, LiteralString, nil},
			Include("escape-sequence"),
			{`\\\n`, LiteralString, nil},
			{`\n`, LiteralString, nil},
			{`"B?`, LiteralString, Pop(1)},
		},
		"lstring": {
			{`[^"]+`, LiteralString, nil},
			{`\n`, LiteralString, nil},
			{`""`, LiteralString, nil},
			{`"B?`, LiteralString, Pop(1)},
		},
		"tqs": {
			{`[^"]+`, LiteralString, nil},
			{`\n`, LiteralString, nil},
			{`"""B?`, LiteralString, Pop(1)},
			{`"`, LiteralString, nil},
		},
	}
}
