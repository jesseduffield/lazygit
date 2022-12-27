package components

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type View struct {
	// context is prepended to any error messages e.g. 'context: "current view"'
	context string
	getView func() *gocui.View
	input   *Input
}

// asserts that the view has the expected name. This is typically used in tandem with the CurrentView method i.e.;
// input.CurrentView().Name("commits") to assert that the current view is the commits view.
func (self *View) Name(expected string) *View {
	self.input.assertWithRetries(func() (bool, string) {
		actual := self.getView().Name()
		return actual == expected, fmt.Sprintf("%s: Expected view name to be '%s', but got '%s'", self.context, expected, actual)
	})

	return self
}

// asserts that the view has the expected title
func (self *View) Title(expected *matcher) *View {
	self.input.assertWithRetries(func() (bool, string) {
		actual := self.getView().Title
		return expected.context(fmt.Sprintf("%s title", self.context)).test(actual)
	})

	return self
}

// asserts that the view has lines matching the given matchers. So if three matchers
// are passed, we only check the first three lines of the view.
// This method is convenient when you have a list of commits but you only want to
// assert on the first couple of commits.
func (self *View) TopLines(matchers ...*matcher) *View {
	self.input.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) >= len(matchers), fmt.Sprintf("unexpected number of lines in view. Expected at least %d, got %d", len(matchers), len(lines))
	})

	return self.assertLines(matchers...)
}

// asserts that the view has lines matching the given matchers. One matcher must be passed for each line.
// If you only care about the top n lines, use the TopLines method instead.
func (self *View) Lines(matchers ...*matcher) *View {
	self.input.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) == len(matchers), fmt.Sprintf("unexpected number of lines in view. Expected %d, got %d", len(matchers), len(lines))
	})

	return self.assertLines(matchers...)
}

func (self *View) assertLines(matchers ...*matcher) *View {
	view := self.getView()

	for i, matcher := range matchers {
		checkIsSelected, matcher := matcher.checkIsSelected()

		self.input.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", view.Name()),
			func() string {
				return view.BufferLines()[i]
			},
		)

		if checkIsSelected {
			self.input.assertWithRetries(func() (bool, string) {
				lineIdx := view.SelectedLineIdx()
				return lineIdx == i, fmt.Sprintf("Unexpected selected line index in view '%s'. Expected %d, got %d", view.Name(), i, lineIdx)
			})
		}
	}

	return self
}

// asserts on the content of the view i.e. the stuff within the view's frame.
func (self *View) Content(matcher *matcher) *View {
	self.input.matchString(matcher, fmt.Sprintf("%s: Unexpected content.", self.context),
		func() string {
			return self.getView().Buffer()
		},
	)

	return self
}

// asserts on the selected line of the view
func (self *View) SelectedLine(matcher *matcher) *View {
	self.input.matchString(matcher, fmt.Sprintf("%s: Unexpected selected line.", self.context),
		func() string {
			return self.getView().SelectedLine()
		},
	)

	return self
}

// asserts on the index of the selected line. 0 is the first index, representing the line at the top of the view.
func (self *View) SelectedLineIdx(expected int) *View {
	self.input.assertWithRetries(func() (bool, string) {
		actual := self.getView().SelectedLineIdx()
		return expected == actual, fmt.Sprintf("%s: Expected selected line index to be %d, got %d", self.context, expected, actual)
	})

	return self
}

// focus the view (assumes the view is a side-view that can be focused via a keybinding)
func (self *View) Focus() *View {
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
		self.input.fail(fmt.Sprintf("Cannot focus view %s: Focus() method not implemented", viewName))
	}

	self.input.press(self.input.keys.Universal.JumpToBlock[index])

	// assert that we land in the expected view
	self.IsFocused()

	return self
}

// asserts that the view is focused
func (self *View) IsFocused() *View {
	self.input.assertWithRetries(func() (bool, string) {
		expected := self.getView().Name()
		actual := self.input.gui.CurrentContext().GetView().Name()
		return actual == expected, fmt.Sprintf("%s: Unexpected view focused. Expected %s, got %s", self.context, expected, actual)
	})

	return self
}

func (self *View) Press(keyStr string) *View {
	self.IsFocused()

	self.input.press(keyStr)

	return self
}

// i.e. pressing down arrow
func (self *View) SelectNextItem() *View {
	return self.Press(self.input.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *View) SelectPreviousItem() *View {
	return self.Press(self.input.keys.Universal.PrevItem)
}

// i.e. pressing space
func (self *View) PressPrimaryAction() *View {
	return self.Press(self.input.keys.Universal.Select)
}

// i.e. pressing space
func (self *View) PressEnter() *View {
	return self.Press(self.input.keys.Universal.Confirm)
}

// i.e. pressing escape
func (self *View) PressEscape() *View {
	return self.Press(self.input.keys.Universal.Return)
}

func (self *View) NavigateToListItem(matcher *matcher) *View {
	self.IsFocused()

	self.input.navigateToListItem(matcher)

	return self
}
