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
func (self *Input) Press(keyStrs ...string) {
	for _, keyStr := range keyStrs {
		self.press(keyStr)
	}
}

func (self *Input) press(keyStr string) {
	self.Wait(self.pushKeyDelay)

	self.gui.PressKey(keyStr)
}

func (self *Input) SwitchToStatusWindow() {
	self.press(self.keys.Universal.JumpToBlock[0])
	self.assert.CurrentWindowName("status")
}

// switch to status window and assert that the status view is on top
func (self *Input) SwitchToStatusView() {
	self.SwitchToStatusWindow()
	self.assert.CurrentView().Name("status")
}

func (self *Input) SwitchToFilesWindow() {
	self.press(self.keys.Universal.JumpToBlock[1])
	self.assert.CurrentWindowName("files")
}

// switch to files window and assert that the files view is on top
func (self *Input) SwitchToFilesView() {
	self.SwitchToFilesWindow()
	self.assert.CurrentView().Name("files")
}

func (self *Input) SwitchToBranchesWindow() {
	self.press(self.keys.Universal.JumpToBlock[2])
	self.assert.CurrentWindowName("localBranches")
}

// switch to branches window and assert that the branches view is on top
func (self *Input) SwitchToBranchesView() {
	self.SwitchToBranchesWindow()
	self.assert.CurrentView().Name("localBranches")
}

func (self *Input) SwitchToCommitsWindow() {
	self.press(self.keys.Universal.JumpToBlock[3])
	self.assert.CurrentWindowName("commits")
}

// switch to commits window and assert that the commits view is on top
func (self *Input) SwitchToCommitsView() {
	self.SwitchToCommitsWindow()
	self.assert.CurrentView().Name("commits")
}

func (self *Input) SwitchToStashWindow() {
	self.press(self.keys.Universal.JumpToBlock[4])
	self.assert.CurrentWindowName("stash")
}

// switch to stash window and assert that the stash view is on top
func (self *Input) SwitchToStashView() {
	self.SwitchToStashWindow()
	self.assert.CurrentView().Name("stash")
}

func (self *Input) Type(content string) {
	for _, char := range content {
		self.press(string(char))
	}
}

// i.e. pressing enter
func (self *Input) Confirm() {
	self.press(self.keys.Universal.Confirm)
}

// i.e. same as Confirm
func (self *Input) Enter() {
	self.press(self.keys.Universal.Confirm)
}

// i.e. pressing escape
func (self *Input) Cancel() {
	self.press(self.keys.Universal.Return)
}

// i.e. pressing space
func (self *Input) PrimaryAction() {
	self.press(self.keys.Universal.Select)
}

// i.e. pressing down arrow
func (self *Input) NextItem() {
	self.press(self.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *Input) PreviousItem() {
	self.press(self.keys.Universal.PrevItem)
}

func (self *Input) ContinueMerge() {
	self.Press(self.keys.Universal.CreateRebaseOptionsMenu)
	self.assert.CurrentView().SelectedLine(Contains("continue"))
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
func (self *Input) NavigateToListItem(matcher *matcher) {
	self.assert.InListContext()

	currentContext := self.gui.CurrentContext().(types.IListContext)

	view := currentContext.GetView()

	var matchIndex int

	self.assert.assertWithRetries(func() (bool, string) {
		matchIndex = -1
		var matches []string
		lines := view.ViewBufferLines()
		// first we look for a duplicate on the current screen. We won't bother looking beyond that though.
		for i, line := range lines {
			ok, _ := matcher.test(line)
			if ok {
				matches = append(matches, line)
				matchIndex = i
			}
		}
		if len(matches) > 1 {
			return false, fmt.Sprintf("Found %d matches for `%s`, expected only a single match. Matching lines:\n%s", len(matches), matcher.name(), strings.Join(matches, "\n"))
		} else if len(matches) == 0 {
			return false, fmt.Sprintf("Could not find item matching: %s. Lines:\n%s", matcher.name(), strings.Join(lines, "\n"))
		} else {
			return true, ""
		}
	})

	selectedLineIdx := view.SelectedLineIdx()
	if selectedLineIdx == matchIndex {
		self.assert.CurrentView().SelectedLine(matcher)
		return
	}
	if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.NextItem()
		}
		self.assert.CurrentView().SelectedLine(matcher)
		return
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.PreviousItem()
		}
		self.assert.CurrentView().SelectedLine(matcher)
		return
	}
}

func (self *Input) InConfirm() *ConfirmationAsserter {
	self.assert.InConfirm()

	return &ConfirmationAsserter{assert: self.assert, input: self}
}

func (self *Input) Prompt(title *matcher, textToType string) {
	self.assert.InPrompt()
	self.assert.CurrentView().Title(title)
	self.Type(textToType)
	self.Confirm()
}

// type some text into a prompt, then switch to the suggestions panel and expect the first
// item to match the given matcher, then confirm that item.
func (self *Input) Typeahead(title *matcher, textToType string, expectedFirstOption *matcher) {
	self.assert.InPrompt()
	self.assert.CurrentView().Title(title)
	self.Type(textToType)
	self.Press(self.keys.Universal.TogglePanel)
	self.assert.CurrentView().Name("suggestions")
	self.assert.CurrentView().SelectedLine(expectedFirstOption)
	self.Confirm()
}

func (self *Input) Menu(title *matcher, optionToSelect *matcher) {
	self.assert.InMenu()
	self.assert.CurrentView().Title(title)
	self.NavigateToListItem(optionToSelect)
	self.Confirm()
}

func (self *Input) Alert(title *matcher, content *matcher) {
	self.assert.InListContext()
	self.assert.CurrentView().Title(title)
	self.assert.CurrentView().Content(content)
	self.Confirm()
}
