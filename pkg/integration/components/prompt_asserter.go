package components

type PromptAsserter struct {
	assert          *Assert
	input           *Input
	hasCheckedTitle bool
}

func (self *PromptAsserter) getViewAsserter() *ViewAsserter {
	return self.assert.Views().ByName("confirmation")
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
	self.input.Type(value)

	return self
}

func (self *PromptAsserter) Clear() *PromptAsserter {
	panic("Clear method not yet implemented!")
}

func (self *PromptAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.input.Confirm()
}

func (self *PromptAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.input.Press(self.input.keys.Universal.Return)
}

func (self *PromptAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.assert.Fail("You must check the title of a prompt popup by calling Title() before calling Confirm()/Cancel().")
	}
}

func (self *PromptAsserter) SuggestionLines(matchers ...*matcher) *PromptAsserter {
	self.assert.Views().ByName("suggestions").Lines(matchers...)

	return self
}

func (self *PromptAsserter) SuggestionTopLines(matchers ...*matcher) *PromptAsserter {
	self.assert.Views().ByName("suggestions").TopLines(matchers...)

	return self
}

func (self *PromptAsserter) SelectFirstSuggestion() *PromptAsserter {
	self.input.Press(self.input.keys.Universal.TogglePanel)
	self.assert.Views().Current().Name("suggestions")

	return self
}

func (self *PromptAsserter) SelectSuggestion(matcher *matcher) *PromptAsserter {
	self.input.Press(self.input.keys.Universal.TogglePanel)
	self.assert.Views().Current().Name("suggestions")

	self.input.NavigateToListItem(matcher)

	return self
}
