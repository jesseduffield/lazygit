package styles

import (
	"github.com/alecthomas/chroma"
)

// Perldoc style.
var Perldoc = Register(chroma.MustNewStyle("perldoc", chroma.StyleEntries{
	chroma.TextWhitespace:       "#bbbbbb",
	chroma.Comment:              "#228B22",
	chroma.CommentPreproc:       "#1e889b",
	chroma.CommentSpecial:       "#8B008B bold",
	chroma.LiteralString:        "#CD5555",
	chroma.LiteralStringHeredoc: "#1c7e71 italic",
	chroma.LiteralStringRegex:   "#1c7e71",
	chroma.LiteralStringOther:   "#cb6c20",
	chroma.LiteralNumber:        "#B452CD",
	chroma.OperatorWord:         "#8B008B",
	chroma.Keyword:              "#8B008B bold",
	chroma.KeywordType:          "#00688B",
	chroma.NameClass:            "#008b45 bold",
	chroma.NameException:        "#008b45 bold",
	chroma.NameFunction:         "#008b45",
	chroma.NameNamespace:        "#008b45 underline",
	chroma.NameVariable:         "#00688B",
	chroma.NameConstant:         "#00688B",
	chroma.NameDecorator:        "#707a7c",
	chroma.NameTag:              "#8B008B bold",
	chroma.NameAttribute:        "#658b00",
	chroma.NameBuiltin:          "#658b00",
	chroma.GenericHeading:       "bold #000080",
	chroma.GenericSubheading:    "bold #800080",
	chroma.GenericDeleted:       "#aa0000",
	chroma.GenericInserted:      "#00aa00",
	chroma.GenericError:         "#aa0000",
	chroma.GenericEmph:          "italic",
	chroma.GenericStrong:        "bold",
	chroma.GenericPrompt:        "#555555",
	chroma.GenericOutput:        "#888888",
	chroma.GenericTraceback:     "#aa0000",
	chroma.GenericUnderline:     "underline",
	chroma.Error:                "bg:#e3d2d2 #a61717",
	chroma.Background:           " bg:#eeeedd",
}))
