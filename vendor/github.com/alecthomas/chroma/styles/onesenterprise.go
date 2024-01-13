package styles

import (
	"github.com/alecthomas/chroma"
)

// 1S:Designer color palette
var OnesEnterprise = Register(chroma.MustNewStyle("onesenterprise", chroma.StyleEntries{
	chroma.Text:           "#000000",
	chroma.Comment:        "#008000",
	chroma.CommentPreproc: "#963200",
	chroma.Operator:       "#FF0000",
	chroma.Keyword:        "#FF0000",
	chroma.Punctuation:    "#FF0000",
	chroma.LiteralString:  "#000000",
	chroma.Name:           "#0000FF",
}))
