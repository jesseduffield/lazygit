package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// through this struct we assert on the state of the lazygit gui

type Assert struct {
	gui integrationTypes.GuiAdapter
}

func NewAssert(gui integrationTypes.GuiAdapter) *Assert {
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

func (self *Assert) HeadCommitMessage(expectedMessage string) {
	self.assertWithRetries(func() (bool, string) {
		if len(self.gui.Model().Commits) == 0 {
			return false, "Expected at least one commit to be present"
		}

		headCommit := self.gui.Model().Commits[0]
		if headCommit.Name != expectedMessage {
			return false, fmt.Sprintf(
				"Expected commit message to be '%s', but got '%s'",
				expectedMessage, headCommit.Name,
			)
		}

		return true, ""
	})
}

func (self *Assert) CurrentViewName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CurrentContext().GetView().Name()
		return actual == expectedViewName, fmt.Sprintf("Expected current view name to be '%s', but got '%s'", expectedViewName, actual)
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

func (self *Assert) SelectedLineContains(text string) {
	self.assertWithRetries(func() (bool, string) {
		line := self.gui.CurrentContext().GetView().SelectedLine()
		return strings.Contains(line, text), fmt.Sprintf("Expected selected line to contain '%s', but got '%s'", text, line)
	})
}

func (self *Assert) assertWithRetries(test func() (bool, string)) {
	waitTimes := []int{0, 1, 5, 10, 200, 500, 1000}

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
