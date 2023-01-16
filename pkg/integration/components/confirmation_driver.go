package components

type ConfirmationDriver struct {
	t                 *TestDriver
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *ConfirmationDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Confirmation()
}

// asserts that the confirmation view has the expected title
func (self *ConfirmationDriver) Title(expected *matcher) *ConfirmationDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts that the confirmation view has the expected content
func (self *ConfirmationDriver) Content(expected *matcher) *ConfirmationDriver {
	self.getViewDriver().Content(expected)

	self.hasCheckedContent = true

	return self
}

func (self *ConfirmationDriver) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEnter()
}

func (self *ConfirmationDriver) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEscape()
}

func (self *ConfirmationDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.t.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
