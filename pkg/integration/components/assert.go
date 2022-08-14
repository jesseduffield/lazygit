package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"golang.org/x/exp/constraints"
)

// through this struct we assert on the state of the lazygit gui

type Assert struct {
	gui integrationTypes.GuiDriver
}

func NewAssert(gui integrationTypes.GuiDriver) *Assert {
	return &Assert{gui: gui}
}

// for making assertions on string values
type matcher[T any] struct {
	testFn func(T) (bool, string)
	prefix string
}

func (self *matcher[T]) test(value T) (bool, string) {
	ok, message := self.testFn(value)
	if ok {
		return true, ""
	}

	if self.prefix != "" {
		return false, self.prefix + " " + message
	}

	return false, message
}

func (self *matcher[T]) context(prefix string) *matcher[T] {
	self.prefix = prefix

	return self
}

func Contains(target string) *matcher[string] {
	return &matcher[string]{testFn: func(value string) (bool, string) {
		return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to contain '%s'", value, target)
	}}
}

func Equals[T constraints.Ordered](target T) *matcher[T] {
	return &matcher[T]{testFn: func(value T) (bool, string) {
		return target == value, fmt.Sprintf("Expected '%T' to equal '%T'", value, target)
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

func (self *Assert) MatchHeadCommitMessage(matcher *matcher[string]) {
	self.assertWithRetries(func() (bool, string) {
		return len(self.gui.Model().Commits) == 0, "Expected at least one commit to be present"
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

func (self *Assert) MatchSelectedLine(matcher *matcher[string]) {
	self.matchString(matcher, "Unexpected selected line.",
		func() string {
			return self.gui.CurrentContext().GetView().SelectedLine()
		},
	)
}

func (self *Assert) InPrompt() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, fmt.Sprintf("Expected prompt popup to be focused")
	})
}

func (self *Assert) InConfirm() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, fmt.Sprintf("Expected confirmation popup to be focused")
	})
}

func (self *Assert) InAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, fmt.Sprintf("Expected alert popup to be focused")
	})
}

func (self *Assert) InMenu() {
	self.assertWithRetries(func() (bool, string) {
		return self.gui.CurrentContext().GetView().Name() == "menu", fmt.Sprintf("Expected popup menu to be focused")
	})
}

func (self *Assert) MatchCurrentViewTitle(matcher *matcher[string]) {
	self.matchString(matcher, "Unexpected current view title.",
		func() string {
			return self.gui.CurrentContext().GetView().Title
		},
	)
}

func (self *Assert) MatchMainViewContent(matcher *matcher[string]) {
	self.matchString(matcher, "Unexpected main view content.",
		func() string {
			return self.gui.MainView().Buffer()
		},
	)
}

func (self *Assert) MatchSecondaryViewContent(matcher *matcher[string]) {
	self.matchString(matcher, "Unexpected secondary view title.",
		func() string {
			return self.gui.SecondaryView().Buffer()
		},
	)
}

func (self *Assert) matchString(matcher *matcher[string], context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
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
