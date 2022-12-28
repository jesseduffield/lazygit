package components

type ConfirmationAsserter struct {
	t                 *TestDriver
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *ConfirmationAsserter) getViewAsserter() *View {
	return self.t.Views().Confirmation()
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

	self.getViewAsserter().PressEnter()
}

func (self *ConfirmationAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEscape()
}

func (self *ConfirmationAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.t.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
