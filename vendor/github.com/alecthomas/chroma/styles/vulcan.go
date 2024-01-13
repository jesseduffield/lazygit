package styles

import (
	"github.com/alecthomas/chroma"
)

var (
	// inspired by Doom Emacs's One Doom Theme
	black  = "#282C34"
	grey   = "#3E4460"
	grey2  = "#43454f"
	white  = "#C9C9C9"
	red    = "#CF5967"
	yellow = "#ECBE7B"
	green  = "#82CC6A"
	cyan   = "#56B6C2"
	blue   = "#7FBAF5"
	blue2  = "#57C7FF"
	purple = "#BC74C4"
)

var Vulcan = Register(chroma.MustNewStyle("vulcan", chroma.StyleEntries{
	chroma.Comment:                  grey,
	chroma.CommentHashbang:          grey + " italic",
	chroma.CommentMultiline:         grey,
	chroma.CommentPreproc:           blue,
	chroma.CommentSingle:            grey,
	chroma.CommentSpecial:           purple + " italic",
	chroma.Generic:                  white,
	chroma.GenericDeleted:           red,
	chroma.GenericEmph:              white + " underline",
	chroma.GenericError:             red + " bold",
	chroma.GenericHeading:           yellow + " bold",
	chroma.GenericInserted:          yellow,
	chroma.GenericOutput:            grey2,
	chroma.GenericPrompt:            white,
	chroma.GenericStrong:            red + " bold",
	chroma.GenericSubheading:        red + " italic",
	chroma.GenericTraceback:         white,
	chroma.GenericUnderline:         "underline",
	chroma.Error:                    red,
	chroma.Keyword:                  blue,
	chroma.KeywordConstant:          red + " bg:" + grey2,
	chroma.KeywordDeclaration:       blue,
	chroma.KeywordNamespace:         purple,
	chroma.KeywordPseudo:            purple,
	chroma.KeywordReserved:          blue,
	chroma.KeywordType:              blue2 + " bold",
	chroma.Literal:                  white,
	chroma.LiteralDate:              blue2,
	chroma.Name:                     white,
	chroma.NameAttribute:            purple,
	chroma.NameBuiltin:              blue,
	chroma.NameBuiltinPseudo:        blue,
	chroma.NameClass:                yellow,
	chroma.NameConstant:             yellow,
	chroma.NameDecorator:            yellow,
	chroma.NameEntity:               white,
	chroma.NameException:            red,
	chroma.NameFunction:             blue2,
	chroma.NameLabel:                red,
	chroma.NameNamespace:            white,
	chroma.NameOther:                white,
	chroma.NameTag:                  purple,
	chroma.NameVariable:             purple + " italic",
	chroma.NameVariableClass:        blue2 + " bold",
	chroma.NameVariableGlobal:       yellow,
	chroma.NameVariableInstance:     blue2,
	chroma.LiteralNumber:            cyan,
	chroma.LiteralNumberBin:         blue2,
	chroma.LiteralNumberFloat:       cyan,
	chroma.LiteralNumberHex:         blue2,
	chroma.LiteralNumberInteger:     cyan,
	chroma.LiteralNumberIntegerLong: cyan,
	chroma.LiteralNumberOct:         blue2,
	chroma.Operator:                 purple,
	chroma.OperatorWord:             purple,
	chroma.Other:                    white,
	chroma.Punctuation:              cyan,
	chroma.LiteralString:            green,
	chroma.LiteralStringBacktick:    blue2,
	chroma.LiteralStringChar:        blue2,
	chroma.LiteralStringDoc:         green,
	chroma.LiteralStringDouble:      green,
	chroma.LiteralStringEscape:      cyan,
	chroma.LiteralStringHeredoc:     cyan,
	chroma.LiteralStringInterpol:    green,
	chroma.LiteralStringOther:       green,
	chroma.LiteralStringRegex:       blue2,
	chroma.LiteralStringSingle:      green,
	chroma.LiteralStringSymbol:      green,
	chroma.Text:                     white,
	chroma.TextWhitespace:           white,
	chroma.Background:               " bg: " + black,
}))
