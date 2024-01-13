package styles

import (
	"github.com/alecthomas/chroma"
)

// SwapOff theme.
var SwapOff = Register(chroma.MustNewStyle("swapoff", chroma.StyleEntries{
	chroma.Background:        "#lightgray bg:#black",
	chroma.Number:            "bold #ansiyellow",
	chroma.Comment:           "#ansiteal",
	chroma.CommentPreproc:    "bold #ansigreen",
	chroma.String:            "bold #ansiturquoise",
	chroma.Keyword:           "bold #ansiwhite",
	chroma.NameKeyword:       "bold #ansiwhite",
	chroma.NameBuiltin:       "bold #ansiwhite",
	chroma.GenericHeading:    "bold",
	chroma.GenericSubheading: "bold",
	chroma.GenericStrong:     "bold",
	chroma.GenericUnderline:  "underline",
	chroma.NameTag:           "bold",
	chroma.NameAttribute:     "#ansiteal",
	chroma.Error:             "#ansired",
	chroma.LiteralDate:       "bold #ansiyellow",
}))
