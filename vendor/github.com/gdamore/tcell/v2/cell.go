// Copyright 2025 The TCell Authors
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

import (
	"github.com/rivo/uniseg"
)

type cell struct {
	currStr   string
	lastStr   string
	currStyle Style
	lastStyle Style
	width     int
	lock      bool
}

func (c *cell) setDirty(dirty bool) {
	if dirty {
		c.lastStr = ""
	} else {
		if c.currStr == "" {
			c.currStr = " "
		}
		c.lastStr = c.currStr
		c.lastStyle = c.currStyle
	}
}

// CellBuffer represents a two-dimensional array of character cells.
// This is primarily intended for use by Screen implementors; it
// contains much of the common code they need.  To create one, just
// declare a variable of its type; no explicit initialization is necessary.
//
// CellBuffer is not thread safe.
type CellBuffer struct {
	w     int
	h     int
	cells []cell
}

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.  If the background or
// foreground of the style is set to ColorNone, then the respective
// color is left un changed.
//
// Deprecated: Use Put instead, which this is implemented in terms of.
func (cb *CellBuffer) SetContent(x int, y int, mainc rune, combc []rune, style Style) {
	cb.Put(x, y, string(append([]rune{mainc}, combc...)), style)
}

// Put a single styled grapheme using the given string and style
// at the same location.  Note that only the first grapheme in the string
// will bre displayed, using only the 1 or 2 (depending on width) cells
// located at x, y. It returns the rest of the string, and the width used.
func (cb *CellBuffer) Put(x int, y int, str string, style Style) (string, int) {
	var width int = 0
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		var cl string
		c := &cb.cells[(y*cb.w)+x]
		state := -1
		for width == 0 && str != "" {
			var g string
			g, str, width, state = uniseg.FirstGraphemeClusterInString(str, state)
			cl += g
			if g == "" {
				break
			}
		}

		// Wide characters: we want to mark the "wide" cells
		// dirty as well as the base cell, to make sure we consider
		// both cells as dirty together.  We only need to do this
		// if we're changing content
		if width > 0 && cl != c.currStr {
			// Prevent unnecessary boundchecks for first cell, since we already
			// received that one.
			c.setDirty(true)
			for i := 1; i < width; i++ {
				cb.SetDirty(x+i, y, true)
			}
		}

		c.currStr = cl
		c.width = width

		if style.fg == ColorNone {
			style.fg = c.currStyle.fg
		}
		if style.bg == ColorNone {
			style.bg = c.currStyle.bg
		}
		c.currStyle = style
	}
	return str, width
}

// Get the contents of a character cell (or two adjacent cells), including the
// the style and the display width in cells.  (The width can be either 1, normally,
// or 2 for East Asian full-width characters.  If the width is 0, then the cell is
// is empty.)
func (cb *CellBuffer) Get(x, y int) (string, Style, int) {
	var style Style
	var width int
	var str string
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := &cb.cells[(y*cb.w)+x]
		str, style = c.currStr, c.currStyle
		if width = c.width; width == 0 || str == "" {
			width = 1
			str = " "
		}
	}
	return str, style, width
}

// GetContent returns the contents of a character cell, including the
// primary rune, any combining character runes (which will usually be
// nil), the style, and the display width in cells.  (The width can be
// either 1, normally, or 2 for East Asian full-width characters.)
//
// Deprecated: Use Get, which this implemented in terms of.
func (cb *CellBuffer) GetContent(x, y int) (rune, []rune, Style, int) {
	var style Style
	var width int
	var mainc rune
	var combc []rune
	str, style, width := cb.Get(x, y)
	for i, r := range str {
		if i == 0 {
			mainc = r
		} else {
			combc = append(combc, r)
		}
	}
	return mainc, combc, style, width
}

// Size returns the (width, height) in cells of the buffer.
func (cb *CellBuffer) Size() (int, int) {
	return cb.w, cb.h
}

// Invalidate marks all characters within the buffer as dirty.
func (cb *CellBuffer) Invalidate() {
	for i := range cb.cells {
		cb.cells[i].lastStr = ""
	}
}

// Dirty checks if a character at the given location needs to be
// refreshed on the physical display.  This returns true if the cell
// content is different since the last time it was marked clean.
func (cb *CellBuffer) Dirty(x, y int) bool {
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := &cb.cells[(y*cb.w)+x]
		if c.lock {
			return false
		}
		if c.lastStyle != c.currStyle {
			return true
		}
		if c.lastStr != c.currStr {
			return true
		}
	}
	return false
}

// SetDirty is normally used to indicate that a cell has
// been displayed (in which case dirty is false), or to manually
// force a cell to be marked dirty.
func (cb *CellBuffer) SetDirty(x, y int, dirty bool) {
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := &cb.cells[(y*cb.w)+x]
		c.setDirty(dirty)
	}
}

// LockCell locks a cell from being drawn, effectively marking it "clean" until
// the lock is removed. This can be used to prevent tcell from drawing a given
// cell, even if the underlying content has changed. For example, when drawing a
// sixel graphic directly to a TTY screen an implementer must lock the region
// underneath the graphic to prevent tcell from drawing on top of the graphic.
func (cb *CellBuffer) LockCell(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= cb.w || y >= cb.h {
		return
	}
	c := &cb.cells[(y*cb.w)+x]
	c.lock = true
}

// UnlockCell removes a lock from the cell and marks it as dirty
func (cb *CellBuffer) UnlockCell(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= cb.w || y >= cb.h {
		return
	}
	c := &cb.cells[(y*cb.w)+x]
	c.lock = false
	cb.SetDirty(x, y, true)
}

// Resize is used to resize the cells array, with different dimensions,
// while preserving the original contents.  The cells will be invalidated
// so that they can be redrawn.
func (cb *CellBuffer) Resize(w, h int) {
	if cb.h == h && cb.w == w {
		return
	}

	newc := make([]cell, w*h)
	for y := 0; y < h && y < cb.h; y++ {
		for x := 0; x < w && x < cb.w; x++ {
			oc := &cb.cells[(y*cb.w)+x]
			nc := &newc[(y*w)+x]
			nc.currStr = oc.currStr
			nc.currStyle = oc.currStyle
			nc.width = oc.width
			nc.lastStr = ""
		}
	}
	cb.cells = newc
	cb.h = h
	cb.w = w
}

// Fill fills the entire cell buffer array with the specified character
// and style.  Normally choose ' ' to clear the screen.  This API doesn't
// support combining characters, or characters with a width larger than one.
// If either the foreground or background are ColorNone, then the respective
// color is unchanged.
func (cb *CellBuffer) Fill(r rune, style Style) {
	for i := range cb.cells {
		c := &cb.cells[i]
		c.currStr = string(r)
		cs := style
		if cs.fg == ColorNone {
			cs.fg = c.currStyle.fg
		}
		if cs.bg == ColorNone {
			cs.bg = c.currStyle.bg
		}
		c.currStyle = cs
		c.width = 1
	}
}
