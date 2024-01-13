package styles

import (
	"github.com/alecthomas/chroma"
)

// ParaisoDark style.
var ParaisoDark = Register(chroma.MustNewStyle("paraiso-dark", chroma.StyleEntries{
	chroma.Text:                  "#e7e9db",
	chroma.Error:                 "#ef6155",
	chroma.Comment:               "#776e71",
	chroma.Keyword:               "#815ba4",
	chroma.KeywordNamespace:      "#5bc4bf",
	chroma.KeywordType:           "#fec418",
	chroma.Operator:              "#5bc4bf",
	chroma.Punctuation:           "#e7e9db",
	chroma.Name:                  "#e7e9db",
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
	chroma.LiteralStringChar:     "#e7e9db",
	chroma.LiteralStringDoc:      "#776e71",
	chroma.LiteralStringEscape:   "#f99b15",
	chroma.LiteralStringInterpol: "#f99b15",
	chroma.GenericDeleted:        "#ef6155",
	chroma.GenericEmph:           "italic",
	chroma.GenericHeading:        "bold #e7e9db",
	chroma.GenericInserted:       "#48b685",
	chroma.GenericPrompt:         "bold #776e71",
	chroma.GenericStrong:         "bold",
	chroma.GenericSubheading:     "bold #5bc4bf",
	chroma.Background:            "bg:#2f1e2e",
}))
