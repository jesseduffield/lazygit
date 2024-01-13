package styles

import (
	"github.com/alecthomas/chroma"
)

// Theme based on HackerRank Dark Editor theme
var HrDark = Register(chroma.MustNewStyle("hrdark", chroma.StyleEntries{
	chroma.Comment:      "italic #828b96",
	chroma.Keyword:      "#ff636f",
	chroma.OperatorWord: "#ff636f",
	chroma.Name:         "#58a1dd",
	chroma.Literal:      "#a6be9d",
	chroma.Operator:     "#ff636f",
	chroma.Background:   "#1d2432",
	chroma.Other:        "#fff",
}))
