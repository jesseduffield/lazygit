package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type Impl struct {
	gui          integrationTypes.GuiAdapter
	keys         config.KeybindingConfig
	assert       *Assert
	pushKeyDelay int
}

func NewInput(gui integrationTypes.GuiAdapter, keys config.KeybindingConfig, assert *Assert, pushKeyDelay int) *Impl {
	return &Impl{
		gui:          gui,
		keys:         keys,
		assert:       assert,
		pushKeyDelay: pushKeyDelay,
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *Impl) PressKeys(keyStrs ...string) {
	for _, keyStr := range keyStrs {
		self.pressKey(keyStr)
	}
}

func (self *Impl) pressKey(keyStr string) {
	self.Wait(self.pushKeyDelay)

	self.gui.PressKey(keyStr)
}

func (self *Impl) SwitchToStatusWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[0])
}

func (self *Impl) SwitchToFilesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[1])
}

func (self *Impl) SwitchToBranchesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[2])
}

func (self *Impl) SwitchToCommitsWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[3])
}

func (self *Impl) SwitchToStashWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[4])
}

func (self *Impl) Type(content string) {
	for _, char := range content {
		self.pressKey(string(char))
	}
}

// i.e. pressing enter
func (self *Impl) Confirm() {
	self.pressKey(self.keys.Universal.Confirm)
}

// i.e. pressing escape
func (self *Impl) Cancel() {
	self.pressKey(self.keys.Universal.Return)
}

// i.e. pressing space
func (self *Impl) Select() {
	self.pressKey(self.keys.Universal.Select)
}

// i.e. pressing down arrow
func (self *Impl) NextItem() {
	self.pressKey(self.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *Impl) PreviousItem() {
	self.pressKey(self.keys.Universal.PrevItem)
}

func (self *Impl) ContinueMerge() {
	self.PressKeys(self.keys.Universal.CreateRebaseOptionsMenu)
	self.assert.SelectedLineContains("continue")
	self.Confirm()
}

func (self *Impl) ContinueRebase() {
	self.ContinueMerge()
}

// for when you want to allow lazygit to process something before continuing
func (self *Impl) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func (self *Impl) LogUI(message string) {
	self.gui.LogUI(message)
}

func (self *Impl) Log(message string) {
	self.gui.LogUI(message)
}

// this will look for a list item in the current panel and if it finds it, it will
// enter the keypresses required to navigate to it.
// The test will fail if:
//  - the user is not in a list item
//  - no list item is found containing the given text
//  - multiple list items are found containing the given text in the initial page of items
//
// NOTE: this currently assumes that ViewBufferLines returns all the lines that can be accessed.
// If this changes in future, we'll need to update this code to first attempt to find the item
// in the current page and failing that, jump to the top of the view and iterate through all of it,
// looking for the item.
func (self *Impl) NavigateToListItemContainingText(text string) {
	self.assert.InListContext()

	currentContext := self.gui.CurrentContext().(types.IListContext)

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
