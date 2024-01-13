package chroma

import (
	"fmt"
	"strings"
)

// Trilean value for StyleEntry value inheritance.
type Trilean uint8

// Trilean states.
const (
	Pass Trilean = iota
	Yes
	No
)

func (t Trilean) String() string {
	switch t {
	case Yes:
		return "Yes"
	case No:
		return "No"
	default:
		return "Pass"
	}
}

// Prefix returns s with "no" as a prefix if Trilean is no.
func (t Trilean) Prefix(s string) string {
	if t == Yes {
		return s
	} else if t == No {
		return "no" + s
	}
	return ""
}

// A StyleEntry in the Style map.
type StyleEntry struct {
	// Hex colours.
	Colour     Colour
	Background Colour
	Border     Colour

	Bold      Trilean
	Italic    Trilean
	Underline Trilean
	NoInherit bool
}

func (s StyleEntry) String() string {
	out := []string{}
	if s.Bold != Pass {
		out = append(out, s.Bold.Prefix("bold"))
	}
	if s.Italic != Pass {
		out = append(out, s.Italic.Prefix("italic"))
	}
	if s.Underline != Pass {
		out = append(out, s.Underline.Prefix("underline"))
	}
	if s.NoInherit {
		out = append(out, "noinherit")
	}
	if s.Colour.IsSet() {
		out = append(out, s.Colour.String())
	}
	if s.Background.IsSet() {
		out = append(out, "bg:"+s.Background.String())
	}
	if s.Border.IsSet() {
		out = append(out, "border:"+s.Border.String())
	}
	return strings.Join(out, " ")
}

// Sub subtracts e from s where elements match.
func (s StyleEntry) Sub(e StyleEntry) StyleEntry {
	out := StyleEntry{}
	if e.Colour != s.Colour {
		out.Colour = s.Colour
	}
	if e.Background != s.Background {
		out.Background = s.Background
	}
	if e.Bold != s.Bold {
		out.Bold = s.Bold
	}
	if e.Italic != s.Italic {
		out.Italic = s.Italic
	}
	if e.Underline != s.Underline {
		out.Underline = s.Underline
	}
	if e.Border != s.Border {
		out.Border = s.Border
	}
	return out
}

// Inherit styles from ancestors.
//
// Ancestors should be provided from oldest to newest.
func (s StyleEntry) Inherit(ancestors ...StyleEntry) StyleEntry {
	out := s
	for i := len(ancestors) - 1; i >= 0; i-- {
		if out.NoInherit {
			return out
		}
		ancestor := ancestors[i]
		if !out.Colour.IsSet() {
			out.Colour = ancestor.Colour
		}
		if !out.Background.IsSet() {
			out.Background = ancestor.Background
		}
		if !out.Border.IsSet() {
			out.Border = ancestor.Border
		}
		if out.Bold == Pass {
			out.Bold = ancestor.Bold
		}
		if out.Italic == Pass {
			out.Italic = ancestor.Italic
		}
		if out.Underline == Pass {
			out.Underline = ancestor.Underline
		}
	}
	return out
}

func (s StyleEntry) IsZero() bool {
	return s.Colour == 0 && s.Background == 0 && s.Border == 0 && s.Bold == Pass && s.Italic == Pass &&
		s.Underline == Pass && !s.NoInherit
}

// A StyleBuilder is a mutable structure for building styles.
//
// Once built, a Style is immutable.
type StyleBuilder struct {
	entries map[TokenType]string
	name    string
	parent  *Style
}

func NewStyleBuilder(name string) *StyleBuilder {
	return &StyleBuilder{name: name, entries: map[TokenType]string{}}
}

func (s *StyleBuilder) AddAll(entries StyleEntries) *StyleBuilder {
	for ttype, entry := range entries {
		s.entries[ttype] = entry
	}
	return s
}

func (s *StyleBuilder) Get(ttype TokenType) StyleEntry {
	// This is less than ideal, but it's the price for having to check errors on each Add().
	entry, _ := ParseStyleEntry(s.entries[ttype])
	return entry.Inherit(s.parent.Get(ttype))
}

// Add an entry to the Style map.
//
// See http://pygments.org/docs/styles/#style-rules for details.
func (s *StyleBuilder) Add(ttype TokenType, entry string) *StyleBuilder { // nolint: gocyclo
	s.entries[ttype] = entry
	return s
}

func (s *StyleBuilder) AddEntry(ttype TokenType, entry StyleEntry) *StyleBuilder {
	s.entries[ttype] = entry.String()
	return s
}

func (s *StyleBuilder) Build() (*Style, error) {
	style := &Style{
		Name:    s.name,
		entries: map[TokenType]StyleEntry{},
		parent:  s.parent,
	}
	for ttype, descriptor := range s.entries {
		entry, err := ParseStyleEntry(descriptor)
		if err != nil {
			return nil, fmt.Errorf("invalid entry for %s: %s", ttype, err)
		}
		style.entries[ttype] = entry
	}
	return style, nil
}

// StyleEntries mapping TokenType to colour definition.
type StyleEntries map[TokenType]string

// NewStyle creates a new style definition.
func NewStyle(name string, entries StyleEntries) (*Style, error) {
	return NewStyleBuilder(name).AddAll(entries).Build()
}

// MustNewStyle creates a new style or panics.
func MustNewStyle(name string, entries StyleEntries) *Style {
	style, err := NewStyle(name, entries)
	if err != nil {
		panic(err)
	}
	return style
}

// A Style definition.
//
// See http://pygments.org/docs/styles/ for details. Semantics are intended to be identical.
type Style struct {
	Name    string
	entries map[TokenType]StyleEntry
	parent  *Style
}

// Types that are styled.
func (s *Style) Types() []TokenType {
	dedupe := map[TokenType]bool{}
	for tt := range s.entries {
		dedupe[tt] = true
	}
	if s.parent != nil {
		for _, tt := range s.parent.Types() {
			dedupe[tt] = true
		}
	}
	out := make([]TokenType, 0, len(dedupe))
	for tt := range dedupe {
		out = append(out, tt)
	}
	return out
}

// Builder creates a mutable builder from this Style.
//
// The builder can then be safely modified. This is a cheap operation.
func (s *Style) Builder() *StyleBuilder {
	return &StyleBuilder{
		name:    s.Name,
		entries: map[TokenType]string{},
		parent:  s,
	}
}

// Has checks if an exact style entry match exists for a token type.
//
// This is distinct from Get() which will merge parent tokens.
func (s *Style) Has(ttype TokenType) bool {
	return !s.get(ttype).IsZero() || s.synthesisable(ttype)
}

// Get a style entry. Will try sub-category or category if an exact match is not found, and
// finally return the Background.
func (s *Style) Get(ttype TokenType) StyleEntry {
	return s.get(ttype).Inherit(
		s.get(Background),
		s.get(Text),
		s.get(ttype.Category()),
		s.get(ttype.SubCategory()))
}

func (s *Style) get(ttype TokenType) StyleEntry {
	out := s.entries[ttype]
	if out.IsZero() && s.parent != nil {
		return s.parent.get(ttype)
	}
	if out.IsZero() && s.synthesisable(ttype) {
		out = s.synthesise(ttype)
	}
	return out
}

func (s *Style) synthesise(ttype TokenType) StyleEntry {
	bg := s.get(Background)
	text := StyleEntry{Colour: bg.Colour}
	text.Colour = text.Colour.BrightenOrDarken(0.5)

	switch ttype {
	// If we don't have a line highlight colour, make one that is 10% brighter/darker than the background.
	case LineHighlight:
		return StyleEntry{Background: bg.Background.BrightenOrDarken(0.1)}

	// If we don't have line numbers, use the text colour but 20% brighter/darker
	case LineNumbers, LineNumbersTable:
		return text

	default:
		return StyleEntry{}
	}
}

func (s *Style) synthesisable(ttype TokenType) bool {
	return ttype == LineHighlight || ttype == LineNumbers || ttype == LineNumbersTable
}

// ParseStyleEntry parses a Pygments style entry.
func ParseStyleEntry(entry string) (StyleEntry, error) { // nolint: gocyclo
	out := StyleEntry{}
	parts := strings.Fields(entry)
	for _, part := range parts {
		switch {
		case part == "italic":
			out.Italic = Yes
		case part == "noitalic":
			out.Italic = No
		case part == "bold":
			out.Bold = Yes
		case part == "nobold":
			out.Bold = No
		case part == "underline":
			out.Underline = Yes
		case part == "nounderline":
			out.Underline = No
		case part == "inherit":
			out.NoInherit = false
		case part == "noinherit":
			out.NoInherit = true
		case part == "bg:":
			out.Background = 0
		case strings.HasPrefix(part, "bg:#"):
			out.Background = ParseColour(part[3:])
			if !out.Background.IsSet() {
				return StyleEntry{}, fmt.Errorf("invalid background colour %q", part)
			}
		case strings.HasPrefix(part, "border:#"):
			out.Border = ParseColour(part[7:])
			if !out.Border.IsSet() {
				return StyleEntry{}, fmt.Errorf("invalid border colour %q", part)
			}
		case strings.HasPrefix(part, "#"):
			out.Colour = ParseColour(part)
			if !out.Colour.IsSet() {
				return StyleEntry{}, fmt.Errorf("invalid colour %q", part)
			}
		default:
			return StyleEntry{}, fmt.Errorf("unknown style element %q", part)
		}
	}
	return out, nil
}
