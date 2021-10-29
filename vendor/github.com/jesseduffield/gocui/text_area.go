package gocui

import "github.com/mattn/go-runewidth"

type TextArea struct {
	content   []rune
	cursor    int
	overwrite bool
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

	self.cursor++
}

func (self *TextArea) BackSpaceChar() {
	if self.cursor == 0 {
		return
	}

	self.content = append(self.content[:self.cursor-1], self.content[self.cursor:]...)
	self.cursor--
}

func (self *TextArea) DeleteChar() {
	if self.atEnd() {
		return
	}

	self.content = append(self.content[:self.cursor], self.content[self.cursor+1:]...)
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

func (self *TextArea) MoveCursorUp() {
	x, y := self.GetCursorXY()
	self.SetCursor2D(x, y-1)
}

func (self *TextArea) MoveCursorDown() {
	x, y := self.GetCursorXY()
	self.SetCursor2D(x, y+1)
}

func (self *TextArea) GetContent() string {
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
		return
	}

	// otherwise, you delete everything up to the start of the current line, without
	// deleting the newline character
	newlineIndex := self.closestNewlineOnLeft()
	self.content = append(self.content[:newlineIndex+1], self.content[self.cursor:]...)
	self.cursor = newlineIndex + 1
}

func (self *TextArea) GoToStartOfLine() {
	if self.atLineStart() {
		return
	}

	// otherwise, you delete everything up to the start of the current line, without
	// deleting the newline character
	newlineIndex := self.closestNewlineOnLeft()
	self.cursor = newlineIndex + 1
}

func (self *TextArea) closestNewlineOnLeft() int {
	newlineIndex := -1

	for i, r := range self.content[0:self.cursor] {
		if r == '\n' {
			newlineIndex = i
		}
	}

	return newlineIndex
}

func (self *TextArea) GoToEndOfLine() {
	if self.atEnd() {
		return
	}

	self.cursor = self.closestNewlineOnRight()
}

func (self *TextArea) closestNewlineOnRight() int {
	for i, r := range self.content[self.cursor:] {
		if r == '\n' {
			return self.cursor + i
		}
	}

	return len(self.content)
}

func (self *TextArea) atLineStart() bool {
	return self.cursor == 0 ||
		(len(self.content) > self.cursor-1 && self.content[self.cursor-1] == '\n')
}

func (self *TextArea) GetCursorXY() (int, int) {
	cursorX := 0
	cursorY := 0
	for _, r := range self.content[0:self.cursor] {
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
	for _, r := range self.content {
		if x <= 0 && y == 0 {
			self.cursor = newCursor
			return
		}

		if r == '\n' {
			if y == 0 {
				self.cursor = newCursor
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

	self.cursor = newCursor
}

func (self *TextArea) Clear() {
	self.content = []rune{}
	self.cursor = 0
}

func (self *TextArea) TypeString(str string) {
	for _, r := range str {
		self.TypeRune(r)
	}
}
