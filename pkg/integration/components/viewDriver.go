package components

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

type ViewDriver struct {
	// context is prepended to any error messages e.g. 'context: "current view"'
	context string
	getView func() *gocui.View
	t       *TestDriver
}

// asserts that the view has the expected title
func (self *ViewDriver) Title(expected *matcher) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		actual := self.getView().Title
		return expected.context(fmt.Sprintf("%s title", self.context)).test(actual)
	})

	return self
}

// asserts that the view has lines matching the given matchers. So if three matchers
// are passed, we only check the first three lines of the view.
// This method is convenient when you have a list of commits but you only want to
// assert on the first couple of commits.
func (self *ViewDriver) TopLines(matchers ...*matcher) *ViewDriver {
	if len(matchers) < 1 {
		self.t.fail("TopLines method requires at least one matcher. If you are trying to assert that there are no lines, use .IsEmpty()")
	}

	self.t.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) >= len(matchers), fmt.Sprintf("unexpected number of lines in view. Expected at least %d, got %d", len(matchers), len(lines))
	})

	return self.assertLines(matchers...)
}

// asserts that the view has lines matching the given matchers. One matcher must be passed for each line.
// If you only care about the top n lines, use the TopLines method instead.
func (self *ViewDriver) Lines(matchers ...*matcher) *ViewDriver {
	self.LineCount(len(matchers))

	return self.assertLines(matchers...)
}

func (self *ViewDriver) assertLines(matchers ...*matcher) *ViewDriver {
	view := self.getView()

	for i, matcher := range matchers {
		checkIsSelected, matcher := matcher.checkIsSelected()

		self.t.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", view.Name()),
			func() string {
				return view.BufferLines()[i]
			},
		)

		if checkIsSelected {
			self.t.assertWithRetries(func() (bool, string) {
				lineIdx := view.SelectedLineIdx()
				return lineIdx == i, fmt.Sprintf("Unexpected selected line index in view '%s'. Expected %d, got %d", view.Name(), i, lineIdx)
			})
		}
	}

	return self
}

// asserts on the content of the view i.e. the stuff within the view's frame.
func (self *ViewDriver) Content(matcher *matcher) *ViewDriver {
	self.t.matchString(matcher, fmt.Sprintf("%s: Unexpected content.", self.context),
		func() string {
			return self.getView().Buffer()
		},
	)

	return self
}

// asserts on the selected line of the view
func (self *ViewDriver) SelectedLine(matcher *matcher) *ViewDriver {
	self.t.matchString(matcher, fmt.Sprintf("%s: Unexpected selected line.", self.context),
		func() string {
			return self.getView().SelectedLine()
		},
	)

	return self
}

// asserts on the index of the selected line. 0 is the first index, representing the line at the top of the view.
func (self *ViewDriver) SelectedLineIdx(expected int) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		actual := self.getView().SelectedLineIdx()
		return expected == actual, fmt.Sprintf("%s: Expected selected line index to be %d, got %d", self.context, expected, actual)
	})

	return self
}

// focus the view (assumes the view is a side-view that can be focused via a keybinding)
func (self *ViewDriver) Focus() *ViewDriver {
	// we can easily change focus by switching to the view's window, but this assumes that the desired view
	// is at the top of that window. So for now we'll switch to the window then assert that the desired
	// view is on top (i.e. that it's the current view).
	// If we want to support other views e.g. the tags view, we'll need to add more logic here.
	viewName := self.getView().Name()

	// using a map rather than a slice because we might add other views which share a window index later
	windowIndexMap := map[string]int{
		"status":        0,
		"files":         1,
		"localBranches": 2,
		"commits":       3,
		"stash":         4,
	}

	index, ok := windowIndexMap[viewName]
	if !ok {
		self.t.fail(fmt.Sprintf("Cannot focus view %s: Focus() method not implemented", viewName))
	}

	self.t.press(self.t.keys.Universal.JumpToBlock[index])

	// assert that we land in the expected view
	self.IsFocused()

	return self
}

// asserts that the view is focused
func (self *ViewDriver) IsFocused() *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		expected := self.getView().Name()
		actual := self.t.gui.CurrentContext().GetView().Name()
		return actual == expected, fmt.Sprintf("%s: Unexpected view focused. Expected %s, got %s", self.context, expected, actual)
	})

	return self
}

func (self *ViewDriver) Press(keyStr string) *ViewDriver {
	self.IsFocused()

	self.t.press(keyStr)

	return self
}

// i.e. pressing down arrow
func (self *ViewDriver) SelectNextItem() *ViewDriver {
	return self.Press(self.t.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *ViewDriver) SelectPreviousItem() *ViewDriver {
	return self.Press(self.t.keys.Universal.PrevItem)
}

// i.e. pressing space
func (self *ViewDriver) PressPrimaryAction() *ViewDriver {
	return self.Press(self.t.keys.Universal.Select)
}

// i.e. pressing space
func (self *ViewDriver) PressEnter() *ViewDriver {
	return self.Press(self.t.keys.Universal.Confirm)
}

// i.e. pressing escape
func (self *ViewDriver) PressEscape() *ViewDriver {
	return self.Press(self.t.keys.Universal.Return)
}

func (self *ViewDriver) NavigateToListItem(matcher *matcher) *ViewDriver {
	self.IsFocused()

	self.t.navigateToListItem(matcher)

	return self
}

// returns true if the view is a list view and it contains no items
func (self *ViewDriver) IsEmpty() *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		actual := strings.TrimSpace(self.getView().Buffer())
		return actual == "", fmt.Sprintf("%s: Unexpected content in view: expected no content. Content: %s", self.context, actual)
	})

	return self
}

func (self *ViewDriver) LineCount(expectedCount int) *ViewDriver {
	if expectedCount == 0 {
		return self.IsEmpty()
	}

	self.t.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) == expectedCount, fmt.Sprintf("unexpected number of lines in view. Expected %d, got %d", expectedCount, len(lines))
	})

	self.t.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()

		// if the view has a single blank line (often the case) we want to treat that as having no lines
		if len(lines) == 1 && expectedCount == 1 {
			actual := strings.TrimSpace(self.getView().Buffer())
			return actual != "", "unexpected number of lines in view. Expected 1, got 0"
		}

		return len(lines) == expectedCount, fmt.Sprintf("unexpected number of lines in view. Expected %d, got %d", expectedCount, len(lines))
	})

	return self
}

// for when you want to make some assertion unrelated to the current view
// without breaking the method chain
func (self *ViewDriver) Tap(f func()) *ViewDriver {
	f()

	return self
}
