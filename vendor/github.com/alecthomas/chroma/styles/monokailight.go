package styles

import (
	"github.com/alecthomas/chroma"
)

// MonokaiLight style.
var MonokaiLight = Register(chroma.MustNewStyle("monokailight", chroma.StyleEntries{
	chroma.Text:                "#272822",
	chroma.Error:               "#960050 bg:#1e0010",
	chroma.Comment:             "#75715e",
	chroma.Keyword:             "#00a8c8",
	chroma.KeywordNamespace:    "#f92672",
	chroma.Operator:            "#f92672",
	chroma.Punctuation:         "#111111",
	chroma.Name:                "#111111",
	chroma.NameAttribute:       "#75af00",
	chroma.NameClass:           "#75af00",
	chroma.NameConstant:        "#00a8c8",
	chroma.NameDecorator:       "#75af00",
	chroma.NameException:       "#75af00",
	chroma.NameFunction:        "#75af00",
	chroma.NameOther:           "#75af00",
	chroma.NameTag:             "#f92672",
	chroma.LiteralNumber:       "#ae81ff",
	chroma.Literal:             "#ae81ff",
	chroma.LiteralDate:         "#d88200",
	chroma.LiteralString:       "#d88200",
	chroma.LiteralStringEscape: "#8045FF",
	chroma.GenericEmph:         "italic",
	chroma.GenericStrong:       "bold",
	chroma.Background:          " bg:#fafafa",
}))
