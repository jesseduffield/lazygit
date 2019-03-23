// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"github.com/go-errors/errors"

	"github.com/mattn/go-runewidth"
)

const maxInt = int(^uint(0) >> 1)

// Editor interface must be satisfied by gocui editors.
type Editor interface {
	Edit(v *View, key Key, ch rune, mod Modifier)
}

// The EditorFunc type is an adapter to allow the use of ordinary functions as
// Editors. If f is a function with the appropriate signature, EditorFunc(f)
// is an Editor object that calls f.
type EditorFunc func(v *View, key Key, ch rune, mod Modifier)

// Edit calls f(v, key, ch, mod)
func (f EditorFunc) Edit(v *View, key Key, ch rune, mod Modifier) {
	f(v, key, ch, mod)
}

// DefaultEditor is the default editor.
var DefaultEditor Editor = EditorFunc(simpleEditor)

// simpleEditor is used as the default gocui editor.
func simpleEditor(v *View, key Key, ch rune, mod Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == KeySpace:
		v.EditWrite(' ')
	case key == KeyBackspace || key == KeyBackspace2:
		v.EditDelete(true)
	case key == KeyDelete:
		v.EditDelete(false)
	case key == KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == KeyEnter:
		v.EditNewLine()
	case key == KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

// EditWrite writes a rune at the cursor position.
func (v *View) EditWrite(ch rune) {
	w := runewidth.RuneWidth(ch)
	v.writeRune(v.cx, v.cy, ch)
	v.moveCursor(w, 0, true)
}

// EditDelete deletes a rune at the cursor position. back determines the
// direction.
func (v *View) EditDelete(back bool) {
	x, y := v.ox+v.cx, v.oy+v.cy
	if y < 0 {
		return
	} else if y >= len(v.viewLines) {
		v.MoveCursor(-1, 0, true)
		return
	}

	maxX, _ := v.Size()
	if back {
		if x == 0 { // start of the line
			if y < 1 {
				return
			}

			var maxPrevWidth int
			if v.Wrap {
				maxPrevWidth = maxX
			} else {
				maxPrevWidth = maxInt
			}

			if v.viewLines[y].linesX == 0 { // regular line
				v.mergeLines(v.cy - 1)
				if len(v.viewLines[y-1].line) < maxPrevWidth {
					v.MoveCursor(-1, 0, true)
				}
			} else { // wrapped line
				n, _ := v.deleteRune(len(v.viewLines[y-1].line)-1, v.cy-1)
				v.MoveCursor(-n, 0, true)
			}
		} else { // middle/end of the line
			n, _ := v.deleteRune(v.cx-1, v.cy)
			v.MoveCursor(-n, 0, true)
		}
	} else {
		if x == len(v.viewLines[y].line) { // end of the line
			v.mergeLines(v.cy)
		} else { // start/middle of the line
			v.deleteRune(v.cx, v.cy)
		}
	}
}

// EditNewLine inserts a new line under the cursor.
func (v *View) EditNewLine() {
	v.breakLine(v.cx, v.cy)
	v.ox = 0
	v.cy = v.cy + 1
	v.cx = 0
}

// MoveCursor moves the cursor taking into account the width of the line/view,
// displacing the origin if necessary.
func (v *View) MoveCursor(dx, dy int, writeMode bool) {
	ox, oy := v.cx+v.ox, v.cy+v.oy
	x, y := ox+dx, oy+dy

	if y < 0 || y >= len(v.viewLines) {
		v.moveCursor(dx, dy, writeMode)
		return
	}

	// Removing newline.
	if x < 0 {
		var prevLen int
		if y-1 >= 0 && y-1 < len(v.viewLines) {
			prevLen = lineWidth(v.viewLines[y-1].line)
		}

		v.MoveCursor(prevLen, -1, writeMode)
		return
	}

	line := v.viewLines[y].line
	var col int
	var prevCol int
	for i := range line {
		prevCol = col
		col += runewidth.RuneWidth(line[i].chr)
		if dx > 0 {
			if x <= col {
				x = col
				break
			}
			continue
		}

		if x < col {
			x = prevCol
			break
		}
	}

	v.moveCursor(x-ox, y-oy, writeMode)
}

func (v *View) moveCursor(dx, dy int, writeMode bool) {
	maxX, maxY := v.Size()
	cx, cy := v.cx+dx, v.cy+dy
	x, y := v.ox+cx, v.oy+cy

	var curLineWidth, prevLineWidth int
	// get the width of the current line
	curLineWidth = maxInt
	if v.Wrap {
		curLineWidth = maxX - 1
	}

	if !writeMode {
		curLineWidth = 0
		if y >= 0 && y < len(v.viewLines) {
			curLineWidth = lineWidth(v.viewLines[y].line)
			if v.Wrap && curLineWidth >= maxX {
				curLineWidth = maxX - 1
			}
		}
	}
	// get the width of the previous line
	prevLineWidth = 0
	if y-1 >= 0 && y-1 < len(v.viewLines) {
		prevLineWidth = lineWidth(v.viewLines[y-1].line)
	}
	// adjust cursor's x position and view's x origin
	if x > curLineWidth { // move to next line
		if dx > 0 { // horizontal movement
			cy++
			if writeMode || v.oy+cy < len(v.viewLines) {
				if !v.Wrap {
					v.ox = 0
				}
				v.cx = 0
			}
		} else { // vertical movement
			if curLineWidth > 0 { // move cursor to the EOL
				if v.Wrap {
					v.cx = curLineWidth
				} else {
					ncx := curLineWidth - v.ox
					if ncx < 0 {
						v.ox += ncx
						if v.ox < 0 {
							v.ox = 0
						}
						v.cx = 0
					} else {
						v.cx = ncx
					}
				}
			} else {
				if writeMode || v.oy+cy < len(v.viewLines) {
					if !v.Wrap {
						v.ox = 0
					}
					v.cx = 0
				}
			}
		}
	} else if cx < 0 {
		if !v.Wrap && v.ox > 0 { // move origin to the left
			v.ox += cx
			v.cx = 0
		} else { // move to previous line
			cy--
			if prevLineWidth > 0 {
				if !v.Wrap { // set origin so the EOL is visible
					nox := prevLineWidth - maxX + 1
					if nox < 0 {
						nox = 0
					}
					v.ox = nox
				}
				v.cx = prevLineWidth
			} else {
				if !v.Wrap {
					v.ox = 0
				}
				v.cx = 0
			}
		}
	} else { // stay on the same line
		if v.Wrap {
			v.cx = cx
		} else {
			if cx >= maxX {
				v.ox += cx - maxX + 1
				v.cx = maxX
			} else {
				v.cx = cx
			}
		}
	}

	// adjust cursor's y position and view's y origin
	if cy < 0 {
		if v.oy > 0 {
			v.oy--
		}
	} else if writeMode || v.oy+cy < len(v.viewLines) {
		if cy >= maxY {
			v.oy++
		} else {
			v.cy = cy
		}
	}
}

// writeRune writes a rune into the view's internal buffer, at the
// position corresponding to the point (x, y). The length of the internal
// buffer is increased if the point is out of bounds. Overwrite mode is
// governed by the value of View.overwrite.
func (v *View) writeRune(x, y int, ch rune) error {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return err
	}

	if x < 0 || y < 0 {
		return errors.New("invalid point")
	}

	if y >= len(v.lines) {
		s := make([][]cell, y-len(v.lines)+1)
		v.lines = append(v.lines, s...)
	}

	olen := len(v.lines[y])

	var s []cell
	if x >= len(v.lines[y]) {
		s = make([]cell, x-len(v.lines[y])+1)
	} else if !v.Overwrite {
		s = make([]cell, 1)
	}
	v.lines[y] = append(v.lines[y], s...)

	if !v.Overwrite || (v.Overwrite && x >= olen-1) {
		copy(v.lines[y][x+1:], v.lines[y][x:])
	}
	v.lines[y][x] = cell{
		fgColor: v.FgColor,
		bgColor: v.BgColor,
		chr:     ch,
	}

	return nil
}

// deleteRune removes a rune from the view's internal buffer, at the
// position corresponding to the point (x, y).
// returns the amount of columns that where removed.
func (v *View) deleteRune(x, y int) (int, error) {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return 0, err
	}

	if x < 0 || y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
		return 0, errors.New("invalid point")
	}

	var tw int
	for i := range v.lines[y] {
		w := runewidth.RuneWidth(v.lines[y][i].chr)
		tw += w
		if tw > x {
			v.lines[y] = append(v.lines[y][:i], v.lines[y][i+1:]...)
			return w, nil
		}

	}

	return 0, nil
}

// mergeLines merges the lines "y" and "y+1" if possible.
func (v *View) mergeLines(y int) error {
	v.tainted = true

	_, y, err := v.realPosition(0, y)
	if err != nil {
		return err
	}

	if y < 0 || y >= len(v.lines) {
		return errors.New("invalid point")
	}

	if y < len(v.lines)-1 { // otherwise we don't need to merge anything
		v.lines[y] = append(v.lines[y], v.lines[y+1]...)
		v.lines = append(v.lines[:y+1], v.lines[y+2:]...)
	}
	return nil
}

// breakLine breaks a line of the internal buffer at the position corresponding
// to the point (x, y).
func (v *View) breakLine(x, y int) error {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return err
	}

	if y < 0 || y >= len(v.lines) {
		return errors.New("invalid point")
	}

	var left, right []cell
	if x < len(v.lines[y]) { // break line
		left = make([]cell, len(v.lines[y][:x]))
		copy(left, v.lines[y][:x])
		right = make([]cell, len(v.lines[y][x:]))
		copy(right, v.lines[y][x:])
	} else { // new empty line
		left = v.lines[y]
	}

	lines := make([][]cell, len(v.lines)+1)
	lines[y] = left
	lines[y+1] = right
	copy(lines, v.lines[:y])
	copy(lines[y+2:], v.lines[y+1:])
	v.lines = lines
	return nil
}
