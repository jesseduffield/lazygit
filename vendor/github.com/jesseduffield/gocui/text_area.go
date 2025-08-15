package gocui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	WHITESPACES     = " \t"
	WORD_SEPARATORS = "*?_+-.[]~=/&;!#$%^(){}<>"
)

type CursorMapping struct {
	Orig    int
	Wrapped int
}

type TextArea struct {
	content        []rune
	wrappedContent []rune
	cursorMapping  []CursorMapping
	cursor         int
	overwrite      bool
	clipboard      string
	AutoWrap       bool
	AutoWrapWidth  int
}

func AutoWrapContent(content []rune, autoWrapWidth int) ([]rune, []CursorMapping) {
	estimatedNumberOfSoftLineBreaks := len(content) / autoWrapWidth
	cursorMapping := make([]CursorMapping, 0, estimatedNumberOfSoftLineBreaks)
	wrappedContent := make([]rune, 0, len(content)+estimatedNumberOfSoftLineBreaks)
	startOfLine := 0
	indexOfLastWhitespace := -1

	for currentPos, r := range content {
		if r == '\n' {
			wrappedContent = append(wrappedContent, content[startOfLine:currentPos+1]...)
			startOfLine = currentPos + 1
			indexOfLastWhitespace = -1
		} else {
			if r == ' ' {
				indexOfLastWhitespace = currentPos + 1
			} else if currentPos-startOfLine >= autoWrapWidth && indexOfLastWhitespace >= 0 {
				wrapAt := indexOfLastWhitespace
				wrappedContent = append(wrappedContent, content[startOfLine:wrapAt]...)
				wrappedContent = append(wrappedContent, '\n')
				cursorMapping = append(cursorMapping, CursorMapping{wrapAt, len(wrappedContent)})
				startOfLine = wrapAt
				indexOfLastWhitespace = -1
			}
		}
	}

	wrappedContent = append(wrappedContent, content[startOfLine:]...)

	return wrappedContent, cursorMapping
}

func (self *TextArea) autoWrapContent() {
	if self.AutoWrap {
		self.wrappedContent, self.cursorMapping = AutoWrapContent(self.content, self.AutoWrapWidth)
	} else {
		self.wrappedContent, self.cursorMapping = self.content, []CursorMapping{}
	}
}

func (self *TextArea) TypeRune(r rune) {
	if self.overwrite && !self.atEnd() {
		self.content[self.cursor] = r
	} else {
		self.content = append(
			self.content[:self.cursor],
			append([]rune{r}, self.content[self.cursor:]...)...,
		)
	}
	self.autoWrapContent()

	self.cursor++
}

func (self *TextArea) BackSpaceChar() {
	if self.cursor == 0 {
		return
	}

	self.content = append(self.content[:self.cursor-1], self.content[self.cursor:]...)
	self.autoWrapContent()
	self.cursor--
}

func (self *TextArea) DeleteChar() {
	if self.atEnd() {
		return
	}

	self.content = append(self.content[:self.cursor], self.content[self.cursor+1:]...)
	self.autoWrapContent()
}

func (self *TextArea) MoveCursorLeft() {
	if self.cursor == 0 {
		return
	}

	self.cursor--
}

func (self *TextArea) MoveCursorRight() {
	if self.cursor == len(self.content) {
		return
	}

	self.cursor++
}

func (self *TextArea) MoveLeftWord() {
	if self.cursor == 0 {
		return
	}
	if self.atLineStart() {
		self.cursor--
		return
	}

	for !self.atLineStart() && strings.ContainsRune(WHITESPACES, self.content[self.cursor-1]) {
		self.cursor--
	}
	separators := false
	for !self.atLineStart() && strings.ContainsRune(WORD_SEPARATORS, self.content[self.cursor-1]) {
		self.cursor--
		separators = true
	}
	if !separators {
		for !self.atLineStart() && !strings.ContainsRune(WHITESPACES+WORD_SEPARATORS, self.content[self.cursor-1]) {
			self.cursor--
		}
	}
}

func (self *TextArea) MoveRightWord() {
	if self.atEnd() {
		return
	}
	if self.atLineEnd() {
		self.cursor++
		return
	}

	for !self.atLineEnd() && strings.ContainsRune(WHITESPACES, self.content[self.cursor]) {
		self.cursor++
	}
	separators := false
	for !self.atLineEnd() && strings.ContainsRune(WORD_SEPARATORS, self.content[self.cursor]) {
		self.cursor++
		separators = true
	}
	if !separators {
		for !self.atLineEnd() && !strings.ContainsRune(WHITESPACES+WORD_SEPARATORS, self.content[self.cursor]) {
			self.cursor++
		}
	}
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
	return string(self.wrappedContent)
}

func (self *TextArea) GetUnwrappedContent() string {
	return string(self.content)
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

		self.content = append(self.content[:self.cursor-1], self.content[self.cursor:]...)
		self.cursor--
		self.autoWrapContent()
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
	self.clipboard = string(self.content[newlineIndex+1 : self.cursor])
	self.content = append(self.content[:newlineIndex+1], self.content[self.cursor:]...)
	self.autoWrapContent()
	self.cursor = newlineIndex + 1
}

func (self *TextArea) DeleteToEndOfLine() {
	if self.atEnd() {
		return
	}

	// if we're at the end of the line, delete just the newline character
	if self.atLineEnd() {
		self.content = append(self.content[:self.cursor], self.content[self.cursor+1:]...)
		self.autoWrapContent()
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
	self.clipboard = string(self.content[self.cursor:lineEndIndex])
	self.content = append(self.content[:self.cursor], self.content[lineEndIndex:]...)
	self.autoWrapContent()
}

func (self *TextArea) GoToStartOfLine() {
	if self.atSoftLineStart() {
		return
	}

	// otherwise, you delete everything up to the start of the current line, without
	// deleting the newline character
	newlineIndex := self.closestNewlineOnLeft()
	self.cursor = newlineIndex + 1
}

func (self *TextArea) closestNewlineOnLeft() int {
	wrappedCursor := self.origCursorToWrappedCursor(self.cursor)

	newlineIndex := -1

	for i, r := range self.wrappedContent[0:wrappedCursor] {
		if r == '\n' {
			newlineIndex = i
		}
	}

	unwrappedNewlineIndex := self.wrappedCursorToOrigCursor(newlineIndex)
	if unwrappedNewlineIndex >= 0 && self.content[unwrappedNewlineIndex] != '\n' {
		unwrappedNewlineIndex--
	}
	return unwrappedNewlineIndex
}

func (self *TextArea) GoToEndOfLine() {
	if self.atEnd() {
		return
	}

	self.cursor = self.closestNewlineOnRight()

	self.moveLeftFromSoftLineBreak()
}

func (self *TextArea) closestNewlineOnRight() int {
	wrappedCursor := self.origCursorToWrappedCursor(self.cursor)

	for i, r := range self.wrappedContent[wrappedCursor:] {
		if r == '\n' {
			return self.wrappedCursorToOrigCursor(wrappedCursor + i)
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

func (self *TextArea) atSoftLineStart() bool {
	wrappedCursor := self.origCursorToWrappedCursor(self.cursor)
	return wrappedCursor == 0 ||
		(len(self.wrappedContent) > wrappedCursor-1 && self.wrappedContent[wrappedCursor-1] == '\n')
}

func (self *TextArea) atLineEnd() bool {
	return self.atEnd() ||
		(len(self.content) > self.cursor && self.content[self.cursor] == '\n')
}

func (self *TextArea) atSoftLineEnd() bool {
	wrappedCursor := self.origCursorToWrappedCursor(self.cursor)
	return wrappedCursor == len(self.wrappedContent) ||
		(len(self.wrappedContent) > wrappedCursor+1 && self.wrappedContent[wrappedCursor+1] == '\n')
}

func (self *TextArea) BackSpaceWord() {
	if self.cursor == 0 {
		return
	}
	if self.atLineStart() {
		self.BackSpaceChar()
		return
	}

	right := self.cursor
	for !self.atLineStart() && strings.ContainsRune(WHITESPACES, self.content[self.cursor-1]) {
		self.cursor--
	}
	separators := false
	for !self.atLineStart() && strings.ContainsRune(WORD_SEPARATORS, self.content[self.cursor-1]) {
		self.cursor--
		separators = true
	}
	if !separators {
		for !self.atLineStart() && !strings.ContainsRune(WHITESPACES+WORD_SEPARATORS, self.content[self.cursor-1]) {
			self.cursor--
		}
	}

	self.clipboard = string(self.content[self.cursor:right])
	self.content = append(self.content[:self.cursor], self.content[right:]...)
	self.autoWrapContent()
}

func (self *TextArea) Yank() {
	self.TypeString(self.clipboard)
}

func origCursorToWrappedCursor(origCursor int, cursorMapping []CursorMapping) int {
	prevMapping := CursorMapping{0, 0}
	for _, mapping := range cursorMapping {
		if origCursor < mapping.Orig {
			break
		}
		prevMapping = mapping
	}

	return origCursor + prevMapping.Wrapped - prevMapping.Orig
}

func (self *TextArea) origCursorToWrappedCursor(origCursor int) int {
	return origCursorToWrappedCursor(origCursor, self.cursorMapping)
}

func wrappedCursorToOrigCursor(wrappedCursor int, cursorMapping []CursorMapping) int {
	prevMapping := CursorMapping{0, 0}
	for _, mapping := range cursorMapping {
		if wrappedCursor < mapping.Wrapped {
			break
		}
		prevMapping = mapping
	}

	return wrappedCursor + prevMapping.Orig - prevMapping.Wrapped
}

func (self *TextArea) wrappedCursorToOrigCursor(wrappedCursor int) int {
	return wrappedCursorToOrigCursor(wrappedCursor, self.cursorMapping)
}

func (self *TextArea) GetCursorXY() (int, int) {
	cursorX := 0
	cursorY := 0
	wrappedCursor := self.origCursorToWrappedCursor(self.cursor)
	for _, r := range self.wrappedContent[0:wrappedCursor] {
		if r == '\n' {
			cursorY++
			cursorX = 0
		} else {
			chWidth := runewidth.RuneWidth(r)
			cursorX += chWidth
		}
	}

	return cursorX, cursorY
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
	for _, r := range self.wrappedContent {
		if x <= 0 && y == 0 {
			self.cursor = self.wrappedCursorToOrigCursor(newCursor)
			if self.wrappedContent[newCursor] == '\n' {
				self.moveLeftFromSoftLineBreak()
			}
			return
		}

		if r == '\n' {
			if y == 0 {
				self.cursor = self.wrappedCursorToOrigCursor(newCursor)
				self.moveLeftFromSoftLineBreak()
				return
			}
			y--
		} else if y == 0 {
			chWidth := runewidth.RuneWidth(r)
			x -= chWidth
		}

		newCursor++
	}

	// if we weren't able to run-down our arg, the user is trying to move out of
	// bounds so we'll just return
	if y > 0 {
		return
	}

	self.cursor = self.wrappedCursorToOrigCursor(newCursor)
}

func (self *TextArea) Clear() {
	self.content = []rune{}
	self.wrappedContent = []rune{}
	self.cursor = 0
}

func (self *TextArea) TypeString(str string) {
	for _, r := range str {
		self.TypeRune(r)
	}
}
