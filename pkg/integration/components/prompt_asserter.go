package components

type PromptAsserter struct {
	assert          *Assert
	input           *Input
	hasCheckedTitle bool
}

func (self *PromptAsserter) getViewAsserter() *ViewAsserter {
	return self.assert.View("confirmation")
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

func (self *PromptAsserter) Confirm() *PromptAsserter {
	self.checkNecessaryChecksCompleted()

	self.input.Confirm()

	return self
}

func (self *PromptAsserter) Cancel() *PromptAsserter {
	self.checkNecessaryChecksCompleted()

	self.input.Press(self.input.keys.Universal.Return)

	return self
}

func (self *PromptAsserter) Type(value string) *PromptAsserter {
	self.input.Type(value)

	return self
}

func (self *PromptAsserter) Clear() *PromptAsserter {
	panic("Clear method not yet implemented!")
}

func (self *PromptAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.assert.Fail("You must check the title of a prompt popup by calling Title() before calling Confirm()/Cancel().")
	}
}
