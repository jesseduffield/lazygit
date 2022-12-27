package components

type ConfirmationAsserter struct {
	assert            *Assert
	input             *Input
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *ConfirmationAsserter) getViewAsserter() *ViewAsserter {
	return self.assert.Views().ByName("confirmation")
}

// asserts that the confirmation view has the expected title
func (self *ConfirmationAsserter) Title(expected *matcher) *ConfirmationAsserter {
	self.getViewAsserter().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts that the confirmation view has the expected content
func (self *ConfirmationAsserter) Content(expected *matcher) *ConfirmationAsserter {
	self.getViewAsserter().Content(expected)

	self.hasCheckedContent = true

	return self
}

func (self *ConfirmationAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.input.Confirm()
}

func (self *ConfirmationAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.input.Press(self.input.keys.Universal.Return)
}

func (self *ConfirmationAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.assert.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
