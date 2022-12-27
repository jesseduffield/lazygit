package components

type PromptAsserter struct {
	input           *Input
	hasCheckedTitle bool
}

func (self *PromptAsserter) getViewAsserter() *View {
	return self.input.Views().Confirmation()
}

// asserts that the popup has the expected title
func (self *PromptAsserter) Title(expected *matcher) *PromptAsserter {
	self.getViewAsserter().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts on the text initially present in the prompt
func (self *PromptAsserter) InitialText(expected *matcher) *PromptAsserter {
	self.getViewAsserter().Content(expected)

	return self
}

func (self *PromptAsserter) Type(value string) *PromptAsserter {
	self.input.typeContent(value)

	return self
}

func (self *PromptAsserter) Clear() *PromptAsserter {
	panic("Clear method not yet implemented!")
}

func (self *PromptAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEnter()
}

func (self *PromptAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEscape()
}

func (self *PromptAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.input.Fail("You must check the title of a prompt popup by calling Title() before calling Confirm()/Cancel().")
	}
}

func (self *PromptAsserter) SuggestionLines(matchers ...*matcher) *PromptAsserter {
	self.input.Views().Suggestions().Lines(matchers...)

	return self
}

func (self *PromptAsserter) SuggestionTopLines(matchers ...*matcher) *PromptAsserter {
	self.input.Views().Suggestions().TopLines(matchers...)

	return self
}

func (self *PromptAsserter) SelectFirstSuggestion() *PromptAsserter {
	self.input.press(self.input.keys.Universal.TogglePanel)
	self.input.Views().Suggestions().
		IsFocused().
		SelectedLineIdx(0)

	return self
}

func (self *PromptAsserter) SelectSuggestion(matcher *matcher) *PromptAsserter {
	self.input.press(self.input.keys.Universal.TogglePanel)
	self.input.Views().Suggestions().
		IsFocused().
		NavigateToListItem(matcher)

	return self
}
