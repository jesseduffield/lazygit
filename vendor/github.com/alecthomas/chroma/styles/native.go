package styles

import (
	"github.com/alecthomas/chroma"
)

// Native style.
var Native = Register(chroma.MustNewStyle("native", chroma.StyleEntries{
	chroma.Background:         "#d0d0d0 bg:#202020",
	chroma.TextWhitespace:     "#666666",
	chroma.Comment:            "italic #999999",
	chroma.CommentPreproc:     "noitalic bold #cd2828",
	chroma.CommentSpecial:     "noitalic bold #e50808 bg:#520000",
	chroma.Keyword:            "bold #6ab825",
	chroma.KeywordPseudo:      "nobold",
	chroma.OperatorWord:       "bold #6ab825",
	chroma.LiteralString:      "#ed9d13",
	chroma.LiteralStringOther: "#ffa500",
	chroma.LiteralNumber:      "#3677a9",
	chroma.NameBuiltin:        "#24909d",
	chroma.NameVariable:       "#40ffff",
	chroma.NameConstant:       "#40ffff",
	chroma.NameClass:          "underline #447fcf",
	chroma.NameFunction:       "#447fcf",
	chroma.NameNamespace:      "underline #447fcf",
	chroma.NameException:      "#bbbbbb",
	chroma.NameTag:            "bold #6ab825",
	chroma.NameAttribute:      "#bbbbbb",
	chroma.NameDecorator:      "#ffa500",
	chroma.GenericHeading:     "bold #ffffff",
	chroma.GenericSubheading:  "underline #ffffff",
	chroma.GenericDeleted:     "#d22323",
	chroma.GenericInserted:    "#589819",
	chroma.GenericError:       "#d22323",
	chroma.GenericEmph:        "italic",
	chroma.GenericStrong:      "bold",
	chroma.GenericPrompt:      "#aaaaaa",
	chroma.GenericOutput:      "#cccccc",
	chroma.GenericTraceback:   "#d22323",
	chroma.GenericUnderline:   "underline",
	chroma.Error:              "bg:#e3d2d2 #a61717",
}))
