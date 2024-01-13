package styles

import (
	"github.com/alecthomas/chroma"
)

// RainbowDash style.
var RainbowDash = Register(chroma.MustNewStyle("rainbow_dash", chroma.StyleEntries{
	chroma.Comment:             "italic #0080ff",
	chroma.CommentPreproc:      "noitalic",
	chroma.CommentSpecial:      "bold",
	chroma.Error:               "bg:#cc0000 #ffffff",
	chroma.GenericDeleted:      "border:#c5060b bg:#ffcccc",
	chroma.GenericEmph:         "italic",
	chroma.GenericError:        "#ff0000",
	chroma.GenericHeading:      "bold #2c5dcd",
	chroma.GenericInserted:     "border:#00cc00 bg:#ccffcc",
	chroma.GenericOutput:       "#aaaaaa",
	chroma.GenericPrompt:       "bold #2c5dcd",
	chroma.GenericStrong:       "bold",
	chroma.GenericSubheading:   "bold #2c5dcd",
	chroma.GenericTraceback:    "#c5060b",
	chroma.GenericUnderline:    "underline",
	chroma.Keyword:             "bold #2c5dcd",
	chroma.KeywordPseudo:       "nobold",
	chroma.KeywordType:         "#5918bb",
	chroma.NameAttribute:       "italic #2c5dcd",
	chroma.NameBuiltin:         "bold #5918bb",
	chroma.NameClass:           "underline",
	chroma.NameConstant:        "#318495",
	chroma.NameDecorator:       "bold #ff8000",
	chroma.NameEntity:          "bold #5918bb",
	chroma.NameException:       "bold #5918bb",
	chroma.NameFunction:        "bold #ff8000",
	chroma.NameTag:             "bold #2c5dcd",
	chroma.LiteralNumber:       "bold #5918bb",
	chroma.Operator:            "#2c5dcd",
	chroma.OperatorWord:        "bold",
	chroma.LiteralString:       "#00cc66",
	chroma.LiteralStringDoc:    "italic",
	chroma.LiteralStringEscape: "bold #c5060b",
	chroma.LiteralStringOther:  "#318495",
	chroma.LiteralStringSymbol: "bold #c5060b",
	chroma.Text:                "#4d4d4d",
	chroma.TextWhitespace:      "#cbcbcb",
	chroma.Background:          " bg:#ffffff",
}))
