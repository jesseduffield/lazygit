package styles

import (
	"github.com/alecthomas/chroma"
)

// Arduino style.
var Arduino = Register(chroma.MustNewStyle("arduino", chroma.StyleEntries{
	chroma.Error:           "#a61717",
	chroma.Comment:         "#95a5a6",
	chroma.CommentPreproc:  "#728E00",
	chroma.Keyword:         "#728E00",
	chroma.KeywordConstant: "#00979D",
	chroma.KeywordPseudo:   "#00979D",
	chroma.KeywordReserved: "#00979D",
	chroma.KeywordType:     "#00979D",
	chroma.Operator:        "#728E00",
	chroma.Name:            "#434f54",
	chroma.NameBuiltin:     "#728E00",
	chroma.NameFunction:    "#D35400",
	chroma.NameOther:       "#728E00",
	chroma.LiteralNumber:   "#8A7B52",
	chroma.LiteralString:   "#7F8C8D",
	chroma.Background:      " bg:#ffffff",
}))
