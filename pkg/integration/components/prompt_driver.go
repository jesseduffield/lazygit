package components

type PromptDriver struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *PromptDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Confirmation()
}

// asserts that the popup has the expected title
func (self *PromptDriver) Title(expected *TextMatcher) *PromptDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts on the text initially present in the prompt
func (self *PromptDriver) InitialText(expected *TextMatcher) *PromptDriver {
	self.getViewDriver().Content(expected)

	return self
}

func (self *PromptDriver) Type(value string) *PromptDriver {
	self.t.typeContent(value)

	return self
}

func (self *PromptDriver) Clear() *PromptDriver {
	self.t.press(ClearKey)

	return self
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

func (self *PromptDriver) SuggestionLines(matchers ...*TextMatcher) *PromptDriver {
	self.t.Views().Suggestions().Lines(matchers...)

	return self
}

func (self *PromptDriver) SuggestionTopLines(matchers ...*TextMatcher) *PromptDriver {
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

func (self *PromptDriver) ConfirmSuggestion(matcher *TextMatcher) {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		NavigateToLine(matcher).
		PressEnter()
}

func (self *PromptDriver) DeleteSuggestion(matcher *TextMatcher) *PromptDriver {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		NavigateToLine(matcher)
	self.t.press(self.t.keys.Universal.Remove)
	return self
}

func (self *PromptDriver) EditSuggestion(matcher *TextMatcher) *PromptDriver {
	self.t.press(self.t.keys.Universal.TogglePanel)
	self.t.Views().Suggestions().
		IsFocused().
		NavigateToLine(matcher)
	self.t.press(self.t.keys.Universal.Edit)
	return self
}
