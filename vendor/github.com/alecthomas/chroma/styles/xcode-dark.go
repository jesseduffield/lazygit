package styles

import (
	"github.com/alecthomas/chroma"
)

var (
	// Inspired by Apple's Xcode "Default (Dark)" Theme
	background                  = "#1F1F24"
	plainText                   = "#FFFFFF"
	comments                    = "#6C7986"
	strings                     = "#FC6A5D"
	numbers                     = "#D0BF69"
	keywords                    = "#FC5FA3"
	preprocessorStatements      = "#FD8F3F"
	typeDeclarations            = "#5DD8FF"
	otherDeclarations           = "#41A1C0"
	otherFunctionAndMethodNames = "#A167E6"
	otherTypeNames              = "#D0A8FF"
)

// Xcode dark style
var XcodeDark = Register(chroma.MustNewStyle("xcode-dark", chroma.StyleEntries{
	chroma.Background: plainText + " bg:" + background,

	chroma.Comment:          comments,
	chroma.CommentMultiline: comments,
	chroma.CommentPreproc:   preprocessorStatements,
	chroma.CommentSingle:    comments,
	chroma.CommentSpecial:   comments + " italic",

	chroma.Error: "#960050",

	chroma.Keyword:            keywords,
	chroma.KeywordConstant:    keywords,
	chroma.KeywordDeclaration: keywords,
	chroma.KeywordReserved:    keywords,

	chroma.LiteralNumber:        numbers,
	chroma.LiteralNumberBin:     numbers,
	chroma.LiteralNumberFloat:   numbers,
	chroma.LiteralNumberHex:     numbers,
	chroma.LiteralNumberInteger: numbers,
	chroma.LiteralNumberOct:     numbers,

	chroma.LiteralString:         strings,
	chroma.LiteralStringEscape:   strings,
	chroma.LiteralStringInterpol: plainText,

	chroma.Name:              plainText,
	chroma.NameBuiltin:       otherTypeNames,
	chroma.NameBuiltinPseudo: otherFunctionAndMethodNames,
	chroma.NameClass:         typeDeclarations,
	chroma.NameFunction:      otherDeclarations,
	chroma.NameVariable:      otherDeclarations,

	chroma.Operator: plainText,

	chroma.Punctuation: plainText,

	chroma.Text: plainText,
}))
