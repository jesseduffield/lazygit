// Copyright 2026 The TCell Authors
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

package tcell

import (
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3/color"
)

// Style represents a complete text style, including both foreground color,
// background color, and additional attributes such as "bold" or "underline".
//
// Note that not all terminals can display all colors or attributes, and
// many might have specific incompatibilities between specific attributes
// and color combinations.
//
// To use Style, just declare a variable of its type.
type Style struct {
	fg      color.Color
	bg      color.Color
	ulColor color.Color
	attrs   AttrMask
	ulStyle UnderlineStyle
	url     *urlInfo
}

type urlInfo struct {
	url string
	id  string
}

func stripOSCControls(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			c := s[i]
			if c <= 0x1f || c == 0x7f || (c >= 0x80 && c <= 0x9f) {
				i++
				continue
			}
			_ = b.WriteByte(c)
			i++
			continue
		}
		if r <= 0x1f || r == 0x7f || (r >= 0x80 && r <= 0x9f) {
			i += size
			continue
		}
		b.WriteString(s[i : i+size])
		i += size
	}
	return b.String()
}

// StyleDefault represents a default style, based upon the context.
// It is the zero value.
var StyleDefault Style

// styleInvalid is just an arbitrary invalid style used internally.
var styleInvalid = Style{attrs: AttrInvalid}

// Foreground returns a new style based on s, with the foreground color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Foreground(c color.Color) Style {
	s2 := s
	s2.fg = c
	return s2
}

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c color.Color) Style {
	s2 := s
	s2.bg = c
	return s2
}

func (s Style) setAttrs(attrs AttrMask, on bool) Style {
	s2 := s
	if on {
		s2.attrs |= attrs
	} else {
		s2.attrs &^= attrs
	}
	return s2
}

// Normal returns the style with all attributes disabled.
func (s Style) Normal() Style {
	return Style{
		fg: s.fg,
		bg: s.bg,
	}
}

// Bold returns a new style based on s, with the bold attribute set
// as requested.
func (s Style) Bold(on bool) Style {
	return s.setAttrs(AttrBold, on)
}

// Blink returns a new style based on s, with the blink attribute set
// as requested.
func (s Style) Blink(on bool) Style {
	return s.setAttrs(AttrBlink, on)
}

// Dim returns a new style based on s, with the dim attribute set
// as requested.
func (s Style) Dim(on bool) Style {
	return s.setAttrs(AttrDim, on)
}

// Italic returns a new style based on s, with the italic attribute set
// as requested.
func (s Style) Italic(on bool) Style {
	return s.setAttrs(AttrItalic, on)
}

// Reverse returns a new style based on s, with the reverse attribute set
// as requested.  (Reverse usually changes the foreground and background
// colors.)
func (s Style) Reverse(on bool) Style {
	return s.setAttrs(AttrReverse, on)
}

// StrikeThrough sets strike-through mode.
func (s Style) StrikeThrough(on bool) Style {
	return s.setAttrs(AttrStrikeThrough, on)
}

// Underline style.  Modern terminals have the option of rendering the
// underline using different styles, and even different colors.
type UnderlineStyle uint8

const (
	UnderlineStyleNone = UnderlineStyle(iota)
	UnderlineStyleSolid
	UnderlineStyleDouble
	UnderlineStyleCurly
	UnderlineStyleDotted
	UnderlineStyleDashed
)

// Underline returns a new style based on s, with the underline attribute set
// as requested.  The parameters can be:
//
// bool: on / off - enables just a simple underline
// UnderlineStyle: sets a specific style (should not coexist with the bool)
// Color: the color to use
func (s Style) Underline(params ...any) Style {
	s2 := s
	for _, param := range params {
		switch v := param.(type) {
		case bool:
			if v {
				s2.ulStyle = UnderlineStyleSolid
			} else {
				s2.ulStyle = UnderlineStyleNone
			}
		case UnderlineStyle:
			s2.ulStyle = v
		case Color:
			s2.ulColor = v
		default:
			panic("Bad type for underline")
		}
	}
	return s2
}

// GetForeground returns the foreground (text) color.
func (s Style) GetForeground() color.Color {
	return s.fg
}

// GetBackground returns the background color.
func (s Style) GetBackground() color.Color {
	return s.bg
}

// GetUnderlineStyle returns the underline style for the style.
func (s Style) GetUnderlineStyle() UnderlineStyle {
	return s.ulStyle
}

// GetUnderlineColor returns the underline color for the style.
func (s Style) GetUnderlineColor() color.Color {
	return s.ulColor
}

// Attributes returns a new style based on s, with its attributes set as
// specified.
//
// Deprecated: Use direct functions instead.
func (s Style) Attributes(attrs AttrMask) Style {
	s2 := s
	s2.attrs = attrs
	return s2
}

// GetAttributes gets the attributes for a style.
// Deprecated: Use individual properties instead.
func (s Style) GetAttributes() AttrMask {
	return s.attrs
}

// Url returns a style with the Url set.  If the provided Url is not empty,
// and the terminal supports it, text will typically be marked up as a clickable
// link to that Url.  If the Url is empty, then this mode is turned off.
func (s Style) Url(url string) Style {

	s2 := s
	s2.url = &urlInfo{url: stripOSCControls(url)}
	if s.url != nil {
		s2.url.id = s.url.id
	}
	return s2
}

// UrlId returns a style with the UrlId set. If the provided UrlId is not empty,
// any marked up Url with this style will be given the UrlId also. If the
// terminal supports it, any text with the same UrlId will be grouped as if it
// were one Url, even if it spans multiple lines.
func (s Style) UrlId(id string) Style {
	s2 := s
	s2.url = &urlInfo{
		id: "id=" + stripOSCControls(id),
	}
	if s.url != nil {
		s2.url.url = s.url.url
	}
	return s2
}

// GetUrl returns the URL (id and actual URL) associated with the style.
// This is a hyper link that will be used for cells marked up with this style.
func (s Style) GetUrl() (id string, url string) {
	if s.url != nil {
		return strings.TrimPrefix(s.url.id, "id="), s.url.url
	}
	return "", ""
}

// HasBold returns true if the style indicates bold text.
// Note that on some terminals bold text is simply brighter.
func (s Style) HasBold() bool {
	return s.attrs&AttrBold != 0
}

// HasBlink returns true if the style indicates blinking text.
func (s Style) HasBlink() bool {
	return s.attrs&AttrBlink != 0
}

// HasReverse returns true if the style indicates reverse video text.
func (s Style) HasReverse() bool {
	return s.attrs&AttrReverse != 0
}

// HasItalic returns true if the style indicates italicized text.
func (s Style) HasItalic() bool {
	return s.attrs&AttrItalic != 0
}

// HasDim returns true if the style indicates dim or faint text.
func (s Style) HasDim() bool {
	return s.attrs&AttrDim != 0
}

// HasStrikeThrough returns true if the style indicates crossed-out text.
func (s Style) HasStrikeThrough() bool {
	return s.attrs&AttrStrikeThrough != 0
}

// HasUnderline returns true if any underline style is set.
// Note that more detail is available via the GetUnderlineStyle
// and GetUnderlineColor methods.
func (s Style) HasUnderline() bool {
	return s.ulStyle != UnderlineStyleNone
}
