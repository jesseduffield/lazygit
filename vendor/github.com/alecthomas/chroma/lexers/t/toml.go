package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var TOML = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "TOML",
		Aliases:   []string{"toml"},
		Filenames: []string{"*.toml"},
		MimeTypes: []string{"text/x-toml"},
	},
	tomlRules,
))

func tomlRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`#.*`, Comment, nil},
			{Words(``, `\b`, `true`, `false`), KeywordConstant, nil},
			{`\d\d\d\d-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d\+)?(Z|[+-]\d{2}:\d{2})`, LiteralDate, nil},
			{`[+-]?[0-9](_?\d)*\.\d+`, LiteralNumberFloat, nil},
			{`[+-]?[0-9](_?\d)*`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, StringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, StringSingle, nil},
			{`[.,=\[\]{}]`, Punctuation, nil},
			{`[A-Za-z0-9_-]+`, NameOther, nil},
		},
	}
}
