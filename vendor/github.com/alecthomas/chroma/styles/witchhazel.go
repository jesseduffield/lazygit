// Copyright 2018 Alethea Katherine Flowers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package styles

import (
	"github.com/alecthomas/chroma"
)

// WitchHazel Style
var WitchHazel = Register(chroma.MustNewStyle("witchhazel", chroma.StyleEntries{
	chroma.Text:              "#F8F8F2",
	chroma.Whitespace:        "#A8757B",
	chroma.Error:             "#960050 bg:#1e0010",
	chroma.Comment:           "#b0bec5",
	chroma.Keyword:           "#C2FFDF",
	chroma.KeywordNamespace:  "#FFB8D1",
	chroma.Operator:          "#FFB8D1",
	chroma.Punctuation:       "#F8F8F2",
	chroma.Name:              "#F8F8F2",
	chroma.NameAttribute:     "#ceb1ff",
	chroma.NameBuiltinPseudo: "#80cbc4",
	chroma.NameClass:         "#ceb1ff",
	chroma.NameConstant:      "#C5A3FF",
	chroma.NameDecorator:     "#ceb1ff",
	chroma.NameException:     "#ceb1ff",
	chroma.NameFunction:      "#ceb1ff",
	chroma.NameProperty:      "#F8F8F2",
	chroma.NameTag:           "#FFB8D1",
	chroma.NameVariable:      "#F8F8F2",
	chroma.Number:            "#C5A3FF",
	chroma.Literal:           "#ae81ff",
	chroma.LiteralDate:       "#e6db74",
	chroma.String:            "#1bc5e0",
	chroma.GenericDeleted:    "#f92672",
	chroma.GenericEmph:       "italic",
	chroma.GenericInserted:   "#a6e22e",
	chroma.GenericStrong:     "bold",
	chroma.GenericSubheading: "#75715e",
	chroma.Background:        " bg:#433e56",
}))
