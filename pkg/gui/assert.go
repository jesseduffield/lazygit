package gui

import (
	"fmt"
	"time"

	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

type AssertImpl struct {
	gui *Gui
}

var _ types.Assert = &AssertImpl{}

func (self *AssertImpl) WorkingTreeFileCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.State.Model.Files)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d changed working tree files, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *AssertImpl) CommitCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.State.Model.Commits)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d commits present, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *AssertImpl) HeadCommitMessage(expectedMessage string) {
	self.assertWithRetries(func() (bool, string) {
		if len(self.gui.State.Model.Commits) == 0 {
			return false, "Expected at least one commit to be present"
		}

		headCommit := self.gui.State.Model.Commits[0]
		if headCommit.Name != expectedMessage {
			return false, fmt.Sprintf(
				"Expected commit message to be '%s', but got '%s'",
				expectedMessage, headCommit.Name,
			)
		}

		return true, ""
	})
}

func (self *AssertImpl) CurrentViewName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.currentViewName()
		return actual == expectedViewName, fmt.Sprintf("Expected current view name to be '%s', but got '%s'", expectedViewName, actual)
	})
}

func (self *AssertImpl) CurrentBranchName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.helpers.Refs.GetCheckedOutRef().Name
		return actual == expectedViewName, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expectedViewName, actual)
	})
}

func (self *AssertImpl) assertWithRetries(test func() (bool, string)) {
	waitTimes := []int{0, 100, 200, 400, 800, 1600}

	var message string
	for _, waitTime := range waitTimes {
		time.Sleep(time.Duration(waitTime) * time.Millisecond)

		var ok bool
		ok, message = test()
		if ok {
			return
		}
	}

	self.gui.g.Close()
	// need to give the gui time to close
	time.Sleep(time.Millisecond * 100)
	panic(message)
}
