package components

import (
	"fmt"
	"os"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// through this struct we assert on the state of the lazygit gui

type Assert struct {
	gui integrationTypes.GuiDriver
}

func NewAssert(gui integrationTypes.GuiDriver) *Assert {
	return &Assert{gui: gui}
}

func (self *Assert) WorkingTreeFileCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Files)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d changed working tree files, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Assert) CommitCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Commits)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d commits present, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Assert) StashCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().StashEntries)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d stash entries, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Assert) AtLeastOneCommit() {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Commits)

		return actualCount > 0, "Expected at least one commit present"
	})
}

func (self *Assert) HeadCommitMessage(matcher *matcher) {
	self.assertWithRetries(func() (bool, string) {
		return len(self.gui.Model().Commits) > 0, "Expected at least one commit to be present"
	})

	self.matchString(matcher, "Unexpected commit message.",
		func() string {
			return self.gui.Model().Commits[0].Name
		},
	)
}

func (self *Assert) CurrentWindowName(expectedWindowName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CurrentContext().GetView().Name()
		return actual == expectedWindowName, fmt.Sprintf("Expected current window name to be '%s', but got '%s'", expectedWindowName, actual)
	})
}

func (self *Assert) CurrentBranchName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CheckedOutRef().Name
		return actual == expectedViewName, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expectedViewName, actual)
	})
}

func (self *Assert) InListContext() {
	self.assertWithRetries(func() (bool, string) {
		currentContext := self.gui.CurrentContext()
		_, ok := currentContext.(types.IListContext)
		return ok, fmt.Sprintf("Expected current context to be a list context, but got %s", currentContext.GetKey())
	})
}

func (self *Assert) InPrompt() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, "Expected prompt popup to be focused"
	})
}

func (self *Assert) InConfirm() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected confirmation popup to be focused"
	})
}

func (self *Assert) InAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected alert popup to be focused"
	})
}

func (self *Assert) InCommitMessagePanel() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "commitMessage", "Expected commit message panel to be focused"
	})
}

func (self *Assert) InMenu() {
	self.assertWithRetries(func() (bool, string) {
		return self.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Assert) NotInPopup() {
	self.assertWithRetries(func() (bool, string) {
		currentViewName := self.gui.CurrentContext().GetView().Name()
		return currentViewName != "menu" && currentViewName != "confirmation" && currentViewName != "commitMessage", "Expected popup not to be focused"
	})
}

func (self *Assert) matchString(matcher *matcher, context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
	})
}

func (self *Assert) assertWithRetries(test func() (bool, string)) {
	waitTimes := []int{0, 1, 1, 1, 1, 1, 5, 10, 20, 40, 100, 200, 500, 1000, 2000, 4000}

	var message string
	for _, waitTime := range waitTimes {
		time.Sleep(time.Duration(waitTime) * time.Millisecond)

		var ok bool
		ok, message = test()
		if ok {
			return
		}
	}

	self.Fail(message)
}

// for when you just want to fail the test yourself
func (self *Assert) Fail(message string) {
	self.gui.Fail(message)
}

// This does _not_ check the files panel, it actually checks the filesystem
func (self *Assert) FileSystemPathPresent(path string) {
	self.assertWithRetries(func() (bool, string) {
		_, err := os.Stat(path)
		return err == nil, fmt.Sprintf("Expected path '%s' to exist, but it does not", path)
	})
}

// This does _not_ check the files panel, it actually checks the filesystem
func (self *Assert) FileSystemPathNotPresent(path string) {
	self.assertWithRetries(func() (bool, string) {
		_, err := os.Stat(path)
		return os.IsNotExist(err), fmt.Sprintf("Expected path '%s' to not exist, but it does", path)
	})
}

func (self *Assert) CurrentView() *ViewAsserter {
	return &ViewAsserter{
		context: "current view",
		getView: func() *gocui.View { return self.gui.CurrentContext().GetView() },
		assert:  self,
	}
}

func (self *Assert) View(viewName string) *ViewAsserter {
	return &ViewAsserter{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.gui.View(viewName) },
		assert:  self,
	}
}

func (self *Assert) MainView() *ViewAsserter {
	return &ViewAsserter{
		context: "main view",
		getView: func() *gocui.View { return self.gui.MainView() },
		assert:  self,
	}
}

func (self *Assert) SecondaryView() *ViewAsserter {
	return &ViewAsserter{
		context: "secondary view",
		getView: func() *gocui.View { return self.gui.SecondaryView() },
		assert:  self,
	}
}
