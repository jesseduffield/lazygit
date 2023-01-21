package components

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
)

type ViewDriver struct {
	// context is prepended to any error messages e.g. 'context: "current view"'
	context              string
	getView              func() *gocui.View
	t                    *TestDriver
	getSelectedLinesFn   func() ([]string, error)
	getSelectedRangeFn   func() (int, int, error)
	getSelectedLineIdxFn func() (int, error)
}

func (self *ViewDriver) getSelectedLines() ([]string, error) {
	if self.getSelectedLinesFn == nil {
		view := self.t.gui.View(self.getView().Name())

		return []string{view.SelectedLine()}, nil
	}

	return self.getSelectedLinesFn()
}

func (self *ViewDriver) getSelectedRange() (int, int, error) {
	if self.getSelectedRangeFn == nil {
		view := self.t.gui.View(self.getView().Name())
		idx := view.SelectedLineIdx()

		return idx, idx, nil
	}

	return self.getSelectedRangeFn()
}

// even if you have a selected range, there may still be a line within that range
// which the cursor points at. This function returns that line index.
func (self *ViewDriver) getSelectedLineIdx() (int, error) {
	if self.getSelectedLineIdxFn == nil {
		view := self.t.gui.View(self.getView().Name())

		return view.SelectedLineIdx(), nil
	}

	return self.getSelectedLineIdxFn()
}

// asserts that the view has the expected title
func (self *ViewDriver) Title(expected *Matcher) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		actual := self.getView().Title
		return expected.context(fmt.Sprintf("%s title", self.context)).test(actual)
	})

	return self
}

// asserts that the view has lines matching the given matchers. One matcher must be passed for each line.
// If you only care about the top n lines, use the TopLines method instead.
// If you only care about a subset of lines, use the ContainsLines method instead.
func (self *ViewDriver) Lines(matchers ...*Matcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.LineCount(len(matchers))

	return self.assertLines(0, matchers...)
}

// asserts that the view has lines matching the given matchers. So if three matchers
// are passed, we only check the first three lines of the view.
// This method is convenient when you have a list of commits but you only want to
// assert on the first couple of commits.
func (self *ViewDriver) TopLines(matchers ...*Matcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	return self.assertLines(0, matchers...)
}

// asserts that somewhere in the view there are consequetive lines matching the given matchers.
func (self *ViewDriver) ContainsLines(matchers ...*Matcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	self.t.assertWithRetries(func() (bool, string) {
		content := self.getView().Buffer()
		lines := strings.Split(content, "\n")

		startIdx, endIdx, err := self.getSelectedRange()

		for i := 0; i < len(lines)-len(matchers)+1; i++ {
			matches := true
			for j, matcher := range matchers {
				checkIsSelected, matcher := matcher.checkIsSelected() // strip the IsSelected matcher out
				lineIdx := i + j
				ok, _ := matcher.test(lines[lineIdx])
				if !ok {
					matches = false
					break
				}
				if checkIsSelected {
					if err != nil {
						matches = false
						break
					}
					if lineIdx < startIdx || lineIdx > endIdx {
						matches = false
						break
					}
				}
			}
			if matches {
				return true, ""
			}
		}

		expectedContent := expectedContentFromMatchers(matchers)

		return false, fmt.Sprintf(
			"Expected the following to be contained in the staging panel:\n-----\n%s\n-----\nBut got:\n-----\n%s\n-----\nSelected range: %d-%d",
			expectedContent,
			content,
			startIdx,
			endIdx,
		)
	})

	return self
}

// asserts on the lines that are selected in the view. Don't use the `IsSelected` matcher with this because it's redundant.
func (self *ViewDriver) SelectedLines(matchers ...*Matcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	self.t.assertWithRetries(func() (bool, string) {
		selectedLines, err := self.getSelectedLines()
		if err != nil {
			return false, err.Error()
		}

		selectedContent := strings.Join(selectedLines, "\n")
		expectedContent := expectedContentFromMatchers(matchers)

		if len(selectedLines) != len(matchers) {
			return false, fmt.Sprintf("Expected the following to be selected:\n-----\n%s\n-----\nBut got:\n-----\n%s\n-----", expectedContent, selectedContent)
		}

		for i, line := range selectedLines {
			checkIsSelected, matcher := matchers[i].checkIsSelected()
			if checkIsSelected {
				self.t.fail("You cannot use the IsSelected matcher with the SelectedLines method")
			}

			ok, message := matcher.test(line)
			if !ok {
				return false, fmt.Sprintf("Error: %s. Expected the following to be selected:\n-----\n%s\n-----\nBut got:\n-----\n%s\n-----", message, expectedContent, selectedContent)
			}
		}

		return true, ""
	})

	return self
}

func (self *ViewDriver) validateMatchersPassed(matchers []*Matcher) {
	if len(matchers) < 1 {
		self.t.fail("'Lines' methods require at least one matcher to be passed as an argument. If you are trying to assert that there are no lines, use .IsEmpty()")
	}
}

func (self *ViewDriver) validateEnoughLines(matchers []*Matcher) {
	view := self.getView()

	self.t.assertWithRetries(func() (bool, string) {
		lines := view.BufferLines()
		return len(lines) >= len(matchers), fmt.Sprintf("unexpected number of lines in view '%s'. Expected at least %d, got %d", view.Name(), len(matchers), len(lines))
	})
}

func (self *ViewDriver) assertLines(offset int, matchers ...*Matcher) *ViewDriver {
	view := self.getView()

	for matcherIndex, matcher := range matchers {
		lineIdx := matcherIndex + offset
		checkIsSelected, matcher := matcher.checkIsSelected()

		self.t.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", view.Name()),
			func() string {
				return view.BufferLines()[lineIdx]
			},
		)

		if checkIsSelected {
			self.t.assertWithRetries(func() (bool, string) {
				startIdx, endIdx, err := self.getSelectedRange()
				if err != nil {
					return false, err.Error()
				}

				if lineIdx < startIdx || lineIdx > endIdx {
					if startIdx == endIdx {
						return false, fmt.Sprintf("Unexpected selected line index in view '%s'. Expected %d, got %d", view.Name(), lineIdx, startIdx)
					} else {
						lines, err := self.getSelectedLines()
						if err != nil {
							return false, err.Error()
						}
						return false, fmt.Sprintf("Unexpected selected line index in view '%s'. Expected line %d to be in range %d to %d. Selected lines:\n---\n%s\n---\n\nExpected line: '%s'", view.Name(), lineIdx, startIdx, endIdx, strings.Join(lines, "\n"), matcher.name())
					}
				}
				return true, ""
			})
		}
	}

	return self
}

// asserts on the content of the view i.e. the stuff within the view's frame.
func (self *ViewDriver) Content(matcher *Matcher) *ViewDriver {
	self.t.matchString(matcher, fmt.Sprintf("%s: Unexpected content.", self.context),
		func() string {
			return self.getView().Buffer()
		},
	)

	return self
}

// asserts on the selected line of the view. If your view has multiple lines selected,
// but also has a concept of a cursor position, this will assert on the line that
// the cursor is on. Otherwise it will assert on the first line of the selection.
func (self *ViewDriver) SelectedLine(matcher *Matcher) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		selectedLineIdx, err := self.getSelectedLineIdx()
		if err != nil {
			return false, err.Error()
		}

		viewLines := self.getView().BufferLines()

		if selectedLineIdx >= len(viewLines) {
			return false, fmt.Sprintf("%s: Expected view to have at least %d lines, but it only has %d", self.context, selectedLineIdx+1, len(viewLines))
		}

		value := viewLines[selectedLineIdx]

		return matcher.context(fmt.Sprintf("%s: Unexpected selected line.", self.context)).test(value)
	})

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

// focus the view (assumes the view is a side-view)
func (self *ViewDriver) Focus() *ViewDriver {
	viewName := self.getView().Name()

	type window struct {
		name      string
		viewNames []string
	}
	windows := []window{
		{name: "status", viewNames: []string{"status"}},
		{name: "files", viewNames: []string{"files", "submodules"}},
		{name: "branches", viewNames: []string{"localBranches", "remotes", "tags"}},
		{name: "commits", viewNames: []string{"commits", "reflogCommits"}},
		{name: "stash", viewNames: []string{"stash"}},
	}

	for windowIndex, window := range windows {
		if lo.Contains(window.viewNames, viewName) {
			tabIndex := lo.IndexOf(window.viewNames, viewName)
			// jump to the desired window
			self.t.press(self.t.keys.Universal.JumpToBlock[windowIndex])

			// assert we're in the window before continuing
			self.t.assertWithRetries(func() (bool, string) {
				currentWindowName := self.t.gui.CurrentContext().GetWindowName()
				// by convention the window is named after the first view in the window
				return currentWindowName == window.name, fmt.Sprintf("Expected to be in window '%s', but was in '%s'", window.name, currentWindowName)
			})

			// switch to the desired tab
			currentViewName := self.t.gui.CurrentContext().GetViewName()
			currentViewTabIndex := lo.IndexOf(window.viewNames, currentViewName)
			if tabIndex > currentViewTabIndex {
				for i := 0; i < tabIndex-currentViewTabIndex; i++ {
					self.t.press(self.t.keys.Universal.NextTab)
				}
			} else if tabIndex < currentViewTabIndex {
				for i := 0; i < currentViewTabIndex-tabIndex; i++ {
					self.t.press(self.t.keys.Universal.PrevTab)
				}
			}

			// assert that we're now in the expected view
			self.IsFocused()

			return self
		}
	}

	self.t.fail(fmt.Sprintf("Cannot focus view %s: Focus() method not implemented", viewName))

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

// i.e. pressing tab
func (self *ViewDriver) PressTab() *ViewDriver {
	return self.Press(self.t.keys.Universal.TogglePanel)
}

// i.e. pressing escape
func (self *ViewDriver) PressEscape() *ViewDriver {
	return self.Press(self.t.keys.Universal.Return)
}

// this will look for a list item in the current panel and if it finds it, it will
// enter the keypresses required to navigate to it.
// The test will fail if:
// - the user is not in a list item
// - no list item is found containing the given text
// - multiple list items are found containing the given text in the initial page of items
//
// NOTE: this currently assumes that BufferLines returns all the lines that can be accessed.
// If this changes in future, we'll need to update this code to first attempt to find the item
// in the current page and failing that, jump to the top of the view and iterate through all of it,
// looking for the item.
func (self *ViewDriver) NavigateToLine(matcher *Matcher) *ViewDriver {
	self.IsFocused()

	view := self.getView()

	var matchIndex int

	self.t.assertWithRetries(func() (bool, string) {
		matchIndex = -1
		var matches []string
		lines := view.BufferLines()
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

	selectedLineIdx, err := self.getSelectedLineIdx()
	if err != nil {
		self.t.fail(err.Error())
		return self
	}
	if selectedLineIdx == matchIndex {
		self.SelectedLine(matcher)
	} else if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.SelectNextItem()
		}
		self.SelectedLine(matcher)
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.SelectPreviousItem()
		}
		self.SelectedLine(matcher)
	}

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

	view := self.getView()

	self.t.assertWithRetries(func() (bool, string) {
		lines := view.BufferLines()
		return len(lines) == expectedCount, fmt.Sprintf("unexpected number of lines in view '%s'. Expected %d, got %d", view.Name(), expectedCount, len(lines))
	})

	self.t.assertWithRetries(func() (bool, string) {
		lines := view.BufferLines()

		// if the view has a single blank line (often the case) we want to treat that as having no lines
		if len(lines) == 1 && expectedCount == 1 {
			actual := strings.TrimSpace(view.Buffer())
			return actual != "", fmt.Sprintf("unexpected number of lines in view '%s'. Expected 1, got 0", view.Name())
		}

		return len(lines) == expectedCount, fmt.Sprintf("unexpected number of lines in view '%s'. Expected %d, got %d", view.Name(), expectedCount, len(lines))
	})

	return self
}

// for when you want to make some assertion unrelated to the current view
// without breaking the method chain
func (self *ViewDriver) Tap(f func()) *ViewDriver {
	f()

	return self
}

// This purely exists as a convenience method for those who hate the trailing periods in multi-line method chains
func (self *ViewDriver) Self() *ViewDriver {
	return self
}

func expectedContentFromMatchers(matchers []*Matcher) string {
	return strings.Join(lo.Map(matchers, func(matcher *Matcher, _ int) string {
		return matcher.name()
	}), "\n")
}
