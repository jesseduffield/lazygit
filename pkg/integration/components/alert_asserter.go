package components

type AlertAsserter struct {
	assert            *Assert
	input             *Input
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *AlertAsserter) getViewAsserter() *ViewAsserter {
	return self.assert.Views().ByName("confirmation")
}

// asserts that the alert view has the expected title
func (self *AlertAsserter) Title(expected *matcher) *AlertAsserter {
	self.getViewAsserter().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts that the alert view has the expected content
func (self *AlertAsserter) Content(expected *matcher) *AlertAsserter {
	self.getViewAsserter().Content(expected)

	self.hasCheckedContent = true

	return self
}

func (self *AlertAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.input.Confirm()
}

func (self *AlertAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.input.Press(self.input.keys.Universal.Return)
}

func (self *AlertAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.assert.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
