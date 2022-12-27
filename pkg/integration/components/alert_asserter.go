package components

type AlertAsserter struct {
	t                 *TestDriver
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *AlertAsserter) getViewAsserter() *View {
	return self.t.Views().Confirmation()
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

	self.getViewAsserter().PressEnter()
}

func (self *AlertAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEscape()
}

func (self *AlertAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.t.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
