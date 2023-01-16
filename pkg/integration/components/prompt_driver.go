package components

type PromptDriver struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *PromptDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Confirmation()
}

// asserts that the popup has the expected title
func (self *PromptDriver) Title(expected *matcher) *PromptDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts on the text initially present in the prompt
func (self *PromptDriver) InitialText(expected *matcher) *PromptDriver {
	self.getViewDriver().Content(expected)

	return self
}

func (self *PromptDriver) Type(value string) *PromptDriver {
	self.t.typeContent(value)

	return self
}

func (self *PromptDriver) Clear() *PromptDriver {
	panic("Clear method not yet implemented!")
}

func (self *PromptDriver) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEnter()
}

func (self *PromptDriver) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEscape()
}

func (self *PromptDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.t.Fail("You must check the title of a prompt popup by calling Title() before calling Confirm()/Cancel().")
	}
}

func (self *PromptDriver) SuggestionLines(matchers ...*matcher) *PromptDriver {
	self.t.Views().Suggestions().Lines(matchers...)

	return self
}

func (self *PromptDriver) SuggestionTopLines(matchers ...*matcher) *PromptDriver {
	self.t.Views().Suggestions().TopLines(matchers...)

	return self
}

func (self *PromptDriver) ConfirmFirstSuggestion() {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		SelectedLineIdx(0).
		PressEnter()
}

func (self *PromptDriver) ConfirmSuggestion(matcher *matcher) {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		NavigateToListItem(matcher).
		PressEnter()
}
