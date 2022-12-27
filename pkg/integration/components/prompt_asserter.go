package components

type PromptAsserter struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *PromptAsserter) getViewAsserter() *View {
	return self.t.Views().Confirmation()
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
	self.t.typeContent(value)

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
		self.t.Fail("You must check the title of a prompt popup by calling Title() before calling Confirm()/Cancel().")
	}
}

func (self *PromptAsserter) SuggestionLines(matchers ...*matcher) *PromptAsserter {
	self.t.Views().Suggestions().Lines(matchers...)

	return self
}

func (self *PromptAsserter) SuggestionTopLines(matchers ...*matcher) *PromptAsserter {
	self.t.Views().Suggestions().TopLines(matchers...)

	return self
}

func (self *PromptAsserter) ConfirmFirstSuggestion() {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		SelectedLineIdx(0).
		PressEnter()
}

func (self *PromptAsserter) ConfirmSuggestion(matcher *matcher) {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		NavigateToListItem(matcher).
		PressEnter()
}
