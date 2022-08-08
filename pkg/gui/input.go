package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	guiTypes "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

type InputImpl struct {
	gui          *Gui
	keys         config.KeybindingConfig
	assert       types.Assert
	pushKeyDelay int
}

func NewInputImpl(gui *Gui, keys config.KeybindingConfig, assert types.Assert, pushKeyDelay int) *InputImpl {
	return &InputImpl{
		gui:          gui,
		keys:         keys,
		assert:       assert,
		pushKeyDelay: pushKeyDelay,
	}
}

var _ types.Input = &InputImpl{}

func (self *InputImpl) PressKeys(keyStrs ...string) {
	for _, keyStr := range keyStrs {
		self.pressKey(keyStr)
	}
}

func (self *InputImpl) pressKey(keyStr string) {
	self.Wait(self.pushKeyDelay)

	key := keybindings.GetKey(keyStr)

	var r rune
	var tcellKey tcell.Key
	switch v := key.(type) {
	case rune:
		r = v
		tcellKey = tcell.KeyRune
	case gocui.Key:
		tcellKey = tcell.Key(v)
	}

	self.gui.g.ReplayedEvents.Keys <- gocui.NewTcellKeyEventWrapper(
		tcell.NewEventKey(tcellKey, r, tcell.ModNone),
		0,
	)
}

func (self *InputImpl) SwitchToStatusWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[0])
}

func (self *InputImpl) SwitchToFilesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[1])
}

func (self *InputImpl) SwitchToBranchesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[2])
}

func (self *InputImpl) SwitchToCommitsWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[3])
}

func (self *InputImpl) SwitchToStashWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[4])
}

func (self *InputImpl) Type(content string) {
	for _, char := range content {
		self.pressKey(string(char))
	}
}

func (self *InputImpl) Confirm() {
	self.pressKey(self.keys.Universal.Confirm)
}

func (self *InputImpl) Cancel() {
	self.pressKey(self.keys.Universal.Return)
}

func (self *InputImpl) Select() {
	self.pressKey(self.keys.Universal.Select)
}

func (self *InputImpl) NextItem() {
	self.pressKey(self.keys.Universal.NextItem)
}

func (self *InputImpl) PreviousItem() {
	self.pressKey(self.keys.Universal.PrevItem)
}

func (self *InputImpl) ContinueMerge() {
	self.PressKeys(self.keys.Universal.CreateRebaseOptionsMenu)
	self.assert.SelectedLineContains("continue")
	self.Confirm()
}

func (self *InputImpl) ContinueRebase() {
	self.ContinueMerge()
}

func (self *InputImpl) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func (self *InputImpl) log(message string) {
	self.gui.c.LogAction(message)
}

// NOTE: this currently assumes that ViewBufferLines returns all the lines that can be accessed.
// If this changes in future, we'll need to update this code to first attempt to find the item
// in the current page and failing that, jump to the top of the view and iterate through all of it,
// looking for the item.
func (self *InputImpl) NavigateToListItemContainingText(text string) {
	self.assert.InListContext()

	currentContext := self.gui.currentContext().(guiTypes.IListContext)
	view := currentContext.GetView()

	// first we look for a duplicate on the current screen. We won't bother looking beyond that though.
	matchCount := 0
	matchIndex := -1
	for i, line := range view.ViewBufferLines() {
		if strings.Contains(line, text) {
			matchCount++
			matchIndex = i
		}
	}
	if matchCount > 1 {
		self.assert.Fail(fmt.Sprintf("Found %d matches for %s, expected only a single match", matchCount, text))
	}
	if matchCount == 1 {
		selectedLineIdx := view.SelectedLineIdx()
		if selectedLineIdx == matchIndex {
			return
		}
		if selectedLineIdx < matchIndex {
			for i := selectedLineIdx; i < matchIndex; i++ {
				self.NextItem()
			}
			return
		} else {
			for i := selectedLineIdx; i > matchIndex; i-- {
				self.PreviousItem()
			}
			return
		}
	}

	self.assert.Fail(fmt.Sprintf("Could not find item containing text: %s", text))
}
