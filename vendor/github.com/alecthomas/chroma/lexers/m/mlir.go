package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// MLIR lexer.
var Mlir = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "MLIR",
		Aliases:   []string{"mlir"},
		Filenames: []string{"*.mlir"},
		MimeTypes: []string{"text/x-mlir"},
	},
	mlirRules,
))

func mlirRules() Rules {
	return Rules{
		"root": {
			Include("whitespace"),
			{`c?"[^"]*?"`, LiteralString, nil},
			{`\^([-a-zA-Z$._][\w\-$.0-9]*)\s*`, NameLabel, nil},
			{`([\w\d_$.]+)\s*=`, NameLabel, nil},
			Include("keyword"),
			{`->`, Punctuation, nil},
			{`@([\w_][\w\d_$.]*)`, NameFunction, nil},
			{`[%#][\w\d_$.]+`, NameVariable, nil},
			{`([1-9?][\d?]*\s*x)+`, LiteralNumber, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumber, nil},
			{`-?\d+(?:[.]\d+)?(?:[eE][-+]?\d+(?:[.]\d+)?)?`, LiteralNumber, nil},
			{`[=<>{}\[\]()*.,!:]|x\b`, Punctuation, nil},
			{`[\w\d]+`, Text, nil},
		},
		"whitespace": {
			{`(\n|\s)+`, Text, nil},
			{`//.*?\n`, Comment, nil},
		},
		"keyword": {
			{Words(``, ``, `constant`, `return`), KeywordType, nil},
			{Words(``, ``, `func`, `loc`, `memref`, `tensor`, `vector`), KeywordType, nil},
			{`bf16|f16|f32|f64|index`, Keyword, nil},
			{`i[1-9]\d*`, Keyword, nil},
		},
	}
}
