// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

// Style represents a complete text style, including both foreground color,
// background color, and additional attributes such as "bold" or "underline".
//
// Note that not all terminals can display all colors or attributes, and
// many might have specific incompatibilities between specific attributes
// and color combinations.
//
// To use Style, just declare a variable of its type.
type Style struct {
	fg    Color
	bg    Color
	attrs AttrMask
}

// StyleDefault represents a default style, based upon the context.
// It is the zero value.
var StyleDefault Style

// styleInvalid is just an arbitrary invalid style used internally.
var styleInvalid = Style{attrs: AttrInvalid}

// Foreground returns a new style based on s, with the foreground color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Foreground(c Color) Style {
	return Style{
		fg:    c,
		bg:    s.bg,
		attrs: s.attrs,
	}
}

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c Color) Style {
	return Style{
		fg:    s.fg,
		bg:    c,
		attrs: s.attrs,
	}
}

// Decompose breaks a style up, returning the foreground, background,
// and other attributes.
func (s Style) Decompose() (fg Color, bg Color, attr AttrMask) {
	return s.fg, s.bg, s.attrs
}

func (s Style) setAttrs(attrs AttrMask, on bool) Style {
	if on {
		return Style{
			fg:    s.fg,
			bg:    s.bg,
			attrs: s.attrs | attrs,
		}
	}
	return Style{
		fg:    s.fg,
		bg:    s.bg,
		attrs: s.attrs &^ attrs,
	}
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

// Underline returns a new style based on s, with the underline attribute set
// as requested.
func (s Style) Underline(on bool) Style {
	return s.setAttrs(AttrUnderline, on)
}

// StrikeThrough sets strikethrough mode.
func (s Style) StrikeThrough(on bool) Style {
	return s.setAttrs(AttrStrikeThrough, on)
}

// Attributes returns a new style based on s, with its attributes set as
// specified.
func (s Style) Attributes(attrs AttrMask) Style {
	return Style{
		fg:    s.fg,
		bg:    s.bg,
		attrs: attrs,
	}
}
