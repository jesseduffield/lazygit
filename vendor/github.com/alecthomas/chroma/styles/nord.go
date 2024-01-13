package styles

import (
	"github.com/alecthomas/chroma"
)

var (
	// colors and palettes based on https://www.nordtheme.com/docs/colors-and-palettes
	nord0  = "#2e3440"
	nord1  = "#3b4252" // nolint
	nord2  = "#434c5e" // nolint
	nord3  = "#4c566a"
	nord3b = "#616e87"

	nord4 = "#d8dee9"
	nord5 = "#e5e9f0" // nolint
	nord6 = "#eceff4"

	nord7  = "#8fbcbb"
	nord8  = "#88c0d0"
	nord9  = "#81a1c1"
	nord10 = "#5e81ac"

	nord11 = "#bf616a"
	nord12 = "#d08770"
	nord13 = "#ebcb8b"
	nord14 = "#a3be8c"
	nord15 = "#b48ead"
)

// Nord, an arctic, north-bluish color palette
var Nord = Register(chroma.MustNewStyle("nord", chroma.StyleEntries{
	chroma.TextWhitespace:        nord4,
	chroma.Comment:               "italic " + nord3b,
	chroma.CommentPreproc:        nord10,
	chroma.Keyword:               "bold " + nord9,
	chroma.KeywordPseudo:         "nobold " + nord9,
	chroma.KeywordType:           "nobold " + nord9,
	chroma.Operator:              nord9,
	chroma.OperatorWord:          "bold " + nord9,
	chroma.Name:                  nord4,
	chroma.NameBuiltin:           nord9,
	chroma.NameFunction:          nord8,
	chroma.NameClass:             nord7,
	chroma.NameNamespace:         nord7,
	chroma.NameException:         nord11,
	chroma.NameVariable:          nord4,
	chroma.NameConstant:          nord7,
	chroma.NameLabel:             nord7,
	chroma.NameEntity:            nord12,
	chroma.NameAttribute:         nord7,
	chroma.NameTag:               nord9,
	chroma.NameDecorator:         nord12,
	chroma.Punctuation:           nord6,
	chroma.LiteralString:         nord14,
	chroma.LiteralStringDoc:      nord3b,
	chroma.LiteralStringInterpol: nord14,
	chroma.LiteralStringEscape:   nord13,
	chroma.LiteralStringRegex:    nord13,
	chroma.LiteralStringSymbol:   nord14,
	chroma.LiteralStringOther:    nord14,
	chroma.LiteralNumber:         nord15,
	chroma.GenericHeading:        "bold " + nord8,
	chroma.GenericSubheading:     "bold " + nord8,
	chroma.GenericDeleted:        nord11,
	chroma.GenericInserted:       nord14,
	chroma.GenericError:          nord11,
	chroma.GenericEmph:           "italic",
	chroma.GenericStrong:         "bold",
	chroma.GenericPrompt:         "bold " + nord3,
	chroma.GenericOutput:         nord4,
	chroma.GenericTraceback:      nord11,
	chroma.Error:                 nord11,
	chroma.Background:            nord4 + " bg:" + nord0,
}))
