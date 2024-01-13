package styles

import (
	"github.com/alecthomas/chroma"
)

// SolarizedDark256 style.
var SolarizedDark256 = Register(chroma.MustNewStyle("solarized-dark256", chroma.StyleEntries{
	chroma.Keyword:               "#5f8700",
	chroma.KeywordConstant:       "#d75f00",
	chroma.KeywordDeclaration:    "#0087ff",
	chroma.KeywordNamespace:      "#d75f00",
	chroma.KeywordReserved:       "#0087ff",
	chroma.KeywordType:           "#af0000",
	chroma.NameAttribute:         "#8a8a8a",
	chroma.NameBuiltin:           "#0087ff",
	chroma.NameBuiltinPseudo:     "#0087ff",
	chroma.NameClass:             "#0087ff",
	chroma.NameConstant:          "#d75f00",
	chroma.NameDecorator:         "#0087ff",
	chroma.NameEntity:            "#d75f00",
	chroma.NameException:         "#af8700",
	chroma.NameFunction:          "#0087ff",
	chroma.NameTag:               "#0087ff",
	chroma.NameVariable:          "#0087ff",
	chroma.LiteralString:         "#00afaf",
	chroma.LiteralStringBacktick: "#4e4e4e",
	chroma.LiteralStringChar:     "#00afaf",
	chroma.LiteralStringDoc:      "#00afaf",
	chroma.LiteralStringEscape:   "#af0000",
	chroma.LiteralStringHeredoc:  "#00afaf",
	chroma.LiteralStringRegex:    "#af0000",
	chroma.LiteralNumber:         "#00afaf",
	chroma.Operator:              "#8a8a8a",
	chroma.OperatorWord:          "#5f8700",
	chroma.Comment:               "#4e4e4e",
	chroma.CommentPreproc:        "#5f8700",
	chroma.CommentSpecial:        "#5f8700",
	chroma.GenericDeleted:        "#af0000",
	chroma.GenericEmph:           "italic",
	chroma.GenericError:          "#af0000 bold",
	chroma.GenericHeading:        "#d75f00",
	chroma.GenericInserted:       "#5f8700",
	chroma.GenericStrong:         "bold",
	chroma.GenericSubheading:     "#0087ff",
	chroma.Background:            "#8a8a8a bg:#1c1c1c",
	chroma.Other:                 "#d75f00",
}))
