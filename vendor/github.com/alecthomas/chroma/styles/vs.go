package styles

import (
	"github.com/alecthomas/chroma"
)

// VisualStudio style.
var VisualStudio = Register(chroma.MustNewStyle("vs", chroma.StyleEntries{
	chroma.Comment:           "#008000",
	chroma.CommentPreproc:    "#0000ff",
	chroma.Keyword:           "#0000ff",
	chroma.OperatorWord:      "#0000ff",
	chroma.KeywordType:       "#2b91af",
	chroma.NameClass:         "#2b91af",
	chroma.LiteralString:     "#a31515",
	chroma.GenericHeading:    "bold",
	chroma.GenericSubheading: "bold",
	chroma.GenericEmph:       "italic",
	chroma.GenericStrong:     "bold",
	chroma.GenericPrompt:     "bold",
	chroma.Error:             "border:#FF0000",
	chroma.Background:        " bg:#ffffff",
}))
