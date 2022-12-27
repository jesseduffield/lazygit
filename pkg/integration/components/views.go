package components

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type Views struct {
	// context is prepended to any error messages e.g. 'context: "current view"'
	context string
	getView func() *gocui.View
	assert  *Assert
}

// asserts that the view has the expected name. This is typically used in tandem with the CurrentView method i.e.;
// assert.CurrentView().Name("commits") to assert that the current view is the commits view.
func (self *Views) Name(expected string) *Views {
	self.assert.assertWithRetries(func() (bool, string) {
		actual := self.getView().Name()
		return actual == expected, fmt.Sprintf("%s: Expected view name to be '%s', but got '%s'", self.context, expected, actual)
	})

	return self
}

// asserts that the view has the expected title
func (self *Views) Title(expected *matcher) *Views {
	self.assert.assertWithRetries(func() (bool, string) {
		actual := self.getView().Title
		return expected.context(fmt.Sprintf("%s title", self.context)).test(actual)
	})

	return self
}

// asserts that the view has lines matching the given matchers. So if three matchers
// are passed, we only check the first three lines of the view.
// This method is convenient when you have a list of commits but you only want to
// assert on the first couple of commits.
func (self *Views) TopLines(matchers ...*matcher) *Views {
	self.assert.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) >= len(matchers), fmt.Sprintf("unexpected number of lines in view. Expected at least %d, got %d", len(matchers), len(lines))
	})

	return self.assertLines(matchers...)
}

// asserts that the view has lines matching the given matchers. One matcher must be passed for each line.
// If you only care about the top n lines, use the TopLines method instead.
func (self *Views) Lines(matchers ...*matcher) *Views {
	self.assert.assertWithRetries(func() (bool, string) {
		lines := self.getView().BufferLines()
		return len(lines) == len(matchers), fmt.Sprintf("unexpected number of lines in view. Expected %d, got %d", len(matchers), len(lines))
	})

	return self.assertLines(matchers...)
}

func (self *Views) assertLines(matchers ...*matcher) *Views {
	view := self.getView()

	for i, matcher := range matchers {
		checkIsSelected, matcher := matcher.checkIsSelected()

		self.assert.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", view.Name()),
			func() string {
				return view.BufferLines()[i]
			},
		)

		if checkIsSelected {
			self.assert.assertWithRetries(func() (bool, string) {
				lineIdx := view.SelectedLineIdx()
				return lineIdx == i, fmt.Sprintf("Unexpected selected line index in view '%s'. Expected %d, got %d", view.Name(), i, lineIdx)
			})
		}
	}

	return self
}

// asserts on the content of the view i.e. the stuff within the view's frame.
func (self *Views) Content(matcher *matcher) *Views {
	self.assert.matchString(matcher, fmt.Sprintf("%s: Unexpected content.", self.context),
		func() string {
			return self.getView().Buffer()
		},
	)

	return self
}

// asserts on the selected line of the view
func (self *Views) SelectedLine(matcher *matcher) *Views {
	self.assert.matchString(matcher, fmt.Sprintf("%s: Unexpected selected line.", self.context),
		func() string {
			return self.getView().SelectedLine()
		},
	)

	return self
}

// asserts on the index of the selected line. 0 is the first index, representing the line at the top of the view.
func (self *Views) SelectedLineIdx(expected int) *Views {
	self.assert.assertWithRetries(func() (bool, string) {
		actual := self.getView().SelectedLineIdx()
		return expected == actual, fmt.Sprintf("%s: Expected selected line index to be %d, got %d", self.context, expected, actual)
	})

	return self
}

type ViewAsserterGetter struct {
	assert *Assert
}

func (self *ViewAsserterGetter) Current() *Views {
	return &Views{
		context: "current view",
		getView: func() *gocui.View { return self.assert.gui.CurrentContext().GetView() },
		assert:  self.assert,
	}
}

func (self *ViewAsserterGetter) Main() *Views {
	return &Views{
		context: "main view",
		getView: func() *gocui.View { return self.assert.gui.MainView() },
		assert:  self.assert,
	}
}

func (self *ViewAsserterGetter) Secondary() *Views {
	return &Views{
		context: "secondary view",
		getView: func() *gocui.View { return self.assert.gui.SecondaryView() },
		assert:  self.assert,
	}
}

func (self *ViewAsserterGetter) ByName(viewName string) *Views {
	return &Views{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.assert.gui.View(viewName) },
		assert:  self.assert,
	}
}
