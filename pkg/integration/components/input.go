package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/samber/lo"
)

type Input struct {
	gui          integrationTypes.GuiDriver
	keys         config.KeybindingConfig
	pushKeyDelay int
	*assertionHelper
}

func NewInput(gui integrationTypes.GuiDriver, keys config.KeybindingConfig, pushKeyDelay int) *Input {
	return &Input{
		gui:             gui,
		keys:            keys,
		pushKeyDelay:    pushKeyDelay,
		assertionHelper: &assertionHelper{gui: gui},
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *Input) press(keyStr string) {
	self.Wait(self.pushKeyDelay)

	self.gui.PressKey(keyStr)
}

func (self *Input) typeContent(content string) {
	for _, char := range content {
		self.press(string(char))
	}
}

func (self *Input) ContinueMerge() {
	self.Views().current().Press(self.keys.Universal.CreateRebaseOptionsMenu)

	self.ExpectMenu().
		Title(Equals("Rebase Options")).
		Select(Contains("continue")).
		Confirm()
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
func (self *Input) navigateToListItem(matcher *matcher) {
	self.inListContext()

	currentContext := self.gui.CurrentContext().(types.IListContext)

	view := currentContext.GetView()

	var matchIndex int

	self.assertWithRetries(func() (bool, string) {
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
		self.Views().current().SelectedLine(matcher)
		return
	}
	if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.Views().current().SelectNextItem()
		}
		self.Views().current().SelectedLine(matcher)
		return
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.Views().current().SelectPreviousItem()
		}
		self.Views().current().SelectedLine(matcher)
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

func (self *Input) ExpectConfirmation() *ConfirmationAsserter {
	self.inConfirm()

	return &ConfirmationAsserter{input: self}
}

func (self *Input) inConfirm() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected confirmation popup to be focused"
	})
}

func (self *Input) ExpectPrompt() *PromptAsserter {
	self.inPrompt()

	return &PromptAsserter{input: self}
}

func (self *Input) inPrompt() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, "Expected prompt popup to be focused"
	})
}

func (self *Input) ExpectAlert() *AlertAsserter {
	self.inAlert()

	return &AlertAsserter{input: self}
}

func (self *Input) inAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected alert popup to be focused"
	})
}

func (self *Input) ExpectMenu() *MenuAsserter {
	self.inMenu()

	return &MenuAsserter{input: self}
}

func (self *Input) inMenu() {
	self.assertWithRetries(func() (bool, string) {
		return self.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Input) ExpectCommitMessagePanel() *CommitMessagePanelAsserter {
	self.inCommitMessagePanel()

	return &CommitMessagePanelAsserter{input: self}
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

// for making assertions on lazygit views
func (self *Input) Views() *Views {
	return &Views{input: self}
}

// for making assertions on the lazygit model
func (self *Input) Model() *Model {
	return &Model{assertionHelper: self.assertionHelper, gui: self.gui}
}

// for making assertions on the file system
func (self *Input) FileSystem() *FileSystem {
	return &FileSystem{assertionHelper: self.assertionHelper}
}

// for when you just want to fail the test yourself.
// This runs callbacks to ensure we render the error after closing the gui.
func (self *Input) Fail(message string) {
	self.assertionHelper.fail(message)
}

func (self *Input) NotInPopup() {
	self.assertWithRetries(func() (bool, string) {
		viewName := self.gui.CurrentContext().GetView().Name()
		return !lo.Contains([]string{"menu", "confirmation", "commitMessage"}, viewName), fmt.Sprintf("Unexpected popup view present: %s view", viewName)
	})
}
