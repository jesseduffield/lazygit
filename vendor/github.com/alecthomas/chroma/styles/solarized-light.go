package styles

import (
	"github.com/alecthomas/chroma"
)

// SolarizedLight style.
var SolarizedLight = Register(chroma.MustNewStyle("solarized-light", chroma.StyleEntries{
	chroma.Text:             "bg: #eee8d5 #586e75",
	chroma.Keyword:          "#859900",
	chroma.KeywordConstant:  "bold",
	chroma.KeywordNamespace: "#dc322f bold",
	chroma.KeywordType:      "bold",
	chroma.Name:             "#268bd2",
	chroma.NameBuiltin:      "#cb4b16",
	chroma.NameClass:        "#cb4b16",
	chroma.NameTag:          "bold",
	chroma.Literal:          "#2aa198",
	chroma.LiteralNumber:    "bold",
	chroma.OperatorWord:     "#859900",
	chroma.Comment:          "#93a1a1 italic",
	chroma.Generic:          "#d33682",
	chroma.Background:       " bg:#eee8d5",
}))
