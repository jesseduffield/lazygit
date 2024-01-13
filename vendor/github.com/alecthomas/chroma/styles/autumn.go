package styles

import (
	"github.com/alecthomas/chroma"
)

// Autumn style.
var Autumn = Register(chroma.MustNewStyle("autumn", chroma.StyleEntries{
	chroma.TextWhitespace:      "#bbbbbb",
	chroma.Comment:             "italic #aaaaaa",
	chroma.CommentPreproc:      "noitalic #4c8317",
	chroma.CommentSpecial:      "italic #0000aa",
	chroma.Keyword:             "#0000aa",
	chroma.KeywordType:         "#00aaaa",
	chroma.OperatorWord:        "#0000aa",
	chroma.NameBuiltin:         "#00aaaa",
	chroma.NameFunction:        "#00aa00",
	chroma.NameClass:           "underline #00aa00",
	chroma.NameNamespace:       "underline #00aaaa",
	chroma.NameVariable:        "#aa0000",
	chroma.NameConstant:        "#aa0000",
	chroma.NameEntity:          "bold #800",
	chroma.NameAttribute:       "#1e90ff",
	chroma.NameTag:             "bold #1e90ff",
	chroma.NameDecorator:       "#888888",
	chroma.LiteralString:       "#aa5500",
	chroma.LiteralStringSymbol: "#0000aa",
	chroma.LiteralStringRegex:  "#009999",
	chroma.LiteralNumber:       "#009999",
	chroma.GenericHeading:      "bold #000080",
	chroma.GenericSubheading:   "bold #800080",
	chroma.GenericDeleted:      "#aa0000",
	chroma.GenericInserted:     "#00aa00",
	chroma.GenericError:        "#aa0000",
	chroma.GenericEmph:         "italic",
	chroma.GenericStrong:       "bold",
	chroma.GenericPrompt:       "#555555",
	chroma.GenericOutput:       "#888888",
	chroma.GenericTraceback:    "#aa0000",
	chroma.GenericUnderline:    "underline",
	chroma.Error:               "#F00 bg:#FAA",
	chroma.Background:          " bg:#ffffff",
}))
