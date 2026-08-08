package gocui

import "strings"

type VimMode int

const (
	VimModeInsert VimMode = iota
	VimModeNormal
	VimModeVisual
)

// VimRegister is the yank/delete buffer. It is separate from TextArea's
// readline clipboard so that a shared register can be passed to several
// VimEditors, letting text yanked in one editable view be pasted in another
// (e.g. commit summary <-> description).
type VimRegister struct {
	text     string
	linewise bool
}

func (self *VimRegister) set(text string, linewise bool) {
	if text == "" {
		return
	}
	self.text = text
	self.linewise = linewise
}

type vimSnapshot struct {
	content string
	cursor  int
}

const vimUndoLimit = 200

type vimFind struct {
	target  string
	forward bool
	till    bool
}

const (
	motionExclusive = iota
	motionInclusive
)

// VimEditor is a modal Editor: insert mode delegates to SimpleEditor, normal
// mode interprets vim commands against the view's TextArea. One instance
// holds the state for one view.
type VimEditor struct {
	allowMultiline bool
	register       *VimRegister

	mode         VimMode
	count        int
	opCount      int
	op           string
	pending      string
	lastFind     vimFind
	visualAnchor int
	undo         []vimSnapshot
	redo         []vimSnapshot
}

func NewVimEditor(register *VimRegister, allowMultiline bool) *VimEditor {
	if register == nil {
		register = &VimRegister{}
	}
	return &VimEditor{
		allowMultiline: allowMultiline,
		register:       register,
		mode:           VimModeInsert,
	}
}

func (self *VimEditor) Mode() VimMode {
	return self.mode
}

func (self *VimEditor) Reset(mode VimMode) {
	self.mode = mode
	self.clearPending()
	self.undo = nil
	self.redo = nil
}

// HandleEscape must be called by the Escape keybinding of a view using this
// editor, because keybindings fire before the Editor sees a key: leaving
// insert mode has to win over closing the panel. Returns true if the editor
// consumed the key.
func (self *VimEditor) HandleEscape(v *View) bool {
	if !self.escape(v.TextArea) {
		return false
	}
	self.syncSelection(v)
	v.RenderTextArea()
	return true
}

func (self *VimEditor) escape(ta *TextArea) bool {
	if self.mode == VimModeInsert {
		self.enterNormalMode(ta)
		return true
	}
	if self.mode == VimModeVisual {
		self.mode = VimModeNormal
		self.clearPending()
		self.clampCursor(ta)
		return true
	}
	if self.op != "" || self.count != 0 || self.pending != "" {
		self.clearPending()
		return true
	}
	return false
}

func (self *VimEditor) Edit(v *View, key Key) bool {
	if self.mode == VimModeInsert {
		if key.Equals(NewKeyName(KeyEsc)) {
			self.enterNormalMode(v.TextArea)
			v.RenderTextArea()
			return true
		}
		return SimpleEditor(v, key)
	}

	handled := self.normalKey(v.TextArea, key)
	if handled {
		self.syncSelection(v)
		v.RenderTextArea()
	}
	return handled
}

func (self *VimEditor) syncSelection(v *View) {
	if self.mode != VimModeVisual {
		v.SetVimSelection(nil)
		return
	}
	start, end := self.visualRange(v.TextArea)
	v.SetVimSelection(v.TextArea.selectionPositions(start, end))
}

func (self *VimEditor) visualRange(ta *TextArea) (int, int) {
	start, end := self.visualAnchor, ta.cursor
	if end < start {
		start, end = end, start
	}
	return start, ta.nextGraphemeStart(end)
}

func (self *VimEditor) normalKey(ta *TextArea, key Key) bool {
	if key.Equals(NewKeyStrMod("r", ModCtrl)) {
		self.restoreRedo(ta)
		return true
	}
	if key.Equals(NewKeyName(KeyEsc)) {
		return self.escape(ta)
	}

	ch := ""
	if key.Mod() == 0 && key.Str() != "" {
		ch = key.Str()
	}
	if ch == "" {
		switch {
		case key.Equals(NewKeyName(KeyArrowLeft)):
			ch = "h"
		case key.Equals(NewKeyName(KeyArrowRight)):
			ch = "l"
		case key.Equals(NewKeyName(KeyArrowUp)):
			ch = "k"
		case key.Equals(NewKeyName(KeyArrowDown)):
			ch = "j"
		case key.Equals(NewKeyName(KeyHome)):
			ch = "0"
		case key.Equals(NewKeyName(KeyEnd)):
			ch = "$"
		case key.Equals(NewKeyName(KeyDelete)):
			ch = "x"
		case key.Equals(NewKeyName(KeyBackspace)):
			ch = "h"
		case key.Equals(NewKeyName(KeyEnter)):
			ch = "j"
		default:
			return false
		}
	}

	self.normalChar(ta, ch)
	return true
}

// normalChar handles one normal-mode key. Unrecognized characters are
// swallowed: on an editable view a plain character that fell through to a
// global keybinding (q, ...) would be disastrous.
func (self *VimEditor) normalChar(ta *TextArea, ch string) {
	if self.pending != "" {
		self.resolvePending(ta, ch)
		return
	}

	if (ch >= "1" && ch <= "9") || (ch == "0" && self.count != 0) {
		self.count = self.count*10 + int(ch[0]-'0')
		return
	}

	if self.mode == VimModeVisual {
		self.visualChar(ta, ch)
		return
	}

	switch ch {
	case "d", "c", "y":
		self.startOperator(ta, ch)
		return
	case "g", "f", "F", "t", "T":
		self.pending = ch
		return
	case "r":
		if self.op == "" {
			self.pending = ch
			return
		}
	case "i", "a":
		if self.op != "" {
			self.pending = "obj" + ch
			return
		}
	}

	if self.motion(ta, ch) {
		return
	}

	if self.op != "" {
		self.clearPending()
		return
	}

	self.command(ta, ch)
}

// visualChar handles one visual-mode key: operators act on the selection,
// motions move its free end, everything else is swallowed.
func (self *VimEditor) visualChar(ta *TextArea, ch string) {
	switch ch {
	case "v":
		self.mode = VimModeNormal
		self.clearPending()
		self.clampCursor(ta)
	case "o":
		self.visualAnchor, ta.cursor = ta.cursor, self.visualAnchor
	case "d", "x", "c", "s", "y":
		op := ch
		switch ch {
		case "x":
			op = "d"
		case "s":
			op = "c"
		}
		start, end := self.visualRange(ta)
		self.mode = VimModeNormal
		self.op = op
		self.applyRange(ta, start, end)
	case "g", "f", "F", "t", "T":
		self.pending = ch
	default:
		if !self.motion(ta, ch) {
			self.clearPending()
		}
	}
}

func (self *VimEditor) motion(ta *TextArea, ch string) bool {
	cur := ta.cursor
	n := self.repeatCount()

	var target int
	kind := motionExclusive

	switch ch {
	case "h":
		target = cur
		lineStart := ta.hardLineStart(cur)
		for range n {
			if target == lineStart {
				break
			}
			target = ta.prevGraphemeStart(target)
		}
	case "l", " ":
		target = cur
		lineEnd := ta.hardLineEnd(cur)
		for range n {
			next := ta.nextGraphemeStart(target)
			if next > lineEnd {
				break
			}
			target = next
		}
	case "j", "k":
		self.motionUpDown(ta, ch, n)
		return true
	case "w":
		if self.op == "c" && vimCharClass(ta.graphemeAt(cur)) != charClassWhitespace {
			// vim's documented quirk: cw on a word acts like ce
			target = ta.vimWordEnd(cur, n)
			kind = motionInclusive
		} else {
			target = ta.vimWordForward(cur, n)
			if self.op != "" {
				// as an operator target, w stops at the end of the
				// cursor's line rather than eating the newline
				lineEnd := ta.hardLineEnd(cur)
				if target > lineEnd && lineEnd > cur {
					target = lineEnd
				}
			}
		}
	case "e":
		target = ta.vimWordEnd(cur, n)
		kind = motionInclusive
	case "b":
		target = ta.vimWordBack(cur, n)
	case "0":
		target = ta.hardLineStart(cur)
	case "^":
		target = ta.vimFirstNonBlank(cur)
	case "$":
		if self.op != "" {
			target = ta.hardLineEnd(cur)
		} else {
			target = ta.vimLastCharOfLine(cur)
		}
	case "G":
		line := ta.vimLineCount() - 1
		if self.count != 0 || self.opCount != 0 {
			line = n - 1
		}
		self.motionToLine(ta, line)
		return true
	case ";", ",":
		if self.lastFind.target == "" {
			self.clearPending()
			return true
		}
		find := self.lastFind
		if ch == "," {
			find.forward = !find.forward
		}
		self.applyFind(ta, find, n)
		return true
	default:
		return false
	}

	self.applyMotion(ta, target, kind)
	return true
}

func (self *VimEditor) motionUpDown(ta *TextArea, ch string, n int) {
	if self.op != "" {
		delta := n
		if ch == "k" {
			delta = -n
		}
		line := ta.vimLineNumber(ta.cursor)
		self.applyLinewise(ta, line, line+delta)
		return
	}
	for range n {
		if ch == "j" {
			ta.MoveCursorDown()
		} else {
			ta.MoveCursorUp()
		}
	}
	self.clearPending()
	self.clampCursor(ta)
}

func (self *VimEditor) motionToLine(ta *TextArea, line int) {
	target := ta.vimGotoLine(max(line, 0))
	if self.op != "" {
		self.applyLinewise(ta, ta.vimLineNumber(ta.cursor), ta.vimLineNumber(target))
		return
	}
	ta.cursor = ta.vimFirstNonBlank(target)
	self.clearPending()
	self.clampCursor(ta)
}

func (self *VimEditor) applyMotion(ta *TextArea, target, kind int) {
	if target < 0 {
		self.clearPending()
		return
	}
	if self.op == "" {
		ta.cursor = target
		self.clearPending()
		self.clampCursor(ta)
		return
	}

	start, end := ta.cursor, target
	if end < start {
		start, end = end, start
	}
	if kind == motionInclusive {
		end = ta.nextGraphemeStart(end)
	}
	self.applyRange(ta, start, end)
}

func (self *VimEditor) applyRange(ta *TextArea, start, end int) {
	op := self.op
	self.clearPending()

	self.register.set(ta.getRange(start, end), false)
	switch op {
	case "y":
		ta.cursor = start
	case "d":
		self.pushUndo(ta)
		ta.deleteRange(start, end)
	case "c":
		self.pushUndo(ta)
		ta.deleteRange(start, end)
		self.mode = VimModeInsert
	}
	self.clampCursor(ta)
}

func (self *VimEditor) applyLinewise(ta *TextArea, lineA, lineB int) {
	op := self.op
	self.clearPending()

	if lineB < lineA {
		lineA, lineB = lineB, lineA
	}
	lineA = max(lineA, 0)
	lineB = min(lineB, ta.vimLineCount()-1)

	start := ta.vimGotoLine(lineA)
	end := ta.hardLineEnd(ta.vimGotoLine(lineB))

	self.register.set(ta.getRange(start, end), true)
	switch op {
	case "y":
	case "d":
		self.pushUndo(ta)
		delStart, delEnd := start, end
		if delEnd < len(ta.content) {
			delEnd++
		} else if delStart > 0 {
			delStart--
		}
		ta.deleteRange(delStart, delEnd)
		ta.cursor = ta.vimFirstNonBlank(ta.cursor)
	case "c":
		self.pushUndo(ta)
		ta.deleteRange(start, end)
		self.mode = VimModeInsert
	}
	self.clampCursor(ta)
}

func (self *VimEditor) startOperator(ta *TextArea, ch string) {
	if self.op == ch {
		line := ta.vimLineNumber(ta.cursor)
		self.applyLinewise(ta, line, line+self.repeatCount()-1)
		return
	}
	if self.op != "" {
		self.clearPending()
		return
	}
	self.op = ch
	self.opCount = self.count
	self.count = 0
}

func (self *VimEditor) resolvePending(ta *TextArea, ch string) {
	pending := self.pending
	self.pending = ""

	switch pending {
	case "g":
		if ch == "g" {
			line := 0
			if self.count != 0 || self.opCount != 0 {
				line = self.repeatCount() - 1
			}
			self.motionToLine(ta, line)
			return
		}
	case "r":
		self.replaceChars(ta, ch)
		return
	case "f", "F", "t", "T":
		find := vimFind{
			target:  ch,
			forward: pending == "f" || pending == "t",
			till:    pending == "t" || pending == "T",
		}
		self.lastFind = find
		self.applyFind(ta, find, self.repeatCount())
		return
	case "obji", "obja":
		if ch == "w" {
			start, end := ta.vimWordBounds(ta.cursor, pending == "obja")
			self.applyRange(ta, start, end)
			return
		}
	}
	self.clearPending()
}

func (self *VimEditor) applyFind(ta *TextArea, find vimFind, n int) {
	target := ta.vimFindOnLine(ta.cursor, find.target, find.forward, find.till, n)
	if target == -1 {
		self.clearPending()
		return
	}
	kind := motionExclusive
	if find.forward {
		kind = motionInclusive
	}
	self.applyMotion(ta, target, kind)
}

func (self *VimEditor) replaceChars(ta *TextArea, ch string) {
	n := self.repeatCount()
	self.clearPending()
	self.pushUndo(ta)

	pos := ta.cursor
	last := ta.cursor
	for range n {
		g := ta.graphemeAt(pos)
		if g == "" || g == "\n" {
			break
		}
		ta.replaceGraphemeAt(pos, ch)
		last = pos
		pos += len(ch)
	}
	ta.cursor = last
}

func (self *VimEditor) command(ta *TextArea, ch string) {
	n := self.repeatCount()
	self.clearPending()

	switch ch {
	case "i":
		self.startInsert(ta)
	case "I":
		ta.cursor = ta.vimFirstNonBlank(ta.cursor)
		self.startInsert(ta)
	case "a":
		ta.cursor = min(ta.nextGraphemeStart(ta.cursor), ta.hardLineEnd(ta.cursor))
		self.startInsert(ta)
	case "A":
		ta.cursor = ta.hardLineEnd(ta.cursor)
		self.startInsert(ta)
	case "o", "O":
		if !self.allowMultiline {
			return
		}
		self.pushUndo(ta)
		if ch == "o" {
			ta.insertAt(ta.hardLineEnd(ta.cursor), "\n")
		} else {
			start := ta.hardLineStart(ta.cursor)
			ta.insertAt(start, "\n")
			ta.cursor = start
		}
		self.mode = VimModeInsert
	case "x":
		self.deleteForward(ta, n)
	case "X":
		self.deleteBackward(ta, n)
	case "s":
		self.deleteForward(ta, n)
		self.startInsert(ta)
	case "S":
		self.op = "c"
		line := ta.vimLineNumber(ta.cursor)
		self.applyLinewise(ta, line, line+n-1)
	case "D":
		self.op = "d"
		self.applyRange(ta, ta.cursor, ta.hardLineEnd(ta.cursor))
	case "C":
		self.op = "c"
		self.applyRange(ta, ta.cursor, ta.hardLineEnd(ta.cursor))
	case "Y":
		self.op = "y"
		line := ta.vimLineNumber(ta.cursor)
		self.applyLinewise(ta, line, line+n-1)
	case "v":
		self.visualAnchor = ta.cursor
		self.mode = VimModeVisual
	case "p":
		self.paste(ta, true, n)
	case "P":
		self.paste(ta, false, n)
	case "u":
		for range n {
			self.restoreUndo(ta)
		}
	}
}

func (self *VimEditor) deleteForward(ta *TextArea, n int) {
	end := ta.cursor
	lineEnd := ta.hardLineEnd(ta.cursor)
	for range n {
		next := ta.nextGraphemeStart(end)
		if next > lineEnd {
			break
		}
		end = next
	}
	if end == ta.cursor {
		return
	}
	self.register.set(ta.getRange(ta.cursor, end), false)
	self.pushUndo(ta)
	ta.deleteRange(ta.cursor, end)
	self.clampCursor(ta)
}

func (self *VimEditor) deleteBackward(ta *TextArea, n int) {
	start := ta.cursor
	lineStart := ta.hardLineStart(ta.cursor)
	for range n {
		if start == lineStart {
			break
		}
		start = ta.prevGraphemeStart(start)
	}
	if start == ta.cursor {
		return
	}
	self.register.set(ta.getRange(start, ta.cursor), false)
	self.pushUndo(ta)
	ta.deleteRange(start, ta.cursor)
}

func (self *VimEditor) paste(ta *TextArea, after bool, n int) {
	if self.register.text == "" {
		return
	}
	self.pushUndo(ta)

	if self.register.linewise && self.allowMultiline {
		block := strings.Repeat(self.register.text+"\n", n)
		if after {
			lineEnd := ta.hardLineEnd(ta.cursor)
			if lineEnd == len(ta.content) {
				ta.insertAt(lineEnd, "\n"+strings.TrimSuffix(block, "\n"))
				ta.cursor = min(lineEnd+1, len(ta.content))
			} else {
				ta.insertAt(lineEnd+1, block)
				ta.cursor = lineEnd + 1
			}
		} else {
			start := ta.hardLineStart(ta.cursor)
			ta.insertAt(start, block)
			ta.cursor = start
		}
		ta.cursor = ta.vimFirstNonBlank(ta.cursor)
		self.clampCursor(ta)
		return
	}

	text := self.register.text
	if self.register.linewise {
		// pasting a linewise register into a single-line view degrades
		// to charwise so the paste isn't silently dropped
		text = strings.ReplaceAll(text, "\n", " ")
	}
	text = strings.Repeat(text, n)
	pos := ta.cursor
	if after {
		pos = min(ta.nextGraphemeStart(pos), ta.hardLineEnd(pos))
	}
	ta.insertAt(pos, text)
	ta.cursor = ta.prevGraphemeStart(pos + len(text))
	self.clampCursor(ta)
}

func (self *VimEditor) pushUndo(ta *TextArea) {
	self.undo = append(self.undo, vimSnapshot{content: ta.content, cursor: ta.cursor})
	if len(self.undo) > vimUndoLimit {
		self.undo = self.undo[1:]
	}
	self.redo = nil
}

func (self *VimEditor) restoreUndo(ta *TextArea) {
	for len(self.undo) > 0 {
		snap := self.undo[len(self.undo)-1]
		self.undo = self.undo[:len(self.undo)-1]
		if snap.content == ta.content {
			continue
		}
		self.redo = append(self.redo, vimSnapshot{content: ta.content, cursor: ta.cursor})
		ta.setState(snap.content, snap.cursor)
		self.clampCursor(ta)
		return
	}
}

func (self *VimEditor) restoreRedo(ta *TextArea) {
	for len(self.redo) > 0 {
		snap := self.redo[len(self.redo)-1]
		self.redo = self.redo[:len(self.redo)-1]
		if snap.content == ta.content {
			continue
		}
		self.undo = append(self.undo, vimSnapshot{content: ta.content, cursor: ta.cursor})
		ta.setState(snap.content, snap.cursor)
		self.clampCursor(ta)
		return
	}
}

func (self *VimEditor) startInsert(ta *TextArea) {
	self.pushUndo(ta)
	self.mode = VimModeInsert
}

func (self *VimEditor) enterNormalMode(ta *TextArea) {
	self.mode = VimModeNormal
	self.clearPending()
	if ta.cursor > ta.hardLineStart(ta.cursor) {
		ta.cursor = ta.prevGraphemeStart(ta.cursor)
	}
}

// vim's normal-mode cursor sits on a character, never on the newline after
// the last one; the insertion-point cursor must be pulled back one grapheme
// whenever a motion or edit leaves it hanging at a line end.
func (self *VimEditor) clampCursor(ta *TextArea) {
	if self.mode != VimModeNormal {
		return
	}
	lineStart := ta.hardLineStart(ta.cursor)
	if ta.cursor > lineStart && ta.cursor == ta.hardLineEnd(ta.cursor) {
		ta.cursor = ta.prevGraphemeStart(ta.cursor)
	}
}

func (self *VimEditor) repeatCount() int {
	return max(self.count, 1) * max(self.opCount, 1)
}

func (self *VimEditor) clearPending() {
	self.count = 0
	self.opCount = 0
	self.op = ""
	self.pending = ""
}
