package styles

import (
	"github.com/alecthomas/chroma"
)

// Rrt style.
var Rrt = Register(chroma.MustNewStyle("rrt", chroma.StyleEntries{
	chroma.CommentPreproc:      "#e5e5e5",
	chroma.Comment:             "#00ff00",
	chroma.KeywordType:         "#ee82ee",
	chroma.Keyword:             "#ff0000",
	chroma.LiteralNumber:       "#ff6600",
	chroma.LiteralStringSymbol: "#ff6600",
	chroma.LiteralString:       "#87ceeb",
	chroma.NameFunction:        "#ffff00",
	chroma.NameConstant:        "#7fffd4",
	chroma.NameVariable:        "#eedd82",
	chroma.Background:          "#f8f8f2 bg:#000000",
}))
