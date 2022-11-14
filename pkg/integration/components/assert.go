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
	gui integrationTypes.GuiDriver
}

func NewAssert(gui integrationTypes.GuiDriver) *Assert {
	return &Assert{gui: gui}
}

// for making assertions on string values
type matcher struct {
	testFn func(string) (bool, string)
	prefix string
}

func (self *matcher) test(value string) (bool, string) {
	ok, message := self.testFn(value)
	if ok {
		return true, ""
	}

	if self.prefix != "" {
		return false, self.prefix + " " + message
	}

	return false, message
}

func (self *matcher) context(prefix string) *matcher {
	self.prefix = prefix

	return self
}

func Contains(target string) *matcher {
	return &matcher{testFn: func(value string) (bool, string) {
		return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to be found in '%s'", target, value)
	}}
}

func NotContains(target string) *matcher {
	return &matcher{testFn: func(value string) (bool, string) {
		return !strings.Contains(value, target), fmt.Sprintf("Expected '%s' to NOT be found in '%s'", target, value)
	}}
}

func Equals(target string) *matcher {
	return &matcher{testFn: func(value string) (bool, string) {
		return target == value, fmt.Sprintf("Expected '%s' to equal '%s'", value, target)
	}}
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

func (self *Assert) MatchHeadCommitMessage(matcher *matcher) {
	self.assertWithRetries(func() (bool, string) {
		return len(self.gui.Model().Commits) > 0, "Expected at least one commit to be present"
	})

	self.matchString(matcher, "Unexpected commit message.",
		func() string {
			return self.gui.Model().Commits[0].Name
		},
	)
}

func (self *Assert) CurrentViewName(expectedViewName string) {
	self.assertWithRetries(func() (bool, string) {
		actual := self.gui.CurrentContext().GetView().Name()
		return actual == expectedViewName, fmt.Sprintf("Expected current view name to be '%s', but got '%s'", expectedViewName, actual)
	})
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

func (self *Assert) MatchSelectedLine(matcher *matcher) {
	self.matchString(matcher, "Unexpected selected line.",
		func() string {
			return self.gui.CurrentContext().GetView().SelectedLine()
		},
	)
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

func (self *Assert) InMenu() {
	self.assertWithRetries(func() (bool, string) {
		return self.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Assert) MatchCurrentViewTitle(matcher *matcher) {
	self.matchString(matcher, "Unexpected current view title.",
		func() string {
			return self.gui.CurrentContext().GetView().Title
		},
	)
}

func (self *Assert) MatchViewContent(viewName string, matcher *matcher) {
	self.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", viewName),
		func() string {
			return self.gui.View(viewName).Buffer()
		},
	)
}

func (self *Assert) MatchCurrentViewContent(matcher *matcher) {
	self.matchString(matcher, "Unexpected content in current view.",
		func() string {
			return self.gui.CurrentContext().GetView().Buffer()
		},
	)
}

func (self *Assert) MatchMainViewContent(matcher *matcher) {
	self.matchString(matcher, "Unexpected main view content.",
		func() string {
			return self.gui.MainView().Buffer()
		},
	)
}

func (self *Assert) MatchSecondaryViewContent(matcher *matcher) {
	self.matchString(matcher, "Unexpected secondary view title.",
		func() string {
			return self.gui.SecondaryView().Buffer()
		},
	)
}

func (self *Assert) matchString(matcher *matcher, context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
	})
}

func (self *Assert) assertWithRetries(test func() (bool, string)) {
	waitTimes := []int{0, 1, 5, 10, 200, 500, 1000, 2000, 4000}

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
