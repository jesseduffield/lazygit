package styles

import (
	"github.com/alecthomas/chroma"
)

// Pygments default theme.
var Pygments = Register(chroma.MustNewStyle("pygments", chroma.StyleEntries{
	chroma.Whitespace:     "#bbbbbb",
	chroma.Comment:        "italic #408080",
	chroma.CommentPreproc: "noitalic #BC7A00",

	chroma.Keyword:       "bold #008000",
	chroma.KeywordPseudo: "nobold",
	chroma.KeywordType:   "nobold #B00040",

	chroma.Operator:     "#666666",
	chroma.OperatorWord: "bold #AA22FF",

	chroma.NameBuiltin:   "#008000",
	chroma.NameFunction:  "#0000FF",
	chroma.NameClass:     "bold #0000FF",
	chroma.NameNamespace: "bold #0000FF",
	chroma.NameException: "bold #D2413A",
	chroma.NameVariable:  "#19177C",
	chroma.NameConstant:  "#880000",
	chroma.NameLabel:     "#A0A000",
	chroma.NameEntity:    "bold #999999",
	chroma.NameAttribute: "#7D9029",
	chroma.NameTag:       "bold #008000",
	chroma.NameDecorator: "#AA22FF",

	chroma.String:         "#BA2121",
	chroma.StringDoc:      "italic",
	chroma.StringInterpol: "bold #BB6688",
	chroma.StringEscape:   "bold #BB6622",
	chroma.StringRegex:    "#BB6688",
	chroma.StringSymbol:   "#19177C",
	chroma.StringOther:    "#008000",
	chroma.Number:         "#666666",

	chroma.GenericHeading:    "bold #000080",
	chroma.GenericSubheading: "bold #800080",
	chroma.GenericDeleted:    "#A00000",
	chroma.GenericInserted:   "#00A000",
	chroma.GenericError:      "#FF0000",
	chroma.GenericEmph:       "italic",
	chroma.GenericStrong:     "bold",
	chroma.GenericPrompt:     "bold #000080",
	chroma.GenericOutput:     "#888",
	chroma.GenericTraceback:  "#04D",
	chroma.GenericUnderline:  "underline",

	chroma.Error: "border:#FF0000",
}))
