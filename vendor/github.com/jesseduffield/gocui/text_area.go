package gocui

import (
	"regexp"
	"slices"
	"strings"

	"github.com/rivo/uniseg"
)

const (
	WHITESPACES     = " \t"
	WORD_SEPARATORS = "*?_+-.[]~=/&;!#$%^(){}<>"
)

type TextAreaCell struct {
	char         string // string because it could be a multi-rune grapheme cluster
	width        int
	x, y         int // cell coordinates
	contentIndex int // byte index into the original content
}

// returns the cursor x,y position after this cell
func (c *TextAreaCell) nextCursorXY() (int, int) {
	if c.char == "\n" {
		return 0, c.y + 1
	}
	return c.x + c.width, c.y
}

type TextArea struct {
	content       string
	cells         []TextAreaCell
	cursor        int // position in content, as an index into the byte array
	overwrite     bool
	clipboard     string
	AutoWrap      bool
	AutoWrapWidth int
}

func stringToTextAreaCells(str string) []TextAreaCell {
	result := make([]TextAreaCell, 0, len(str))

	contentIndex := 0
	state := -1
	for len(str) > 0 {
		var c string
		var w int
		c, str, w, state = uniseg.FirstGraphemeClusterInString(str, state)
		// only set char, width, and contentIndex; x and y will be set later
		result = append(result, TextAreaCell{char: c, width: w, contentIndex: contentIndex})
		contentIndex += len(c)
	}
	return result
}

// Returns the indices in content where soft line breaks occur due to auto-wrapping to the given width.
func AutoWrapContent(content string, autoWrapWidth int) []int {
	_, softLineBreakIndices := contentToCells(content, autoWrapWidth)
	return softLineBreakIndices
}

func contentToCells(content string, autoWrapWidth int) ([]TextAreaCell, []int) {
	estimatedNumberOfSoftLineBreaks := 0
	if autoWrapWidth > 0 {
		estimatedNumberOfSoftLineBreaks = len(content) / autoWrapWidth
	}
	softLineBreakIndices := make([]int, 0, estimatedNumberOfSoftLineBreaks)
	result := make([]TextAreaCell, 0, len(content)+estimatedNumberOfSoftLineBreaks)
	startOfLine := 0
	currentLineWidth := 0
	indexOfLastWhitespace := -1
	var footNoteMatcher footNoteMatcher

	cells := stringToTextAreaCells(content)
	y := 0

	appendCellsSinceLineStart := func(to int) {
		x := 0
		for i := startOfLine; i < to; i++ {
			cells[i].x = x
			cells[i].y = y
			x += cells[i].width
		}

		result = append(result, cells[startOfLine:to]...)
	}

	for currentPos, c := range cells {
		if c.char == "\n" {
			appendCellsSinceLineStart(currentPos + 1)
			y++
			startOfLine = currentPos + 1
			indexOfLastWhitespace = -1
			currentLineWidth = 0
			footNoteMatcher.reset()
		} else {
			currentLineWidth += c.width
			if c.char == " " && !footNoteMatcher.isFootNote() {
				indexOfLastWhitespace = currentPos + 1
			} else if autoWrapWidth > 0 && currentLineWidth > autoWrapWidth && indexOfLastWhitespace >= 0 {
				wrapAt := indexOfLastWhitespace
				appendCellsSinceLineStart(wrapAt)
				contentIndex := cells[wrapAt].contentIndex
				y++
				result = append(result, TextAreaCell{char: "\n", width: 1, contentIndex: contentIndex, x: 0, y: y})
				softLineBreakIndices = append(softLineBreakIndices, contentIndex)
				startOfLine = wrapAt
				indexOfLastWhitespace = -1
				currentLineWidth = 0
				for _, c1 := range cells[startOfLine : currentPos+1] {
					currentLineWidth += c1.width
				}
				footNoteMatcher.reset()
			}

			footNoteMatcher.addCharacter(c.char)
		}
	}

	appendCellsSinceLineStart(len(cells))

	return result, softLineBreakIndices
}

var footNoteRe = regexp.MustCompile(`^\[\d+\]:\s*$`)

type footNoteMatcher struct {
	lineStr        strings.Builder
	didFailToMatch bool
}

func (self *footNoteMatcher) addCharacter(chr string) {
	if self.didFailToMatch {
		// don't bother tracking the rune if we know it can't possibly match any more
		return
	}

	if self.lineStr.Len() == 0 && chr != "[" {
		// fail early if the first rune of a line isn't a '['; this is mainly to avoid a (possibly
		// expensive) regex match
		self.didFailToMatch = true
		return
	}

	self.lineStr.WriteString(chr)
}

func (self *footNoteMatcher) isFootNote() bool {
	if self.didFailToMatch {
		return false
	}

	if footNoteRe.MatchString(self.lineStr.String()) {
		// it's a footnote, so treat spaces as non-breaking. It's important not to reset the matcher
		// here, because there could be multiple spaces after a footnote.
		return true
	}

	// no need to check again for this line
	self.didFailToMatch = true
	return false
}

func (self *footNoteMatcher) reset() {
	self.lineStr.Reset()
	self.didFailToMatch = false
}

func (self *TextArea) updateCells() {
	width := self.AutoWrapWidth
	if !self.AutoWrap {
		width = -1
	}

	self.cells, _ = contentToCells(self.content, width)
}

func (self *TextArea) typeCharacter(ch string) {
	widthToDelete := 0
	if self.overwrite && !self.atEnd() {
		s, _, _, _ := uniseg.FirstGraphemeClusterInString(self.content[self.cursor:], -1)
		widthToDelete = len(s)
	}

	self.content = self.content[:self.cursor] + ch + self.content[self.cursor+widthToDelete:]
	self.cursor += len(ch)
}

func (self *TextArea) TypeCharacter(ch string) {
	self.typeCharacter(ch)
	self.updateCells()
}

func (self *TextArea) BackSpaceChar() {
	if self.cursor == 0 {
		return
	}

	cellCursor := self.contentCursorToCellCursor(self.cursor)
	widthToDelete := len(self.cells[cellCursor-1].char)

	oldCursor := self.cursor
	self.cursor -= widthToDelete
	self.content = self.content[:self.cursor] + self.content[oldCursor:]

	self.updateCells()
}

func (self *TextArea) DeleteChar() {
	if self.atEnd() {
		return
	}

	s, _, _, _ := uniseg.FirstGraphemeClusterInString(self.content[self.cursor:], -1)
	widthToDelete := len(s)
	self.content = self.content[:self.cursor] + self.content[self.cursor+widthToDelete:]
	self.updateCells()
}

func (self *TextArea) MoveCursorLeft() {
	if self.cursor == 0 {
		return
	}

	cellCursor := self.contentCursorToCellCursor(self.cursor)
	self.cursor -= len(self.cells[cellCursor-1].char)
}

func (self *TextArea) MoveCursorRight() {
	if self.cursor == len(self.content) {
		return
	}

	s, _, _, _ := uniseg.FirstGraphemeClusterInString(self.content[self.cursor:], -1)
	self.cursor += len(s)
}

func (self *TextArea) newCursorForMoveLeftWord() int {
	if self.cursor == 0 {
		return 0
	}
	if self.atLineStart() {
		return self.cursor - 1
	}

	cellCursor := self.contentCursorToCellCursor(self.cursor)
	for cellCursor > 0 && (self.isSoftLineBreak(cellCursor-1) || strings.Contains(WHITESPACES, self.cells[cellCursor-1].char)) {
		cellCursor--
	}
	separators := false
	for cellCursor > 0 && strings.Contains(WORD_SEPARATORS, self.cells[cellCursor-1].char) {
		cellCursor--
		separators = true
	}
	if !separators {
		for cellCursor > 0 && self.cells[cellCursor-1].char != "\n" && !strings.Contains(WHITESPACES+WORD_SEPARATORS, self.cells[cellCursor-1].char) {
			cellCursor--
		}
	}

	return self.cellCursorToContentCursor(cellCursor)
}

func (self *TextArea) MoveLeftWord() {
	self.cursor = self.newCursorForMoveLeftWord()
}

func (self *TextArea) MoveRightWord() {
	if self.atEnd() {
		return
	}
	if self.atLineEnd() {
		self.cursor++
		return
	}

	cellCursor := self.contentCursorToCellCursor(self.cursor)
	for cellCursor < len(self.cells) && (self.isSoftLineBreak(cellCursor) || strings.Contains(WHITESPACES, self.cells[cellCursor].char)) {
		cellCursor++
	}
	separators := false
	for cellCursor < len(self.cells) && strings.Contains(WORD_SEPARATORS, self.cells[cellCursor].char) {
		cellCursor++
		separators = true
	}
	if !separators {
		for cellCursor < len(self.cells) && self.cells[cellCursor].char != "\n" && !strings.Contains(WHITESPACES+WORD_SEPARATORS, self.cells[cellCursor].char) {
			cellCursor++
		}
	}

	self.cursor = self.cellCursorToContentCursor(cellCursor)
}

func (self *TextArea) MoveCursorUp() {
	x, y := self.GetCursorXY()
	self.SetCursor2D(x, y-1)
}

func (self *TextArea) MoveCursorDown() {
	x, y := self.GetCursorXY()
	self.SetCursor2D(x, y+1)
}

func (self *TextArea) GetContent() string {
	var b strings.Builder
	for _, c := range self.cells {
		b.WriteString(c.char)
	}
	return b.String()
}

func (self *TextArea) GetUnwrappedContent() string {
	return self.content
}

func (self *TextArea) ToggleOverwrite() {
	self.overwrite = !self.overwrite
}

func (self *TextArea) atEnd() bool {
	return self.cursor == len(self.content)
}

func (self *TextArea) DeleteToStartOfLine() {
	// copying vim's logic: if you're at the start of the line, you delete the newline
	// character and go to the end of the previous line
	if self.atLineStart() {
		if self.cursor == 0 {
			return
		}

		self.content = self.content[:self.cursor-1] + self.content[self.cursor:]
		self.cursor--
		self.updateCells()
		return
	}

	// otherwise, if we're at a soft line start, skip left past the soft line
	// break, so we'll end up deleting the previous line. This seems like the
	// only reasonable behavior in this case, as you can't delete just the soft
	// line break.
	if self.atSoftLineStart() {
		self.cursor--
	}

	// otherwise, you delete everything up to the start of the current line, without
	// deleting the newline character
	newlineIndex := self.closestNewlineOnLeft()
	self.clipboard = self.content[newlineIndex+1 : self.cursor]
	self.content = self.content[:newlineIndex+1] + self.content[self.cursor:]
	self.updateCells()
	self.cursor = newlineIndex + 1
}

func (self *TextArea) DeleteToEndOfLine() {
	if self.atEnd() {
		return
	}

	// if we're at the end of the line, delete just the newline character
	if self.atLineEnd() {
		self.content = self.content[:self.cursor] + self.content[self.cursor+1:]
		self.updateCells()
		return
	}

	// otherwise, if we're at a soft line end, skip right past the soft line
	// break, so we'll end up deleting the next line. This seems like the
	// only reasonable behavior in this case, as you can't delete just the soft
	// line break.
	if self.atSoftLineEnd() {
		self.cursor++
	}

	lineEndIndex := self.closestNewlineOnRight()
	self.clipboard = self.content[self.cursor:lineEndIndex]
	self.content = self.content[:self.cursor] + self.content[lineEndIndex:]
	self.updateCells()
}

func (self *TextArea) GoToStartOfLine() {
	if self.atSoftLineStart() {
		return
	}

	newlineIndex := self.closestNewlineOnLeft()
	self.cursor = newlineIndex + 1
}

func (self *TextArea) closestNewlineOnLeft() int {
	cellCursor := self.contentCursorToCellCursor(self.cursor)

	newlineCellIndex := -1

	for i, c := range self.cells[0:cellCursor] {
		if c.char == "\n" {
			newlineCellIndex = i
		}
	}

	if newlineCellIndex == -1 {
		return -1
	}

	newlineContentIndex := self.cells[newlineCellIndex].contentIndex
	if self.content[newlineContentIndex] != '\n' {
		newlineContentIndex--
	}
	return newlineContentIndex
}

func (self *TextArea) GoToEndOfLine() {
	if self.atEnd() {
		return
	}

	self.cursor = self.closestNewlineOnRight()

	self.moveLeftFromSoftLineBreak()
}

func (self *TextArea) closestNewlineOnRight() int {
	cellCursor := self.contentCursorToCellCursor(self.cursor)

	for i, c := range self.cells[cellCursor:] {
		if c.char == "\n" {
			return self.cellCursorToContentCursor(cellCursor + i)
		}
	}

	return len(self.content)
}

func (self *TextArea) moveLeftFromSoftLineBreak() {
	// If the end of line is a soft line break, we need to move left by one so
	// that we end up at the last whitespace before the line break. Otherwise
	// we'd be at the start of the next line, since the newline character
	// doesn't really exist in the real content.
	if self.cursor < len(self.content) && self.content[self.cursor] != '\n' {
		self.cursor--
	}
}

func (self *TextArea) atLineStart() bool {
	return self.cursor == 0 ||
		(len(self.content) > self.cursor-1 && self.content[self.cursor-1] == '\n')
}

func (self *TextArea) isSoftLineBreak(cellCursor int) bool {
	cell := self.cells[cellCursor]
	return cell.char == "\n" && self.content[cell.contentIndex] != '\n'
}

func (self *TextArea) atSoftLineStart() bool {
	cellCursor := self.contentCursorToCellCursor(self.cursor)
	return cellCursor == 0 ||
		(len(self.cells) > cellCursor-1 && self.cells[cellCursor-1].char == "\n")
}

func (self *TextArea) atLineEnd() bool {
	return self.atEnd() ||
		(len(self.content) > self.cursor && self.content[self.cursor] == '\n')
}

func (self *TextArea) atSoftLineEnd() bool {
	cellCursor := self.contentCursorToCellCursor(self.cursor)
	return cellCursor == len(self.cells) ||
		(len(self.cells) > cellCursor+1 && self.cells[cellCursor+1].char == "\n")
}

func (self *TextArea) BackSpaceWord() {
	newCursor := self.newCursorForMoveLeftWord()
	if newCursor == self.cursor {
		return
	}

	clipboard := self.content[newCursor:self.cursor]
	if clipboard != "\n" {
		self.clipboard = clipboard
	}
	self.content = self.content[:newCursor] + self.content[self.cursor:]
	self.cursor = newCursor
	self.updateCells()
}

func (self *TextArea) Yank() {
	self.TypeString(self.clipboard)
}

func (self *TextArea) contentCursorToCellCursor(origCursor int) int {
	idx, _ := slices.BinarySearchFunc(self.cells, origCursor, func(cell TextAreaCell, cursor int) int {
		return cell.contentIndex - cursor
	})
	for idx < len(self.cells)-1 && self.cells[idx+1].contentIndex == origCursor {
		idx++
	}
	return idx
}

func (self *TextArea) cellCursorToContentCursor(cellCursor int) int {
	if cellCursor >= len(self.cells) {
		return len(self.content)
	}

	return self.cells[cellCursor].contentIndex
}

func (self *TextArea) GetCursorXY() (int, int) {
	if len(self.cells) == 0 {
		return 0, 0
	}
	cellCursor := self.contentCursorToCellCursor(self.cursor)
	if cellCursor >= len(self.cells) {
		return self.cells[len(self.cells)-1].nextCursorXY()
	}
	if cellCursor > 0 && self.cells[cellCursor].char == "\n" {
		return self.cells[cellCursor-1].nextCursorXY()
	}
	cell := self.cells[cellCursor]
	return cell.x, cell.y
}

// takes an x,y position and maps it to a 1D cursor position
func (self *TextArea) SetCursor2D(x int, y int) {
	if y < 0 {
		y = 0
	}
	if x < 0 {
		x = 0
	}

	newCursor := 0
	for _, c := range self.cells {
		if x <= 0 && y == 0 {
			self.cursor = self.cellCursorToContentCursor(newCursor)
			if self.cells[newCursor].char == "\n" {
				self.moveLeftFromSoftLineBreak()
			}
			return
		}

		if c.char == "\n" {
			if y == 0 {
				self.cursor = self.cellCursorToContentCursor(newCursor)
				self.moveLeftFromSoftLineBreak()
				return
			}
			y--
		} else if y == 0 {
			x -= c.width
		}

		newCursor++
	}

	// if we weren't able to run-down our arg, the user is trying to move out of
	// bounds so we'll just return
	if y > 0 {
		return
	}

	self.cursor = self.cellCursorToContentCursor(newCursor)
}

func (self *TextArea) Clear() {
	self.content = ""
	self.cells = nil
	self.cursor = 0
}

func (self *TextArea) TypeString(str string) {
	state := -1
	for str != "" {
		var chr string
		chr, str, _, state = uniseg.FirstGraphemeClusterInString(str, state)
		self.typeCharacter(chr)
	}

	self.updateCells()
}
