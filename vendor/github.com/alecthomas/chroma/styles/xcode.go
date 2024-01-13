package styles

import (
	"github.com/alecthomas/chroma"
)

// Xcode style.
var Xcode = Register(chroma.MustNewStyle("xcode", chroma.StyleEntries{
	chroma.Comment:           "#177500",
	chroma.CommentPreproc:    "#633820",
	chroma.LiteralString:     "#C41A16",
	chroma.LiteralStringChar: "#2300CE",
	chroma.Operator:          "#000000",
	chroma.Keyword:           "#A90D91",
	chroma.Name:              "#000000",
	chroma.NameAttribute:     "#836C28",
	chroma.NameClass:         "#3F6E75",
	chroma.NameFunction:      "#000000",
	chroma.NameBuiltin:       "#A90D91",
	chroma.NameBuiltinPseudo: "#5B269A",
	chroma.NameVariable:      "#000000",
	chroma.NameTag:           "#000000",
	chroma.NameDecorator:     "#000000",
	chroma.NameLabel:         "#000000",
	chroma.Literal:           "#1C01CE",
	chroma.LiteralNumber:     "#1C01CE",
	chroma.Error:             "#000000",
	chroma.Background:        " bg:#ffffff",
}))
