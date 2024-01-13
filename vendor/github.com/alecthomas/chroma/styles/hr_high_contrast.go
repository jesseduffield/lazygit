package styles

import (
	"github.com/alecthomas/chroma"
)

// Theme based on HackerRank High Contrast Editor Theme
var HrHighContrast = Register(chroma.MustNewStyle("hr_high_contrast", chroma.StyleEntries{
	chroma.Comment:              "#5a8349",
	chroma.Keyword:              "#467faf",
	chroma.OperatorWord:         "#467faf",
	chroma.Name:                 "#ffffff",
	chroma.LiteralString:        "#a87662",
	chroma.LiteralNumber:        "#fff",
	chroma.LiteralStringBoolean: "#467faf",
	chroma.Operator:             "#e4e400",
	chroma.Background:           "#000",
	chroma.Other:                "#d5d500",
}))
