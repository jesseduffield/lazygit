package components

import (
	"fmt"

	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type Model struct {
	*assertionHelper
	gui integrationTypes.GuiDriver
}

func (self *Model) WorkingTreeFileCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Files)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d changed working tree files, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Model) CommitCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Commits)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d commits present, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Model) StashCount(expectedCount int) {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().StashEntries)

		return actualCount == expectedCount, fmt.Sprintf(
			"Expected %d stash entries, but got %d",
			expectedCount, actualCount,
		)
	})
}

func (self *Model) AtLeastOneCommit() {
	self.assertWithRetries(func() (bool, string) {
		actualCount := len(self.gui.Model().Commits)

		return actualCount > 0, "Expected at least one commit present"
	})
}

func (self *Model) HeadCommitMessage(matcher *matcher) {
	self.assertWithRetries(func() (bool, string) {
		return len(self.gui.Model().Commits) > 0, "Expected at least one commit to be present"
	})

	self.matchString(matcher, "Unexpected commit message.",
		func() string {
			return self.gui.Model().Commits[0].Name
		},
	)
}

func (self *Model) CurrentBranchName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CheckedOutRef().Name
		return actual == expectedViewName, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expectedViewName, actual)
	})
}
