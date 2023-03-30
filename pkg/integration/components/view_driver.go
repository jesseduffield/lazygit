package components

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
)

type ViewDriver struct {
	// context is prepended to any error messages e.g. 'context: "current view"'
	context string
	getView func() *gocui.View
	t       *TestDriver
}

func (self *ViewDriver) getSelectedLines() []string {
	view := self.t.gui.View(self.getView().Name())
	return view.SelectedLines()
}

func (self *ViewDriver) getSelectedRange() (int, int) {
	view := self.t.gui.View(self.getView().Name())
	return view.SelectedLineRange()
}

func (self *ViewDriver) getSelectedLineIdx() int {
	view := self.t.gui.View(self.getView().Name())
	return view.SelectedLineIdx()
}

// asserts that the view has the expected title
func (self *ViewDriver) Title(expected *TextMatcher) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		actual := self.getView().Title
		return expected.context(fmt.Sprintf("%s title", self.context)).test(actual)
	})

	return self
}

// asserts that the view has lines matching the given matchers. One matcher must be passed for each line.
// If you only care about the top n lines, use the TopLines method instead.
// If you only care about a subset of lines, use the ContainsLines method instead.
func (self *ViewDriver) Lines(matchers ...*TextMatcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.LineCount(EqualsInt(len(matchers)))

	return self.assertLines(0, matchers...)
}

// asserts that the view has lines matching the given matchers. So if three matchers
// are passed, we only check the first three lines of the view.
// This method is convenient when you have a list of commits but you only want to
// assert on the first couple of commits.
func (self *ViewDriver) TopLines(matchers ...*TextMatcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	return self.assertLines(0, matchers...)
}

// Asserts on the visible lines of the view.
// Note, this assumes that the view's viewport is filled with lines
func (self *ViewDriver) VisibleLines(matchers ...*TextMatcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateVisibleLineCount(matchers)

	// Get the origin of the view and offset that.
	// Note that we don't do any retrying here so if we want to bring back retry logic
	// we'll need to update this.
	originY := self.getView().OriginY()

	return self.assertLines(originY, matchers...)
}

// asserts that somewhere in the view there are consequetive lines matching the given matchers.
func (self *ViewDriver) ContainsLines(matchers ...*TextMatcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	self.t.assertWithRetries(func() (bool, string) {
		content := self.getView().Buffer()
		lines := strings.Split(content, "\n")

		startIdx, endIdx := self.getSelectedRange()

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

func (self *ViewDriver) ContainsColoredText(fgColorStr string, text string) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		view := self.getView()
		ok := self.getView().ContainsColoredText(fgColorStr, text)
		if !ok {
			return false, fmt.Sprintf("expected view '%s' to contain colored text '%s' but it didn't", view.Name(), text)
		}

		return true, ""
	})

	return self
}

func (self *ViewDriver) DoesNotContainColoredText(fgColorStr string, text string) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		view := self.getView()
		ok := !self.getView().ContainsColoredText(fgColorStr, text)
		if !ok {
			return false, fmt.Sprintf("expected view '%s' to NOT contain colored text '%s' but it didn't", view.Name(), text)
		}

		return true, ""
	})

	return self
}

// asserts on the lines that are selected in the view. Don't use the `IsSelected` matcher with this because it's redundant.
func (self *ViewDriver) SelectedLines(matchers ...*TextMatcher) *ViewDriver {
	self.validateMatchersPassed(matchers)
	self.validateEnoughLines(matchers)

	self.t.assertWithRetries(func() (bool, string) {
		selectedLines := self.getSelectedLines()

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

func (self *ViewDriver) validateMatchersPassed(matchers []*TextMatcher) {
	if len(matchers) < 1 {
		self.t.fail("'Lines' methods require at least one matcher to be passed as an argument. If you are trying to assert that there are no lines, use .IsEmpty()")
	}
}

func (self *ViewDriver) validateEnoughLines(matchers []*TextMatcher) {
	view := self.getView()

	self.t.assertWithRetries(func() (bool, string) {
		lines := view.BufferLines()
		return len(lines) >= len(matchers), fmt.Sprintf("unexpected number of lines in view '%s'. Expected at least %d, got %d", view.Name(), len(matchers), len(lines))
	})
}

// assumes the view's viewport is filled with lines
func (self *ViewDriver) validateVisibleLineCount(matchers []*TextMatcher) {
	view := self.getView()

	self.t.assertWithRetries(func() (bool, string) {
		count := view.InnerHeight() + 1
		return count == len(matchers), fmt.Sprintf("unexpected number of visible lines in view '%s'. Expected exactly %d, got %d", view.Name(), len(matchers), count)
	})
}

func (self *ViewDriver) assertLines(offset int, matchers ...*TextMatcher) *ViewDriver {
	view := self.getView()

	var expectedStartIdx, expectedEndIdx int
	foundSelectionStart := false
	foundSelectionEnd := false
	expectedSelectedLines := []string{}

	for matcherIndex, matcher := range matchers {
		lineIdx := matcherIndex + offset

		checkIsSelected, matcher := matcher.checkIsSelected()

		if checkIsSelected {
			if foundSelectionEnd {
				self.t.fail("The IsSelected matcher can only be used on a contiguous range of lines.")
			}
			if !foundSelectionStart {
				expectedStartIdx = lineIdx
				foundSelectionStart = true
			}
			expectedSelectedLines = append(expectedSelectedLines, matcher.name())
			expectedEndIdx = lineIdx
		} else if foundSelectionStart {
			foundSelectionEnd = true
		}
	}

	for matcherIndex, matcher := range matchers {
		lineIdx := matcherIndex + offset
		expectSelected, matcher := matcher.checkIsSelected()

		self.t.matchString(matcher, fmt.Sprintf("Unexpected content in view '%s'.", view.Name()),
			func() string {
				return view.BufferLines()[lineIdx]
			},
		)

		// If any of the matchers care about the selection, we need to
		// assert on the selection for each matcher.
		if foundSelectionStart {
			self.t.assertWithRetries(func() (bool, string) {
				startIdx, endIdx := self.getSelectedRange()

				selected := lineIdx >= startIdx && lineIdx <= endIdx

				if (selected && expectSelected) || (!selected && !expectSelected) {
					return true, ""
				}

				lines := self.getSelectedLines()

				return false, fmt.Sprintf(
					"Unexpected selection in view '%s'. Expected %s to be selected but got %s.\nExpected selected lines:\n---\n%s\n---\n\nActual selected lines:\n---\n%s\n---\n",
					view.Name(),
					formatLineRange(expectedStartIdx, expectedEndIdx),
					formatLineRange(startIdx, endIdx),
					strings.Join(expectedSelectedLines, "\n"),
					strings.Join(lines, "\n"),
				)
			})
		}
	}

	return self
}

func formatLineRange(from int, to int) string {
	if from == to {
		return "line " + fmt.Sprintf("%d", from)
	}

	return "lines " + fmt.Sprintf("%d-%d", from, to)
}

// asserts on the content of the view i.e. the stuff within the view's frame.
func (self *ViewDriver) Content(matcher *TextMatcher) *ViewDriver {
	self.t.matchString(matcher, fmt.Sprintf("%s: Unexpected content.", self.context),
		func() string {
			return self.getView().Buffer()
		},
	)

	return self
}

// asserts on the selected line of the view. If you are selecting a range,
// you should use the SelectedLines method instead.
func (self *ViewDriver) SelectedLine(matcher *TextMatcher) *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		selectedLineIdx := self.getSelectedLineIdx()

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
		{name: "files", viewNames: []string{"files", "worktrees", "submodules"}},
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

func (self *ViewDriver) Delay() *ViewDriver {
	self.t.Wait(self.t.inputDelay)

	return self
}

// for use when typing or navigating, because in demos we want that to happen
// faster
func (self *ViewDriver) PressFast(keyStr string) *ViewDriver {
	self.IsFocused()

	self.t.pressFast(keyStr)

	return self
}

func (self *ViewDriver) Click(x, y int) *ViewDriver {
	offsetX, offsetY, _, _ := self.getView().Dimensions()

	self.t.click(offsetX+1+x, offsetY+1+y)

	return self
}

// i.e. pressing down arrow
func (self *ViewDriver) SelectNextItem() *ViewDriver {
	return self.PressFast(self.t.keys.Universal.NextItem)
}

// i.e. pressing up arrow
func (self *ViewDriver) SelectPreviousItem() *ViewDriver {
	return self.PressFast(self.t.keys.Universal.PrevItem)
}

// i.e. pressing '<'
func (self *ViewDriver) GotoTop() *ViewDriver {
	return self.PressFast(self.t.keys.Universal.GotoTop)
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
func (self *ViewDriver) NavigateToLine(matcher *TextMatcher) *ViewDriver {
	self.IsFocused()

	view := self.getView()
	lines := view.BufferLines()

	matchIndex := -1

	self.t.assertWithRetries(func() (bool, string) {
		var matches []string
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
		}
		return true, ""
	})

	// If no match was found, it could be that this is a view that renders only
	// the visible lines. In that case, we jump to the top and then press
	// down-arrow until we found the match. We simply return the first match we
	// find, so we have no way to assert that there are no duplicates.
	if matchIndex == -1 {
		self.GotoTop()
		matchIndex = len(lines)
	}

	selectedLineIdx := self.getSelectedLineIdx()
	if selectedLineIdx == matchIndex {
		return self.SelectedLine(matcher)
	}

	// At this point we can't just take the difference of selected and matched
	// index and press up or down arrow this many times. The reason is that
	// there might be section headers between those lines, and these will be
	// skipped when pressing up or down arrow. So we must keep pressing the
	// arrow key in a loop, and check after each one whether we now reached the
	// target line.
	var maxNumKeyPresses int
	var keyPress func()
	if selectedLineIdx < matchIndex {
		maxNumKeyPresses = matchIndex - selectedLineIdx
		keyPress = func() { self.SelectNextItem() }
	} else {
		maxNumKeyPresses = selectedLineIdx - matchIndex
		keyPress = func() { self.SelectPreviousItem() }
	}

	for i := 0; i < maxNumKeyPresses; i++ {
		keyPress()
		idx := self.getSelectedLineIdx()
		// It is important to use view.BufferLines() here and not lines, because it
		// could change with every keypress.
		if ok, _ := matcher.test(view.BufferLines()[idx]); ok {
			return self
		}
	}

	self.t.fail(fmt.Sprintf("Could not navigate to item matching: %s. Lines:\n%s", matcher.name(), strings.Join(view.BufferLines(), "\n")))
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

func (self *ViewDriver) LineCount(matcher *IntMatcher) *ViewDriver {
	view := self.getView()

	self.t.assertWithRetries(func() (bool, string) {
		lineCount := self.getLineCount()
		ok, _ := matcher.test(lineCount)
		return ok, fmt.Sprintf("unexpected number of lines in view '%s'. Expected %s, got %d", view.Name(), matcher.name(), lineCount)
	})

	return self
}

func (self *ViewDriver) getLineCount() int {
	// can't rely entirely on view.BufferLines because it returns 1 even if there's nothing in the view
	if strings.TrimSpace(self.getView().Buffer()) == "" {
		return 0
	}

	view := self.getView()
	return len(view.BufferLines())
}

func (self *ViewDriver) IsVisible() *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		return self.getView().Visible, fmt.Sprintf("%s: Expected view to be visible, but it was not", self.context)
	})

	return self
}

func (self *ViewDriver) IsInvisible() *ViewDriver {
	self.t.assertWithRetries(func() (bool, string) {
		return !self.getView().Visible, fmt.Sprintf("%s: Expected view to be invisible, but it was not", self.context)
	})

	return self
}

// will filter or search depending on whether the view supports filtering/searching
func (self *ViewDriver) FilterOrSearch(text string) *ViewDriver {
	self.IsFocused()

	self.Press(self.t.keys.Universal.StartSearch).
		Tap(func() {
			self.t.ExpectSearch().
				Clear().
				Type(text).
				Confirm()

			self.t.Views().Search().IsVisible().Content(Contains(fmt.Sprintf("matches for '%s'", text)))
		})

	return self
}

func (self *ViewDriver) SetCaption(caption string) *ViewDriver {
	self.t.gui.SetCaption(caption)

	return self
}

func (self *ViewDriver) SetCaptionPrefix(prefix string) *ViewDriver {
	self.t.gui.SetCaptionPrefix(prefix)

	return self
}

func (self *ViewDriver) Wait(milliseconds int) *ViewDriver {
	if !self.t.gui.Headless() {
		self.t.Wait(milliseconds)
	}

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

func expectedContentFromMatchers(matchers []*TextMatcher) string {
	return strings.Join(lo.Map(matchers, func(matcher *TextMatcher, _ int) string {
		return matcher.name()
	}), "\n")
}
