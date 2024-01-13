package styles

import (
	"github.com/alecthomas/chroma"
)

// BlackWhite style.
var BlackWhite = Register(chroma.MustNewStyle("bw", chroma.StyleEntries{
	chroma.Comment:               "italic",
	chroma.CommentPreproc:        "noitalic",
	chroma.Keyword:               "bold",
	chroma.KeywordPseudo:         "nobold",
	chroma.KeywordType:           "nobold",
	chroma.OperatorWord:          "bold",
	chroma.NameClass:             "bold",
	chroma.NameNamespace:         "bold",
	chroma.NameException:         "bold",
	chroma.NameEntity:            "bold",
	chroma.NameTag:               "bold",
	chroma.LiteralString:         "italic",
	chroma.LiteralStringInterpol: "bold",
	chroma.LiteralStringEscape:   "bold",
	chroma.GenericHeading:        "bold",
	chroma.GenericSubheading:     "bold",
	chroma.GenericEmph:           "italic",
	chroma.GenericStrong:         "bold",
	chroma.GenericPrompt:         "bold",
	chroma.Error:                 "border:#FF0000",
	chroma.Background:            " bg:#ffffff",
}))
