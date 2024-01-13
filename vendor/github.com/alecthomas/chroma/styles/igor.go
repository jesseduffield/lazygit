package styles

import (
	"github.com/alecthomas/chroma"
)

// Igor style.
var Igor = Register(chroma.MustNewStyle("igor", chroma.StyleEntries{
	chroma.Comment:       "italic #FF0000",
	chroma.Keyword:       "#0000FF",
	chroma.NameFunction:  "#C34E00",
	chroma.NameDecorator: "#CC00A3",
	chroma.NameClass:     "#007575",
	chroma.LiteralString: "#009C00",
	chroma.Background:    " bg:#ffffff",
}))
