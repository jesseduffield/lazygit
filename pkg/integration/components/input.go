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
	gui    integrationTypes.GuiDriver
	keys   config.KeybindingConfig
	assert *Assert
	*assertionHelper
	pushKeyDelay int
}

func NewInput(gui integrationTypes.GuiDriver, keys config.KeybindingConfig, pushKeyDelay int) *Input {
	return &Input{
		gui:             gui,
		keys:            keys,
		pushKeyDelay:    pushKeyDelay,
		assertionHelper: assert.assertionHelper,
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
	self.currentWindowName("status")
}

// switch to status window and assert that the status view is on top
func (self *Input) SwitchToStatusView() {
	self.SwitchToStatusWindow()
	self.assert.Views().Current().Name("status")
}

func (self *Input) switchToFilesWindow() {
	self.press(self.keys.Universal.JumpToBlock[1])
	self.currentWindowName("files")
}

// switch to files window and assert that the files view is on top
func (self *Input) SwitchToFilesView() {
	self.switchToFilesWindow()
	self.assert.Views().Current().Name("files")
}

func (self *Input) SwitchToBranchesWindow() {
	self.press(self.keys.Universal.JumpToBlock[2])
	self.currentWindowName("localBranches")
}

// switch to branches window and assert that the branches view is on top
func (self *Input) SwitchToBranchesView() {
	self.SwitchToBranchesWindow()
	self.assert.Views().Current().Name("localBranches")
}

func (self *Input) SwitchToCommitsWindow() {
	self.press(self.keys.Universal.JumpToBlock[3])
	self.currentWindowName("commits")
}

// switch to commits window and assert that the commits view is on top
func (self *Input) SwitchToCommitsView() {
	self.SwitchToCommitsWindow()
	self.assert.Views().Current().Name("commits")
}

func (self *Input) SwitchToStashWindow() {
	self.press(self.keys.Universal.JumpToBlock[4])
	self.currentWindowName("stash")
}

// switch to stash window and assert that the stash view is on top
func (self *Input) SwitchToStashView() {
	self.SwitchToStashWindow()
	self.assert.Views().Current().Name("stash")
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
	self.assert.Views().Current().SelectedLine(Contains("continue"))
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
	self.inListContext()

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
		self.assert.Views().Current().SelectedLine(matcher)
		return
	}
	if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.NextItem()
		}
		self.assert.Views().Current().SelectedLine(matcher)
		return
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.PreviousItem()
		}
		self.assert.Views().Current().SelectedLine(matcher)
		return
	}
}

func (self *Input) inListContext() {
	self.assertWithRetries(func() (bool, string) {
		currentContext := self.gui.CurrentContext()
		_, ok := currentContext.(types.IListContext)
		return ok, fmt.Sprintf("Expected current context to be a list context, but got %s", currentContext.GetKey())
	})
}

func (self *Input) Confirmation() *ConfirmationAsserter {
	self.inConfirm()

	return &ConfirmationAsserter{assert: self.assert, input: self}
}

func (self *Input) inConfirm() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected confirmation popup to be focused"
	})
}

func (self *Input) Prompt() *PromptAsserter {
	self.inPrompt()

	return &PromptAsserter{assert: self.assert, input: self}
}

func (self *Input) inPrompt() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, "Expected prompt popup to be focused"
	})
}

func (self *Input) Alert() *AlertAsserter {
	self.inAlert()

	return &AlertAsserter{assert: self.assert, input: self}
}

func (self *Input) inAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected alert popup to be focused"
	})
}

func (self *Input) Menu() *MenuAsserter {
	self.inMenu()

	return &MenuAsserter{assert: self.assert, input: self}
}

func (self *Input) inMenu() {
	self.assertWithRetries(func() (bool, string) {
		return self.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Input) CommitMessagePanel() *CommitMessagePanelAsserter {
	self.inCommitMessagePanel()

	return &CommitMessagePanelAsserter{assert: self.assert, input: self}
}

func (self *Input) inCommitMessagePanel() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "commitMessage", "Expected commit message panel to be focused"
	})
}

func (self *Input) currentWindowName(expectedWindowName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CurrentContext().GetView().Name()
		return actual == expectedWindowName, fmt.Sprintf("Expected current window name to be '%s', but got '%s'", expectedWindowName, actual)
	})
}
