package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Abnf lexer.
var Abnf = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "ABNF",
		Aliases:   []string{"abnf"},
		Filenames: []string{"*.abnf"},
		MimeTypes: []string{"text/x-abnf"},
	},
	abnfRules,
))

func abnfRules() Rules {
	return Rules{
		"root": {
			{`;.*$`, CommentSingle, nil},
			{`(%[si])?"[^"]*"`, Literal, nil},
			{`%b[01]+\-[01]+\b`, Literal, nil},
			{`%b[01]+(\.[01]+)*\b`, Literal, nil},
			{`%d[0-9]+\-[0-9]+\b`, Literal, nil},
			{`%d[0-9]+(\.[0-9]+)*\b`, Literal, nil},
			{`%x[0-9a-fA-F]+\-[0-9a-fA-F]+\b`, Literal, nil},
			{`%x[0-9a-fA-F]+(\.[0-9a-fA-F]+)*\b`, Literal, nil},
			{`\b[0-9]+\*[0-9]+`, Operator, nil},
			{`\b[0-9]+\*`, Operator, nil},
			{`\b[0-9]+`, Operator, nil},
			{`\*`, Operator, nil},
			{Words(``, `\b`, `ALPHA`, `BIT`, `CHAR`, `CR`, `CRLF`, `CTL`, `DIGIT`, `DQUOTE`, `HEXDIG`, `HTAB`, `LF`, `LWSP`, `OCTET`, `SP`, `VCHAR`, `WSP`), Keyword, nil},
			{`[a-zA-Z][a-zA-Z0-9-]+\b`, NameClass, nil},
			{`(=/|=|/)`, Operator, nil},
			{`[\[\]()]`, Punctuation, nil},
			{`\s+`, Text, nil},
			{`.`, Text, nil},
		},
	}
}
