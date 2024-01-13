package styles

import (
	"github.com/alecthomas/chroma"
)

// ParaisoLight style.
var ParaisoLight = Register(chroma.MustNewStyle("paraiso-light", chroma.StyleEntries{
	chroma.Text:                  "#2f1e2e",
	chroma.Error:                 "#ef6155",
	chroma.Comment:               "#8d8687",
	chroma.Keyword:               "#815ba4",
	chroma.KeywordNamespace:      "#5bc4bf",
	chroma.KeywordType:           "#fec418",
	chroma.Operator:              "#5bc4bf",
	chroma.Punctuation:           "#2f1e2e",
	chroma.Name:                  "#2f1e2e",
	chroma.NameAttribute:         "#06b6ef",
	chroma.NameClass:             "#fec418",
	chroma.NameConstant:          "#ef6155",
	chroma.NameDecorator:         "#5bc4bf",
	chroma.NameException:         "#ef6155",
	chroma.NameFunction:          "#06b6ef",
	chroma.NameNamespace:         "#fec418",
	chroma.NameOther:             "#06b6ef",
	chroma.NameTag:               "#5bc4bf",
	chroma.NameVariable:          "#ef6155",
	chroma.LiteralNumber:         "#f99b15",
	chroma.Literal:               "#f99b15",
	chroma.LiteralDate:           "#48b685",
	chroma.LiteralString:         "#48b685",
	chroma.LiteralStringChar:     "#2f1e2e",
	chroma.LiteralStringDoc:      "#8d8687",
	chroma.LiteralStringEscape:   "#f99b15",
	chroma.LiteralStringInterpol: "#f99b15",
	chroma.GenericDeleted:        "#ef6155",
	chroma.GenericEmph:           "italic",
	chroma.GenericHeading:        "bold #2f1e2e",
	chroma.GenericInserted:       "#48b685",
	chroma.GenericPrompt:         "bold #8d8687",
	chroma.GenericStrong:         "bold",
	chroma.GenericSubheading:     "bold #5bc4bf",
	chroma.Background:            "bg:#e7e9db",
}))
