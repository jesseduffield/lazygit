package gocui

import (
	"strings"

	"github.com/rivo/uniseg"
)

// Vim-style motions and range edits for TextArea. These are the primitives
// that VimEditor composes into normal-mode commands; they are pure with
// respect to the cursor (they return a new cursor position) unless the name
// says otherwise.
//
// Character classes follow vim's small-word rules: whitespace, punctuation
// (WORD_SEPARATORS), and word characters are distinct classes, and a motion
// stops at every class boundary. Newlines are treated as whitespace; vim's
// "an empty line is a word" special case is deliberately not implemented, as
// it doesn't pull its weight in single-message editing.

const (
	charClassWhitespace = iota
	charClassSeparator
	charClassWord
)

func vimCharClass(ch string) int {
	if ch == "\n" || strings.Contains(WHITESPACES, ch) {
		return charClassWhitespace
	}
	if strings.Contains(WORD_SEPARATORS, ch) {
		return charClassSeparator
	}
	return charClassWord
}

// cells excluding the synthetic soft-line-break cells, i.e. one cell per
// grapheme of real content, with strictly increasing contentIndex.
func (self *TextArea) realCells() []TextAreaCell {
	result := make([]TextAreaCell, 0, len(self.cells))
	for i, cell := range self.cells {
		if cell.char == "\n" && self.isSoftLineBreak(i) {
			continue
		}
		result = append(result, cell)
	}
	return result
}

func realCellIndex(cells []TextAreaCell, cursor int) int {
	for i, cell := range cells {
		if cell.contentIndex >= cursor {
			return i
		}
	}
	return len(cells)
}

func (self *TextArea) vimWordForward(cursor int, count int) int {
	cells := self.realCells()
	i := realCellIndex(cells, cursor)
	for range count {
		if i >= len(cells) {
			break
		}
		class := vimCharClass(cells[i].char)
		if class != charClassWhitespace {
			for i < len(cells) && vimCharClass(cells[i].char) == class {
				i++
			}
		}
		for i < len(cells) && vimCharClass(cells[i].char) == charClassWhitespace {
			i++
		}
	}
	if i >= len(cells) {
		return len(self.content)
	}
	return cells[i].contentIndex
}

func (self *TextArea) vimWordEnd(cursor int, count int) int {
	cells := self.realCells()
	i := realCellIndex(cells, cursor)
	for range count {
		i++
		for i < len(cells) && vimCharClass(cells[i].char) == charClassWhitespace {
			i++
		}
		if i >= len(cells) {
			break
		}
		class := vimCharClass(cells[i].char)
		for i+1 < len(cells) && vimCharClass(cells[i+1].char) == class {
			i++
		}
	}
	if i >= len(cells) {
		i = len(cells) - 1
	}
	if i < 0 {
		return cursor
	}
	return cells[i].contentIndex
}

func (self *TextArea) vimWordBack(cursor int, count int) int {
	cells := self.realCells()
	i := realCellIndex(cells, cursor)
	for range count {
		i--
		for i >= 0 && vimCharClass(cells[i].char) == charClassWhitespace {
			i--
		}
		if i < 0 {
			break
		}
		class := vimCharClass(cells[i].char)
		for i > 0 && vimCharClass(cells[i-1].char) == class {
			i--
		}
	}
	if i < 0 {
		return 0
	}
	return cells[i].contentIndex
}

func (self *TextArea) hardLineStart(cursor int) int {
	return strings.LastIndexByte(self.content[:cursor], '\n') + 1
}

func (self *TextArea) hardLineEnd(cursor int) int {
	if idx := strings.IndexByte(self.content[cursor:], '\n'); idx != -1 {
		return cursor + idx
	}
	return len(self.content)
}

func (self *TextArea) vimFirstNonBlank(cursor int) int {
	start := self.hardLineStart(cursor)
	end := self.hardLineEnd(cursor)
	for i := start; i < end; i++ {
		if !strings.Contains(WHITESPACES, string(self.content[i])) {
			return i
		}
	}
	return start
}

// start of the last grapheme of the line containing cursor; the line start
// itself if the line is empty
func (self *TextArea) vimLastCharOfLine(cursor int) int {
	start := self.hardLineStart(cursor)
	end := self.hardLineEnd(cursor)
	if start == end {
		return start
	}
	return self.prevGraphemeStart(end)
}

// start of the grapheme preceding cursor (0 if there is none)
func (self *TextArea) prevGraphemeStart(cursor int) int {
	cells := self.realCells()
	i := realCellIndex(cells, cursor)
	if i == 0 {
		return 0
	}
	return cells[i-1].contentIndex
}

// start of the grapheme following cursor (end of content if there is none)
func (self *TextArea) nextGraphemeStart(cursor int) int {
	if cursor >= len(self.content) {
		return len(self.content)
	}
	s, _, _, _ := uniseg.FirstGraphemeClusterInString(self.content[cursor:], -1)
	return cursor + len(s)
}

func (self *TextArea) graphemeAt(cursor int) string {
	if cursor >= len(self.content) {
		return ""
	}
	s, _, _, _ := uniseg.FirstGraphemeClusterInString(self.content[cursor:], -1)
	return s
}

// vimFindOnLine returns the cursor for f/F/t/T motions, or -1 if the target
// does not occur count times on the cursor's line.
func (self *TextArea) vimFindOnLine(cursor int, target string, forward bool, till bool, count int) int {
	cells := self.realCells()
	i := realCellIndex(cells, cursor)
	lineStart := self.hardLineStart(cursor)
	lineEnd := self.hardLineEnd(cursor)

	step := -1
	if forward {
		step = 1
	}
	found := i
	for range count {
		j := found + step
		for {
			if j < 0 || j >= len(cells) {
				return -1
			}
			idx := cells[j].contentIndex
			if idx < lineStart || idx >= lineEnd {
				return -1
			}
			if cells[j].char == target {
				break
			}
			j += step
		}
		found = j
	}
	if till {
		found -= step
	}
	return cells[found].contentIndex
}

// line number (0-based) and total line count, both over hard lines
func (self *TextArea) vimLineNumber(cursor int) int {
	return strings.Count(self.content[:cursor], "\n")
}

func (self *TextArea) vimLineCount() int {
	return strings.Count(self.content, "\n") + 1
}

// start of the (0-based) line'th hard line, clamped to the last line
func (self *TextArea) vimGotoLine(line int) int {
	cursor := 0
	for range line {
		idx := strings.IndexByte(self.content[cursor:], '\n')
		if idx == -1 {
			break
		}
		cursor += idx + 1
	}
	return cursor
}

// vimWordBounds returns the [start, end) range of the word (or whitespace
// run) under cursor. With around=true it extends over trailing whitespace,
// or leading whitespace if there is none trailing, per vim's aw.
func (self *TextArea) vimWordBounds(cursor int, around bool) (int, int) {
	cells := self.realCells()
	if len(cells) == 0 {
		return 0, 0
	}
	i := realCellIndex(cells, cursor)
	if i >= len(cells) {
		i = len(cells) - 1
	}
	if cells[i].char == "\n" {
		return cells[i].contentIndex, cells[i].contentIndex
	}

	class := vimCharClass(cells[i].char)
	start := i
	for start > 0 && cells[start-1].char != "\n" && vimCharClass(cells[start-1].char) == class {
		start--
	}
	end := i
	for end+1 < len(cells) && cells[end+1].char != "\n" && vimCharClass(cells[end+1].char) == class {
		end++
	}

	if around && class != charClassWhitespace {
		trailed := false
		for end+1 < len(cells) && cells[end+1].char != "\n" && vimCharClass(cells[end+1].char) == charClassWhitespace {
			end++
			trailed = true
		}
		if !trailed {
			for start > 0 && cells[start-1].char != "\n" && vimCharClass(cells[start-1].char) == charClassWhitespace {
				start--
			}
		}
	}

	startIdx := cells[start].contentIndex
	if end+1 < len(cells) {
		return startIdx, cells[end+1].contentIndex
	}
	return startIdx, len(self.content)
}

func (self *TextArea) getRange(start, end int) string {
	start = max(start, 0)
	end = min(end, len(self.content))
	if start >= end {
		return ""
	}
	return self.content[start:end]
}

func (self *TextArea) deleteRange(start, end int) {
	start = max(start, 0)
	end = min(end, len(self.content))
	if start >= end {
		return
	}
	self.content = self.content[:start] + self.content[end:]
	self.cursor = start
	self.updateCells()
}

func (self *TextArea) insertAt(cursor int, str string) {
	self.cursor = min(max(cursor, 0), len(self.content))
	self.TypeString(str)
}

func (self *TextArea) replaceGraphemeAt(cursor int, ch string) {
	old := self.graphemeAt(cursor)
	if old == "" || old == "\n" {
		return
	}
	self.content = self.content[:cursor] + ch + self.content[cursor+len(old):]
	self.cursor = cursor
	self.updateCells()
}

func (self *TextArea) setState(content string, cursor int) {
	self.content = content
	self.updateCells()
	self.cursor = min(max(cursor, 0), len(content))
}

func (self *TextArea) SetCursor(cursor int) {
	self.cursor = min(max(cursor, 0), len(self.content))
}

func (self *TextArea) GetCursor() int {
	return self.cursor
}
