package styles

import (
	"github.com/alecthomas/chroma"
)

// Vim style.
var Vim = Register(chroma.MustNewStyle("vim", chroma.StyleEntries{
	chroma.Background:         "#cccccc bg:#000000",
	chroma.Comment:            "#000080",
	chroma.CommentSpecial:     "bold #cd0000",
	chroma.Keyword:            "#cdcd00",
	chroma.KeywordDeclaration: "#00cd00",
	chroma.KeywordNamespace:   "#cd00cd",
	chroma.KeywordType:        "#00cd00",
	chroma.Operator:           "#3399cc",
	chroma.OperatorWord:       "#cdcd00",
	chroma.NameClass:          "#00cdcd",
	chroma.NameBuiltin:        "#cd00cd",
	chroma.NameException:      "bold #666699",
	chroma.NameVariable:       "#00cdcd",
	chroma.LiteralString:      "#cd0000",
	chroma.LiteralNumber:      "#cd00cd",
	chroma.GenericHeading:     "bold #000080",
	chroma.GenericSubheading:  "bold #800080",
	chroma.GenericDeleted:     "#cd0000",
	chroma.GenericInserted:    "#00cd00",
	chroma.GenericError:       "#FF0000",
	chroma.GenericEmph:        "italic",
	chroma.GenericStrong:      "bold",
	chroma.GenericPrompt:      "bold #000080",
	chroma.GenericOutput:      "#888",
	chroma.GenericTraceback:   "#04D",
	chroma.GenericUnderline:   "underline",
	chroma.Error:              "border:#FF0000",
}))
