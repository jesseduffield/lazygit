package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type Input struct {
	gui          integrationTypes.GuiDriver
	keys         config.KeybindingConfig
	assert       *Assert
	pushKeyDelay int
}

func NewInput(gui integrationTypes.GuiDriver, keys config.KeybindingConfig, assert *Assert, pushKeyDelay int) *Input {
	return &Input{
		gui:          gui,
		keys:         keys,
		assert:       assert,
		pushKeyDelay: pushKeyDelay,
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *Input) PressKeys(keyStrs ...string) {
	for _, keyStr := range keyStrs {
		self.pressKey(keyStr)
	}
}

func (self *Input) pressKey(keyStr string) {
	self.Wait(self.pushKeyDelay)

	self.gui.PressKey(keyStr)
}

func (self *Input) SwitchToStatusWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[0])
	self.assert.CurrentWindowName("status")
}

func (self *Input) SwitchToFilesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[1])
	self.assert.CurrentWindowName("files")
}

func (self *Input) SwitchToBranchesWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[2])
	self.assert.CurrentWindowName("localBranches")
}

func (self *Input) SwitchToCommitsWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[3])
	self.assert.CurrentWindowName("commits")
}

func (self *Input) SwitchToStashWindow() {
	self.pressKey(self.keys.Universal.JumpToBlock[4])
	self.assert.CurrentWindowName("stash")
}

func (self *Input) Type(content string) {
	for _, char := range content {
		self.pressKey(string(char))
	}
}

// i.e. pressing enter
func (self *Input) Confirm() {
	self.pressKey(self.keys.Universal.Confirm)
}

func (self *Input) ProceedWhenAsked(matcher *matcher) {
	self.assert.InConfirm()
	self.assert.MatchCurrentViewContent(matcher)
	self.Confirm()
}

// i.e. same as Confirm
func (self *Input) Enter() {
	self.pressKey(self.keys.Universal.Confirm)
}

// i.e. pressing escape
func (self *Input) Cancel() {
	self.pressKey(self.keys.Universal.Return)
}

// i.e. pressing space
func (self *Input) PrimaryAction() {
	self.pressKey(self.keys.Universal.Select)
}

// i.e. pressing down arrow
func (self *Input) NextItem() {
	self.pressKey(self.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *Input) PreviousItem() {
	self.pressKey(self.keys.Universal.PrevItem)
}

func (self *Input) ContinueMerge() {
	self.PressKeys(self.keys.Universal.CreateRebaseOptionsMenu)
	self.assert.MatchSelectedLine(Contains("continue"))
	self.Confirm()
}

func (self *Input) ContinueRebase() {
	self.ContinueMerge()
}

// for when you want to allow lazygit to process something before continuing
func (self *Input) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func (self *Input) LogUI(message string) {
	self.gui.LogUI(message)
}

func (self *Input) Log(message string) {
	self.gui.LogUI(message)
}

// this will look for a list item in the current panel and if it finds it, it will
// enter the keypresses required to navigate to it.
// The test will fail if:
// - the user is not in a list item
// - no list item is found containing the given text
// - multiple list items are found containing the given text in the initial page of items
//
// NOTE: this currently assumes that ViewBufferLines returns all the lines that can be accessed.
// If this changes in future, we'll need to update this code to first attempt to find the item
// in the current page and failing that, jump to the top of the view and iterate through all of it,
// looking for the item.
func (self *Input) NavigateToListItemContainingText(text string) {
	self.assert.InListContext()

	currentContext := self.gui.CurrentContext().(types.IListContext)

	view := currentContext.GetView()

	var matchIndex int

	self.assert.assertWithRetries(func() (bool, string) {
		matchCount := 0
		matchIndex = -1
		// first we look for a duplicate on the current screen. We won't bother looking beyond that though.
		for i, line := range view.ViewBufferLines() {
			if strings.Contains(line, text) {
				matchCount++
				matchIndex = i
			}
		}
		if matchCount > 1 {
			return false, fmt.Sprintf("Found %d matches for %s, expected only a single match", matchCount, text)
		} else if matchCount == 0 {
			return false, fmt.Sprintf("Could not find item containing text: %s", text)
		} else {
			return true, ""
		}
	})

	selectedLineIdx := view.SelectedLineIdx()
	if selectedLineIdx == matchIndex {
		self.assert.MatchSelectedLine(Contains(text))
		return
	}
	if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.NextItem()
		}
		self.assert.MatchSelectedLine(Contains(text))
		return
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.PreviousItem()
		}
		self.assert.MatchSelectedLine(Contains(text))
		return
	}
}
